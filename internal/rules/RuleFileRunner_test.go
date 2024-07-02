package rules

import (
	"fmt"
	"irules/internal/data"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRule(t *testing.T) {
	fmt.Println("Started Exec")
	var builderPtr *strings.Builder
	builderPtr = new(strings.Builder)
	inDataObjWithArray2 := CreateCustTreeWithoutArrays()

	dir, err := filepath.Abs("..\\..\\..\\R.csv")
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
			name: "2nd Grid",
			args: args{
				ruleFile:   "C:\\Users\\Srav\\Documents\\code\\irules\\R.csv",
				ruleName:   "LicenseEligibility",
				version:    "1.5",
				inData:     inDataObjWithArray2,
				logBuilder: builderPtr,
			},
			want: true,
		},
		{
			name: "1st Grid",
			args: args{
				ruleFile:   "C:\\Users\\Srav\\Documents\\code\\irules\\R.csv",
				ruleName:   "LicenseEligibility",
				version:    "1.3",
				inData:     inDataObjWithArray2,
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

func CreateCustTreeWithoutArrays() map[string]interface{} {
	inDataObjWithArray := map[string]interface{}{
		"currentyear": 2024,
		"splState":    "TX",
		"cust": map[string]interface{}{
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
	}
	return inDataObjWithArray
}
