package rotas

import (
	"api/src/controllers"
	"api/src/middlewares"
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
	rotas := rotasPublicacoes(publicacoesController)
	rotas = append(rotas, rotaLogin(usuarioController))

	for _, rota := range rotas {
		if rota.RequerAutenticacao {
			r.HandleFunc(rota.URI, middlewares.Logger(middlewares.Autenticar(rota.Funcao))).Methods(rota.Metodo)
		} else {
			r.HandleFunc(rota.URI, middlewares.Logger(rota.Funcao)).Methods(rota.Metodo)
		}
	}

	return r
}
