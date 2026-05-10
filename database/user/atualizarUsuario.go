// Pacote "user" — contém todas as funções que interagem com a tabela de usuários no banco.
// Essa camada é chamada de "repository" ou "data layer" — ela só sabe falar com o banco,
// não sabe nada sobre HTTP, rotas ou o Gin.
package user

import (
	// Pacote interno que provavelmente contém a função de conectar ao banco de dados.
	// Separar a conexão em um pacote próprio evita repetir esse código em todo lugar.
	"database-api/database"
	"database-api/utils"

	// Pacote nativo do Go para formatar strings e criar erros customizados.
	// fmt.Errorf() cria um erro novo com uma mensagem que você define.
	"fmt"
)

// updateUserModel é uma struct que define o formato dos dados que essa função retorna.
// Ela representa apenas os campos que fazem sentido devolver ao cliente após atualizar —
// repare que não inclui "old_password", pois isso nunca deve ser exposto na resposta.
// As tags `json:"..."` definem como os campos aparecem quando viram JSON.
type updateUserModel struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUser recebe o id do usuário e os novos dados, valida a senha antiga
// e atualiza o registro no banco. Retorna o usuário atualizado ou um erro.
// Retornar a struct vazia "updateUserModel{}" junto com o erro é o padrão Go:
// quando há erro, você retorna o "zero value" do tipo e o erro preenchido.
func UpdateUser(id int, name, email, password, old_password string) (updateUserModel, error) {
	db, err := database.ConectarBanco()
	if err != nil {
		return updateUserModel{}, err
	}

	var password_old string
	err = db.QueryRow(`SELECT password FROM "users" WHERE id=$1`, id).Scan(&password_old)
	if err != nil {
		return updateUserModel{}, err
	}

	if !utils.VerificarSenha(password_old, old_password) {
		return updateUserModel{}, fmt.Errorf("senha incorreta")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return updateUserModel{}, err
	}

	var updated_user updateUserModel
	err = db.QueryRow(
		`UPDATE "users" SET name=$1, email=$2, password=$3 WHERE id=$4 RETURNING name, email, password`,
		name, email, hashedPassword, id,
	).Scan(&updated_user.Name, &updated_user.Email, &updated_user.Password)
	if err != nil {
		return updateUserModel{}, err
	}

	return updated_user, nil
}
