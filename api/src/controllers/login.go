package controllers

import (
	"api/src/autenticacao"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type UsuarioController struct {
	Repositorio repositorios.UsuarioRepositorio
}

func NovoUsuarioController(repositorio repositorios.UsuarioRepositorio) *UsuarioController {
	return &UsuarioController{Repositorio: repositorio}
}

func (uc *UsuarioController) Login(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario

	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioSalvoNoBanco, erro := uc.Repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusNotFound, erro)
		return
	}

	if erro = seguranca.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); erro != nil {
		if erro != nil {
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}
	}

	token, erro := autenticacao.CriarToken(usuarioSalvoNoBanco.ID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	w.Write([]byte(token))
}
