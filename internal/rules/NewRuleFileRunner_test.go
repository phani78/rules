package rules

import (
	"fmt"
	"irules/internal/data"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRuleNew(t *testing.T) {
	fmt.Println("Started Exec")
	var builderPtr *strings.Builder
	builderPtr = new(strings.Builder)
	inDataObj1 := CreateTestObj1()
	inDataObj2 := CreateTestObj2()

	dir, err := filepath.Abs("..\\..\\")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Dir : ", dir)

	type args struct {
		ruleFile   string
		ruleName   string
		version    string
		inData     map[string]interface{}
		logBuilder *strings.Builder
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   map[string]interface{}
		wantErr bool
	}{
		{
			name: "New Test",
			args: args{
				ruleFile:   "..\\..\\NewTest.csv",
				ruleName:   "New Rule",
				version:    "0.99",
				inData:     inDataObj1,
				logBuilder: builderPtr,
			},
			want: true,
		},
		{
			name: "1st Grid",
			args: args{
				ruleFile:   "..\\..\\NewTest.csv",
				ruleName:   "New Rule",
				version:    "0.99",
				inData:     inDataObj2,
				logBuilder: builderPtr,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			successCode, dataGrid, outColNames, err := RunRuleFromFile(true, tt.args.ruleFile, tt.args.ruleName, tt.args.version, tt.args.inData, tt.args.logBuilder)
			fmt.Println("\n=================== Get Back an OutArray :Start ===============")
			data.PrintObjectRecursively(outColNames, 0)
			data.PrintObjectRecursively(dataGrid, 0)
			fmt.Println("\n--------------- End of OutArray ---------------")
			fmt.Println("\n=================== Updated InData :Start ===============")
			data.PrintObjectRecursively(tt.args.inData, 0)
			fmt.Println("\n--------------- End of Updated InData ---------------")

			if err != nil {
				t.Errorf("RunRule() error = %v", err)
				return
			}
			if successCode != tt.want {
				t.Errorf("RunRule() got = %v, want %v", successCode, tt.want)
			} else {
				fmt.Printf("\nPassed the test with name : %s args : %v", tt.name, tt.args)
			}
			if dataGrid == nil {
				t.Errorf("RunRule() got OurData as null")
			}

			//if !reflect.DeepEqual(got1, tt.want1) {
			//	t.Errorf("RunRule() got1 = %v, want %v", got1, tt.want1)
			//}
		})
	}

	fmt.Println("\n=================== Start of Decision Log ===============")
	fmt.Println(builderPtr.String())
	fmt.Println("\n--------------- End of Decision Log -----------------")

}

func CreateTestObj1() map[string]interface{} {
	inDataObj := map[string]interface{}{
		"Fld2":  "040",
		"Fld4":  300,
		"Fld6":  "IN2",
		"Fld8":  "In3",
		"Fld10": "09",

		"Fld11": "Y",
		"Fld12": "Y",
		"Fld13": "xyz1",
	}
	return inDataObj
}

func CreateTestObj2() map[string]interface{} {
	inDataObj2 := map[string]interface{}{
		"Fld2":  "050",
		"Fld4":  300,
		"Fld6":  "IN2",
		"Fld8":  "In3",
		"Fld10": "09",

		"Fld11": "Y",
		"Fld12": "N",
		"Fld13": "xyz2",
	}
	return inDataObj2
}
