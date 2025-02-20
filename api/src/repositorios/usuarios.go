package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type UsuarioRepositorio interface {
	BuscarPorEmail(email string) (modelos.Usuario, error)
}

type usuarioRepositorio struct {
	db *sql.DB
}

func NovoRepositorioDeUsuarios(db *sql.DB) UsuarioRepositorio {
	return &usuarioRepositorio{db}
}

func (repositorio *usuarioRepositorio) BuscarPorEmail(email string) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"SELECT id, senha FROM usuarios WHERE email = $1", email,
	)

	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario modelos.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Senha,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}
