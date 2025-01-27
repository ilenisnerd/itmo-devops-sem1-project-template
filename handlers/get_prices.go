package handlers

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"log"
	"net/http"
	"project-sem-1/models"
	"strconv"
	"time"
)

func GetPrices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query(`
            SELECT id, created_at, name, category, price 
            FROM prices
        `)

		if err != nil {
			log.Printf("Error querying database: %v", err)
			http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
			return
		}

		var allPrices []models.PriceStruct

		for rows.Next() {
			var (
				idInt     int
				createdAt time.Time
				name      string
				category  string
				priceVal  float64
			)
			if err := rows.Scan(&idInt, &createdAt, &name, &category, &priceVal); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			allPrices = append(allPrices, models.PriceStruct{
				ID:        strconv.Itoa(idInt),
				CreatedAt: createdAt,
				Name:      name,
				Category:  category,
				Price:     priceVal,
			})
		}
		if rows.Err() != nil {
			log.Printf("Error after rows.Next(): %v", rows.Err())
			http.Error(w, "Failed to read rows", http.StatusInternalServerError)
			return
		}
		rows.Close()

		csvBuffer := &bytes.Buffer{}
		writer := csv.NewWriter(csvBuffer)

		writer.Write([]string{"id", "name", "category", "price", "create_date"})

		for _, p := range allPrices {
			record := []string{
				p.ID,
				p.Name,
				p.Category,
				strconv.FormatFloat(p.Price, 'f', 2, 64),
				p.CreatedAt.Format("2006-01-02"),
			}
			writer.Write(record)
		}
		writer.Flush()

		if err := writer.Error(); err != nil {
			log.Printf("Error finalizing CSV: %v", err)
			http.Error(w, "Failed to write CSV", http.StatusInternalServerError)
			return
		}

		zipBuffer := &bytes.Buffer{}
		zipWriter := zip.NewWriter(zipBuffer)

		csvFile, err := zipWriter.Create("data.csv")
		if err != nil {
			log.Printf("Error creating file in ZIP: %v", err)
			http.Error(w, "Failed to create ZIP", http.StatusInternalServerError)
			return
		}

		if _, err := csvFile.Write(csvBuffer.Bytes()); err != nil {
			log.Printf("Error writing CSV to ZIP: %v", err)
			http.Error(w, "Failed to write ZIP", http.StatusInternalServerError)
			return
		}

		if err := zipWriter.Close(); err != nil {
			log.Printf("Error closing ZIP writer: %v", err)
			http.Error(w, "Failed to close ZIP", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=data.zip")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(zipBuffer.Bytes()); err != nil {
			log.Printf("Error sending ZIP file: %v", err)
		}
	}
}
