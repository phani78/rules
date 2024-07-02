package rules

import (
	"fmt"
	"strings"
)

func RunRuleFromFile(returnOnFirstPass bool, ruleFile string, ruleName string, version string, inData map[string]interface{}, logBuilder *strings.Builder) (bool, [][]string, []string, error) {
	logBuilder.WriteString(fmt.Sprintf("Starting the Rules in the file %s with rule name %s and version # %s \n", ruleFile, ruleName, version))

	RuleDataArray, err := ReadCSVData(ruleFile) //".\\R.csv"
	numRules := len(RuleDataArray)
	if err == nil {
		for i := 0; i < numRules; i++ {
			if RuleDataArray[i].RuleName == ruleName && RuleDataArray[i].RuleVersion == version {

				ruleSuccess, outGrid, outColNames := RunRuleTable(returnOnFirstPass, RuleDataArray, i, inData, logBuilder)

				logBuilder.WriteString(fmt.Sprintf("Completed the Rules in the file %s with rule name %s and version # %s with a result : %v \n\n\n", ruleFile, ruleName, version, ruleSuccess))

				return ruleSuccess, outGrid, outColNames, nil
			}
		}
	}

	logBuilder.WriteString(fmt.Sprintf("Completed the Rules in the file %s with rule name %s and version # %s with an error : %v \n\n\n", ruleFile, ruleName, version, err))

	return false, nil, nil, err
}

func RunRuleTable(returnOnFirstPass bool, RuleDataArray []*RuleData, index int, inData map[string]interface{}, logBuilder *strings.Builder) (bool, [][]string, []string) {

	PreRuleDataArray := getRulesByNames(RuleDataArray[index].PreRules, RuleDataArray)
	for _, currPreRule := range PreRuleDataArray {
		pre_ruleSuccess, pre_outGrid, _ := EvaluateRuleTable(returnOnFirstPass, currPreRule.RuleExprRecs, currPreRule.RuleHeaderNames, currPreRule.RuleDirStr, inData, logBuilder)
		fmt.Printf("ran the prerule %s-%s with : result = %v, outgrid=%v \n", currPreRule.RuleName, currPreRule.RuleVersion, pre_ruleSuccess, pre_outGrid)
	}

	ruleSuccess, outGrid, outColNames := EvaluateRuleTable(returnOnFirstPass, RuleDataArray[index].RuleExprRecs, RuleDataArray[index].RuleHeaderNames, RuleDataArray[index].RuleDirStr, inData, logBuilder)

	PostRuleDataArray := getRulesByNames(RuleDataArray[index].PostRules, RuleDataArray)
	for _, currPostRule := range PostRuleDataArray {
		post_ruleSuccess, post_outGrid, _ := EvaluateRuleTable(returnOnFirstPass, currPostRule.RuleExprRecs, currPostRule.RuleHeaderNames, currPostRule.RuleDirStr, inData, logBuilder)
		fmt.Printf("ran the postrule %s-%s with : result = %v, outgrid=%v \n", currPostRule.RuleName, currPostRule.RuleVersion, post_ruleSuccess, post_outGrid)
	}

	return ruleSuccess, outGrid, outColNames
}

func getRulesByNames(ruleListStr string, rulesArray []*RuleData) []*RuleData {
	RuleDataArray := []*RuleData{}

	if ruleListStr == "" {
		return RuleDataArray
	}

	fullRuleNames := strings.Split(ruleListStr, ",")

	for _, fullRuleName := range fullRuleNames {
		ruleParts := strings.Split(fullRuleName, "-")
		for _, currRuleData := range rulesArray {
			if currRuleData.RuleName == ruleParts[0] && currRuleData.RuleVersion == ruleParts[1] {
				RuleDataArray = append(RuleDataArray, currRuleData)
			}
		}

	}
	return RuleDataArray
}
