package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

type PublicacoesRepositorio interface {
	Criar(publicacao modelos.Publicacao) (uint64, error)
	BuscarPorId(publicacaoID uint64) (modelos.Publicacao, error)
	BuscarPublicacoes(usuarioID uint64) ([]modelos.Publicacao, error)
	Atualizar(publicacaoID uint64, publicacao modelos.Publicacao) error
	DeletarPublicacao(publicacaoID uint64) error
	BuscarPorUsuario(usuarioID uint64) ([]modelos.Publicacao, error)
	Curtir(publicacaoID uint64) error
	Descurtir(publicacaoID uint64) error
}

type publicacoesRepositorio struct {
	db *sql.DB
}

func NovoRepositorioDePublicacoes(db *sql.DB) PublicacoesRepositorio {
	return &publicacoesRepositorio{db}
}

func (repositorio *publicacoesRepositorio) Criar(publicacao modelos.Publicacao) (uint64, error) {
	query := "INSERT INTO publicacoes (titulo, conteudo, autor_id) VALUES ($1, $2, $3) RETURNING id"
	var ultimoIDInserido uint64

	erro := repositorio.db.QueryRow(query, publicacao.Titulo, publicacao.Conteudo, publicacao.AutorID).Scan(&ultimoIDInserido)
	if erro != nil {
		return 0, erro
	}

	return ultimoIDInserido, nil
}

func (repositorio *publicacoesRepositorio) BuscarPorId(publicacaoID uint64) (modelos.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT p.*, u.nick FROM 
		publicacoes p INNER JOIN usuarios u 
		ON u.id = p.autor_id WHERE p.id = $1
	`, publicacaoID)

	if erro != nil {
		return modelos.Publicacao{}, erro
	}
	defer linhas.Close()

	var publicacao modelos.Publicacao

	if linhas.Next() {
		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadoEm,
			&publicacao.AutorNick,
		); erro != nil {
			return modelos.Publicacao{}, erro
		}
	}

	return publicacao, nil
}

func (repositorio *publicacoesRepositorio) BuscarPublicacoes(usuarioID uint64) ([]modelos.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
	SELECT DISTINCT p.*, u.nick FROM publicacoes p 
	INNER JOIN usuarios u ON u.id = p.autor_id 
	INNER JOIN seguidores s ON p.autor_id = s.usuario_id 
	WHERE u.id = $1 or s.seguidor_id = $2
	ORDER BY 1 DESC`,
		usuarioID, usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var publicacoes []modelos.Publicacao

	for linhas.Next() {
		var publicacao modelos.Publicacao
		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadoEm,
			&publicacao.AutorNick,
		); erro != nil {
			return nil, erro
		}
		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func (repositorio *publicacoesRepositorio) Atualizar(publicacaoID uint64, publicacao modelos.Publicacao) error {
	statement, erro := repositorio.db.Prepare("UPDATE publicacoes set titulo = $1, conteudo = $2 WHERE id = $3")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacaoID); erro != nil {
		return erro
	}

	return nil
}

func (repositorio *publicacoesRepositorio) DeletarPublicacao(publicacaoID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM publicacoes WHERE id = $1")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}

	return nil
}

func (repositorio *publicacoesRepositorio) BuscarPorUsuario(usuarioID uint64) ([]modelos.Publicacao, error) {
	linhas, erro := repositorio.db.Query(`
	SELECT p.*, u.nick FROM publicacoes p 
	INNER JOIN usuarios u ON u.id = p.autor_id 
	WHERE p.autor_id = $1`,
		usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var publicacoes []modelos.Publicacao

	for linhas.Next() {
		var publicacao modelos.Publicacao
		if erro = linhas.Scan(
			&publicacao.ID,
			&publicacao.Titulo,
			&publicacao.Conteudo,
			&publicacao.AutorID,
			&publicacao.Curtidas,
			&publicacao.CriadoEm,
			&publicacao.AutorNick,
		); erro != nil {
			return nil, erro
		}
		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func (repositorio *publicacoesRepositorio) Curtir(publicacaoID uint64) error {
	statement, erro := repositorio.db.Prepare("UPDATE publicacoes set curtidas = curtidas + 1 WHERE id = $1")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}

	return nil
}

func (repositorio *publicacoesRepositorio) Descurtir(publicacaoID uint64) error {
	statement, erro := repositorio.db.Prepare(`
	UPDATE publicacoes set curtidas = 
	CASE WHEN curtidas > 0 THEN curtidas - 1 
	ELSE curtidas END 
	WHERE id = $1`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publicacaoID); erro != nil {
		return erro
	}

	return nil
}
