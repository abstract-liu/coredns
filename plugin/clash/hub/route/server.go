package route

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net"
	"net/http"
)

var log = clog.NewWithPlugin(constant.PluginName)

func router() *chi.Mux {
	r := chi.NewRouter()
	corsM := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:         300,
	})
	r.Use(corsM.Handler)
	r.Group(func(r chi.Router) {
		r.Mount("/config", configRouter())
	})
	return r
}

func Start(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("External controller listen error: %s", err)
		return
	}
	serverAddr := l.Addr().String()
	log.Infof("RESTful API listening at: %s", serverAddr)

	if err = http.Serve(l, router()); err != nil {
		log.Errorf("External controller serve error: %s", err)
	}
}
