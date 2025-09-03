# GoBid 🚀

![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

## Sobre o Projeto

GoBid é o backend para um sistema de leilões em tempo real, desenvolvido durante a formação Go da [Rocketseat](https://www.rocketseat.com.br/).

O projeto implementa as funcionalidades centrais de uma plataforma de leilões, utilizando Go para o backend e WebSockets para a comunicação em tempo real. Ele serve como um estudo prático sobre desenvolvimento backend, aplicando conceitos como arquitetura em camadas, acesso a banco de dados, autenticação e concorrência, com algumas mudanças pontuais sobre a estrutura ensinada no curso.

## Funcionalidades

* **Autenticação de Usuários:** Sistema de cadastro, login e logout com gerenciamento de sessão.
* **Criação de Leilões:** Usuários autenticados podem cadastrar produtos, o que automaticamente inicia um leilão.
* **Salas de Leilão em Tempo Real:** Cada produto em leilão possui uma "sala" para onde os eventos são transmitidos via WebSockets.
* **Lances em Tempo Real:** Os lances são enviados e recebidos instantaneamente por todos os participantes do leilão.
* **Gerenciamento de Estado:** Validação de lances e o ciclo de vida do leilão são gerenciados pelo servidor.

## Tecnologias Utilizadas

* **Linguagem:** [Go](https://golang.org/)
* **Banco de Dados:** [PostgreSQL](https://www.postgresql.org/)
* **Comunicação Real-time:** [WebSockets](https://github.com/gorilla/websocket)
* **Router HTTP:** [Chi/v5](https://github.com/go-chi/chi)
* **Driver do Banco de Dados:** [pgx/v5](https://github.com/jackc/pgx)
* **Geração de Código SQL:** [sqlc](https://github.com/sqlc-dev/sqlc)
* **Migrations de Banco de Dados:** [tern](https://github.com/jackc/tern)
* **Gerenciamento de Sessão:** [alexedwards/scs](https://github.com/alexedwards/scs)
* **Live Reloading (Dev):** [Air](https://github.com/cosmtrek/air)

## Estrutura do Projeto

O projeto segue uma arquitetura em camadas para uma boa separação de responsabilidades:

.
├── cmd/
│   ├── api/            # Ponto de entrada da aplicação principal.
│   └── terndotenv/     # Utilitário para rodar as migrations com .env.
├── internal/
│   ├── api/            # Handlers HTTP, rotas (Chi) e middlewares.
│   ├── services/       # Lógica de negócio (leilão, lances, usuários).
│   ├── store/pgstore/  # Camada de acesso a dados.
│   │   ├── migrations/ # Arquivos de migration (tern).
│   │   ├── queries/    # Arquivos .sql com as queries (sqlc).
│   │   └── *.sql.go    # Código Go gerado pelo sqlc.
│   ├── usecase/        # Structs de requisição e sua validação.
│   └── validator/      # Utilitários para validação de dados.
├── .air.toml           # Configuração da ferramenta Air.
├── go.mod
└── tern.conf           # Configuração da ferramenta Tern.

## 🚀 Começando

Siga os passos abaixo para configurar e rodar o projeto localmente.

### Pré-requisitos

* [Go](https://go.dev/doc/install) (versão 1.21 ou superior)
* [PostgreSQL](https://www.postgresql.org/download/) (sugestão: rodar via Docker)
* [tern](https://github.com/jackc/tern?tab=readme-ov-file#installation) (precisa estar instalado e no `PATH` do sistema)
* [Air](https://github.com/cosmtrek/air#installation) (para desenvolvimento)

### Instalação

1.  **Clone o repositório:**
    ```bash
    git clone [https://github.com/gregoryAlvim/gobid.git](https://github.com/gregoryAlvim/gobid.git)
    cd gobid
    ```

2.  **Configure as Variáveis de Ambiente:**
    Crie um arquivo `.env` na raiz do projeto para armazenar as credenciais do banco de dados.
    ```env
    # Exemplo de conteúdo para o arquivo .env
    GOBID_DATABASE_HOST=localhost
    GOBID_DATABASE_PORT=5432
    GOBID_DATABASE_USER=seu_usuario
    GOBID_DATABASE_NAME=gobid
    GOBID_DATABASE_PASSWORD=sua_senha
    ```

3.  **Configure o Banco de Dados:**
    Inicie seu servidor PostgreSQL e crie o banco de dados (`gobid` ou o nome que você definiu no `.env`).

4.  **Rode as Migrations:**
    O projeto inclui um programa Go que carrega as variáveis do arquivo `.env` e executa o `tern` para você, aplicando as migrations e criando as tabelas.
    ```bash
    go run ./cmd/terndotenv/main.go
    ```

5.  **Instale as Dependências do Go:**
    ```bash
    go mod tidy
    ```

6.  **Execute a Aplicação Principal:**
    * **Para desenvolvimento (com live reload):**
        ```bash
        air
        ```
    * **Para executar manualmente:**
        ```bash
        go run ./cmd/api
        ```
    O servidor estará rodando na porta `3080`.

## Endpoints da API

| Método | Endpoint                                         | Descrição                                      | Autenticação |
| :----- | :----------------------------------------------- | :--------------------------------------------- | :----------- |
| `POST` | `/api/v1/users/signup`                           | Cadastra um novo usuário.                      | Nenhuma      |
| `POST` | `/api/v1/users/login`                            | Autentica um usuário e cria uma sessão.        | Nenhuma      |
| `POST` | `/api/v1/users/logout`                           | Invalida a sessão do usuário.                  | Requerida    |
| `POST` | `/api/v1/products`                               | Cria um novo produto e inicia seu leilão.      | Requerida    |
| `GET`  | `/api/v1/products/ws/subscribe/{product_id}`     | Inscreve o usuário no leilão via WebSocket.    | Requerida    |

## Origem do Projeto

Este projeto foi desenvolvido com base nos conhecimentos e desafios propostos na formação de Go da Rocketseat. Algumas alterações e adições foram implementadas sobre a estrutura original do curso para explorar diferentes conceitos e aprofundar o aprendizado.