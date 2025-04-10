// countries-service/main.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Country struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Continent  string `json:"continent"`
	Population string `json:"population"`
}

func main() {
	port := getEnv("PORT", "8080")
	csvPath := getEnv("CSV_PATH", "countries.csv")

	countries, err := loadCountriesFromCSV(csvPath)
	if err != nil {
		log.Fatalf("Failed to load countries: %v", err)
	}

	http.HandleFunc("/countries", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(countries)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting countries service on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func loadCountriesFromCSV(filePath string) ([]Country, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read CSV: %v", err)
	}

	var countries []Country
	// Skip header row
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}
		if len(record) >= 4 {
			country := Country{
				Code:       record[0],
				Name:       record[1],
				Continent:  record[2],
				Population: record[3],
			}
			countries = append(countries, country)
		}
	}

	return countries, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
