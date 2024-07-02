package rules

import (
	"fmt"
	"irules/internal/data"
	"reflect"
	"strings"
	"testing"
)

func TestCheckNumericRange(t *testing.T) {

	var builderPtr *strings.Builder
	builderPtr = new(strings.Builder)

	type args struct {
		val        string
		rangeStr   string
		logBuilder *strings.Builder
	}
	tests := []struct {
		name string
		vals args
		want bool
	}{
		{name: "Int Positive 1", vals: args{val: "30", rangeStr: "30..60", logBuilder: builderPtr}, want: true},
		{name: "Int Positive 2", vals: args{val: "60", rangeStr: "30..60", logBuilder: builderPtr}, want: true},
		{name: "Int Negative 1", vals: args{val: "20", rangeStr: "30..60", logBuilder: builderPtr}, want: false},
		{name: "Int Negative 2", vals: args{val: "81", rangeStr: "30..60", logBuilder: builderPtr}, want: false},
		{name: "Float Positive 1", vals: args{val: "30.2", rangeStr: "30.2..60.8", logBuilder: builderPtr}, want: true},
		{name: "Float Positive 2", vals: args{val: "60.8", rangeStr: "30.2..60.8", logBuilder: builderPtr}, want: true},
		{name: "Float Negative 1", vals: args{val: "30.1", rangeStr: "30.2..60.8", logBuilder: builderPtr}, want: false},
		{name: "Float Negative 2", vals: args{val: "60.9", rangeStr: "30.2..60.8", logBuilder: builderPtr}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckNumericRange(tt.vals.val, tt.vals.rangeStr, tt.vals.logBuilder); got != tt.want {
				t.Errorf("CheckNumericRange() = %v, want %v", got, tt.want)
			}
		})
	}
	fmt.Println(builderPtr.String())
}

func TestEvaluateRuleRecordWithArray(t *testing.T) {
	var builderPtr *strings.Builder
	builderPtr = new(strings.Builder)
	/*
		outMap0 := make(map[string]interface{})
		outMap1 := make(map[string]interface{})


		inDataObjWithArray := CreateCustTreeWithArrays()

		inDataObjWithoutArray := map[string]interface{}{
			"cust": map[string]interface{}{
				"primaryaddress": map[string]interface{}{
					"street": "123 Elm St",
					"city":   "Somewhere",
					"state":  "TX",
					"zip":    "75000",
				},
				"age": 30,
				"examPassed": map[string]interface{}{
					"status": true,
					"year":   "2022",
					"score": 3.5
				},
			},
		}

		SetMapValueByJSONPath(outMap0, "cust.adrs.street1", "100 main st")
		fmt.Println("\n=================== Set an obj without array :Start ===============")
		PrintObjectRecursively(outMap0, 0)
		fmt.Println("\n--------------- End of Obj ---------------")
		SetMapValueByJSONPath(outMap1, "cust[0].adrs.street1", "200 main st")
		SetMapValueByJSONPath(outMap1, "cust[1].adrs.street1", "300 main st")
		fmt.Println("\n=================== Set an obj with array :Start ===============")
		PrintObjectRecursively(outMap1, 0)
		fmt.Println("\n--------------- End of Obj ---------------")

		obj0, _ := GetValueFromMap(inDataObjWithArray, "cust[]")
		fmt.Println("\n=================== Get an obj from array without index :Start ===============")
		PrintObjectRecursively(obj0, 0)
		fmt.Println("\n--------------- End of Obj ---------------")

		obj1, _ := GetValueFromMap(inDataObjWithArray, "cust[0]")
		fmt.Println("\n=================== Get an obj from array with index :Start ===============")
		PrintObjectRecursively(obj1, 0)
		fmt.Println("\n--------------- End of Obj ---------------")

		obj2, _ := GetValueFromMap(inDataObjWithoutArray, "cust")
		fmt.Println("\n=================== Get an obj without array :Start ===============")
		PrintObjectRecursively(obj2, 0)
		fmt.Println("\n--------------- End of Obj ---------------")
	*/
	inDataObjWithArray2 := CreateCustTreeWithArrays()
	outMap2 := make(map[string]interface{})
	type args struct {
		ruleString   []string
		headerString []string
		dirString    []string
		inData       map[string]interface{}
		logBuilder   *strings.Builder
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "1st Rec",
			args: args{
				ruleString:   []string{"TX", "18..99", ">currentyear-5", "true", "texas resident, adult, license not older than 5 years"},
				headerString: []string{"cust[*].primaryaddress.state", "cust[*].age", "cust[*].examPassed.year", "cust[*].examPassed.status", "cust[*].result", "cust[*].resultText"},
				dirString:    []string{"in", "in", "in", "in", "out", "out"},
				inData:       inDataObjWithArray2,
				logBuilder:   builderPtr,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rulesPassed := EvaluateRuleRecord(tt.args.ruleString, tt.args.headerString, tt.args.dirString, tt.args.inData, tt.args.logBuilder, outMap2)
			fmt.Println("\n=================== Get Back an OutMap :Start ===============")
			data.PrintObjectRecursively(outMap2, 0)
			fmt.Println("\n--------------- End of OutMap ---------------")
			fmt.Println("\n=================== Updated InData :Start ===============")
			data.PrintObjectRecursively(inDataObjWithArray2, 0)
			fmt.Println("\n--------------- End of Updated InData ---------------")
			if !reflect.DeepEqual(rulesPassed, true) {
				t.Errorf("EvaluateRuleRecord() got = %v, want %v", rulesPassed, true)
			}
		})
	}
	fmt.Println("\n=================== Start of Decision Log ===============")
	fmt.Println(builderPtr.String())
	fmt.Println("\n--------------- End of Decision Log -----------------")

}

func CreateCustTreeWithArrays() map[string]interface{} {
	inDataObjWithArray := map[string]interface{}{
		"cust": []interface{}{
			map[string]interface{}{
				"primaryaddress": map[string]interface{}{
					"street": "123 Elm St",
					"city":   "City1",
					"state":  "TX",
					"zip":    "75000",
				},
				"age": 30,
				"examPassed": map[string]interface{}{
					"status": true,
					"year":   "2022",
					"score":  3.5,
				},
			},
			map[string]interface{}{
				"primaryaddress": map[string]interface{}{
					"street": "456 Main St",
					"city":   "City2",
					"state":  "TX",
					"zip":    "75001",
				},
				"age": 50,
				"examPassed": map[string]interface{}{
					"status": true,
					"year":   "2023",
					"score":  3.5,
				},
			},
		},
	}
	return inDataObjWithArray
}

func TestEvaluateRuleRecordSimple(t *testing.T) {
	var builderPtr *strings.Builder
	builderPtr = new(strings.Builder)

	outMap := make(map[string]interface{})

	inDataObj := map[string]interface{}{
		"cust": map[string]interface{}{
			"primaryaddress": map[string]interface{}{
				"street": "123 Elm St",
				"city":   "Somewhere",
				"state":  "TX",
				"zip":    "75000",
			},
			"age": 30,
			"examPassed": map[string]interface{}{
				"status": true,
				"year":   "2022",
				"score":  3.5,
			},
		},
	}

	type args struct {
		ruleString   []string
		headerString []string
		dirString    []string
		inData       map[string]interface{}
		logBuilder   *strings.Builder
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "1st Rec",
			args: args{
				ruleString:   []string{"TX", "18..99", ">currentyear-5", "true", "texas resident, adult, license not older than 5 years"},
				headerString: []string{"cust.primaryaddress.state", "cust.age", "cust.examPassed.year", "cust.examPassed.score", "cust.result", "cust.resultText"},
				dirString:    []string{"in", "in", "in", "in", "out", "out"},
				inData:       inDataObj,
				logBuilder:   builderPtr,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rulesPassed := EvaluateRuleRecord(tt.args.ruleString, tt.args.headerString, tt.args.dirString, tt.args.inData, tt.args.logBuilder, outMap)
			fmt.Println("\n=================== Get Back an OutMap :Start ===============")
			data.PrintObjectRecursively(outMap, 0)
			fmt.Println("\n--------------- End of OutMap ---------------")
			fmt.Println("\n=================== Updated InData :Start ===============")
			data.PrintObjectRecursively(inDataObj, 0)
			fmt.Println("\n--------------- End of Updated InData ---------------")
			if !reflect.DeepEqual(rulesPassed, true) {
				t.Errorf("EvaluateRuleRecord() got = %v, want %v", rulesPassed, true)
			}
		})
	}
	fmt.Println("\n=================== Start of Decision Log ===============")
	fmt.Println(builderPtr.String())
	fmt.Println("\n--------------- End of Decision Log -----------------")

}
