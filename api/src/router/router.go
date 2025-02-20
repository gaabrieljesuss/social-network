package router

import (
	"api/src/controllers"
	"api/src/router/rotas"
	"github.com/gorilla/mux"
)

func Gerar(usuarioController *controllers.UsuarioController, publicacoesController *controllers.PublicacoesController) *mux.Router {
	r := mux.NewRouter()
	return rotas.Configurar(r, usuarioController, publicacoesController)
}
