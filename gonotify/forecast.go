package gonotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ForecastResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationTimeMS     float64 `json:"generationtime_ms"`
	UTCOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	HourlyUnits          struct {
		Time          string `json:"time"`
		Temperature2m string `json:"temperature_2m"`
	} `json:"hourly_units"`
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

func GetForecast() {
	url := "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m&forecast_days=3"

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error contacting forecast API::", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading the forecast API response body:", err)
			return
		}

		var forecast ForecastResponse

		err = json.Unmarshal(body, &forecast)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		const (
			inputTimeFormat  = "2006-01-02T15:04"
			outputTimeFormat = "January 2nd 3:04 PM"
		)

		for i, t := range forecast.Hourly.Time {
			// Parse the time string into a time.Time object
			timeParsed, err := time.Parse(inputTimeFormat, t)
			if err != nil {
				fmt.Println("For timestamp:", t)
				fmt.Println("Error parsing time:", err)
				return
			}

			// Print the time and temperature for each entry in a human-readable format
			fmt.Printf("Time: %s, Temperature: %.1f°C\n", timeParsed.Format(outputTimeFormat), forecast.Hourly.Temperature2m[i])
		}

	} else {
		fmt.Printf("Failed to retrieve data. Status code: %d\n", response.StatusCode)
	}
}
