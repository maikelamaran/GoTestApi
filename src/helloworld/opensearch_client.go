package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type Alarma struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

func crearClienteOpenSearch() *opensearch.Client {
	cfg := opensearch.Config{
		Addresses: []string{"https://localhost:9200"},
		Username:  "admin",
		Password:  "admin",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creando cliente OpenSearch: %v", err)
	}

	return client
}

func BuscarAlarmasPorNivel(nivel string) []Alarma {
	client := crearClienteOpenSearch()

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"level": nivel,
			},
		},
	}

	queryBody, _ := json.Marshal(query) //aqui serializo el query

	req := opensearchapi.SearchRequest{ //armo la peticion
		Index: []string{"alarmas"},        //es como el nombre de la tabla en sql
		Body:  bytes.NewReader(queryBody), //necesita ioreader
	}

	res, err := req.Do(context.Background(), client) //mando la peticion
	if err != nil {
		log.Fatalf("Error buscando en OpenSearch: %v", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body) //leo el body que trajo como en getlogs

	//por lo que vi siempre vienen en esta estructura de opensearch
	//hits -> hits -> _source
	//de ahi que el for tambien sea diferente
	var resultado struct {
		Hits struct {
			Hits []struct {
				Source Alarma `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(body, &resultado); err != nil {
		log.Fatalf("Error parseando resultado: %v", err)
	}

	alarmas := []Alarma{}
	for _, hit := range resultado.Hits.Hits {
		alarmas = append(alarmas, hit.Source)
	}

	return alarmas
}
func saveToOpenSearch(logs []Alarma) {
	// Configurar cliente de OpenSearch (conexión insegura para desarrollo)
	// cfg := opensearch.Config{
	// 	Addresses: []string{"https://localhost:9200"},
	// 	Username:  "admin",
	// 	Password:  "admin",
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{
	// 			InsecureSkipVerify: true, // Solo para entornos de desarrollo
	// 		},
	// 	},
	// }

	// // Crear cliente OpenSearch
	// client, err := opensearch.NewClient(cfg)
	// if err != nil {
	// 	log.Fatalf("Error creando cliente OpenSearch: %v", err)
	// }
	client := crearClienteOpenSearch()
	// Verificar conexión
	info, err := client.Info()
	if err != nil {
		log.Fatalf("Error conectando a OpenSearch: %v", err)
	}
	defer info.Body.Close()
	fmt.Println("\n✅ Conexión exitosa a OpenSearch")

	// Nombre del índice donde se guardarán las alarmas
	indexName := "alarmas"

	// Indexar cada alarma
	fmt.Println("\n⏳ Subiendo alarmas a OpenSearch...")
	for _, alarma := range logs {
		// Convertir a JSON
		document, _ := json.Marshal(alarma)

		// Configurar petición de indexado
		req := opensearchapi.IndexRequest{
			Index: indexName,
			Body:  bytes.NewReader(document),
		}

		// Ejecutar petición
		res, err := req.Do(context.Background(), client)
		if err != nil {
			log.Printf("Error indexando documento: %v", err)
			continue
		}
		defer res.Body.Close()

		// Manejar errores de respuesta
		if res.IsError() {
			log.Printf("Error en respuesta: %s", res.String())
		} else {
			fmt.Printf("📄 Documento indexado: %s\n", alarma.Name)
		}
	}
	fmt.Println("✅ Todas las alarmas han sido guardadas en OpenSearch")
}
