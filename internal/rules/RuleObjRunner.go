package rules

import (
	"fmt"
	"irules/internal/data"
	"strconv"
	"strings"
)

type RuleData struct {
	RuleName        string
	RuleVersion     string
	PreRules        string
	PostRules       string
	IsRuleTable     bool
	RuleDirStr      []string
	RuleHeaderNames []string
	RuleExprRecs    [][]string
}

//func (r *RuleData) IsRuleTable() bool {
//	return r.isRuleTable
//}
//
//func (r *RuleData) SetIsRuleTable(isRuleTable bool) {
//	r.isRuleTable = isRuleTable
//}

func AppendRuleString(ruleString []string, ruleDataPtr *RuleData) {
	if !data.IsEmptyStringArray(ruleString) {
		ruleDataPtr.RuleExprRecs = append(ruleDataPtr.RuleExprRecs, ruleString)
	}
}

func IsStartOfRule(ruleString []string) bool {

	return (strings.Compare("Rule", ruleString[0]) == 0)
}

func EvaluateRuleTable(returnOnFirstPass bool, ruleGrid [][]string, headerString []string, dirString []string, inData map[string]interface{}, logBuilder *strings.Builder) (bool, [][]string, []string) {
	RanSuccesfully := true
	numRules := len(ruleGrid)
	numOutObjs := numRules

	numOutCols, outColNames := ExtractOutColData(dirString, headerString)

	if returnOnFirstPass {
		numOutObjs = 1
	}
	outDataArr := make([][]string, numOutObjs)
	for l := range outDataArr {
		outDataArr[l] = make([]string, numOutCols)
	}

	//lenJsonPath = prefix+suffix
	currObjIndex := 0
	for i := 0; i < numRules; i++ {
		//MyTODO : call this func recursively , pass numObjs to SetMapValueByJSONPath to control the num array elements

		if !returnOnFirstPass {
			currObjIndex = i
		}

		currentIterSuccess := EvaluateSingleRuleRecord(ruleGrid[i], headerString, dirString, inData, logBuilder, nil, outDataArr, currObjIndex)
		RanSuccesfully = RanSuccesfully || currentIterSuccess
		if returnOnFirstPass && currentIterSuccess {
			break
		}
		//RanSuccesfully = RanSuccesfully || currentIterSuccess
	}

	return RanSuccesfully, outDataArr, outColNames

	//var outStrings []string = []string{}
	//var logBuilder strings.Builder

	//return EvaluateSingleRuleRecord(ruleString, headerString, dirString, inData, logBuilder, outMap)
}

func ExtractOutColData(dirString []string, headerString []string) (int, []string) {
	numOutCols := -1
	minIndex := -1
	maxIndex := -1
	for pos, item := range dirString {
		if item == "out" {
			if minIndex == -1 {
				minIndex = pos
			}
			maxIndex = pos
		}
	}
	numOutCols = maxIndex - minIndex + 1
	outColNames := make([]string, numOutCols)
	outColIndex := 0
	for pos, item := range dirString {
		if item == "out" {
			outColNames[outColIndex] = headerString[pos]
			outColIndex++
		}
	}
	return numOutCols, outColNames
}

func EvaluateRuleRecord(ruleString []string, headerString []string, dirString []string, inData map[string]interface{}, logBuilder *strings.Builder, outMap map[string]interface{}) bool {
	//outMap := make(map[string]interface{})
	hasIndexString, prefix, _ := data.HasSubstringInArray(headerString, "[*]")
	if hasIndexString {
		RanSuccesfully := true

		//lenJsonPath = prefix+suffix
		objs, _ := data.GetMapValueByJSONPath(inData, prefix+"[]")
		numObjs := len(objs.([]interface{}))
		for i := 0; i < numObjs; i++ {
			headerStringWithCurrIndex := data.ReplaceInStringArray(headerString, "*", strconv.Itoa(i))
			fmt.Printf("\nWhen evaluating array # %d, headerStringWithCurrIndex = %v", i, headerStringWithCurrIndex)
			//MyTODO Immediate : call this func recursively , pass numObjs to SetMapValueByJSONPath to control the num array elements
			currentIterSuccess := EvaluateSingleRuleRecord(ruleString, headerStringWithCurrIndex, dirString, inData, logBuilder, outMap, nil, -1)
			RanSuccesfully = RanSuccesfully && currentIterSuccess
		}
		return RanSuccesfully

	} else {
		return EvaluateSingleRuleRecord(ruleString, headerString, dirString, inData, logBuilder, outMap, nil, -1)
	}

	//var outStrings []string = []string{}
	//var logBuilder strings.Builder

	//return EvaluateSingleRuleRecord(ruleString, headerString, dirString, inData, logBuilder, outMap)
}

func EvaluateSingleRuleRecord(ruleString []string, headerString []string, dirString []string, inData map[string]interface{}, logBuilder *strings.Builder, outMap map[string]interface{}, outDataArr [][]string, recIndex int) bool {
	logBuilder.WriteString(fmt.Sprintf("Starting the evaluation of Rule Record %v \n", ruleString))
	allInCondPassed := true
	outColNum := 0
	for i, ruleExpr := range ruleString {
		if len(strings.TrimSpace(dirString[i])) > 0 {
			if dirString[i] == "in" {
				evalResult := true
				if ruleExpr != "" {
					evalResult, _ = evalRuleExpr(ruleExpr, headerString[i], inData, logBuilder, true)
				}
				fmt.Printf("\nevalResult : %v", evalResult)
				allInCondPassed = allInCondPassed && evalResult
				if !evalResult {
					fmt.Printf("\nCond #%d failed with false, ruleExpr=%s", i, ruleExpr)
					return false
				}

			} else if dirString[i] == "out" {
				outputValue, error := data.GetMapValueByJSONPath(inData, ruleExpr)
				var outString string
				if error == nil {
					outString = fmt.Sprintf("%v", outputValue)
					logBuilder.WriteString(fmt.Sprintf("	Taking the calculated value for expression %s using the expression for header : %s\n", ruleExpr, headerString[i]))
				} else {
					if ruleExpr == "true" || ruleExpr == "false" {
						if !allInCondPassed {
							outString = "false"
						} else {
							outString = ruleExpr
						}
					} else {
						outString = ruleExpr
						logBuilder.WriteString(fmt.Sprintf("	Taking the constant value %s for header : %s\n", outString, headerString[i]))
					}
				}

				data.SetMapValueByJSONPath(inData, headerString[i], outString)
				if outMap != nil {
					data.SetMapValueByJSONPath(outMap, headerString[i], outString)
				} else {
					outDataArr[recIndex][outColNum] = outString
					outColNum++
				}

			} else {
				logBuilder.WriteString(fmt.Sprintf("	Ignoring the column with header : %s\n", headerString[i]))
			}
		}

		//if len(str) <= 8 {
		//	result = append(result, str)
		//	logBuilder.WriteString(fmt.Sprintf("'%s' is selected because its length is %d.\n", str, len(str)))
		//} else {
		//	logBuilder.WriteString(fmt.Sprintf("'%s' is not selected because its length is %d.\n", str, len(str)))
		//}
	}
	logBuilder.WriteString(fmt.Sprintf("completed the evaluation of Rule Record %v with a result : %v \n", ruleString, allInCondPassed))
	return allInCondPassed
}

func evalFunCall(ruleExpr string, header string, inData map[string]interface{}, logBuilder *strings.Builder) (bool, interface{}) {
	var result bool
	var value interface{}

	withinParanthesis, funcName := data.ExtractTextInParanthesis(ruleExpr)

	if data.StartsWithFunctionCall(withinParanthesis) {
		result, value = evalFunCall(withinParanthesis, header, inData, logBuilder)
	} else {
		result, value = evalRuleExpr(withinParanthesis, header, inData, logBuilder, true)
	}

	//result, value = evalRuleExpr(withinParanthesis, header, inData, logBuilder)
	switch funcName {
	case "Not":
		result, value = ProcessNotEquals(result)
		fmt.Printf("Func Not() with param %s has returned result=%v and return value=%v", withinParanthesis, result, value)
		return result, nil
	}

	return false, nil

}

func evalRuleExpr(ruleExpr string, header string, inData map[string]interface{}, logBuilder *strings.Builder, topLevelCall bool) (bool, interface{}) {
	var result bool

	switch {
	case data.StartsWithFunctionCall(ruleExpr):
		result, _ := evalFunCall(ruleExpr, header, inData, logBuilder)
		return result, nil

	case data.EnclosedWithParanthesis(ruleExpr):
		withinParanthesis, _ := data.ExtractTextInParanthesis(ruleExpr)
		result, _ = evalRuleExpr(withinParanthesis, header, inData, logBuilder, false)
		return result, nil
	case strings.Contains(ruleExpr, ".."):
		inputValue, _ := data.GetMapValueByJSONPath(inData, header)
		inputStr := fmt.Sprintf("%v", inputValue)
		result = CheckNumericRange(inputStr, ruleExpr, logBuilder)
		return result, nil
	case strings.Contains(ruleExpr, ","):
		tokens := strings.Split(ruleExpr, ",")
		for _, token := range tokens {
			succ, _ := evalRuleExpr(token, header, inData, logBuilder, false)
			if succ {
				return true, nil
			}
		}
		return false, nil
	case strings.ContainsAny(ruleExpr, "<>="):
		result = EvalLogicalExpr(ruleExpr, header, inData, logBuilder)
		//logBuilder.WriteString(fmt.Sprintf("Not Yet Implemented : %s\n", ruleExpr))
		logBuilder.WriteString(fmt.Sprintf("	Result of %s is determined as : %v\n", ruleExpr, result))
		return result, nil
	default:
		inputTypedValue := data.GetPipelineOrStaticValue(inData, ruleExpr)
		if topLevelCall {
			headerValue, _ := data.GetMapValueByJSONPath(inData, header)
			headerValue = data.DetectAndConvert(headerValue)

			result = (headerValue == inputTypedValue)
			logBuilder.WriteString(fmt.Sprintf("	evaluated result for string equality for header %s and rule value ='%s' : %v \n", header, ruleExpr, result))
			return result, nil
		} else {
			return true, inputTypedValue
		}
	}

}

func CheckNumericRange(val string, rangeStr string, logBuilder *strings.Builder) bool {

	minVal, _ := data.GetStringBeforeSubstring(rangeStr, "..")
	maxVal, _ := data.GetStringAfterSubstring(rangeStr, "..")

	// Try converting to integer
	if valInt, err := strconv.Atoi(val); err == nil {
		if minValInt, err := strconv.Atoi(minVal); err == nil {
			if maxValInt, err := strconv.Atoi(maxVal); err == nil {
				result := valInt >= minValInt && valInt <= maxValInt
				logBuilder.WriteString(fmt.Sprintf("	evaluated result for int range check for rangeStr='%s' and rule value='%s' : %v \n", rangeStr, val, result))
				return (result)
			}
		}
	}

	// Try converting to float
	if valFloat, err := strconv.ParseFloat(val, 64); err == nil {
		if minValFloat, err := strconv.ParseFloat(minVal, 64); err == nil {
			if maxValFloat, err := strconv.ParseFloat(maxVal, 64); err == nil {
				result := valFloat >= minValFloat && valFloat <= maxValFloat
				logBuilder.WriteString(fmt.Sprintf("	evaluated float range check to %v for rangeStr='%s' and value='%s'\n", result, rangeStr, val))
				return (result)
			}
		}
	}

	logBuilder.WriteString(fmt.Sprintf("	evaluated range result for rangeStr='%s' and value='%s' : false.\n", rangeStr, val))
	return false
}
