package main

import (
	"api/src/banco"
	"api/src/config"
	"api/src/controllers"
	"api/src/metrics"
	"api/src/repositorios"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Carregar()

	metrics.Init()

	db, err := banco.Conectar()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	repositorioUsuarios := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioController := controllers.NovoUsuarioController(repositorioUsuarios)

	repositorioPublicacoes := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoesController := controllers.NovoPublicacoesController(repositorioPublicacoes)

	fmt.Printf("Rodando a API na porta %d", config.Porta)

	r := router.Gerar(usuarioController, publicacoesController)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
