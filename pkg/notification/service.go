package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SendNotification(request NotificationRequest) error {
	url := "https://util.devi.tools/api/v1/notify"

	log.Printf("Iniciando requisição para enviar a notificação para o serviço: %s", url)

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Erro ao codificar a requisição para JSON: %v", err)
		return fmt.Errorf("erro ao codificar a requisição: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Erro ao criar a requisição: %v", err)
		return fmt.Errorf("erro ao criar a requisição: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao enviar a requisição para %s: %v", url, err)
		return fmt.Errorf("erro ao enviar a requisição: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Resposta do serviço de notificação, status code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro ao enviar a notificação, status code: %d", resp.StatusCode)
		return fmt.Errorf("erro ao enviar a notificação, status code: %d", resp.StatusCode)
	}

	log.Println("Notificação enviada com sucesso")

	return nil
}
