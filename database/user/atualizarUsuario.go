// Pacote "user" — contém todas as funções que interagem com a tabela de usuários no banco.
// Essa camada é chamada de "repository" ou "data layer" — ela só sabe falar com o banco,
// não sabe nada sobre HTTP, rotas ou o Gin.
package user

import (
	// Pacote interno que provavelmente contém a função de conectar ao banco de dados.
	// Separar a conexão em um pacote próprio evita repetir esse código em todo lugar.
	"database-api/database"

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

	// Chama a função do pacote "database" que abre (ou reutiliza) a conexão com o banco.
	// "db" é um objeto que representa essa conexão — é por ele que você executa queries.
	// Se a conexão falhar, retorna a struct vazia e o erro para quem chamou essa função.
	db, err := database.ConectarBanco()
	if err != nil {
		return updateUserModel{}, err
	}

	// Declara uma variável para guardar a senha atual do usuário que está no banco.
	// Ela começa como string vazia "" (zero value de string em Go).
	var password_old string

	// db.QueryRow() executa uma query SQL que espera retornar apenas UMA linha.
	// É ideal para buscas por ID ou campos únicos, pois é mais eficiente que Query().
	// O $1 é um placeholder — o valor real (id) é passado como argumento separado.
	// Isso evita SQL Injection, pois o banco trata o valor como dado, não como código SQL.
	// .Scan() lê o resultado da linha e coloca na variável &password_old.
	// O "&" passa o ponteiro para que Scan() consiga escrever o valor dentro da variável.
	err = db.QueryRow(`SELECT password FROM "user" WHERE id=$1`, id).Scan(&password_old)
	if err != nil {
		// Se o usuário não existir, db.QueryRow().Scan() retorna "sql: no rows in result set".
		// Se houver erro de conexão ou de SQL, também cai aqui.
		return updateUserModel{}, err
	}

	// Compara a senha enviada pelo cliente com a senha que veio do banco.
	// Em Go, "==" em strings compara o conteúdo caractere por caractere.
	// Nota importante: em produção, senhas devem ser armazenadas como hash (bcrypt, argon2)
	// e comparadas com funções específicas — nunca como texto puro como está aqui.
	passwordMatches := old_password == password_old

	// Se as senhas não batem, retorna um erro customizado com fmt.Errorf().
	// O "!" é o operador de negação em Go — "!passwordMatches" = "se NÃO combina".
	if !passwordMatches {
		// fmt.Errorf() cria um erro com mensagem personalizada.
		// Esse erro vai subir até o handler, que vai devolver 500 ao cliente.
		return updateUserModel{}, fmt.Errorf("senha incorreta")
	}

	// Declara a variável que vai receber os dados do usuário após a atualização.
	var updated_user updateUserModel

	// Executa a query de UPDATE no banco.
	// $1, $2, $3, $4 são os placeholders — são preenchidos em ordem pelos argumentos
	// passados depois da string SQL: name, email, password, id.
	// RETURNING é um recurso do PostgreSQL que faz o banco devolver os dados
	// da linha atualizada logo após o UPDATE, evitando uma segunda query de SELECT.
	// .Scan() lê os três campos retornados e preenche diretamente nos campos da struct.
	err = db.QueryRow(
		`UPDATE "user" SET name=$1, email=$2, password=$3 WHERE id=$4 RETURNING name, email, password`,
		name, email, password, id,
	).Scan(&updated_user.Name, &updated_user.Email, &updated_user.Password)

	if err != nil {
		return updateUserModel{}, err
	}

	// Retorna o usuário atualizado e "err" — que aqui vale nil (sem erro).
	// Retornar "err" em vez de "nil" explicitamente é um hábito comum em Go,
	// mas neste caso os dois são equivalentes, pois já verificamos o erro acima.
	return updated_user, err
}
