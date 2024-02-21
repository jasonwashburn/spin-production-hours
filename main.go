package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
)

type ProductionHour struct {
	ModelRun time.Time `json:"modelRun"`
	Forecast int       `json:"forecast"`
}

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		// Assuming URL pattern is /productionHours/yyyy/mm/dd/hh

		pathParts := splitPath(r.URL.Path)
		if len(pathParts) != 6 { // 6 parts: "", "productionHours", "yyyy", "mm", "dd", "hh"
			http.Error(w, "Invalid URL format. Expecting /productionHours/yyyy/mm/dd/hh", http.StatusBadRequest)
			return
		}

		year, _ := strconv.Atoi(pathParts[2])
		month, _ := strconv.Atoi(pathParts[3])
		day, _ := strconv.Atoi(pathParts[4])
		hour, _ := strconv.Atoi(pathParts[5])

		validTime := time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC)

		subtractedTime := validTime.Add(-42 * time.Hour)

		startModel := time.Date(subtractedTime.Year(), subtractedTime.Month(), subtractedTime.Day(), 0, 0, 0, 0, time.UTC)

		var productionHours []ProductionHour
		modelRun := startModel

		productionForecasts := []int{36, 39, 42}
		for modelRun.Before(validTime) {
			for _, forecast := range productionForecasts {
				if modelRun.Add(time.Duration(forecast)*time.Hour) == validTime {
					productionHours = append(productionHours, ProductionHour{ModelRun: modelRun, Forecast: forecast})
				}
			}

			modelRun = modelRun.Add(6 * time.Hour)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(productionHours); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}

func main() {}
