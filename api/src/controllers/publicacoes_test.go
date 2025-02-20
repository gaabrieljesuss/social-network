package controllers

import (
	"api/src/autenticacao"
	"api/src/modelos"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPublicacoesRepositorio struct {
	mock.Mock
}

func (m *MockPublicacoesRepositorio) Criar(publicacao modelos.Publicacao) (uint64, error) {
	args := m.Called(publicacao)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockPublicacoesRepositorio) BuscarPorId(publicacaoID uint64) (modelos.Publicacao, error) {
	args := m.Called(publicacaoID)
	return args.Get(0).(modelos.Publicacao), args.Error(1)
}

func (m *MockPublicacoesRepositorio) BuscarPublicacoes(usuarioID uint64) ([]modelos.Publicacao, error) {
	args := m.Called(usuarioID)
	return args.Get(0).([]modelos.Publicacao), args.Error(1)
}

func (m *MockPublicacoesRepositorio) Atualizar(id uint64, publicacao modelos.Publicacao) error {
	args := m.Called(id, publicacao)
	return args.Error(0)
}

func (m *MockPublicacoesRepositorio) DeletarPublicacao(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPublicacoesRepositorio) BuscarPorUsuario(usuarioID uint64) ([]modelos.Publicacao, error) {
	args := m.Called(usuarioID)
	return args.Get(0).([]modelos.Publicacao), args.Error(1)
}

func (m *MockPublicacoesRepositorio) Curtir(publicacaoID uint64) error {
	args := m.Called(publicacaoID)
	return args.Error(0)
}

func (m *MockPublicacoesRepositorio) Descurtir(publicacaoID uint64) error {
	args := m.Called(publicacaoID)
	return args.Error(0)
}

func setupConfig(t *testing.T, repositorio *MockPublicacoesRepositorio) (*PublicacoesController, *httptest.ResponseRecorder) {
	controller := NovoPublicacoesController(repositorio)
	recorder := httptest.NewRecorder()
	return controller, recorder
}

func TestCriarPublicacao(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	publicacao := modelos.Publicacao{Titulo: "Teste", Conteudo: "Conteudo", AutorID: usuarioID}
	mockRepo.On("Criar", publicacao).Return(uint64(1), nil)

	tokenString, err := autenticacao.CriarToken(usuarioID)
	if err != nil {
		t.Fatalf("Erro ao gerar token JWT: %v", err)
	}
	body, _ := json.Marshal(publicacao)
	r := httptest.NewRequest("POST", "/publicacoes", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+tokenString)

	controller.CriarPublicacao(recorder, r)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestCriarPublicacao_Falha(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	publicacao := modelos.Publicacao{}
	mockRepo.On("Criar", publicacao).Return(uint64(0), errors.New("erro ao criar publicação"))

	tokenString, err := autenticacao.CriarToken(usuarioID)
	if err != nil {
		t.Fatalf("Erro ao gerar token JWT: %v", err)
	}
	body, _ := json.Marshal(publicacao)
	r := httptest.NewRequest("POST", "/publicacoes", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+tokenString)

	controller.CriarPublicacao(recorder, r)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestAtualizarPublicacao(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := "1"
	usuarioID := uint64(1)
	publicacao := modelos.Publicacao{Titulo: "Atualizado", Conteudo: "Novo Conteudo"}
	mockRepo.On("Atualizar", uint64(1), publicacao).Return(nil)

	tokenString, err := autenticacao.CriarToken(usuarioID)
	if err != nil {
		t.Fatalf("Erro ao gerar token JWT: %v", err)
	}
	body, _ := json.Marshal(publicacao)
	r := httptest.NewRequest("PUT", "/publicacoes/"+publicacaoID, bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+tokenString)
	vars := map[string]string{"publicacaoId": publicacaoID}
	r = mux.SetURLVars(r, vars)

	controller.AtualizarPublicacao(recorder, r)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestAtualizarPublicacao_Falha(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := "1"
	usuarioID := uint64(1)
	publicacao := modelos.Publicacao{Titulo: "Atualizado", Conteudo: "Novo Conteudo"}
	mockRepo.On("Atualizar", uint64(1), publicacao).Return(errors.New("erro ao atualizar publicação"))

	tokenString, err := autenticacao.CriarToken(usuarioID)
	if err != nil {
		t.Fatalf("Erro ao gerar token JWT: %v", err)
	}
	body, _ := json.Marshal(publicacao)
	r := httptest.NewRequest("PUT", "/publicacoes/"+publicacaoID, bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+tokenString)
	vars := map[string]string{"publicacaoId": publicacaoID}
	r = mux.SetURLVars(r, vars)

	controller.AtualizarPublicacao(recorder, r)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeletarPublicacao(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := "1"
	publicacaoExistente := modelos.Publicacao{
		ID:       1,
		Titulo:   "Publicação Existente",
		Conteudo: "Conteúdo da publicação",
		AutorID:  1,
	}
	usuarioID := uint64(1)
	mockRepo.On("DeletarPublicacao", uint64(1)).Return(nil)
	mockRepo.On("BuscarPorId", uint64(1)).Return(publicacaoExistente, nil)

	tokenString, err := autenticacao.CriarToken(usuarioID)
	if err != nil {
		t.Fatalf("Erro ao gerar token JWT: %v", err)
	}
	r := httptest.NewRequest("DELETE", "/publicacoes/"+publicacaoID, nil)
	r.Header.Set("Authorization", "Bearer "+tokenString)
	vars := map[string]string{"publicacaoId": publicacaoID}
	r = mux.SetURLVars(r, vars)

	controller.DeletarPublicacao(recorder, r)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeletarPublicacao_Falha(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := "1"
	mockRepo.On("Deletar", uint64(1)).Return(errors.New("erro ao deletar publicação"))

	r := httptest.NewRequest("DELETE", "/publicacoes/"+publicacaoID, nil)
	vars := map[string]string{"publicacaoId": publicacaoID}
	r = mux.SetURLVars(r, vars)

	controller.DeletarPublicacao(recorder, r)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockRepo.AssertExpectations(t)
}
