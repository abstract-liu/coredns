package resource

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/samber/lo"
	"io"
	"net/http"
	"time"
)

var log = clog.NewWithPlugin(constant.PluginName)

type Parser[V any] func([]byte) (V, error)

type Fetcher[V any] struct {
	name     string
	url      string
	interval time.Duration
	done     chan struct{}
	hash     [16]byte
	parser   Parser[V]

	httpClient *http.Client

	UpdatedAt time.Time
	OnUpdate  func(V)
}

func (f *Fetcher[V]) Name() string {
	return f.name
}

func (f *Fetcher[V]) Destroy() error {
	if f.interval > 0 {
		f.done <- struct{}{}
	}
	return nil
}

func (f *Fetcher[V]) Update() (V, bool, error) {
	buf, err := f.read()
	if err != nil {
		return lo.Empty[V](), false, err
	}

	now := time.Now()
	hash := md5.Sum(buf)
	if bytes.Equal(f.hash[:], hash[:]) {
		f.UpdatedAt = now
		return lo.Empty[V](), true, nil
	}

	contents, err := f.parser(buf)
	if err != nil {
		return lo.Empty[V](), false, err
	}

	f.UpdatedAt = now
	f.hash = hash

	return contents, false, nil
}

func (f *Fetcher[V]) read() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	request, _ := http.NewRequest(http.MethodGet, f.url, nil)
	resp, err := f.httpClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.New(resp.Status)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (f *Fetcher[V]) pullLoop() {
	timer := time.NewTimer(f.interval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			timer.Reset(f.interval)
			elm, same, err := f.Update()
			if err != nil {
				log.Errorf("[Fetcher] %s pull error: %s", f.Name(), err.Error())
				continue
			}

			if same {
				log.Debugf("[Fetcher] %s's content doesn't change", f.Name())
				continue
			}

			log.Infof("[Fetcher] %s's content update", f.Name())
			if f.OnUpdate != nil {
				f.OnUpdate(elm)
			}
		case <-f.done:
			return
		}
	}
}

func (f *Fetcher[V]) Initial() (V, error) {
	var (
		preReadBuf []byte
		err        error
	)

	var contents V
	if preReadBuf, err = f.read(); err == nil {
		contents, err = f.parser(preReadBuf)
	}

	if err != nil {
		return lo.Empty[V](), err
	}

	f.hash = md5.Sum(preReadBuf)
	if f.interval > 0 {
		go f.pullLoop()
	}

	return contents, nil
}

func NewFetcher[V any](name, url string, interval time.Duration, parser Parser[V], onUpdate func(V)) *Fetcher[V] {
	return &Fetcher[V]{
		name:       name,
		url:        url,
		interval:   interval,
		done:       make(chan struct{}),
		parser:     parser,
		OnUpdate:   onUpdate,
		httpClient: &http.Client{},
	}
}
