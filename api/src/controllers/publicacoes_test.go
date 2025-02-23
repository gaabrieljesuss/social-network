package controllers

import (
	"api/src/autenticacao"
	"api/src/modelos"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestCreatePublication_WhenAllFieldsArePassed_ExpectedNewPublication(t *testing.T) {
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

func TestCreatePublication_WhenNoTitleIsPassed_ExpectedBadRequestError(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	publicacao := modelos.Publicacao{}

	tokenString, err := autenticacao.CriarToken(usuarioID)
	if err != nil {
		t.Fatalf("Erro ao gerar token JWT: %v", err)
	}
	body, _ := json.Marshal(publicacao)
	r := httptest.NewRequest("POST", "/publicacoes", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+tokenString)

	controller.CriarPublicacao(recorder, r)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreatePublication_WhenCommandFailsInDatabase_ExpectedInternalServerError(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	publicacao := modelos.Publicacao{Titulo: "Teste", Conteudo: "Conteudo", AutorID: usuarioID}
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

func TestUpdatePublication_WhenAllFieldsArePassed_ExpectedUpdatedPublication(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := "1"
	usuarioID := uint64(1)
	publicacao := modelos.Publicacao{Titulo: "Atualizado", Conteudo: "Novo Conteudo"}
	publicacaoExistente := modelos.Publicacao{
		ID:       1,
		Titulo:   "Publicação Existente",
		Conteudo: "Conteúdo da publicação",
		AutorID:  1,
	}
	mockRepo.On("BuscarPorId", uint64(1)).Return(publicacaoExistente, nil)
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

func TestUpdatePublication_WhenPublicationDoesNotExist_ExpectedNotFoundError(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := "1"
	usuarioID := uint64(1)
	publicacao := modelos.Publicacao{Titulo: "Atualizado", Conteudo: "Novo Conteudo"}
	publicacaoVazia := modelos.Publicacao{Titulo: "Atualizado", Conteudo: "Novo Conteudo"}
	mockRepo.On("BuscarPorId", uint64(1)).Return(publicacaoVazia, errors.New("publicação não encontrada"))

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

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePublication_WhenCommandFailsInDatabase_ExpectedInternalServerError(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := "1"
	usuarioID := uint64(1)
	publicacao := modelos.Publicacao{Titulo: "Atualizado", Conteudo: "Novo Conteudo"}
	publicacaoExistente := modelos.Publicacao{
		ID:       1,
		Titulo:   "Publicação Existente",
		Conteudo: "Conteúdo da publicação",
		AutorID:  1,
	}
	mockRepo.On("BuscarPorId", uint64(1)).Return(publicacaoExistente, nil)
	mockRepo.On("Atualizar", uint64(1), publicacao).Return(errors.New("erro atualizando a publicação"))

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

func TestDeletePublication_WhenPublicationExists_ExpectedDeletedPublication(t *testing.T) {
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

func TestDeletePublication_WhenPublicationDoesNotExist_ExpectedNotFoundError(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	publicacaoID := "1"
	publicacaoVazia := modelos.Publicacao{Titulo: "Atualizado", Conteudo: "Novo Conteudo"}
	mockRepo.On("BuscarPorId", uint64(1)).Return(publicacaoVazia, errors.New("publicação não encontrada"))

	tokenString, err := autenticacao.CriarToken(usuarioID)
	if err != nil {
		t.Fatalf("Erro ao gerar token JWT: %v", err)
	}
	r := httptest.NewRequest("DELETE", "/publicacoes/"+publicacaoID, nil)
	r.Header.Set("Authorization", "Bearer "+tokenString)
	vars := map[string]string{"publicacaoId": publicacaoID}
	r = mux.SetURLVars(r, vars)

	controller.DeletarPublicacao(recorder, r)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeletePublication_WhenCommandFailsInDatabase_ExpectedInternalServerError(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	publicacaoID := "1"
	publicacaoExistente := modelos.Publicacao{
		ID:       1,
		Titulo:   "Publicação Existente",
		Conteudo: "Conteúdo da publicação",
		AutorID:  1,
	}
	mockRepo.On("DeletarPublicacao", uint64(1)).Return(errors.New("erro ao deletar a publicacao"))
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

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestFindPublicationById_WhenPublicationExists_ExpectedPublicationReturned(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := uint64(1)
	publicacao := modelos.Publicacao{ID: publicacaoID, Titulo: "Teste", Conteudo: "Conteudo", AutorID: 1}
	mockRepo.On("BuscarPorId", publicacaoID).Return(publicacao, nil)

	r := httptest.NewRequest("GET", fmt.Sprintf("/publicacoes/%d", publicacaoID), nil)
	vars := map[string]string{"publicacaoId": strconv.FormatUint(publicacaoID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.BuscarPublicacao(recorder, r)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestFindPublicationById_WhenPublicationDoesNotExist_ExpectedNotFoundError(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := uint64(1)
	mockRepo.On("BuscarPorId", publicacaoID).Return(modelos.Publicacao{}, errors.New("publicação não encontrada"))

	r := httptest.NewRequest("GET", fmt.Sprintf("/publicacoes/%d", publicacaoID), nil)
	vars := map[string]string{"publicacaoId": strconv.FormatUint(publicacaoID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.BuscarPublicacao(recorder, r)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestSearchPublicationByUser_WhenUserHasPublications_ExpectedPublicationsReturned(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	publicacoes := []modelos.Publicacao{{ID: 1, Titulo: "Teste", Conteudo: "Conteudo", AutorID: usuarioID}}
	mockRepo.On("BuscarPorUsuario", usuarioID).Return(publicacoes, nil)

	r := httptest.NewRequest("GET", fmt.Sprintf("/usuarios/%d/publicacoes", usuarioID), nil)
	vars := map[string]string{"usuarioId": strconv.FormatUint(usuarioID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.BuscarPublicacoesPorUsuario(recorder, r)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestSearchPublicationByUser_WhenUserHasNoPublications_ExpectedEmptyList(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	mockRepo.On("BuscarPorUsuario", usuarioID).Return([]modelos.Publicacao{}, nil)

	r := httptest.NewRequest("GET", fmt.Sprintf("/usuarios/%d/publicacoes", usuarioID), nil)
	vars := map[string]string{"usuarioId": strconv.FormatUint(usuarioID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.BuscarPublicacoesPorUsuario(recorder, r)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestSearchPublications_WhenHasPublications_ExpectedPublicationsReturned(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	publicacoes := []modelos.Publicacao{{ID: 1, Titulo: "Teste", Conteudo: "Conteudo", AutorID: usuarioID}}
	mockRepo.On("BuscarPublicacoes", usuarioID).Return(publicacoes, nil)

	tokenString, err := autenticacao.CriarToken(usuarioID)
	if err != nil {
		t.Fatalf("Erro ao gerar token JWT: %v", err)
	}
	r := httptest.NewRequest("GET", fmt.Sprintf("/publicacoes"), nil)
	r.Header.Set("Authorization", "Bearer "+tokenString)
	vars := map[string]string{"usuarioId": strconv.FormatUint(usuarioID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.BuscarPublicacoes(recorder, r)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestSearchPublications_WhenUserHasNoPublications_ExpectedEmptyList(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	usuarioID := uint64(1)
	mockRepo.On("BuscarPublicacoes", usuarioID).Return([]modelos.Publicacao{}, nil)

	tokenString, err := autenticacao.CriarToken(usuarioID)
	if err != nil {
		t.Fatalf("Erro ao gerar token JWT: %v", err)
	}
	r := httptest.NewRequest("GET", fmt.Sprintf("/publicacoes"), nil)
	r.Header.Set("Authorization", "Bearer "+tokenString)
	vars := map[string]string{"usuarioId": strconv.FormatUint(usuarioID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.BuscarPublicacoes(recorder, r)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestLikePublication_WhenPublicationExists_ExpectedSuccess(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := uint64(1)
	mockRepo.On("Curtir", publicacaoID).Return(nil)

	r := httptest.NewRequest("POST", fmt.Sprintf("/publicacoes/%d/curtir", publicacaoID), nil)
	vars := map[string]string{"publicacaoId": strconv.FormatUint(publicacaoID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.CurtirPublicacao(recorder, r)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestLikePublication_WhenDatabaseFails_ExpectedInternalServerError(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := uint64(1)
	mockRepo.On("Curtir", publicacaoID).Return(errors.New("erro ao curtir a publicação"))

	r := httptest.NewRequest("POST", fmt.Sprintf("/publicacoes/%d/curtir", publicacaoID), nil)
	vars := map[string]string{"publicacaoId": strconv.FormatUint(publicacaoID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.CurtirPublicacao(recorder, r)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestUnlikePublication_WhenPublicationExists_ExpectedSuccess(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := uint64(1)
	mockRepo.On("Descurtir", publicacaoID).Return(nil)

	r := httptest.NewRequest("POST", fmt.Sprintf("/publicacoes/%d/descurtir", publicacaoID), nil)
	vars := map[string]string{"publicacaoId": strconv.FormatUint(publicacaoID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.DescurtirPublicacao(recorder, r)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestUnlikePublication_WhenDatabaseFails_ExpectedInternalServerError(t *testing.T) {
	mockRepo := new(MockPublicacoesRepositorio)
	controller, recorder := setupConfig(t, mockRepo)

	publicacaoID := uint64(1)
	mockRepo.On("Descurtir", publicacaoID).Return(errors.New("erro ao descurtir a publicação"))

	r := httptest.NewRequest("POST", fmt.Sprintf("/publicacoes/%d/descurtir", publicacaoID), nil)
	vars := map[string]string{"publicacaoId": strconv.FormatUint(publicacaoID, 10)}
	r = mux.SetURLVars(r, vars)

	controller.DescurtirPublicacao(recorder, r)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	mockRepo.AssertExpectations(t)
}
