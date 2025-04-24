package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {

	url := "https://demo.utmstack.com/management/logs"
	authToken := "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJkZW1vIiwiYXV0aCI6IlJPTEVfQURNSU4sUk9MRV9VU0VSIiwiYXV0aGVudGljYXRlZCI6dHJ1ZSwiZXhwIjoxNzQ3NTczOTY1fQ._x-YvPRQ4yp3Jo6HH0oQjpiq9RglU84wsk76TKVjHnrCTB5pWQZoM7YKZzkynY7yuSnFz_U_QMEMI3nhIhfCzg"

	putAlarma(authToken, url, "🚨🚨🚨🚨🚨🚨 alarma creada por maikel", "ERROR")

	logs := getLogs(authToken, url)

	saveToOpenSearch(logs)

	resultados := BuscarAlarmasPorNivel("ERROR")
	fmt.Println("🔍 Alarmas con nivel ERROR:")
	for _, a := range resultados {
		fmt.Printf("🔸 %s - %s\n", a.Name, a.Level)
	}
}

func putAlarma(authToken, url, name, level string) {

	alarma := Alarma{
		Name:  name,
		Level: level,
	}

	jsonData, err := json.Marshal(alarma)
	if err != nil {
		log.Fatalf("Error serializando JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creando la request PUT: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	// Enviar petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error enviando PUT: %v", err)
	}
	defer resp.Body.Close()

	// Verificar respuesta
	if resp.StatusCode != http.StatusNoContent {
		log.Fatalf("PUT falló. Código: %d", resp.StatusCode)
	}

	fmt.Println("✅ PUT exitoso: alarma enviada")
}

func getLogs(authToken, url string) []Alarma {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creando request GET: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error haciendo GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: código HTTP %d", resp.StatusCode)
	}
	//leer cuerpo o devuelve en []bytes (Representa el JSON crudo, sin interpretar)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error leyendo la respuesta: %v", err)
	}

	// Deserializar JSON
	var logs []Alarma
	//Unmarshal Convierte esos bytes en structs de Go
	if err := json.Unmarshal(body, &logs); err != nil { //guardo en logs, unmarshal body
		log.Fatalf("Error parseando JSON: %v", err)
	}

	// Mostrar resultados en consola
	fmt.Println("📋 Lista de alarmas:")
	for _, a := range logs {
		fmt.Printf("🔸 Nombre: %s | Nivel: %s\n", a.Name, a.Level)
	}

	return logs
}
