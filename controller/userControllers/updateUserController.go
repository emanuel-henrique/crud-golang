// Mesmo pacote do controller anterior — ambos os handlers ficam aqui.
package usercontrollers

import (
	// Pacote interno com as funções que falam com o banco de dados.
	"database-api/database/user"

	// Pacote nativo com as constantes de status HTTP (200, 400, 500...).
	"net/http"

	// Pacote nativo do Go para converter tipos.
	// "strconv" = string conversion — converte strings em números e vice-versa.
	"strconv"

	// Framework Gin para lidar com requisições e respostas HTTP.
	"github.com/gin-gonic/gin"
)

// UpdateUserHandle é o handler da rota PUT /users/:id.
// Ele recebe o ID do usuário pela URL e os novos dados pelo corpo JSON.
func UpdateUserHandle(c *gin.Context) {

	// c.Param("id") lê o valor do parâmetro dinâmico ":id" da URL.
	// Ex: se a rota chamada foi PUT /users/42, c.Param("id") retorna a string "42".
	// Porém, IDs em banco de dados normalmente são números inteiros, não strings.
	// Por isso usamos strconv.Atoi() para converter — "Atoi" significa "ASCII to integer".
	// Assim como sempre em Go, a função retorna o valor convertido E um possível erro.
	id, err := strconv.Atoi(c.Param("id"))

	// Se a conversão falhou, significa que o que veio na URL não é um número válido.
	// Ex: PUT /users/abc — "abc" não pode virar um inteiro, então retorna erro.
	if err != nil {
		// Retorna 400 Bad Request com uma mensagem clara para o cliente.
		// Aqui a mensagem é manual ("id inválido") porque o erro do strconv
		// seria técnico demais para o cliente final entender.
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	// Struct anônima que define os campos esperados no corpo JSON da requisição.
	// Diferente do CreateUser, aqui também pedimos a senha antiga (Old_password)
	// para validar a identidade do usuário antes de permitir a atualização.
	// Nota: por convenção Go usa camelCase (OldPassword), mas old_password
	// também funciona — o que importa para o JSON é a tag `json:"old_password"`.
	var body struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		Password     string `json:"password"`
		Old_password string `json:"old_password"`
	}

	// Tenta converter o corpo JSON da requisição para dentro da struct "body".
	// Funciona igual ao CreateUserHandle — lê o JSON e preenche os campos.
	if err := c.ShouldBindJSON(&body); err != nil {
		// 400 Bad Request — o cliente mandou um JSON inválido ou malformado.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chama a função UpdateUser do pacote de banco de dados, passando:
	// - id: o número inteiro já convertido lá em cima
	// - os campos do body: nome, email, nova senha e senha antiga
	// Retorna o usuário atualizado e um possível erro.
	updatedUser, err := user.UpdateUser(id, body.Name, body.Email, body.Password, body.Old_password)

	// Se qualquer coisa der errado na camada do banco (senha errada,
	// usuário não encontrado, erro de conexão...), retorna 500.
	// Diferente do CreateUser, aqui não há distinção de tipos de erro —
	// tudo cai no mesmo bloco. Em projetos maiores, valeria separar
	// (ex: 404 se usuário não existir, 401 se senha estiver errada).
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Se chegou aqui, a atualização foi bem-sucedida.
	// http.StatusOK = 200 — resposta padrão para operações bem-sucedidas
	// que não criam um recurso novo (criações usam 201, como no CreateUser).
	// Retorna os dados atualizados do usuário como JSON.
	c.JSON(http.StatusOK, updatedUser)
}
