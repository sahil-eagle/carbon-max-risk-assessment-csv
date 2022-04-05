package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/shopspring/decimal"
)

type deps struct {
	ddb dynamodbiface.DynamoDBAPI
}

type RiskAssessmentData struct {
	Pk                                     string `json:"pk"`
	Sk                                     string `json:"sk"`
	ProjectName                            string `json:"projectName"`
	SellerName                             string `json:"sellerName"`
	PlantName                              string `json:"plantName"`
	ProjectID                              string `json:"projectId"`
	ProjectType                            string `json:"projectType"`
	RegisteredMW                           string `json:"registeredMw"`
	VintageStart                           string `json:"vintageStart"`
	VintageEnd                             string `json:"vintageEnd"`
	NumberOfIssuedCredits                  string `json:"numberOfIssuedCredits"`
	CreditingPeriodStartDate               string `json:"creditingPeriodStartDate"`
	CreditingPeriodEndDate                 string `json:"creditingPeriodEndDate"`
	CODOfTheProject                        string `json:"codOfTheProject"`
	CapacityFactor                         string `json:"capacityFactor"`
	AverageAnnualVariationInCapacityFactor string `json:"averageAnnualVariationInCapacityFactor"`
	CurrentStatusOfSale                    string `json:"currentStatusOfSale"`
	DateOfIssuanceToCC                     string `json:"dateOfIssuanceToCC"`
	ContractStatus                         string `json:"contractStatus"`
	ProjectStatus                          string `json:"projectStatus"`
	HostCountryApproval                    string `json:"hostCountryApproval"`
	ApprovalStatus                         string `json:"approvalStatus"`
	Methodology                            string `json:"methodology"`
}

type RiskDefinitionsData struct {
	Pk                string `json:"pk"`
	Sk                string `json:"sk"`
	CategoryWeightage string `json:"categoryWeightage"`
	SubCategoryValue  string `json:"subCategoryValue"`
}

type RiskScoreData struct {
	Pk        string  `json:"pk"`
	Sk        string  `json:"sk"`
	RiskScore float64 `json:"riskScore"`
}

// function to get risk definition values for calculation of risk
func getRiskDefinitionValues(category, subCategory string, d *deps) (float64, string, error) {
	db := d.ddb
	input := &dynamodb.GetItemInput{
		TableName: aws.String("RiskAssessmentTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String("categoryName#" + category),
			},
			"sk": {
				S: aws.String("subCategoryName#" + subCategory),
			},
		},
	}

	result, _ := db.GetItem(input)
	res := RiskDefinitionsData{}
	dynamodbattribute.UnmarshalMap(result.Item, &res)

	categoryWeightageString := res.CategoryWeightage
	subCategoryValue := res.SubCategoryValue

	categoryWeightageSplitValue := strings.Split(categoryWeightageString, "/")

	valueOne, _ := decimal.NewFromString(categoryWeightageSplitValue[0])
	valueTwo, _ := decimal.NewFromString(categoryWeightageSplitValue[1])
	tempStringCompareValue, _ := decimal.NewFromString("0")
	if valueTwo.Equal(tempStringCompareValue) {
		return 0, "", errors.New("cannot divide with 0(zero) values")
	}

	fmt.Println("categoryWeightageString: ", categoryWeightageString)
	categoryWeightage, _ := valueOne.Div(valueTwo).Float64()

	return categoryWeightage, subCategoryValue, nil
}

// risk calculation function
func calculatePlantRiskScore(contractCategoryWeight,
	projectCategoryWeight,
	approvalCategoryWeight,
	hostCountryApprovalCategoryWeight,
	methodologyCategoryWeight float64, contractSubCategoryValue,
	projectSubCategoryValue,
	approvalSubCategoryValue,
	hostCountryApprovalSubCategoryValue,
	methodologySubCategoryValue string) float64 {

	fmt.Println("contractCategoryWeight: ", contractCategoryWeight)
	fmt.Println("projectCategoryWeight: ", projectCategoryWeight)
	fmt.Println("hostCountryApprovalCategoryWeight: ", hostCountryApprovalCategoryWeight)

	contractSubCategoryValuePer, _ := strconv.ParseFloat(strings.Split(contractSubCategoryValue, "%")[0], 64)
	projectSubCategoryValuePer, _ := strconv.ParseFloat(strings.Split(projectSubCategoryValue, "%")[0], 64)
	approvalSubCategoryValuePer, _ := strconv.ParseFloat(strings.Split(approvalSubCategoryValue, "%")[0], 64)
	hostCountryApprovalSubCategoryValuePer, _ := strconv.ParseFloat(strings.Split(hostCountryApprovalSubCategoryValue, "%")[0], 64)
	methodologySubCategoryValuePer, _ := strconv.ParseFloat(strings.Split(methodologySubCategoryValue, "%")[0], 64)

	fmt.Println("contractSubCategoryValuePer: ", contractSubCategoryValuePer)
	fmt.Println("projectSubCategoryValuePer: ", projectSubCategoryValuePer)
	fmt.Println("approvalSubCategoryValuePer: ", approvalSubCategoryValuePer)
	fmt.Println("hostCountryApprovalSubCategoryValuePer: ", hostCountryApprovalSubCategoryValuePer)
	fmt.Println("methodologySubCategoryValuePer: ", methodologySubCategoryValuePer)

	// formula for risk score
	riskScore := (contractCategoryWeight * contractSubCategoryValuePer) +
		(projectCategoryWeight * projectSubCategoryValuePer) +
		(approvalCategoryWeight * approvalSubCategoryValuePer) +
		(hostCountryApprovalCategoryWeight * hostCountryApprovalSubCategoryValuePer) +
		(methodologyCategoryWeight * methodologySubCategoryValuePer)

	// risk score logic for 10%
	if contractSubCategoryValuePer == 10 || projectSubCategoryValuePer == 10 || approvalSubCategoryValuePer == 10 {
		tempRiskScore, _ := decimal.NewFromString("10.0")
		riskScore, _ = tempRiskScore.Float64()
	}
	fmt.Println("riskType: ", reflect.TypeOf(riskScore))

	return riskScore
}

// function for getting all data required for calculation and returnign risk score
func plantRiskComputation(pk string, sk string, d *deps) (float64, error) {
	db := d.ddb

	input := &dynamodb.GetItemInput{
		TableName: aws.String("RiskAssessmentTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(pk),
			},
			"sk": {
				S: aws.String(sk),
			},
		},
	}
	result, _ := db.GetItem(input)

	res := RiskAssessmentData{}
	dynamodbattribute.UnmarshalMap(result.Item, &res)

	contractStatus := res.ContractStatus
	projectStatus := res.ProjectStatus
	hostCountryApproval := res.HostCountryApproval
	approvalStatus := res.ApprovalStatus
	methodology := res.Methodology

	fmt.Println("contractStatus: ", contractStatus)
	fmt.Println("projectStatus: ", projectStatus)
	fmt.Println("hostCountryApproval: ", hostCountryApproval)

	// getting risk definitions value and sub-category value
	contractCategoryWeight, contractSubCategoryValue, contractErrValue := getRiskDefinitionValues("Contract Status", contractStatus, d)
	projectCategoryWeight, projectSubCategoryValue, projectErrValue := getRiskDefinitionValues("Project Status", projectStatus, d)
	approvalCategoryWeight, approvalSubCategoryValue, approvalErrValue := getRiskDefinitionValues("Registration Status", approvalStatus, d)
	hostCountryApprovalCategoryWeight, hostCountryApprovalSubCategoryValue, hostCountryApprovalErrValue := getRiskDefinitionValues("Host Country Approval", hostCountryApproval, d)
	methodologyCategoryWeight, methodologySubCategoryValue, methodologyErrValue := getRiskDefinitionValues("Methodology", methodology, d)

	if contractErrValue != nil || projectErrValue != nil || approvalErrValue != nil || hostCountryApprovalErrValue != nil || methodologyErrValue != nil {
		return 0, errors.New("error, cannot divide with zero(0)")
	}
	riskScore := calculatePlantRiskScore(contractCategoryWeight,
		projectCategoryWeight,
		approvalCategoryWeight,
		hostCountryApprovalCategoryWeight,
		methodologyCategoryWeight, contractSubCategoryValue,
		projectSubCategoryValue,
		approvalSubCategoryValue,
		hostCountryApprovalSubCategoryValue,
		methodologySubCategoryValue)

	fmt.Println("riskScore: ", riskScore)

	return riskScore, nil
}

// updating risk score in the table
func putRiskScore(pk, sk string, riskScore float64, d *deps) {
	db := d.ddb

	upd := expression.Set(expression.Name("riskScore"), expression.Value(riskScore))

	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		fmt.Println("Got error while building update expression")
	}

	tableName := "RiskAssessmentTable"
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(pk),
			},
			"sk": {
				S: aws.String(sk),
			},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	_, errUpdate := db.UpdateItem(input)
	if errUpdate != nil {
		fmt.Println(errUpdate)
	} else {
		fmt.Println("Risk assessed")
	}
}

func (d *deps) handler(e events.DynamoDBEvent) (string, error) {
	var riskScore float64
	var errValue error
	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)
		pk := record.Change.NewImage["pk"].String()
		sk := record.Change.NewImage["sk"].String()
		// logic for running db streams and risk assesses computation for only seller data
		if strings.HasPrefix(pk, "sellerName") {
			riskScore, errValue = plantRiskComputation(pk, sk, d)
			if errValue == nil {
				putRiskScore(pk, sk, riskScore, d)
			} else {
				return "", errValue
			}
		}
		// Print new values for attributes of type String
		// for name, value := range record.Change.NewImage {
		// 	if value.DataType() == events.DataTypeString {
		// 		fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())
		// 		// calculatePlantRiskScore()
		// 	}
		// }
	}
	riskScoreString := strconv.FormatFloat(riskScore, 'f', -1, 64)

	return riskScoreString, nil
}

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// starting a session of dynamodbiface
	d := deps{
		ddb: dynamodb.New(sess),
	}

	lambda.Start(d.handler)
}
