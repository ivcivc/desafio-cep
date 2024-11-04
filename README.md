# API de Consulta de Temperatura por CEP

Esta API permite consultar a temperatura atual de uma cidade a partir de um CEP brasileiro. A temperatura é retornada em três escalas: Celsius, Fahrenheit e Kelvin.

## 🚀 Funcionalidades

- Consulta de CEP via API ViaCEP
- Consulta de temperatura via WeatherAPI
- Conversão automática entre escalas de temperatura
- Tratamento de erros para CEPs inválidos ou não encontrados

## 📋 Pré-requisitos

- Go 1.21 ou superior
- Conta na WeatherAPI (para obter a chave de API)
- Docker (opcional, para containerização)
- Conta no Google Cloud Platform (opcional, para deploy)


## 🔍 Endpoints

### GET /?cep={cep}

Retorna a temperatura atual da cidade correspondente ao CEP informado.

#### Parâmetros
- `cep`: CEP brasileiro (8 dígitos)

#### Respostas

**Sucesso (200 OK)**
Retorna um JSON com as seguintes informações:
- `temp_C`: Temperatura em Celsius
- `temp_F`: Temperatura em Fahrenheit
- `temp_K`: Temperatura em Kelvin

**Erro (400)**
Retorna um JSON com a mensagem de erro:
- `error`: can not find zipcode


**CEP Inválido (422)**
Retorna um JSON com a mensagem de erro:
- `error`: invalid zipcode

## 🧪 Acessar a API (Google Cloud Run)
https://hello-275133257525.us-central1.run.app/?cep=


## 🧪 Testes

Execute os testes automatizados com:
```
go test -v
``` 


## 🐳 Docker

### Construir a imagem
```
docker build -t temperature-api .
```

### Executar o container    
```
docker run -p 8080:8080 temperature-api
``` 


## 🚀 Deploy no Google Cloud Run

1. Build e push da imagem:
```
gcloud builds submit --tag gcr.io/temperature-api-396212/temperature-api
```


## 🛠️ Tecnologias Utilizadas

- [Go](https://golang.org/)
- [ViaCEP API](https://viacep.com.br/)
- [WeatherAPI](https://www.weatherapi.com/)
- [Docker](https://www.docker.com/)
- [Google Cloud Run](https://cloud.google.com/run)

## ✒️ Estrutura do Projeto

- `main.go`: Função principal que gerencia as requisições e respostas da API.
- `main_test.go`: Conjunto de testes automatizados para validar o comportamento da API.
- `go.mod`: Arquivo de módulo Go que define as dependências do projeto.
- `Dockerfile`: Define a configuração do container Docker.
- `README.md`: Documentação do projeto.

