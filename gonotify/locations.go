package gonotify


/*
This file contains logic that will associate locations with specific device tokens 
when users enable frost alerts for a given locations.

In the future, this file will be updated to store tokens and locations in a CSP database offering
encryption.
*/

// Location struct to represent a geographical location
type Location struct {
	Latitude  string
	Longitude string
}

// A map associating tokens with locations to be notified about
var tokenLocationMap = make(map[token][]Location)

func AddLocation(targetToken token, location Location) {
	// Check if the slice for the token exists in the map
	locations, exists := tokenLocationMap[targetToken]

	// If it doesn't exist, create a new slice
	if !exists {
		locations = make([]Location, 0)
	}

	// Append the location to the slice
	locations = append(locations, location)

	// Update the map with the new slice of locations
	tokenLocationMap[targetToken] = locations
}