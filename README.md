# 🚀 CRUD API em Golang

API RESTful desenvolvida em Go utilizando o framework Gin e PostgreSQL.  
O projeto implementa operações completas de CRUD de usuários com foco em arquitetura backend, segurança e boas práticas.

---

## 🛠 Tecnologias utilizadas

- Golang
- Gin Gonic
- PostgreSQL
- SQL nativo
- bcrypt
- godotenv

---

## ✨ Funcionalidades

- Criar usuários
- Buscar usuário por ID
- Atualizar usuários
- Deletar usuários
- Hash de senha com bcrypt
- Validação de senha antiga para atualização
- Variáveis de ambiente com `.env`
- Pool de conexões com PostgreSQL

---

## 📁 Estrutura do Projeto

```bash
database-api/
│
├── controller/
│   └── userControllers/
│
├── database/
│   ├── user/
│   └── connect.go
│
├── utils/
│
├── main/
│   └── main.go
│
├── .env
├── go.mod
└── go.sum
```

---

## ⚙️ Configuração do ambiente

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/seu-repositorio.git
```

---

### 2. Instale as dependências

```bash
go mod tidy
```

---

### 3. Configure o arquivo `.env`

Crie um arquivo `.env` na raiz do projeto:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=sua_senha
DB_NAME=database_api
```

---

### 4. Configure a tabela no PostgreSQL

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);
```

---

## ▶️ Executando o projeto

```bash
go run main/main.go
```

Servidor iniciado em:

```bash
http://localhost:8080
```

---

## 📡 Rotas da API

### 👤 Criar usuário

```http
POST /users
```

#### Body

```json
{
  "name": "Emanuel",
  "email": "emanuel@email.com",
  "password": "123456"
}
```

---

### 🔍 Buscar usuário por ID

```http
GET /users/:id
```

---

### ✏️ Atualizar usuário

```http
PUT /users/:id
```

#### Body

```json
{
  "name": "Novo Nome",
  "email": "novo@email.com",
  "password": "novaSenha",
  "old_password": "senhaAtual"
}
```

---

### 🗑️ Deletar usuário

```http
DELETE /users/:id
```

---

## 🔒 Segurança

As senhas são armazenadas utilizando hash com bcrypt, evitando armazenamento de senhas em texto puro.

---

## 🧠 Conceitos praticados

- APIs RESTful
- Estruturação de projeto backend
- SQL parametrizado
- Hash de senhas
- Tratamento de erros
- Variáveis de ambiente
- Pool de conexões
- Arquitetura em camadas

---

## 🎯 Objetivo do projeto

Este projeto foi desenvolvido com o objetivo de aprofundar conhecimentos em:

- Desenvolvimento Backend com Go
- Framework Gin
- PostgreSQL
- Segurança de autenticação
- Arquitetura backend
- Desenvolvimento de APIs REST

---

## 👨‍💻 Autor

**Emanuel Henrique**

- LinkedIn: www.linkedin.com/in/emanuel-henrique-38b264392
- GitHub: github.com/emanuel-henrique
