package rotas

import (
	"api/src/controllers"
	"net/http"
)

func rotasPublicacoes(publicacoesController *controllers.PublicacoesController) []Rota {
	return []Rota{
		{
			URI:                "/publicacoes",
			Metodo:             http.MethodPost,
			Funcao:             publicacoesController.CriarPublicacao,
			RequerAutenticacao: true,
		},
		{
			URI:                "/publicacoes",
			Metodo:             http.MethodGet,
			Funcao:             publicacoesController.BuscarPublicacoes,
			RequerAutenticacao: true,
		},
		{
			URI:                "/publicacoes/{publicacaoId}",
			Metodo:             http.MethodGet,
			Funcao:             publicacoesController.BuscarPublicacao,
			RequerAutenticacao: true,
		},
		{
			URI:                "/publicacoes/{publicacaoId}",
			Metodo:             http.MethodPut,
			Funcao:             publicacoesController.AtualizarPublicacao,
			RequerAutenticacao: true,
		},
		{
			URI:                "/publicacoes/{publicacaoId}",
			Metodo:             http.MethodDelete,
			Funcao:             publicacoesController.DeletarPublicacao,
			RequerAutenticacao: true,
		},
		{
			URI:                "/usuarios/{usuarioId}/publicacoes",
			Metodo:             http.MethodGet,
			Funcao:             publicacoesController.BuscarPublicacoesPorUsuario,
			RequerAutenticacao: true,
		},
		{
			URI:                "/publicacoes/{publicacaoId}/curtir",
			Metodo:             http.MethodPost,
			Funcao:             publicacoesController.CurtirPublicacao,
			RequerAutenticacao: true,
		},
		{
			URI:                "/publicacoes/{publicacaoId}/descurtir",
			Metodo:             http.MethodPost,
			Funcao:             publicacoesController.DescurtirPublicacao,
			RequerAutenticacao: true,
		},
	}
}
