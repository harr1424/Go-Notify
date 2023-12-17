package gonotify

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const dynamoDBTableName = "DeviceTokensAndLocations"

func UpdateTokenLocationMap(tokenLocationMap map[string][]Location) {
	// Create a DynamoDB svc
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// Iterate over map and store data in DynamoDB
	for token, locations := range tokenLocationMap {
		// Convert locations to DynamoDB AttributeValue
		avList, err := attributeValueList(locations)
		if err != nil {
			fmt.Println("Error marshaling location:", err)
			continue
		}

		// Prepare input for PutItem operation
		input := &dynamodb.PutItemInput{
			TableName: aws.String(dynamoDBTableName),
			Item: map[string]types.AttributeValue{
				"Token":     &types.AttributeValueMemberS{Value: token},
				"Locations": avList,
			},
		}

		// Perform PutItem operation
		_, err = svc.PutItem(context.Background(), input)
		if err != nil {
			fmt.Println("Error putting item into DynamoDB:", err)
		}
	}
}

// Helper function to convert a slice of locations to DynamoDB AttributeValue
func attributeValueList(locations []Location) (types.AttributeValue, error) {
	avList := make([]types.AttributeValue, len(locations))
	for i, loc := range locations {
		avMap := map[string]types.AttributeValue{
			"Latitude":  &types.AttributeValueMemberN{Value: loc.Latitude},
			"Longitude": &types.AttributeValueMemberN{Value: loc.Longitude},
			"Name":      &types.AttributeValueMemberS{Value: loc.Name},
			"Unit":      &types.AttributeValueMemberS{Value: loc.Unit},
		}
		avList[i] = &types.AttributeValueMemberM{Value: avMap}
	}
	return &types.AttributeValueMemberL{Value: avList}, nil
}

func RetrieveTokenLocationMap() (map[string][]Location, bool, error) {
	// Create a DynamoDB svc
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := dynamodb.NewFromConfig(cfg)

	// Prepare input for Scan operation
	input := &dynamodb.ScanInput{
		TableName: aws.String(dynamoDBTableName),
	}

	// Perform Scan operation
	result, err := svc.Scan(context.Background(), input)
	if err != nil {
		fmt.Println("Error scanning DynamoDB table:", err)
		return nil, true, err
	}

	// Convert DynamoDB items to map[string][]Location
	tokenLocationMap := make(map[string][]Location)
	for _, item := range result.Items {
		token := item["Token"].(*types.AttributeValueMemberS).Value
		locationsAttribute := item["Locations"].(*types.AttributeValueMemberL).Value

		var locations []Location
		for _, locationAttr := range locationsAttribute {
			lat := locationAttr.(*types.AttributeValueMemberM).Value["Latitude"].(*types.AttributeValueMemberN).Value
			lon := locationAttr.(*types.AttributeValueMemberM).Value["Longitude"].(*types.AttributeValueMemberN).Value
			name := locationAttr.(*types.AttributeValueMemberM).Value["Name"].(*types.AttributeValueMemberS).Value
			unit := locationAttr.(*types.AttributeValueMemberM).Value["Unit"].(*types.AttributeValueMemberS).Value

			locations = append(locations, Location{
				Latitude:  lat,
				Longitude: lon,
				Name:      name,
				Unit:      unit,
			})
		}

		tokenLocationMap[token] = locations
	}

	isEmpty := len(result.Items) == 0

	return tokenLocationMap, isEmpty, nil
}
