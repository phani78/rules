package rules

import (
	"encoding/csv"
	"fmt"
	"strings"

	"irules/internal/data"
	"os"
)

func ReadCSVData(filename string) ([]*RuleData, error) {
	// if len(os.Args) < 2 {
	// 	fmt.Println("Usage: go run main.go <csv file>")
	// 	return
	// }

	// filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return nil, err
	}

	dataArray := make([][]string, len(records))

	for i, record := range records {
		dataArray[i] = make([]string, len(record))
		for j, value := range record {
			value = strings.ReplaceAll(value, "\"", "")
			dataArray[i][j] = value
		}
	}

	var RuleDataArrayIndex = -1
	RuleDataArray := []*RuleData{}
	fmt.Println("Data Array:")
	var ruleData *RuleData
	lineNum := 0
	for i, row := range dataArray {
		fmt.Printf("Row %d: %v\n", i, row)
		lineNum += 1
		if IsStartOfRule(row) {
			ruleData = &RuleData{RuleName: row[1], RuleVersion: row[3], IsRuleTable: false, PreRules: row[5], PostRules: row[7]}
			//ruleData.RuleName = row[1]
			//ruleData.RuleVersion = row[3]
			lineNum = 0
			fmt.Println("ruleData.RuleName :", ruleData.RuleName)
			RuleDataArray = append(RuleDataArray, ruleData)
			RuleDataArrayIndex += 1
		} else if lineNum == 1 {
			fmt.Println("ruleData=", ruleData, ", row[0] == RuleTable", row[0] == "RuleTable")
			ruleData.IsRuleTable = (row[0] == "RuleTable")
			//ruleData.isRuleTable = (row[0] == "RuleTable")
		} else if lineNum == 2 {
			ruleData.RuleDirStr = row
		} else if lineNum == 3 {
			ruleData.RuleHeaderNames = row
		} else {
			AppendRuleString(row, RuleDataArray[RuleDataArrayIndex])
		}

	}
	fmt.Printf("1===================After CSV Read here is RuleDataArray : ==============")
	data.PrintObjectRecursively(RuleDataArray, 0)
	fmt.Printf("2===================After printing RuleDataArray : ==============")
	return RuleDataArray, nil

	//cust.primaryaddress.state,cust.age,examPassedDate.status,examPassedDate.year,result.isEligible,result.whyEligible

	//var builderPtr *strings.Builder
	//builderPtr = new(strings.Builder)

	/*
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error reading CSV file:", err)
			return
		}

		dataMap := make(map[int]map[string]string)

		if len(records) > 0 {
			headers := records[0]

			for i, row := range records[0:] {
				fmt.Printf("RAW Row %d: %s\n", i, row)
				rowMap := make(map[string]string)
				for j, value := range row {
					rowMap[headers[j]] = value
				}
				dataMap[i] = rowMap
			}
		}

		fmt.Println("Data Map:")
		for rowNum, rowMap := range dataMap {
			fmt.Printf("Row %d:\n", rowNum)
			for colName, value := range rowMap {
				fmt.Printf("  %s: %s\n", colName, value)
			}
		}

	*/
}
