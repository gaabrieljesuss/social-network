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

## Monitoramento da API com Prometheus e Grafana

Este projeto está configurado para permitir o monitoramento de métricas da API utilizando **Prometheus** e **Grafana**.

### **Prometheus**
O Prometheus coleta métricas sobre a API e está configurado para monitorar o tempo de resposta das requisições e o número total de requisições recebidas.

#### **Métricas Disponíveis:**
- **`api_requests_total`**: Número total de requisições recebidas pela API, categorizadas por método HTTP e endpoint.
- **`api_response_time_seconds`**: Tempo de resposta das requisições, categorizado por método HTTP e endpoint.

### **Grafana**
Grafana é utilizado para visualização das métricas coletadas pelo Prometheus. O painel do Grafana exibe gráficos e informações sobre as requisições da API, incluindo tempos de resposta e contagem de requisições.

#### **Como acessar o Grafana:**
1. Após subir os containers com o comando `docker compose up --build`, acesse o Grafana no endereço `http://localhost:3000`.
2. As credenciais padrão para login são:
  - **Usuário:** `admin`
  - **Senha:** A senha configurada na variável de ambiente `GF_SECURITY_ADMIN_PASSWORD` no arquivo `.env`.

#### **Painéis de Monitoramento:**
- O Grafana estará configurado para buscar métricas do Prometheus, e você poderá visualizar gráficos sobre o tempo de resposta das requisições e a quantidade total de requisições feitas para os diferentes endpoints da API.

### **Acessando as Métricas do Prometheus:**
As métricas expostas pelo Prometheus podem ser acessadas diretamente pelo endereço:
- **Prometheus:** `http://localhost:9090/metrics`

### **Acessando as Informações no Prometheus:**
1. Abra o Prometheus em `http://localhost:9090`.
2. Utilize as queries do Prometheus, como:
  - **`api_requests_total`** para visualizar o número total de requisições.
  - **`api_response_time_seconds_sum / api_response_time_seconds_count`** para ver o tempo médio de resposta das requisições.
