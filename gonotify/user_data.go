package gonotify

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Struct to unmarshall JSON object describing a token
type TokenRequest struct {
	Token string `json:"token"`
}

// Struct to represent a geographical Location and assocated attributes
type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Name      string `json:"name"`
	Unit      string `json:"unit"`
}

// Struct used to deserialize a payload sent when adding or removing a location
type LocationAddRequest struct {
	Token    string   `json:"token"`
	Location Location `json:"location"`
}

// A map with token keys and location values
var TokenLocationMap map[string][]Location

// A context used to call DynamoDB methods
var ctx = context.Background()

func ReadRemoteTableContents() error {
	result, isEmpty, err := RetrieveTokenLocationMap(ctx)
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
        log.Println("Could not create new token from register request:", err)
        http.Error(res, "Invalid request body", http.StatusBadRequest)
        return
    }

    newToken := tokenRequest.Token

    if _, exists := TokenLocationMap[newToken]; !exists {
        TokenLocationMap[newToken] = []Location{}
        if err := UpdateTokenLocation(ctx, newToken, TokenLocationMap[newToken]); err != nil {
            http.Error(res, "API failed to register token", http.StatusInternalServerError)
            return
        }
        fmt.Println("Added token:", newToken)
    }

    res.WriteHeader(http.StatusCreated)
}

func HandleLocationAdd(res http.ResponseWriter, req *http.Request) {
    handleLocationUpdate(res, req, true)
}

func HandleLocationRemove(res http.ResponseWriter, req *http.Request) {
    handleLocationUpdate(res, req, false)
}

func handleLocationUpdate(res http.ResponseWriter, req *http.Request, isAdd bool) {
    var requestBody LocationAddRequest
    decoder := json.NewDecoder(req.Body)

    if err := decoder.Decode(&requestBody); err != nil {
        log.Println("Could not parse request:", err)
        http.Error(res, "Invalid request body", http.StatusBadRequest)
        return
    }

    token := requestBody.Token
    location := requestBody.Location

    if locations, exists := TokenLocationMap[token]; !exists {
        if isAdd {
            TokenLocationMap[token] = []Location{location}
            if err := UpdateTokenLocation(ctx, token, TokenLocationMap[token]); err != nil {
                http.Error(res, "API failed to register token", http.StatusInternalServerError)
                return
            }
            fmt.Println("Location added for the token:", token)
        } else {
            fmt.Println("Token not found:", token)
        }
    } else {
        if isAdd {
            if !locationExists(locations, location) {
                TokenLocationMap[token] = append(locations, location)
                if err := UpdateTokenLocation(ctx, token, TokenLocationMap[token]); err != nil {
                    http.Error(res, "API failed to add location", http.StatusInternalServerError)
                    return
                }
                fmt.Println("Location added for the token:", token)
            } else {
                fmt.Println("Location already exists for the token:", token)
            }
        } else {
            if index := findLocationIndex(locations, location); index != -1 {
                TokenLocationMap[token] = append(locations[:index], locations[index+1:]...)
                if err := UpdateTokenLocation(ctx, token, TokenLocationMap[token]); err != nil {
                    http.Error(res, "API failed to remove location", http.StatusInternalServerError)
                    return
                }
                fmt.Println("Location removed for the token:", token)
            } else {
                fmt.Println("Location not found for the token:", token)
            }
        }
    }

    res.WriteHeader(http.StatusCreated)
}

func locationExists(locations []Location, location Location) bool {
    for _, loc := range locations {
        if loc == location {
            return true
        }
    }
    return false
}

func findLocationIndex(locations []Location, location Location) int {
    for i, loc := range locations {
        if loc == location {
            return i
        }
    }
    return -1
}

