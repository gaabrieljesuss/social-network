package rotas

import (
	"api/src/controllers"
	"net/http"
)

func rotaLogin(usuarioController *controllers.UsuarioController) Rota {
	return Rota{
		URI:                "/login",
		Metodo:             http.MethodPost,
		Funcao:             usuarioController.Login,
		RequerAutenticacao: false,
	}
}
