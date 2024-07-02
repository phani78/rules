package rules

import (
	"fmt"
	"irules/internal/data"
	"strings"
)

func EvalLogicalExpr(ruleExpr string, header string, dataPipeline map[string]interface{}, logBuilder *strings.Builder) bool {
	var result bool

	switch {
	case strings.Contains(ruleExpr, data.TOKEN_AND):
		prefix, suffix, err := data.Split(ruleExpr, data.TOKEN_AND)
		if err == nil {
			prefixresult := EvalLogicalExpr(prefix, header, dataPipeline, logBuilder)
			suffixresult := EvalLogicalExpr(suffix, header, dataPipeline, logBuilder)
			return prefixresult && suffixresult
		} else {
			return false
		}
	case strings.Contains(ruleExpr, data.TOKEN_OR):
		prefix, suffix, err := data.Split(ruleExpr, data.TOKEN_OR)
		if err == nil {
			prefixresult := EvalLogicalExpr(prefix, header, dataPipeline, logBuilder)
			suffixresult := EvalLogicalExpr(suffix, header, dataPipeline, logBuilder)
			return prefixresult || suffixresult
		} else {
			return false
		}
	default:
		operator := data.ContainsAny(ruleExpr, data.TOKEN_LOGICAL_OPS)
		if operator != "" {
			suffix, err := data.GetStringAfterSubstring(ruleExpr, operator)
			if err == nil {
				suffixResult := data.EvalNumericExpr(suffix, dataPipeline)
				logicalOperResult := data.ComparePipelineValue(header, dataPipeline, suffixResult, operator)
				//logBuilder.WriteString(fmt.Sprintf("Result of %s is determined as : %v\n", ruleExpr, logicalOperResult))
				return logicalOperResult
			} else {
				logBuilder.WriteString(fmt.Sprintf("An error occured while evaluating : %s(Error : %v)\n", ruleExpr, err))
				return false
			}
		} else {
			logBuilder.WriteString(fmt.Sprintf("Error : could not evaluate expr(unknown operator) : %s\n", ruleExpr))
		}
	}
	return result
}
