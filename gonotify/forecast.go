package gonotify

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
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

func getForecastAndNotify(targetDevice string, location Location) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&hourly=temperature_2m&forecast_days=3", location.Latitude, location.Longitude)

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
			outputTimeFormat = "January 2"
		)

		for i, t := range forecast.Hourly.Time {
			// Parse the time string into a time.Time object
			timeParsed, err := time.Parse(inputTimeFormat, t)
			if err != nil {
				fmt.Println("For timestamp:", t)
				fmt.Println("Error parsing time:", err)
				return
			}

			temp := forecast.Hourly.Temperature2m[i]
			time := timeParsed.Format(outputTimeFormat)

			var tempString string

			if location.Unit == "F" {
				tempString = celsiusToFahrenheitString(temp)
			} else {
				tempString = strconv.FormatFloat(temp, 'f', -1, 64)
			}

			if temp < 3.0 {
				fmt.Printf("Sending frost notification to %s: \n", targetDevice)
				sendPushNotification(targetDevice, location.Name, time, tempString, location.Unit)
				break
			}
		}

		fmt.Printf("Finished analyzing forecast for (%s): \n", location.Name)

	} else {
		fmt.Printf("Failed to retrieve data. Status code: %d\n", response.StatusCode)
	}
}

func CheckAllLocationsForFrost() {
	uniqueLocations := make(map[string]struct{})

	for token, allLocations := range TokenLocationMap {
		for _, location := range allLocations {
			// Generate a unique key for each location
			locationKey := fmt.Sprintf("%s|%s", location.Latitude, location.Longitude)

			// Check if the location has already been processed
			if _, processed := uniqueLocations[locationKey]; !processed {
				// Process the location
				getForecastAndNotify(token, location)

				// Mark the location as processed
				uniqueLocations[locationKey] = struct{}{}
			}
		}
	}
}

func celsiusToFahrenheitString(celsius float64) string {
	fahrenheit := math.Round((celsius * 9 / 5) + 32)
	return fmt.Sprintf("%.0f", fahrenheit)
}
