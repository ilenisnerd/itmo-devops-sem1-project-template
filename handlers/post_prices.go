package handlers

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	_ "os"
	"path/filepath"
	_ "project-sem-1/database"
	"project-sem-1/models"
	"strconv"
	"time"
)

func PostPrices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error retrieving file: %v", err)
			http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		zipBuffer := &bytes.Buffer{}
		if _, err := io.Copy(zipBuffer, file); err != nil {
			log.Printf("Error reading file: %v", err)
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		zipReader, err := zip.NewReader(bytes.NewReader(zipBuffer.Bytes()), int64(zipBuffer.Len()))
		if err != nil {
			log.Printf("Error opening zip: %v", err)
			http.Error(w, "Invalid zip file", http.StatusBadRequest)
			return
		}

		var validRecords []models.PriceStruct

		for _, zf := range zipReader.File {
			if filepath.Ext(zf.Name) != ".csv" {
				continue
			}
			csvFile, err := zf.Open()
			if err != nil {
				log.Printf("Error opening CSV in zip: %v", err)
				continue
			}
			defer csvFile.Close()

			reader := csv.NewReader(csvFile)

			if _, err := reader.Read(); err != nil {
				log.Printf("Error reading header: %v", err)
				continue
			}

			for {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("Error reading CSV record: %v", err)
					continue
				}
				if len(record) < 5 {
					log.Printf("Invalid record: %v", record)
					continue
				}

				idStr := record[0]
				name := record[1]
				category := record[2]
				priceStr := record[3]
				dateStr := record[4]

				priceVal, err := strconv.ParseFloat(priceStr, 64)
				if err != nil {
					log.Printf("Invalid price: %v", priceStr)
					continue
				}

				createdAt, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					log.Printf("Invalid date: %v", dateStr)
					continue
				}

				validRecords = append(validRecords, models.PriceStruct{
					ID:        idStr,
					CreatedAt: createdAt,
					Name:      name,
					Category:  category,
					Price:     priceVal,
				})
			}
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}
		defer func() { _ = tx.Rollback() }()

		var successfullyProcessed int

		for _, rec := range validRecords {
			_, err := tx.Exec(`
                INSERT INTO prices (id, created_at, name, category, price)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT (id) DO NOTHING
            `, rec.ID, rec.CreatedAt, rec.Name, rec.Category, rec.Price)
			if err != nil {
				log.Printf("DB insert error: %v", err)
				continue
			}

			successfullyProcessed++
		}

		var dbCategories int
		var dbTotalPrice float64

		row := tx.QueryRow(`
            SELECT COUNT(DISTINCT category), COALESCE(SUM(price), 0)
            FROM prices
        `)
		if err := row.Scan(&dbCategories, &dbTotalPrice); err != nil {
			log.Printf("Failed to scan totals: %v", err)
			http.Error(w, "Failed to calculate totals", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		resp := models.PostStruct{
			TotalItems:      successfullyProcessed,
			TotalCategories: dbCategories,
			TotalPrice:      math.Round(dbTotalPrice*100) / 100,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding JSON: %v", err)
		}
	}
}
