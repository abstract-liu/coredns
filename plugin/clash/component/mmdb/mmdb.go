package mmdb

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/oschwald/maxminddb-golang"
	"io"
	"net/http"
	"os"
	"sync"
)

type databaseType = uint8

const (
	typeMaxmind databaseType = iota
	typeSing
	typeMetaV0
)

var (
	log      = clog.NewWithPlugin(constant.PluginName)
	IPreader IPReader
	IPonce   sync.Once
	ASNonce  sync.Once
)

func DownloadMMDB() (err error) {
	resp, err := http.Get(constant.MMDB_URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(constant.MMDB_PATH, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)

	return err
}

func IPInstance() IPReader {
	IPonce.Do(func() {
		mmdbPath := constant.MMDB_PATH
		log.Infof("Load MMDB file: %s", mmdbPath)
		mmdb, err := maxminddb.Open(mmdbPath)
		if err != nil {
			log.Errorf("Can't load MMDB: %s", err.Error())
		}
		IPreader = IPReader{Reader: mmdb}
		switch mmdb.Metadata.DatabaseType {
		case "sing-geoip":
			IPreader.databaseType = typeSing
		case "Meta-geoip0":
			IPreader.databaseType = typeMetaV0
		default:
			IPreader.databaseType = typeMaxmind
		}
	})

	return IPreader
}
