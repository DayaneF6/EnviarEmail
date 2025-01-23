package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go/aws"
)

type RequestBody struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

type Response struct {
	Message      string `json:"message"`
	SNSMessageID string `json:"snsMessageId,omitempty"`
	Error        string `json:"error,omitempty"`
}

func handler(ctx context.Context, event json.RawMessage) (Response, error) {
	// Configura o cliente SNS
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Erro ao carregar configuração AWS: %v", err)
		return Response{Error: "Erro interno ao configurar AWS"}, nil
	}

	snsClient := sns.NewFromConfig(cfg)

	// Parseando o corpo do evento
	var reqBody RequestBody
	err = json.Unmarshal(event, &reqBody)
	if err != nil {
		log.Printf("Erro ao processar corpo da requisição: %v", err)
		return Response{Error: "Formato inválido na requisição"}, nil
	}

	// Validação de campos obrigatórios
	if reqBody.Email == "" || reqBody.Message == "" {
		return Response{Error: "Os campos 'email' e 'message' são obrigatórios"}, nil
	}

	topicArn := os.Getenv("arn:aws:sns:us-east-2:124355683655:EnviarMensagem")
	if topicArn == "" {
		log.Println("ARN do tópico SNS não configurado")
		return Response{Error: "ARN do tópico SNS não configurado"}, nil
	}

	// Configurando os parâmetros da publicação
	input := &sns.PublishInput{
		Message:  aws.String(fmt.Sprintf("Mensagem de: %s\n\n%s", reqBody.Email, reqBody.Message)),
		Subject:  aws.String("Nova mensagem recebida"),
		TopicArn: aws.String(topicArn),
	}

	// Publicando no SNS
	output, err := snsClient.Publish(ctx, input)
	if err != nil {
		log.Printf("Erro ao publicar mensagem no SNS: %v", err)
		return Response{Error: "Erro ao enviar mensagem ao SNS"}, nil
	}

	log.Printf("Mensagem publicada com sucesso. ID: %s", *output.MessageId)
	return Response{
		Message:      "Mensagem enviada com sucesso via SNS!",
		SNSMessageID: *output.MessageId,
	}, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var body json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		resp, err := handler(ctx, body)
		if err != nil {
			http.Error(w, "Erro interno", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
