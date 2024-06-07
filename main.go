package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // MySQL driver

)

// Data struct (same as before)
type Data struct {
	ID          int
	Coordenadas string `json:"coordenadas"`
	Descricao   string
}

func main() {
	// Database connection (same as before)
	db, err := sql.Open("mysql", "user:pass8@tcp(url:3306)/mapeamento")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Handle HTTP requests
	http.HandleFunc("/api/random-data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		// Get total number of rows
		//var count int
		//err := db.QueryRow("SELECT 1'").Scan(&count)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}

		// Generate random offset
		//offset := rand.Intn(count) // Random number between 0 and count-1

		// Fetch random row with offset
		var d Data
		err = db.QueryRow("select  id, coordenadas, descricao from (SELECT id, coordenadas, descricao FROM mapeamento.geometriaespacialjson WHERE  coordenadas    LIKE '%MultiPolygon%' limit 100) a ORDER BY RAND() LIMIT 1;").Scan(&d.ID, &d.Coordenadas, &d.Descricao)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := fmt.Sprintf(`
		{
			"type": "Feature",
			"geometry": %v
		}
		`, d.Coordenadas)

		w.Header().Set("Content-Type", "application/json")
		var payload interface{}                     //The interface where we will save the converted JSON data.
		err = json.Unmarshal([]byte(res), &payload) // Convert JSON data into interface{} type
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(payload) // Encode a single Data object
	})

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
