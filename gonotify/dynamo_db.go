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
	for tokenID, locations := range tokenLocationMap {
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
				"TokenID":   &types.AttributeValueMemberS{Value: tokenID},
				"Locations": avList,
			},
		}

		// Perform PutItem operation
		_, err = svc.PutItem(context.Background(), input)
		if err != nil {
			fmt.Println("Error putting item into DynamoDB:", err)
		} else {
			fmt.Printf("Successfully added data for Token ID %s to DynamoDB\n", tokenID)
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
		}
		avList[i] = &types.AttributeValueMemberM{Value: avMap}
	}
	return &types.AttributeValueMemberL{Value: avList}, nil
}
