package route

import (
	"github.com/coredns/coredns/plugin/clash/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

func configRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/UpdateRemoteConfig", updateRemoteConfig)
	return r
}

func updateRemoteConfig(w http.ResponseWriter, r *http.Request) {
	if err := config.UpdateRemoteClashConfig(); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, newError(err.Error()))
		return
	}
	render.NoContent(w, r)
}
