package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

type PageData struct {
	Countries []Country
	Error     string
}

func main() {
	port := getEnv("PORT", "8081")
	countriesServiceURL := getEnv("COUNTRIES_SERVICE_URL", "http://localhost:8080/countries")

	tmpl, err := template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Countries List</title>
    <style>
        table {
            border-collapse: collapse;
            width: 100%;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        .error {
            color: red;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <h1>Countries of the World</h1>
    
    {{if .Error}}
        <p class="error">{{.Error}}</p>
    {{else}}
        <table>
            <tr>
                <th>Code</th>
                <th>Name</th>
                <th>Continent</th>
                <th>Population</th>
            </tr>
            {{range .Countries}}
            <tr>
                <td>{{.Code}}</td>
                <td>{{.Name}}</td>
                <td>{{.Continent}}</td>
                <td>{{.Population}}</td>
            </tr>
            {{end}}
        </table>
    {{end}}
</body>
</html>
`)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{}

		countries, err := fetchCountries(countriesServiceURL)
		if err != nil {
			data.Error = fmt.Sprintf("Failed to fetch countries: %v", err)
		} else {
			data.Countries = countries
		}

		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Template execution error: %v", err)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting client service on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func fetchCountries(url string) ([]Country, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var countries []Country
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return nil, fmt.Errorf("JSON decoding failed: %v", err)
	}

	return countries, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
