package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"crypto/tls"
)

type ViaCEPResponse struct {
	Cep        string `json:"cep"`
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {
	// Desabilitar a verificação do certificado SSL
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/", handleCEP)
	http.ListenAndServe(":8080", nil)
}

func handleCEP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	cep := r.URL.Query().Get("cep")

	// Consultar ViaCEP
	viaCEPResp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		fmt.Println("Erro ao consultar ViaCEP:", err)
		http.Error(w, "Erro ao consultar CEP", http.StatusInternalServerError)
		return
	}
	defer viaCEPResp.Body.Close()

	// Verificar o status code da resposta
	if viaCEPResp.StatusCode == http.StatusBadRequest {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "can not find zipcode")
		return
	}

	var location ViaCEPResponse
	if err := json.NewDecoder(viaCEPResp.Body).Decode(&location); err != nil {
		fmt.Println("Erro ao decodificar resposta do ViaCEP:", err)
		http.Error(w, "Erro ao decodificar resposta do ViaCEP", http.StatusInternalServerError)
		return
	}

	// Verificar se o ViaCEP retornou erro
	if location.Erro == "true" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, "invalid zipcode")
		return
	}

	// Consultar WeatherAPI
	weatherAPIKey := "3a4554c1d5e54545a46225119240311"
	weatherResp, err := http.Get(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", weatherAPIKey, location.Localidade))
	if err != nil {
		fmt.Println("Erro ao consultar temperatura:", err)
		http.Error(w, "Erro ao consultar temperatura", http.StatusInternalServerError)
		return
	}
	defer weatherResp.Body.Close()

	var weather WeatherResponse
	if err := json.NewDecoder(weatherResp.Body).Decode(&weather); err != nil {
		fmt.Println("Erro ao decodificar resposta do WeatherAPI:", err)
		http.Error(w, "Erro ao decodificar resposta do WeatherAPI", http.StatusInternalServerError)
		return
	}

	// Calcular temperaturas
	tempC := weather.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	// Preparar resposta
	response := TemperatureResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
