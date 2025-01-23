package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

		// Exibe os dados no console (para fins de depuração)
		fmt.Printf("Email: %s\nMensagem: %s\n", data.Email, data.Message)

		// Resposta de sucesso
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Dados recebidos com sucesso"))
	} else {
		// Se não for um POST, retorna um erro
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}

}

const AWS = require('aws-sdk');

// Configuração do cliente SNS
const sns = new AWS.SNS();

exports.handler = async (event) => {
    // Parseando o body do evento recebido
    const body = JSON.parse(event.body);
    const email = body.email;
    const message = body.message;

    // Verifica se os campos obrigatórios foram enviados
    if (!email || !message) {
        return {
            statusCode: 400,
            body: JSON.stringify({ error: "Os campos 'email' e 'message' são obrigatórios." }),
        };
    }

    // Configurando os parâmetros para o SNS
    const params = {
        Message: Mensagem de: ${email}\n\n${message}, // Conteúdo da mensagem
        Subject: "Nova mensagem recebida",
        TopicArn: "arn:aws:sns:us-east-2:123456789012:SeuTopicoSNS" // Substitua pelo ARN do seu tópico SNS
    };

    try {
        // Publicando no SNS
        const result = await sns.publish(params).promise();
        console.log("Mensagem publicada no SNS:", result);

        return {
            statusCode: 200,
            body: JSON.stringify({
                message: "Mensagem enviada com sucesso via SNS!",
                snsMessageId: result.MessageId,
            }),
        };
    } catch (error) {
        console.error("Erro ao publicar no SNS:", error);
        return {
            statusCode: 500,
            body: JSON.stringify({ error: "Erro ao enviar mensagem ao SNS." }),
        };
    }
};

func main() {
	// Rota que captura o email e a mensagem
	http.HandleFunc("/submit", handleRequest)

	// Inicia o servidor HTTP na porta 8080
	fmt.Println("Servidor iniciado na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
