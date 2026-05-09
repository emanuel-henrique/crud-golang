// Pacote "database" — responsável exclusivamente por gerenciar a conexão com o banco.
// Todos os outros pacotes que precisam do banco importam daqui, evitando repetição.
package database

import (
	// Pacote nativo do Go para trabalhar com bancos de dados relacionais.
	// Ele define interfaces genéricas (como *sql.DB) que funcionam com qualquer banco —
	// quem implementa os detalhes específicos do banco é o "driver" (no caso, o pq).
	"database/sql"

	// Pacote nativo para formatar strings — usado aqui para montar a string de conexão.
	"fmt"

	// Pacote nativo para logar erros fatais no terminal.
	"log"

	// Pacote nativo para ler variáveis de ambiente do sistema operacional.
	// os.Getenv("NOME") retorna o valor da variável de ambiente "NOME",
	// ou uma string vazia "" se ela não existir.
	"os"

	// godotenv é um pacote externo que lê um arquivo ".env" e carrega
	// as variáveis definidas nele como variáveis de ambiente do processo.
	// Isso permite guardar senhas e configurações fora do código fonte.
	"github.com/joho/godotenv"

	// Driver do PostgreSQL para Go. O "_" na frente significa import anônimo —
	// não vamos chamar nenhuma função desse pacote diretamente.
	// Mas ao importá-lo, ele se auto-registra no pacote "database/sql" por baixo dos panos,
	// habilitando o uso do banco "postgres". Sem esse import, sql.Open("postgres",...) falharia.
	_ "github.com/lib/pq"
)

// ConectarBanco abre uma conexão com o banco PostgreSQL e a retorna.
// Retorna *sql.DB (ponteiro para a conexão) e um error.
func ConectarBanco() (*sql.DB, error) {

	// godotenv.Load() lê o arquivo ".env" na raiz do projeto e registra
	// cada linha como uma variável de ambiente. Ex: HOST=localhost vira os.Getenv("HOST") = "localhost".
	// É fundamental que seja chamado ANTES de qualquer os.Getenv(),
	// caso contrário as variáveis ainda não existem e retornam "" (string vazia).
	err := godotenv.Load()
	if err != nil {
		// log.Fatal imprime a mensagem e encerra o programa imediatamente.
		// Diferente de "return", ele mata o processo — se não tem .env, não tem como continuar.
		log.Fatal("Error loading .env file")
	}

	// As variáveis são lidas AQUI, depois do godotenv.Load(), garantindo que
	// o arquivo .env já foi carregado e os valores estão disponíveis no ambiente.
	// Antes estavam num bloco "var" global, onde eram lidas na inicialização do programa
	// — antes do Load() ser chamado — e por isso sempre retornavam "".
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	// Monta a string de conexão no formato que o driver PostgreSQL espera.
	// fmt.Sprintf preenche cada %s com o argumento correspondente, em ordem.
	// Todos os placeholders são %s pois os.Getenv() SEMPRE retorna string —
	// antes, "port" usava %d (que espera int), gerando um valor inválido
	// como "port=%!d(string=5432)" e corrompendo a string de conexão inteira.
	// sslmode=require exige conexão criptografada — importante para bancos em produção/nuvem.
	stringConexao := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// sql.Open() NÃO abre a conexão de verdade ainda — ele só valida o formato
	// da string de conexão e prepara o objeto *sql.DB.
	// A conexão real só acontece na primeira query ou no db.Ping() abaixo.
	// O primeiro argumento "postgres" diz qual driver usar — foi registrado pelo import do "pq".
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		// panic() encerra o programa abruptamente e imprime o stack trace.
		// É mais agressivo que log.Fatal — normalmente usado para erros
		// que "nunca deveriam acontecer". Em produção, prefira retornar o erro.
		panic(err)
	}

	// db.Ping() tenta de fato abrir e testar a conexão com o banco.
	// Se o banco estiver fora do ar, a senha estiver errada ou o host inacessível,
	// é aqui que você vai descobrir — não no sql.Open().
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Confirma no terminal que a conexão foi bem-sucedida.
	// Útil para debug, mas em produção considere usar um logger estruturado
	// em vez de fmt.Println, pois ele não tem controle de nível (info, warn, error).
	fmt.Println("Conectado ao Banco de Dados: ", dbname)

	// Retorna a conexão ativa para quem chamou a função.
	// err aqui é nil (sem erro), já que passou pelo Ping com sucesso.
	return db, err
}
