package authorization

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AuthorizationService interface {
	CheckAuthorization() (bool, error)
}

type authorizationService struct{}

func NewAuthorizationService() AuthorizationService {
	return &authorizationService{}
}

func (s *authorizationService) CheckAuthorization() (bool, error) {
	url := "https://util.devi.tools/api/v2/authorize"

	log.Printf("Iniciando requisição de autorização para o serviço: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Erro ao criar requisição: %v", err)
		return false, fmt.Errorf("erro ao criar requisição: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao enviar requisição para %s: %v", url, err)
		return false, fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Resposta do serviço de autorização, status code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro ao consultar serviço de autorização, código de status: %d", resp.StatusCode)
		return false, fmt.Errorf("erro ao consultar serviço de autorização, código de status: %d", resp.StatusCode)
	}

	var authorizationResponse Authorization
	if err := json.NewDecoder(resp.Body).Decode(&authorizationResponse); err != nil {
		log.Printf("Erro ao decodificar a resposta da autorização: %v", err)
		return false, fmt.Errorf("erro ao decodificar a resposta: %v", err)
	}

	log.Printf("Resposta do serviço de autorização: %s", authorizationResponse.Status)

	if authorizationResponse.Status != "success" {
		log.Printf("Serviço de autorização falhou, status: %s", authorizationResponse.Status)
		return false, fmt.Errorf("serviço de autorização falhou, status: %s", authorizationResponse.Status)
	}

	log.Println("Serviço de autorização bem-sucedido")

	return authorizationResponse.Data.Authorization, nil
}
