# Social Network
Este projeto consiste em uma API RESTful simulando uma rede social construida com a linguagem Go.

## Requisitos para executar a aplicação
1. A ferramenta git instalada para controle de versões do projeto (caso não tenha, clique [aqui](https://git-scm.com/));
2. A ferramenta docker instalada (caso não tenha, clique [aqui](https://www.docker.com/) para ir para a página de instalação).

## Como executar?
1. Clone o projeto com o (`git clone`) no seu computador e abra um terminal na pasta baixada;
2. Navegue até a pasta `api`;
3. Execute a aplicação utilizando o Docker. Para isso, execute o seguinte comando:
```bash
docker compose up --build
```

Pronto! O projeto está configurado. A partir de agora, toda vez que quiser iniciar o projeto basta executar o comando `docker compose up --build`. Assim, o projeto estará disponível no endereço `http://localhost:8000`.

## Descrição das rotas da API

A API tem como objetivo o gerenciamento de uma rede social. Abaixo estão listadas as rotas disponíveis:

### **Rotas de Autenticação**

- **`POST /login`**: Autentica o usuário no sistema.

OBS: Para realizar a autenticação, pode-se utilizar o seguinte corpo para esta requisição:

```JSON
{
    "email": "usuario1@gmail.com",
    "senha": "123456"
}
```

Após isto, copie o token retornado e adicione-o no cabeçalho de autenticação dos demais endpoints.

### **Rotas de Publicações**

- **`POST /publicacoes`**: Cria uma nova publicação.  
**Autenticação:** Requerida.

- **`GET /publicacoes`**: Retorna todas as publicações disponíveis.  
  **Autenticação:** Requerida.

- **`GET /publicacoes/{publicacaoId}`**: Retorna os detalhes de uma publicação específica.  
  **Autenticação:** Requerida.

- **`PUT /publicacoes/{publicacaoId}`**: Atualiza os dados de uma publicação específica.  
  **Autenticação:** Requerida.

- **`DELETE /publicacoes/{publicacaoId}`**: Deleta uma publicação específica.  
  **Autenticação:** Requerida.

### **Rotas Relacionadas a Usuários**

- **`GET /usuarios/{usuarioId}/publicacoes`**: Retorna todas as publicações criadas por um usuário específico.  
  **Autenticação:** Requerida.

### **Rotas de Curtidas**

- **`POST /publicacoes/{publicacaoId}/curtir`**: Adiciona uma curtida a uma publicação.  
  **Autenticação:** Requerida.

- **`POST /publicacoes/{publicacaoId}/descurtir`**: Remove uma curtida de uma publicação.  
  **Autenticação:** Requerida.  
