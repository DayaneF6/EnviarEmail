package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Estrutura para armazenar os dados do formulário
type FormData struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

// Função para lidar com requisições POST
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Verifica se o método é POST
	if r.Method == http.MethodPost {
		var data FormData

		// Decodifica o corpo da requisição JSON para a estrutura FormData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Erro ao decodificar os dados", http.StatusBadRequest)
			return
		}

		// Verifica se os campos obrigatórios foram enviados
		if data.Email == "" || data.Message == "" {
			http.Error(w, "Os campos 'email' e 'message' são obrigatórios.", http.StatusBadRequest)
			return
		}

		// Exibe os dados no console (para fins de depuração)
		fmt.Printf("Email: %s\nMensagem: %s\n", data.Email, data.Message)

		// Envia a mensagem para o SNS
		err = sendMessageToSNS(data.Email, data.Message)
		if err != nil {
			http.Error(w, "Erro ao enviar mensagem para o SNS", http.StatusInternalServerError)
			return
		}

		// Resposta de sucesso
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Dados recebidos e mensagem enviada com sucesso!"))
	} else {
		// Se não for um POST, retorna um erro
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// Função para enviar a mensagem para o SNS
func sendMessageToSNS(email, message string) error {
	// Cria uma sessão AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"), // Substitua pela sua região AWS
	})
	if err != nil {
		return fmt.Errorf("erro ao criar sessão AWS: %v", err)
	}

	// Cria um cliente SNS
	svc := sns.New(sess)

	// Configura os parâmetros para o SNS
	params := &sns.PublishInput{
		Message:  aws.String(fmt.Sprintf("Mensagem de: %s\n\n%s", email, message)),
		Subject:  aws.String("Nova mensagem recebida"),
		TopicArn: aws.String("arn:aws:sns:us-east-2:124355683655:EnviarMensagem"),
	}

	// Publica no SNS
	_, err = svc.Publish(params)
	if err != nil {
		return fmt.Errorf("erro ao publicar no SNS: %v", err)
	}

	// Se a publicação for bem-sucedida, retorna nil
	return nil
}

func main() {
	// Rota que captura o email e a mensagem
	http.HandleFunc("/submit", handleRequest)

	// Inicia o servidor HTTP na porta 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Servidor iniciado na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
