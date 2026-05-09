// Declara o pacote como "usercontrollers".
// Controllers são responsáveis por receber a requisição HTTP, validar os dados
// e chamar as funções de negócio/banco de dados. Eles são a "porta de entrada" da sua API.
package usercontrollers

import (
	// Pacote interno do seu projeto que contém as funções que falam com o banco de dados.
	// É aqui que fica a lógica de criar, buscar, atualizar usuários no banco.
	"database-api/database/user"

	// Pacote nativo do Go com constantes HTTP prontas.
	// Ex: http.StatusOK = 200, http.StatusBadRequest = 400, http.StatusCreated = 201
	// Usar essas constantes é melhor do que escrever números "mágicos" no código.
	"net/http"

	// Framework Gin para lidar com requisições e respostas HTTP.
	"github.com/gin-gonic/gin"
)

// CreateUserHandle é o handler da rota POST /users.
// Um handler no Gin sempre recebe um ponteiro para gin.Context como parâmetro.
// O *gin.Context (chamado aqui de "c") contém tudo sobre a requisição atual:
// - o corpo (body) que o cliente enviou
// - os parâmetros da URL
// - os headers
// - e os métodos para enviar a resposta de volta
func CreateUserHandle(c *gin.Context) {

	// Declara uma struct anônima (sem nome, usada só aqui) para representar
	// o corpo JSON esperado da requisição.
	// As tags `json:"nome"` dizem ao Go como mapear os campos do JSON para a struct.
	// Ex: o JSON {"name": "João"} vai preencher o campo Name com "João".
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// c.ShouldBindJSON tenta ler o corpo da requisição HTTP e converter
	// o JSON recebido para dentro da variável "body".
	// O "&body" passa o endereço de memória da variável (ponteiro),
	// permitindo que a função preencha os campos diretamente.
	// Se o JSON estiver malformado ou faltar campos obrigatórios, retorna erro.
	if err := c.ShouldBindJSON(&body); err != nil {

		// c.JSON envia uma resposta HTTP com um status code e um corpo JSON.
		// http.StatusBadRequest = 400 — significa que o cliente mandou dados inválidos.
		// gin.H{} é um atalho do Gin para map[string]any{} — cria um JSON simples.
		// err.Error() converte o erro em uma string legível para mostrar ao cliente.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		// "return" encerra a execução do handler aqui.
		// Sem ele, o código continuaria executando mesmo após enviar a resposta de erro.
		return
	}

	// Chama a função CreateUser do pacote "user" (seu banco de dados),
	// passando os dados que vieram do corpo da requisição.
	// Em Go, funções podem retornar múltiplos valores — aqui retorna
	// o usuário criado e um possível erro.
	createdUser, err := user.CreateUser(body.Name, body.Email, body.Password)

	// Verifica se ocorreu algum erro ao tentar criar o usuário no banco.
	if err != nil {

		// Verifica o tipo específico do erro comparando a mensagem.
		// Isso permite retornar status HTTP diferentes dependendo do problema.
		// Nota: comparar strings de erro é simples, mas em projetos maiores
		// é melhor usar erros tipados (errors.Is / errors.As) — mas para
		// começar, essa abordagem funciona bem.
		if err.Error() == "email already in use" {

			// http.StatusConflict = 409 — usado quando o recurso já existe.
			// Faz sentido aqui porque o email já está cadastrado no banco.
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		// Se o erro não foi "email already in use", é um erro inesperado do servidor.
		// http.StatusInternalServerError = 500 — erro interno, culpa do servidor, não do cliente.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Se chegou até aqui, o usuário foi criado com sucesso.
	// http.StatusCreated = 201 — indica que um recurso foi criado com sucesso.
	// Retorna o objeto do usuário criado como resposta JSON para o cliente.
	c.JSON(http.StatusCreated, createdUser)
}
