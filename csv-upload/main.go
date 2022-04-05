// package main

// import (
// 	"context"
// 	"encoding/csv"
// 	"errors"
// 	"fmt"
// 	"os"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/aws/aws-lambda-go/lambda"
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb"
// 	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
// 	"github.com/aws/aws-sdk-go/service/s3"
// 	"github.com/aws/aws-sdk-go/service/s3/s3manager"
// 	"github.com/google/uuid"
// )

// type MyDataFromS3 struct {
// 	Pk                                     string `json:"pk"`
// 	Gsi1Pk                                 string `json:"gsi1Pk"`
// 	Gsi1Sk                                 string `json:"gsi1Sk"`
// 	Gsi2Pk                                 string `json:"gsi2Pk"`
// 	Gsi2Sk                                 string `json:"gsi2Sk"`
// 	ProjectName                            string `json:"projectName"`
// 	SellerName                             string `json:"sellerName"`
// 	PlantName                              string `json:"plantName"`
// 	ProjectID                              string `json:"projectId"`
// 	ProjectType                            string `json:"projectType"`
// 	RegisteredMW                           string `json:"registeredMw"`
// 	VintageStart                           string `json:"vintageStart"`
// 	VintageEnd                             string `json:"vintageEnd"`
// 	NumberOfIssuedCredits                  string `json:"numberOfIssuedCredits"`
// 	CreditingPeriodStartDate               string `json:"creditingPeriodStartDate"`
// 	CreditingPeriodEndDate                 string `json:"creditingPeriodEndDate"`
// 	CODOfTheProject                        string `json:"codOfTheProject"`
// 	CapacityFactor                         string `json:"capacityFactor"`
// 	AverageAnnualVariationInCapacityFactor string `json:"averageAnnualVariationInCapacityFactor"`
// 	CurrentStatusOfSale                    string `json:"currentStatusOfSale"`
// 	DateOfIssuanceToCC                     string `json:"dateOfIssuanceToCC"`
// 	ContractStatus                         string `json:"contractStatus"`
// 	ProjectStatus                          string `json:"projectStatus"`
// 	HostCountryApproval                    string `json:"hostCountryApproval"`
// 	ApprovalStatus                         string `json:"approvalStatus"`
// 	Methodology                            string `json:"methodology"`
// }

// func deleteItemFromS3(sess *session.Session, bucket *string, item *string) error {
// 	svc := s3.New(sess)

// 	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
// 		Bucket: bucket,
// 		Key:    item,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
// 		Bucket: bucket,
// 		Key:    item,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func getDataFromS3File(bucket string, s3File string) [][]string {
// 	//the only writable directory in the lambda is /tmp
// 	file, err := os.Create("/tmp/" + s3File)
// 	if err != nil {
// 		fmt.Println("Unable to open file %q, %v", s3File, err)
// 	}

// 	defer file.Close()

// 	sess, _ := session.NewSession(&aws.Config{
// 		Region: aws.String("ap-south-1")},
// 	)

// 	downloader := s3manager.NewDownloader(sess)

// 	_, err = downloader.Download(file,
// 		&s3.GetObjectInput{
// 			Bucket:                  aws.String(bucket),
// 			Key:                     aws.String(s3File),
// 			ResponseContentType:     aws.String("text/csv"),
// 			ResponseContentEncoding: aws.String("utf-8"),
// 		})
// 	if err != nil {
// 		fmt.Println("Unable to download s3File %q, %v", s3File, err)
// 	}

// 	// fmt.Println(file)
// 	// fmt.Println(file.Name())

// 	recordFile, err := os.Open(file.Name())
// 	if err != nil {
// 		fmt.Println("cannot open file", err)
// 	}

// 	reader := csv.NewReader(recordFile)
// 	dataCSV, err := reader.ReadAll()
// 	if err != nil {
// 		fmt.Println("Cannot read the file", err)
// 	}

// 	return dataCSV
// }

// func extractData(data [][]string) ([]string, [][]string) {
// 	var headers []string
// 	var rows [][]string

// 	headers = data[:][0]
// 	rows = data[1:]

// 	return headers, rows
// }

// func insertIntoDynamoDB(headerDataToInsert []string, rowDataToInsert [][]string, bucketName string, fileName string) error {

// 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// 		SharedConfigState: session.SharedConfigEnable,
// 	}))

// 	svc := dynamodb.New(sess)

// 	for i := 0; i < len(rowDataToInsert); i = i + 1 {
// 		//validation for plantName
// 		if rowDataToInsert[i][3] == "" {
// 			deleteItemFromS3(sess, &bucketName, &fileName)
// 			return errors.New("plantName is a required value in the csv sheet")
// 		}
// 	}

// 	for i := 0; i < len(rowDataToInsert); i = i + 1 {
// 		data := MyDataFromS3{
// 			Pk:                                     uuid.New().String(),
// 			Gsi1Pk:                                 "sellerName#" + rowDataToInsert[i][0],
// 			Gsi1Sk:                                 "projectId#" + rowDataToInsert[i][1] + "#plantName#" + rowDataToInsert[i][3],
// 			Gsi2Pk:                                 "",
// 			Gsi2Sk:                                 "",
// 			SellerName:                             rowDataToInsert[i][0],
// 			ProjectID:                              rowDataToInsert[i][1],
// 			ProjectName:                            rowDataToInsert[i][2],
// 			PlantName:                              rowDataToInsert[i][3],
// 			ProjectType:                            rowDataToInsert[i][4],
// 			RegisteredMW:                           rowDataToInsert[i][5],
// 			VintageStart:                           rowDataToInsert[i][6],
// 			VintageEnd:                             rowDataToInsert[i][7],
// 			NumberOfIssuedCredits:                  rowDataToInsert[i][8],
// 			CreditingPeriodStartDate:               rowDataToInsert[i][9],
// 			CreditingPeriodEndDate:                 rowDataToInsert[i][10],
// 			CODOfTheProject:                        rowDataToInsert[i][11],
// 			CapacityFactor:                         rowDataToInsert[i][12],
// 			AverageAnnualVariationInCapacityFactor: rowDataToInsert[i][13],
// 			CurrentStatusOfSale:                    rowDataToInsert[i][14],
// 			DateOfIssuanceToCC:                     rowDataToInsert[i][15],
// 			ContractStatus:                         rowDataToInsert[i][16],
// 			ProjectStatus:                          rowDataToInsert[i][17],
// 			HostCountryApproval:                    rowDataToInsert[i][18],
// 			ApprovalStatus:                         rowDataToInsert[i][19],
// 			Methodology:                            rowDataToInsert[i][20],
// 		}

// 		av, err := dynamodbattribute.MarshalMap(data)
// 		if err != nil {
// 			fmt.Println("Got error marshalling new movie item:", av, err)
// 		}

// 		tableName := "RiskAssessmentTable"
// 		input := &dynamodb.PutItemInput{
// 			Item:      av,
// 			TableName: aws.String(tableName),
// 		}

// 		_, err = svc.PutItem(input)
// 		if err != nil {
// 			fmt.Println("Got error calling PutItem:", err)
// 		}
// 	}
// 	return nil
// }

// func HandleRequest(ctx context.Context, s3Event events.S3Event) error {
// 	bukcetName := s3Event.Records[0].S3.Bucket.Name
// 	fileName := s3Event.Records[0].S3.Object.Key
// 	for _, record := range s3Event.Records {
// 		s3 := record.S3
// 		csvData := getDataFromS3File(s3.Bucket.Name, s3.Object.Key)
// 		headerDataToInsert, rowDataToInsert := extractData(csvData)

// 		fmt.Println("headerDataToInsert: ", headerDataToInsert)
// 		fmt.Println("rowDataToInsert: ", rowDataToInsert)

// 		msg := insertIntoDynamoDB(headerDataToInsert, rowDataToInsert, bukcetName, fileName)
// 		if msg != nil {
// 			return msg
// 		}
// 	}

// 	fmt.Println(bukcetName, fileName)
// 	return nil
// }

// func main() {
// 	lambda.Start(HandleRequest)
// }
// *********************************

package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

//Event incoming event
type Event struct {
	Records []Record
}

type Record struct {
	EventSource    string
	EventSourceArn string
	AWSRegion      string
	S3             events.S3Entity
	SQS            events.SQSMessage
}

type TaskData struct {
	Pk           string            `json:"pk"`
	Sk           string            `json:"sk"`
	CreateTime   string            `json:"createTime"`
	UpdateTime   string            `json:"updateTime"`
	TaskProgress string            `json:"taskProgress"`
	ErrorLogs    map[string]string `json:"errorLogs"`
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

func getDataFromS3File(bucket string, s3File string) [][]string {
	//the only writable directory in the lambda is /tmp
	file, err := os.Create("/tmp/" + s3File)
	if err != nil {
		fmt.Println("Unable to open file %q, %v", s3File, err)
	}

	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket:                  aws.String(bucket),
			Key:                     aws.String(s3File),
			ResponseContentType:     aws.String("text/csv"),
			ResponseContentEncoding: aws.String("utf-8"),
		})
	if err != nil {
		fmt.Println("Unable to download s3File %q, %v", s3File, err)
	}

	recordFile, err := os.Open(file.Name())
	if err != nil {
		fmt.Println("cannot open file", err)
	}

	reader := csv.NewReader(recordFile)
	dataCSV, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Cannot read the file", err)
	}

	return dataCSV
}

func extractData(data [][]string) ([]string, [][]string) {
	var headers []string
	var rows [][]string

	headers = data[:][0]
	rows = data[1:]

	return headers, rows
}

func checkPlantExists(pk string, sk string) bool {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	db := dynamodb.New(sess)

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

	return len(result.Item) != 0
}

func insertIntoDynamoDB(headerDataToInsert []string, rowDataToInsert [][]string, bucketName string, fileName string) map[string]map[string]string {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	taskErrorLogsMap := make(map[string]map[string]string)

	for i := 0; i < len(rowDataToInsert); i = i + 1 {
		currentProjectName := rowDataToInsert[i][2]
		pk := "sellerName#" + rowDataToInsert[i][0]
		sk := "registryProjectName#" + currentProjectName + "#plantName#" + rowDataToInsert[i][3]

		taskErrorLogs := make(map[string]string)
		if checkPlantExists(pk, sk) {
			taskErrorLogs["plant already exist"] = fmt.Sprintf("the plant in row %d of csv already exist in database", i+2)
		}

		if rowDataToInsert[i][3] == "" {
			taskErrorLogs["required attribute plant name is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][4] == "" {
			taskErrorLogs["required attribute project type is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][5] == "" {
			taskErrorLogs["required attribute registered mw is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		// if rowDataToInsert[i][6] == "" {
		// 	taskErrorLogs["required attribute vintage start is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		// }
		// if rowDataToInsert[i][7] == "" {
		// 	taskErrorLogs["required attribute vintage end is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		// }
		// if rowDataToInsert[i][8] == "" {
		// 	taskErrorLogs["required attribute number of issued credits is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		// }
		if rowDataToInsert[i][9] == "" {
			taskErrorLogs["required attribute crediting period start date is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][10] == "" {
			taskErrorLogs["required attribute crediting period end date is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][14] == "" {
			taskErrorLogs["required attribute current status of sale is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][15] == "" {
			taskErrorLogs["required attribute date of issuance to cc is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][16] == "" {
			taskErrorLogs["required attribute contract status is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][17] == "" {
			taskErrorLogs["required attribute project status is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][18] == "" {
			taskErrorLogs["required attribute host country approval is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][19] == "" {
			taskErrorLogs["required attribute approval status is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][20] == "" {
			taskErrorLogs["required attribute methodology is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}

		if len(taskErrorLogs) != 0 {
			taskErrorLogsMap[fmt.Sprintf("missing required attributes from row %d", i+2)] = taskErrorLogs
			continue
		}

		data := RiskAssessmentData{
			Pk:                                     "sellerName#" + rowDataToInsert[i][0],
			Sk:                                     "registryProjectName#" + currentProjectName + "#plantName#" + rowDataToInsert[i][3],
			SellerName:                             rowDataToInsert[i][0],
			ProjectID:                              rowDataToInsert[i][1],
			ProjectName:                            rowDataToInsert[i][2],
			PlantName:                              rowDataToInsert[i][3],
			ProjectType:                            rowDataToInsert[i][4],
			RegisteredMW:                           rowDataToInsert[i][5],
			VintageStart:                           rowDataToInsert[i][6],
			VintageEnd:                             rowDataToInsert[i][7],
			NumberOfIssuedCredits:                  rowDataToInsert[i][8],
			CreditingPeriodStartDate:               rowDataToInsert[i][9],
			CreditingPeriodEndDate:                 rowDataToInsert[i][10],
			CODOfTheProject:                        rowDataToInsert[i][11],
			CapacityFactor:                         rowDataToInsert[i][12],
			AverageAnnualVariationInCapacityFactor: rowDataToInsert[i][13],
			CurrentStatusOfSale:                    rowDataToInsert[i][14],
			DateOfIssuanceToCC:                     rowDataToInsert[i][15],
			ContractStatus:                         rowDataToInsert[i][16],
			ProjectStatus:                          rowDataToInsert[i][17],
			HostCountryApproval:                    rowDataToInsert[i][18],
			ApprovalStatus:                         rowDataToInsert[i][19],
			Methodology:                            rowDataToInsert[i][20],
		}

		if currentProjectName == "" {
			data.Sk = "Not registered"
		}

		av, err := dynamodbattribute.MarshalMap(data)
		if err != nil {
			fmt.Println("Got error marshalling new movie item:", av, err)
		}

		tableName := "RiskAssessmentTable"
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			fmt.Println("Got error calling PutItem:", err)
		}
	}

	return taskErrorLogsMap
}

func updateDynamoDB(headerDataToInsert []string, rowDataToInsert [][]string, bucketName string, fileName string) map[string]map[string]string {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	db := dynamodb.New(sess)

	taskErrorLogsMap := make(map[string]map[string]string)

	for i := 0; i < len(rowDataToInsert); i = i + 1 {
		currentProjectName := rowDataToInsert[i][2]
		pk := "sellerName#" + rowDataToInsert[i][0]
		sk := "registryProjectName#" + currentProjectName + "#plantName#" + rowDataToInsert[i][3]

		taskErrorLogs := make(map[string]string)
		if !checkPlantExists(pk, sk) {
			taskErrorLogs["plant do not exist"] = fmt.Sprintf("the update requested for plant in row %d of csv does not exists in database", i+2)
		}

		if rowDataToInsert[i][3] == "" {
			taskErrorLogs["required attribute plant name is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][4] == "" {
			taskErrorLogs["required attribute project type is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][5] == "" {
			taskErrorLogs["required attribute registered mw is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		// if rowDataToInsert[i][6] == "" {
		// 	taskErrorLogs["required attribute vintage start is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		// }
		// if rowDataToInsert[i][7] == "" {
		// 	taskErrorLogs["required attribute vintage end is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		// }
		// if rowDataToInsert[i][8] == "" {
		// 	taskErrorLogs["required attribute number of issued credits is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		// }
		if rowDataToInsert[i][9] == "" {
			taskErrorLogs["required attribute crediting period start date is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][10] == "" {
			taskErrorLogs["required attribute crediting period end date is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][14] == "" {
			taskErrorLogs["required attribute current status of sale is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][15] == "" {
			taskErrorLogs["required attribute date of issuance to cc is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][16] == "" {
			taskErrorLogs["required attribute contract status is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][17] == "" {
			taskErrorLogs["required attribute project status is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][18] == "" {
			taskErrorLogs["required attribute host country approval is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][19] == "" {
			taskErrorLogs["required attribute approval status is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}
		if rowDataToInsert[i][20] == "" {
			taskErrorLogs["required attribute methodology is missing from csv"] = fmt.Sprintf("missing attribute in row %d", i+2)
		}

		if len(taskErrorLogs) != 0 {
			taskErrorLogsMap[fmt.Sprintf("missing required attributes from row %d", i+2)] = taskErrorLogs
			continue
		}

		upd := expression.
			Set(expression.Name("sellerName"), expression.Value(rowDataToInsert[i][0])).
			Set(expression.Name("projectId"), expression.Value(rowDataToInsert[i][1])).
			Set(expression.Name("projectName"), expression.Value(rowDataToInsert[i][2])).
			Set(expression.Name("plantName"), expression.Value(rowDataToInsert[i][3])).
			Set(expression.Name("projectType"), expression.Value(rowDataToInsert[i][4])).
			Set(expression.Name("registeredMw"), expression.Value(rowDataToInsert[i][5])).
			Set(expression.Name("vintageStart"), expression.Value(rowDataToInsert[i][6])).
			Set(expression.Name("vintageEnd"), expression.Value(rowDataToInsert[i][7])).
			Set(expression.Name("numberOfIssuedCredits"), expression.Value(rowDataToInsert[i][8])).
			Set(expression.Name("creditingPeriodStartDate"), expression.Value(rowDataToInsert[i][9])).
			Set(expression.Name("creditingPeriodEndDate"), expression.Value(rowDataToInsert[i][10])).
			Set(expression.Name("codOfTheProject"), expression.Value(rowDataToInsert[i][11])).
			Set(expression.Name("capacityFactor"), expression.Value(rowDataToInsert[i][12])).
			Set(expression.Name("averageAnnualVariationInCapacityFactor"), expression.Value(rowDataToInsert[i][13])).
			Set(expression.Name("currentStatusOfSale"), expression.Value(rowDataToInsert[i][14])).
			Set(expression.Name("dateOfIssuanceToCC"), expression.Value(rowDataToInsert[i][15])).
			Set(expression.Name("contractStatus"), expression.Value(rowDataToInsert[i][16])).
			Set(expression.Name("projectStatus"), expression.Value(rowDataToInsert[i][17])).
			Set(expression.Name("hostCountryApproval"), expression.Value(rowDataToInsert[i][18])).
			Set(expression.Name("approvalStatus"), expression.Value(rowDataToInsert[i][19])).
			Set(expression.Name("methodology"), expression.Value(rowDataToInsert[i][20]))

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
		}
	}

	return taskErrorLogsMap
}

func createTask(taskType string, createTime time.Time) (string, string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	// emptyErrorMap := make(map[string]string)
	taskPk := "task#" + uuid.New().String()
	taskSk := "taskType#" + taskType
	data := TaskData{
		Pk:           taskPk,
		Sk:           taskSk,
		CreateTime:   createTime.String(),
		TaskProgress: "In progress",
		UpdateTime:   "",
		// ErrorLogs:  emptyErrorMap,
	}
	av, _ := dynamodbattribute.MarshalMap(data)

	tableName := "RiskAssessmentTable"
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err := svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:", err)
	}
	return taskPk, taskSk
}

func updateTask(taskPk string, taskSk string, errMsgMap map[string]map[string]string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	taskProgress := ""
	if len(errMsgMap) == 0 {
		taskProgress = "Completed"
	} else {
		taskProgress = "Error detected"
	}

	// building expression for updating items
	upd := expression.
		Set(expression.Name("updateTime"), expression.Value(time.Now().String())).
		Set(expression.Name("errorLogs"), expression.Value(errMsgMap)).
		Set(expression.Name("taskProgress"), expression.Value(taskProgress))

	expr, errUpdateBuild := expression.NewBuilder().WithUpdate(upd).Build()

	if errUpdateBuild != nil {
		fmt.Println(errUpdateBuild)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("RiskAssessmentTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(taskPk),
			},
			"sk": {
				S: aws.String(taskSk),
			},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	db.UpdateItem(input)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, sqsRecord := range sqsEvent.Records {
		// decode sqs body to s3 sqsEvent
		s3Event := &events.S3Event{}
		err := json.Unmarshal([]byte(sqsRecord.Body), s3Event)
		if err != nil {
			return errors.Wrap(err, "Failed to decode sqs body to an S3 sqsEvent")
		}

		if len(s3Event.Records) == 0 {
			return errors.New("S3 Event Records is empty")
		}

		for _, s3Record := range s3Event.Records {
			taskType := ""
			if strings.Contains(s3Record.S3.Object.Key, "update") {
				taskType = "update"
			} else if strings.Contains(s3Record.S3.Object.Key, "create") {
				taskType = "create"
			}

			fileName := s3Record.S3.Object.Key
			bucketName := s3Record.S3.Bucket.Name

			csvData := getDataFromS3File(bucketName, fileName)
			headerDataToInsert, rowDataToInsert := extractData(csvData)
			if taskType == "create" {
				taskPk, taskSk := createTask(taskType, time.Now())
				errMsgMap := insertIntoDynamoDB(headerDataToInsert, rowDataToInsert, bucketName, fileName)
				updateTask(taskPk, taskSk, errMsgMap)
			} else if taskType == "update" {
				taskPk, taskSk := createTask(taskType, time.Now())
				errMsgMap := updateDynamoDB(headerDataToInsert, rowDataToInsert, bucketName, fileName)
				updateTask(taskPk, taskSk, errMsgMap)
			}

		}
	}

	return nil

}

func main() {
	lambda.Start(handler)
}
