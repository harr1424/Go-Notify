package gonotify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Struct to unmarshall JSON object describing a token
type TokenRequest struct {
	Token string `json:"token"`
}

// Location struct to represent a geographical Location
type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Name      string `json:"name"`
}

// Struct used to deserialize a payload sent when adding a location
type LocationAddRequest struct {
	Token    string   `json:"token"`
	Location Location `json:"Location"`
}

var TokenLocationMap map[string][]Location

func ReadRemoteTableContents() error {
	result, isEmpty, err := RetrieveTokenLocationMap()
	if err != nil {
		return fmt.Errorf("error reading remote table contents: %v", err)
	}

	if isEmpty {
		TokenLocationMap = make(map[string][]Location)
		fmt.Println("DynamoDB table was found to be empty.")
	} else {
		TokenLocationMap = result
		fmt.Println("DynamoDB table was read successfully.")
	}

	return nil
}

// Called when the register endpoint is contacted
// Expects to receive POST data describing an iOS device token
func RegisterToken(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/register" {
		http.NotFound(res, req)
		return
	}

	var tokenRequest TokenRequest

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&tokenRequest); err != nil {
		log.Println("Could not create new token from register request): ", err)
		return
	}

	newToken := tokenRequest.Token

	if _, exists := TokenLocationMap[newToken]; !exists {
		TokenLocationMap[newToken] = []Location{}
		UpdateTokenLocationMap(TokenLocationMap)
		fmt.Println("Added token: ", newToken)
	} else {
		fmt.Println("Token already exists in DeviceTokenLocationMap")
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
	if locations, exists := TokenLocationMap[token]; !exists {
		// If the token doesn't exist, associate it with a new slice containing the new location
		TokenLocationMap[token] = []Location{newLocation}
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
			TokenLocationMap[token] = append(TokenLocationMap[token], newLocation)
			UpdateTokenLocationMap(TokenLocationMap)
			fmt.Println("Location added for the token:", token)
		} else {
			fmt.Println("Location already exists for the token:", token)
		}
	}

	// Print the updated map
	fmt.Printf("Updated TokenLocationMap: %v\n", TokenLocationMap)

	// Respond with success status
	res.WriteHeader(http.StatusCreated)
}

// Called when the remove_location endpoint is contacted
// Expects to receive POST data describing an iOS device token and location
func HandleLocationRemove(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/remove_location" {
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
	locationToRemove := requestBody.Location

	// Check if the token exists in the map
	if locations, exists := TokenLocationMap[token]; !exists {
		fmt.Println("Token not found:", token)
	} else {
		// Token exists, check if the location already exists
		locationIndex := -1
		for i, loc := range locations {
			if loc.Latitude == locationToRemove.Latitude && loc.Longitude == locationToRemove.Longitude {
				locationIndex = i
				break
			}
		}

		// If the location exists, remove it from the slice
		if locationIndex != -1 {
			TokenLocationMap[token] = append(locations[:locationIndex], locations[locationIndex+1:]...)
			UpdateTokenLocationMap(TokenLocationMap)
			fmt.Println("Location removed for the token:", token)
		} else {
			fmt.Println("Location not found for the token:", token)
		}
	}

	// Print the updated map
	fmt.Printf("Updated TokenLocationMap: %v\n", TokenLocationMap)

	// Respond with success status
	res.WriteHeader(http.StatusCreated)
}
