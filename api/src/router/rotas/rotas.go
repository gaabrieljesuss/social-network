package rotas

import (
	"api/src/controllers"
	"api/src/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"github.com/gorilla/mux"
)

type Rota struct {
	URI                string
	Metodo             string
	Funcao             http.HandlerFunc
	RequerAutenticacao bool
}

func Configurar(r *mux.Router, usuarioController *controllers.UsuarioController, publicacoesController *controllers.PublicacoesController) *mux.Router {
	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	rotas := rotasPublicacoes(publicacoesController)
	rotas = append(rotas, rotaLogin(usuarioController))

	for _, rota := range rotas {
		handler := middlewares.PrometheusMiddleware(middlewares.Logger(rota.Funcao))

		if rota.RequerAutenticacao {
			handler = middlewares.PrometheusMiddleware(middlewares.Logger(middlewares.Autenticar(rota.Funcao)))
		}

		r.HandleFunc(rota.URI, handler).Methods(rota.Metodo)
	}

	return r
}
