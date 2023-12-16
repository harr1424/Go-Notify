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
}

func sendPushNotification(targetToken string, time string, value string) {
	newAlert := alert{Time: time, Value: value}

	// load signing key from file
	authKey, err := APNS.AuthKeyFromFile("apnkey.p8")
	if err != nil {
		log.Println("Token Error:", err)
	}

	// Generate JWT used for APNs
	requestToken := &APNS.Token{
		AuthKey: authKey,
		KeyID:   signingKey,
		TeamID:  teamID,
	}

	// Construct alert information from alert struct
	alertSubtitle := fmt.Sprintf("A temperature of %s is expected on %s", newAlert.Value, newAlert.Time)
	payload := PAYLOAD.NewPayload().Alert("Frost Alert").AlertSubtitle(alertSubtitle)

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
