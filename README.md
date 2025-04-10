# Pag-Simples


Este é um projeto desenvolvido em Golang para gerenciamento simplificado de transferências financeiras.


## Descrição

  

O projeto oferece endpoints para o gerenciamento de usuários e a realização de transferências entre contas. Ele foi feito construído pensando em uma arquitetura limpa que permitisse modularização de componentes que facilitam a sua expansão e manutenabilidade.

  

O projeto é composto por uma API REST com endpoints para o gerenciamento de usuários e a realização de transferências entre contas. Foi desenvolvido com uma arquitetura limpa, visando a modularização dos componentes, o que facilita tanto a expansão quanto a manutenção do sistema.

  

## Tecnologias usadas

  

-  **Golang**: A linguagem de programação usada no projeto.

-  **Chi**: Framework para roteamento HTTP em Go.

-  **Docker**: Para facilitar o deploy e a execução do serviço em diferentes ambientes.

  

## Como rodar o projeto

 
### 1. Clonar o repositório

Clone o repositório para a sua máquina local:
```bash

git  clone  https://github.com/laianemuckler/pag-simples.git

cd  pag-simples
```
### 2. Rodar o projeto localmente

Passo 1: Instalar as dependências
```bash
go mod tidy
```
Passo 2: Executar o projeto
```bash
go run cmd/api/main.go
```
### 3. Subir o projeto com Docker

Passo 1: Construir a imagem do Docker
```bash
docker-compose build
```
Passo 2: Subir os containers
```bash
docker-compose up
```
Passo 3: Reconstruir a imagem após modificações
```bash
docker-compose up --build
```
### 4. Rodar os testes

Rodar todos os testes unitários (de forma detalhada)
```bash
go test -v ./...
```
Rodar os testes com cobertura
```bash
go test -cover ./...
```

## Endpoints

### **GET** `/users/{id}` 
Obtém as informações de um usuário pelo ID.

### **GET** `/users` 
Obtém a lista de todos os usuários.

### **POST** `/users` 
Cria um novo usuário com as informações fornecidas no corpo da requisição.


### **POST** `/transfer` 
Realiza uma transferência entre dois usuários, especificando o valor, o pagador e o recebedor.

#### Exemplo de requisição:

```json
{
  "value": 100.0,
  "payer": 4,
  "payee": 15
}
```
#### Exemplo de requisição com cURL

Você pode realizar uma transferência usando o `curl` com o seguinte comando:

```bash
curl -X POST http://localhost:8080/transfer \
-H "Content-Type: application/json" \
-d '{
  "value": 100.0,
  "payer": 4,
  "payee": 15
}'
```

#### Resposta (sucesso)

Você pode realizar uma transferência usando o `curl` com o seguinte comando:

```bash
curl -X POST http://localhost:8080/transfer \
-H "Content-Type: application/json" \
-d '{
  "value": 100.0,
  "payer": 4,
  "payee": 15
}'