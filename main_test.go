package main

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCEP(t *testing.T) {
	tests := []struct {
		name           string
		cep            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "CEP válido",
			cep:            "37130093",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"temp_C":`,
		},
		{
			name:           "CEP inválido",
			cep:            "37009999",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "invalid zipcode",
		},
		{
			name:           "CEP não encontrado",
			cep:            "371300930",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "can not find zipcode",
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/?"+"cep="+tt.cep, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handleCEP)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler retornou status code errado: recebido %v esperado %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				// Para respostas bem-sucedidas, verificar se é um JSON válido com as temperaturas
				var response TemperatureResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Falha ao decodificar JSON: %v", err)
				}

				// Verificar se as temperaturas foram calculadas corretamente
				if response.TempF != response.TempC*1.8+32 {
					t.Errorf("Conversão Fahrenheit incorreta")
				}
				if response.TempK != response.TempC+273.15 {
					t.Errorf("Conversão Kelvin incorreta")
				}
			} else {
				// Para respostas de erro, verificar a mensagem exata
				if rr.Body.String() != tt.expectedBody {
					t.Errorf("handler retornou body inesperado: recebido %v esperado %v",
						rr.Body.String(), tt.expectedBody)
				}
			}
		})
	}
}

// TestTemperatureConversion testa as conversões de temperatura isoladamente
func TestTemperatureConversion(t *testing.T) {
	// Definir uma margem de erro aceitável para comparações de ponto flutuante
	const epsilon = 0.00001

	tests := []struct {
		tempC float64
		tempF float64
		tempK float64
	}{
		{0, 32, 273.15},
		{100, 212, 373.15},
		{-40, -40, 233.15},
	}

	for _, tt := range tests {
		calculatedF := tt.tempC*1.8 + 32
		calculatedK := tt.tempC + 273.15

		// Usar função auxiliar para comparar valores float
		if !almostEqual(calculatedF, tt.tempF, epsilon) {
			t.Errorf("Conversão Fahrenheit incorreta para %v°C: esperado %v°F, recebido %v°F",
				tt.tempC, tt.tempF, calculatedF)
		}

		if !almostEqual(calculatedK, tt.tempK, epsilon) {
			t.Errorf("Conversão Kelvin incorreta para %v°C: esperado %v°K, recebido %v°K",
				tt.tempC, tt.tempK, calculatedK)
		}
	}
}

// Função auxiliar para comparar valores float com margem de erro
func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}
