package data

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var TOKEN_OR = " or "
var TOKEN_AND = " and "
var TOKEN_LTE = "<="
var TOKEN_LT = "<"
var TOKEN_GTE = ">="
var TOKEN_GT = ">"
var TOKEN_LOGICAL_OPS = []string{TOKEN_LTE, TOKEN_LT, TOKEN_GTE, TOKEN_GT}

func GetFloat(value interface{}) float64 {
	switch value.(type) {
	case int:
		return float64(value.(int))
	case float64:
		return value.(float64)
	case string:
		floatValue, _ := strconv.ParseFloat(strings.TrimSpace(value.(string)), 64)
		return floatValue
	default:
		return -999999.99
	}
}

func PrintObjectRecursively(obj interface{}, indent int) {
	spacing := ""
	for i := 0; i < indent; i++ {
		spacing += " "
	}

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	//fmt.Printf(" %sEntered PrintObjectRecursively() %v%s\n", spacing, val, typ)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
		typ = typ.Elem()
		PrintObjectRecursively(val.Interface(), indent)

	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			//fmt.Printf("%sStruct Field  Type=%s Kind=%s Value=%v \n", spacing, field.Type, field.Kind(), typ.Field(i))

			if field.Kind() == reflect.Struct || field.Kind() == reflect.Map || field.Kind() == reflect.Slice || field.Kind() == reflect.Array || field.Kind() == reflect.Ptr {
				fieldType := typ.Field(i)
				fmt.Printf("%sStruct Complex Type %s Value %s \n", spacing, fieldType.Name, field.Interface())
				//fmt.Printf("%s%s: %v\n", spacing, fieldType.Name, field.Interface())
				PrintObjectRecursively(field.Interface(), indent+4)
			} else {
				fmt.Printf("%sStruct Simple Field : %s=%v(%s) \n", spacing, typ.Field(i).Name, field.Interface(), field.Kind())

				//if field.Kind() == reflect.Bool {
				//	fmt.Printf("%sStruct Simple Field : %s=%v(%s) \n", spacing+"    ", typ.Field(i).Name, GetUnexportedField(field), field.Kind())
				//} else {
				//	fmt.Printf("%sStruct Simple Field : %s=%s(%s) \n", spacing+"    ", typ.Field(i).Name, field.String(), field.Kind())
				//}

			}

		}

	case reflect.Map:
		for _, key := range val.MapKeys() {
			fmt.Printf("%s%v:\n", spacing, key)
			PrintObjectRecursively(val.MapIndex(key).Interface(), indent+4)
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			fmt.Printf("%s[%d]\n", spacing, i)
			PrintObjectRecursively(val.Index(i).Interface(), indent+4)
		}

	default:
		fmt.Printf("%s%v\n", spacing, val)
	}
}

/*
func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}
*/

func HasSubstringInArray(arr []string, substr string) (bool, string, string) {
	for _, str := range arr {
		index := strings.Index(str, substr)
		if index > 0 {
			prefix := str[:index]
			suffix := str[index+len(substr):]
			return true, prefix, suffix
		}
	}
	return false, "", ""
}

func ReplaceInStringArray(arr []string, oldSubstr, newSubstr string) []string {
	// Create a new slice to store the modified strings
	modifiedArr := make([]string, len(arr))

	// Iterate over each string in the array
	for i, s := range arr {
		// Replace oldSubstr with newSubstr in each string and store in modifiedArr
		modifiedArr[i] = strings.ReplaceAll(s, oldSubstr, newSubstr)
	}
	return modifiedArr
}

func ExtractArrayParts(str string) (bool, string, string) {
	index1 := strings.Index(str, "[")
	if index1 > 0 {
		arrayName := str[:index1]
		index2 := strings.Index(str, "]")
		if index2 > index1 {
			//suffix := str[index2+1:]
			indexStr := str[index1+1 : index2]
			return true, arrayName, indexStr
		}
	}
	return false, "", ""
}

func DetectAndConvert(value interface{}) interface{} {
	switch value.(type) {
	case string:
		strVal := value.(string)

		// Try to convert to int
		if intValue, err := strconv.Atoi(strVal); err == nil {
			return intValue
		}

		// Try to convert to float64
		if floatValue, err := strconv.ParseFloat(strVal, 64); err == nil {
			return floatValue
		}

		// Try to convert to bool
		if boolValue, err := strconv.ParseBool(strVal); err == nil {
			return boolValue
		}

		// Default to string
		return strVal
	}
	return value
}

var floatPattern = `^-?\d+([.])+\d+?$`
var integerPattern = `^-?\d+$`

func EvalNumericExpr(expr string, dataPipeline map[string]interface{}) interface{} {

	hasOperator, operator, prefix, suffix := SplitIfContainsAny(expr, "+-*/")

	if hasOperator {
		prefixValue := GetPipelineOrStaticValue(dataPipeline, prefix)
		suffixValue := GetPipelineOrStaticValue(dataPipeline, suffix)

		isFloat := IsFloat(prefixValue) || IsFloat(suffixValue)
		if isFloat {
			floatResult := DoFloatMath(operator, prefixValue, suffixValue)
			return floatResult
		} else {
			areInts := IsInt(prefixValue) && IsInt(suffixValue)
			if areInts {
				intResult := DoIntMath(operator, prefixValue, suffixValue)
				return intResult
			}
		}
	} else {
		return DetectAndConvert(expr)
	}
	return nil
}

func DoFloatMath(operator rune, prefixValue interface{}, suffixValue interface{}) float64 {
	switch operator {
	case '+':
		return prefixValue.(float64) + suffixValue.(float64)
	case '-':
		return prefixValue.(float64) - suffixValue.(float64)
	case '*':
		return prefixValue.(float64) * suffixValue.(float64)
	case '/':
		return prefixValue.(float64) / suffixValue.(float64)
	}
	return 0
}

func DoIntMath(operator rune, prefixValue interface{}, suffixValue interface{}) int {
	switch operator {
	case '+':
		return prefixValue.(int) + suffixValue.(int)
	case '-':
		return prefixValue.(int) - suffixValue.(int)
	case '*':
		return prefixValue.(int) * suffixValue.(int)
	case '/':
		return prefixValue.(int) / suffixValue.(int)
	}
	return 0
}

func IsFloat(value interface{}) bool {
	switch value.(type) {
	case float64:
		return true
	case string:
		isFloat, _ := regexp.MatchString(floatPattern, value.(string))
		return isFloat
	}

	return false
}

func IsInt(value interface{}) bool {
	switch value.(type) {
	case int:
		return true
	case string:
		isInt, _ := regexp.MatchString(integerPattern, value.(string))
		return isInt
	}

	return false
}

func AppendString(arr []string, str string) ([]string, error) {
	if arr == nil {
		arr = make([]string, 0)
	}

	// Append the string to the array
	arr = append(arr, str)
	return arr, nil
}
