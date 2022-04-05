package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/events/test"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// defining a struct for creating and storing mock dynamodb session
type mockedGetItem struct {
	dynamodbiface.DynamoDBAPI
	ResponseQuery dynamodb.GetItemOutput
}

type mockedGetItemZeroValueCheck struct {
	dynamodbiface.DynamoDBAPI
	ResponseQuery dynamodb.GetItemOutput
}

type RiskAssessmentTestData struct {
	Pk string `json:"pk"`
	Sk string `json:"sk"`
}

func (d mockedGetItem) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {

	var err = errors.New("item not found")
	val, ok := input.Key["pk"]
	if ok {
		if strings.HasPrefix(*val.S, "sellerName") {
			pk := dynamodb.AttributeValue{}
			pk.SetS("sellerName#Renew Power")
			sk := dynamodb.AttributeValue{}
			sk.SetS("registryProjectName#Renewable Solar Power Project by ReNew Solar Power Private Limited#plantName#Abha Solarfarms Limited")
			contractStatus := dynamodb.AttributeValue{}
			contractStatus.SetS("LOC")
			projectStatus := dynamodb.AttributeValue{}
			projectStatus.SetS("Constructed")
			approvalStatus := dynamodb.AttributeValue{}
			approvalStatus.SetS("Under Registration")
			hostCountryApproval := dynamodb.AttributeValue{}
			hostCountryApproval.SetS("Yes")
			methodology := dynamodb.AttributeValue{}
			methodology.SetS("Solar Power")
			resp := make(map[string]*dynamodb.AttributeValue)
			resp["pk"] = &pk
			resp["sk"] = &sk
			resp["contractStatus"] = &contractStatus
			resp["projectStatus"] = &projectStatus
			resp["approvalStatus"] = &approvalStatus
			resp["hostCountryApproval"] = &hostCountryApproval
			resp["methodology"] = &methodology

			output := &dynamodb.GetItemOutput{
				Item: resp,
			}
			return output, nil
		} else if strings.HasPrefix(*val.S, "categoryName") {

			pk := dynamodb.AttributeValue{}
			sk := dynamodb.AttributeValue{}
			subCategoryValue := dynamodb.AttributeValue{}
			categoryWeightage := dynamodb.AttributeValue{}
			fmt.Println("*val.S: ", *val.S)
			if strings.Contains(*val.S, "Contract") {
				pk.SetS("categoryName#Contract Status")
				sk.SetS("subCategoryName#LOC")
				subCategoryValue.SetS("50%")
				categoryWeightage.SetS("3/13")
			} else if strings.Contains(*val.S, "Project") {
				pk.SetS("categoryName#Project Status")
				sk.SetS("subCategoryName#Constructed")
				subCategoryValue.SetS("100%")
				categoryWeightage.SetS("1/13")
			} else if strings.Contains(*val.S, "Host") {
				pk.SetS("categoryName#Host Country Approval")
				sk.SetS("subCategoryName#Yes")
				subCategoryValue.SetS("100%")
				categoryWeightage.SetS("1/13")
			} else if strings.Contains(*val.S, "Registration") {
				pk.SetS("categoryName#Registration Status")
				sk.SetS("subCategoryName#Under Registration")
				subCategoryValue.SetS("75%")
				categoryWeightage.SetS("6/13")
			} else if strings.Contains(*val.S, "Methodology") {
				pk.SetS("categoryName#Methodology")
				sk.SetS("subCategoryName#Solar Power")
				subCategoryValue.SetS("100%")
				categoryWeightage.SetS("2/13")
			}

			resp := make(map[string]*dynamodb.AttributeValue)
			resp["pk"] = &pk
			resp["sk"] = &sk
			resp["subCategoryValue"] = &subCategoryValue
			resp["categoryWeightage"] = &categoryWeightage

			output := &dynamodb.GetItemOutput{
				Item: resp,
			}
			return output, nil
		}

	}
	return &dynamodb.GetItemOutput{}, err
}

// for 2nd test, to handle zero value divide error
func (d mockedGetItemZeroValueCheck) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {

	var err = errors.New("item not found")
	val, ok := input.Key["pk"]
	if ok {
		if strings.HasPrefix(*val.S, "sellerName") { // structuring getItem for a different getItem function for sellerName in main
			pk := dynamodb.AttributeValue{}
			pk.SetS("sellerName#Renew Power")
			sk := dynamodb.AttributeValue{}
			sk.SetS("registryProjectName#Renewable Solar Power Project by ReNew Solar Power Private Limited#plantName#Abha Solarfarms Limited")
			contractStatus := dynamodb.AttributeValue{}
			contractStatus.SetS("LOC")
			projectStatus := dynamodb.AttributeValue{}
			projectStatus.SetS("Constructed")
			approvalStatus := dynamodb.AttributeValue{}
			approvalStatus.SetS("Under Registration")
			hostCountryApproval := dynamodb.AttributeValue{}
			hostCountryApproval.SetS("Yes")
			methodology := dynamodb.AttributeValue{}
			methodology.SetS("Solar Power")
			resp := make(map[string]*dynamodb.AttributeValue)
			resp["pk"] = &pk
			resp["sk"] = &sk
			resp["contractStatus"] = &contractStatus
			resp["projectStatus"] = &projectStatus
			resp["approvalStatus"] = &approvalStatus
			resp["hostCountryApproval"] = &hostCountryApproval
			resp["methodology"] = &methodology

			output := &dynamodb.GetItemOutput{
				Item: resp,
			}
			return output, nil
		} else if strings.HasPrefix(*val.S, "categoryName") { // structuring getItem for a different getItem function for catergoryName in main
			pk := dynamodb.AttributeValue{}
			sk := dynamodb.AttributeValue{}
			subCategoryValue := dynamodb.AttributeValue{}
			categoryWeightage := dynamodb.AttributeValue{}
			fmt.Println("*val.S: ", *val.S)
			if strings.Contains(*val.S, "Contract") {
				pk.SetS("categoryName#Contract Status")
				sk.SetS("subCategoryName#LOC")
				subCategoryValue.SetS("50%")
				categoryWeightage.SetS("3/13")
			} else if strings.Contains(*val.S, "Project") {
				pk.SetS("categoryName#Project Status")
				sk.SetS("subCategoryName#Constructed")
				subCategoryValue.SetS("100%")
				categoryWeightage.SetS("1/13")
			} else if strings.Contains(*val.S, "Host") {
				pk.SetS("categoryName#Host Country Approval")
				sk.SetS("subCategoryName#Yes")
				subCategoryValue.SetS("100%")
				categoryWeightage.SetS("1/13")
			} else if strings.Contains(*val.S, "Registration") {
				pk.SetS("categoryName#Registration Status")
				sk.SetS("subCategoryName#Under Registration")
				subCategoryValue.SetS("75%")
				categoryWeightage.SetS("6/13")
			} else if strings.Contains(*val.S, "Methodology") {
				pk.SetS("categoryName#Methodology")
				sk.SetS("subCategoryName#Solar Power")
				subCategoryValue.SetS("100%")
				categoryWeightage.SetS("2/0")
			}

			resp := make(map[string]*dynamodb.AttributeValue)
			resp["pk"] = &pk
			resp["sk"] = &sk
			resp["subCategoryValue"] = &subCategoryValue
			resp["categoryWeightage"] = &categoryWeightage

			output := &dynamodb.GetItemOutput{
				Item: resp,
			}
			return output, nil
		}

	}
	return &dynamodb.GetItemOutput{}, err
}

// mocking update item of main function
func (d mockedGetItem) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return &dynamodb.UpdateItemOutput{}, nil // empty because updateItem generally dont necessarily need output
}

func TestLambdaReadHandler(t *testing.T) {

	// reading sample dynamodb stream formatted data(../testData/dynamodb-event.json) from external saved file to send as an input stream to handler
	inputJSON := test.ReadJSONFromFile(t, "../testData/dynamodb-event.json")
	var inputEvent events.DynamoDBEvent
	if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}
	// structuring a mocked dynamodb session for get item for test1
	m := mockedGetItem{
		ResponseQuery: dynamodb.GetItemOutput{},
	}

	d := deps{
		ddb: m, //mock session
	}

	// test1:- for verifying correct risk score
	t.Run("check risk score", func(t *testing.T) {
		riskScoreString, _ := d.handler(inputEvent)
		expectedRiskScore := "76.9230769230769"
		actualRiskScore := riskScoreString

		assert.Equal(t, expectedRiskScore, actualRiskScore)
	})

	// structuring a mocked dynamodb session for get item for test2
	m1 := mockedGetItemZeroValueCheck{
		ResponseQuery: dynamodb.GetItemOutput{},
	}

	d1 := deps{
		ddb: m1, //mock session
	}

	// test2:- checking error for division by zero
	t.Run("check division with zero", func(t *testing.T) {
		_, errString := d1.handler(inputEvent)
		if errString == nil {
			panic("cannot divide by zero(0) error should be there")
		}

		expectedErrorMessage := "error, cannot divide with zero(0)"
		actualErrorMessage := errString.Error()
		assert.Equal(t, expectedErrorMessage, actualErrorMessage)
	})

}
