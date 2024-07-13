package gonotify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const dynamoDBTableName = "DeviceTokensAndLocations"

var svc *dynamodb.Client

func InitializeDynamoDBClient(client *dynamodb.Client) {
	svc = client
}

func UpdateTokenLocation(ctx context.Context, token string, locations []Location) error {
    // Convert locations to DynamoDB AttributeValue
    avList, err := attributeValueList(locations)
    if err != nil {
        return errors.New("error constructing avList: " + err.Error())
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
        return errors.New("error putting item into DynamoDB: " + err.Error())
    }

    return nil
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

func RetrieveTokenLocationMap(ctx context.Context) (map[string][]Location, bool, error) {

	// Prepare input for Scan operation
	input := &dynamodb.ScanInput{
		TableName: aws.String(dynamoDBTableName),
	}

	// Perform Scan operation
	result, err := svc.Scan(context.Background(), input)
	if err != nil {
		return nil, true, errors.New("Error scanning DynamoDB table: " + err.Error())
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
