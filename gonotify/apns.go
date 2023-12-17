package gonotify

import (
	"fmt"
	"log"

	"github.com/sideshow/apns2"
	PAYLOAD "github.com/sideshow/apns2/payload"
	APNS "github.com/sideshow/apns2/token"
)

type alert struct {
	Time  string `json:"time"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

func sendPushNotification(targetToken string, location string, time string, value string, unit string) {
	newAlert := alert{Time: time, Value: value, Unit: unit}

	// load signing key from file
	authKey, err := APNS.AuthKeyFromFile("apnkey.p8")
	if err != nil {
		log.Println("Error sending push notification:", err)
	}

	// Generate JWT used for APNs
	requestToken := &APNS.Token{
		AuthKey: authKey,
		KeyID:   signingKey,
		TeamID:  teamID,
	}

	// Construct alert information from alert struct
	//alertTitle := fmt.Sprintf("Frost Alert %s", location)
	alertSubtitle := fmt.Sprintf("%s: %s°%s on %s", location, newAlert.Value, newAlert.Unit, newAlert.Time)
	payload := PAYLOAD.NewPayload().AlertSubtitle(alertSubtitle)

	notification := &apns2.Notification{
		DeviceToken: targetToken,
		Topic:       bundleId,
		Payload:     payload,
	}

	client := apns2.NewTokenClient(requestToken)
	result, err := client.Push(notification)
	if err != nil {
		log.Println("Error Sending Push Notification:", err)
	}
	log.Println("Sent notification with response:", result)
}
