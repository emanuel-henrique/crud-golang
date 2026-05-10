// Declara que este arquivo pertence ao pacote principal.
// Em Go, todo programa executável precisa ter um pacote chamado "main".
package main

import (
	// Importa o pacote "usercontrollers" de dentro do seu próprio projeto.
	// O apelido "usercontrollers" antes do caminho é um alias — você escolhe
	// como quer chamar o pacote dentro deste arquivo.
	// O caminho "database-api/controller/userControllers" é a pasta real no seu projeto.
	usercontrollers "database-api/controller/userControllers"

	// Pacote nativo do Go para registrar erros e mensagens no terminal.
	"log"

	// Gin é um framework externo para criar APIs HTTP em Go.
	// Ele facilita a criação de rotas, lida com requisições e respostas.
	"github.com/gin-gonic/gin"
)

// Função principal — é aqui que o programa começa a executar.
// Todo executável Go começa pelo main().
func main() {

	// gin.Default() cria uma instância do servidor HTTP com dois middlewares
	// já configurados automaticamente:
	// - Logger: exibe no terminal cada requisição recebida (método, rota, status, tempo)
	// - Recovery: se o servidor travar (panic), ele se recupera sem derrubar tudo
	// "app" é a variável que representa seu servidor — você usa ela para tudo.
	app := gin.Default()

	// Define uma rota do tipo POST no caminho "/users".
	// Quando alguém fizer uma requisição POST para "/users", o Gin vai
	// chamar automaticamente a função "CreateUserHandle" do pacote usercontrollers.
	// POST é usado convencionalmente para CRIAR recursos novos.
	app.POST("/users", usercontrollers.CreateUserHandle)

	// Define uma rota do tipo PUT no caminho "/users/:id".
	// O ":id" é um parâmetro dinâmico de rota — significa que o valor muda.
	// Ex: "/users/1", "/users/42", "/users/abc" — todos caem nessa rota.
	// Dentro do handler você consegue ler esse valor com c.Param("id").
	// PUT é usado convencionalmente para ATUALIZAR um recurso existente.
	app.PUT("/users/:id", usercontrollers.UpdateUserHandle)

	app.GET("/users/:id", usercontrollers.FindUserByIdController)

	// app.Run() inicia o servidor e fica ouvindo requisições.
	// Por padrão, sobe na porta 8080 (http://localhost:8080).
	// Ele retorna um erro se não conseguir subir (ex: porta já em uso).
	// O "if err := ...; err != nil" é um padrão muito comum em Go para
	// executar algo e já verificar se deu erro na mesma linha.
	if err := app.Run(); err != nil {
		// log.Fatalf imprime a mensagem de erro formatada no terminal
		// e encerra o programa com código de saída 1 (indicando falha).
		// O "%v" é o formato genérico do Go para imprimir qualquer valor.
		log.Fatalf("failed to run server: %v", err)
	}
}
