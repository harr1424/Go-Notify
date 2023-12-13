package gonotify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Location struct to represent a geographical Location
type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Struct used to deserialize a payload sent when adding a location
type LocationAddRequest struct {
	Token    string   `json:"token"`
	Location Location `json:"Location"`
}

var tokenLocationMap = make(map[string][]Location)

// Called when the register endpoint is contacted
// Expects to receive POST data describing an iOS device token
func RegisterToken(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/register" {
		http.NotFound(res, req)
		return
	}

	var newToken string

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&newToken); err != nil {
		log.Println("Could not create new token from register request): ", err)
		return
	}

	if _, exists := tokenLocationMap[newToken]; !exists {
		tokenLocationMap[newToken] = []Location{}
	} else {
		fmt.Println("Token already exists in DeviceTokenLocationMap.")
	}

	res.WriteHeader(http.StatusCreated)
}

// Called when the add_location endpoint is contacted
// Expects to receive POST data describing an iOS device token and location
func HandleLocationAdd(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/add_location" {
		http.NotFound(res, req)
		return
	}

	var requestBody LocationAddRequest

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&requestBody); err != nil {
		log.Println("Could not add new token from add request:", err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract token and location from the request payload
	token := requestBody.Token
	newLocation := requestBody.Location

	// Check if the token exists in the map
	if locations, exists := tokenLocationMap[token]; !exists {
		// If the token doesn't exist, associate it with a new slice containing the new location
		tokenLocationMap[token] = []Location{newLocation}
		fmt.Println("Location added for the token:", token)
	} else {
		// Token exists, check if the location already exists
		locationExists := false
		for _, loc := range locations {
			if loc == newLocation {
				locationExists = true
				break
			}
		}

		// If the location doesn't exist, add it to the slice
		if !locationExists {
			tokenLocationMap[token] = append(tokenLocationMap[token], newLocation)
			fmt.Println("Location added for the token:", token)
		} else {
			fmt.Println("Location already exists for the token:", token)
		}
	}

	// Print the updated map
	fmt.Printf("Updated tokenLocationMap: %v\n", tokenLocationMap)

	// Respond with success status
	res.WriteHeader(http.StatusCreated)
}