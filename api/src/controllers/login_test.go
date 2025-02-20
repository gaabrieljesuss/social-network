package controllers

import (
	"api/src/modelos"
	"api/src/seguranca"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepositorio struct {
	mock.Mock
}

func (m *MockRepositorio) BuscarPorEmail(email string) (modelos.Usuario, error) {
	args := m.Called(email)
	return args.Get(0).(modelos.Usuario), args.Error(1)
}

func setup(t *testing.T, repositorio *MockRepositorio) (*UsuarioController, *httptest.ResponseRecorder) {
	controller := NovoUsuarioController(repositorio)
	recorder := httptest.NewRecorder()
	return controller, recorder
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("Erro simulado na leitura dos dados")
}

func (e *errorReader) Close() error {
	return nil
}

func criarRequestLogin(t *testing.T, email, senha string) *http.Request {
	body, _ := json.Marshal(map[string]string{"email": email, "senha": senha})
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestLogin_Sucesso(t *testing.T) {
	mockRepo := new(MockRepositorio)
	controller, recorder := setup(t, mockRepo)

	email := "usuario@teste.com"
	senha := "senhaCorreta"
	hashSenha, _ := seguranca.Hash(senha)

	usuarioMock := modelos.Usuario{ID: 1, Email: email, Senha: string(hashSenha)}

	mockRepo.On("BuscarPorEmail", email).Return(usuarioMock, nil)

	req := criarRequestLogin(t, email, senha)
	controller.Login(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, recorder.Body.String())
	mockRepo.AssertExpectations(t)
}

func TestLogin_EmailNaoEncontrado(t *testing.T) {
	mockRepo := new(MockRepositorio)
	controller, recorder := setup(t, mockRepo)

	email := "naoexiste@teste.com"
	mockRepo.On("BuscarPorEmail", email).Return(modelos.Usuario{}, errors.New("usuário não encontrado"))

	req := criarRequestLogin(t, email, "qualquerSenha")
	controller.Login(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "usuário não encontrado")
	mockRepo.AssertExpectations(t)
}

func TestLogin_SenhaIncorreta(t *testing.T) {
	mockRepo := new(MockRepositorio)
	controller, recorder := setup(t, mockRepo)

	email := "usuario@teste.com"
	senhaCorreta := "senhaCorreta"
	hashSenha, _ := seguranca.Hash(senhaCorreta)
	usuarioMock := modelos.Usuario{ID: 1, Email: email, Senha: string(hashSenha)}

	mockRepo.On("BuscarPorEmail", email).Return(usuarioMock, nil)

	req := criarRequestLogin(t, email, "senhaErrada")
	controller.Login(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "hashedPassword is not the hash of the given password")
	mockRepo.AssertExpectations(t)
}

func TestLogin_ErroAoBuscarNoBanco(t *testing.T) {
	mockRepo := new(MockRepositorio)
	controller, recorder := setup(t, mockRepo)

	email := "usuario@teste.com"
	mockRepo.On("BuscarPorEmail", email).Return(modelos.Usuario{}, errors.New("erro no banco"))

	req := criarRequestLogin(t, email, "senha")
	controller.Login(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "erro no banco")
	mockRepo.AssertExpectations(t)
}

func TestLogin_ErroAoLerCorpoRequisicao(t *testing.T) {
	mockRepo := new(MockRepositorio)
	controller, recorder := setup(t, mockRepo)

	req, err := http.NewRequest(http.MethodPost, "/login", ioutil.NopCloser(&errorReader{}))
	assert.NoError(t, err)

	controller.Login(recorder, req)

	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Erro simulado na leitura dos dados")
}
