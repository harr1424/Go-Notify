package gonotify

import (
	"fmt"
	"log"

	"github.com/sideshow/apns2"
	payload "github.com/sideshow/apns2/payload"
	token "github.com/sideshow/apns2/token"
)

func sendPushNotification(targetToken string, location string) {

	// load signing key from file
	authKey, err := token.AuthKeyFromFile("apnkey.p8")
	if err != nil {
		log.Println("Error sending push notification:", err)
	}

	// Generate JWT used for APNs
	requestToken := &token.Token{
		AuthKey: authKey,
		KeyID:   signingKey,
		TeamID:  teamID,
	}

	// Construct alert information from alert struct
	alertSubtitle := fmt.Sprintf("Frost Alert for %s", location)
	payload := payload.NewPayload().AlertSubtitle(alertSubtitle)

	notification := &apns2.Notification{
		DeviceToken: targetToken,
		Topic:       bundleId,
		Payload:     payload,
	}

	client := apns2.NewTokenClient(requestToken).Production()
	result, err := client.Push(notification)
	if err != nil {
		log.Println("Error Sending Push Notification:", err)
	}
	log.Println("Sent notification with response:", result)
}
