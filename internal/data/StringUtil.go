package data

import (
	"fmt"
	"regexp"
	"strings"
)

func ContainsAny(str string, anyStrArray []string) string {
	for _, currStr := range anyStrArray {
		if strings.Index(str, currStr) >= 0 {
			return currStr
		}
	}

	return ""
}

func Split(input, substring string) (string, string, error) {
	// Find the index of the substring
	index := strings.Index(input, substring)
	if index == -1 {
		return "", "", fmt.Errorf("substring not found")
	}

	// Split the string into two parts: before and after the substring
	before := input[:index]
	after := input[index+len(substring):]

	// Trim the whitespace from both parts
	before = strings.TrimSpace(before)
	after = strings.TrimSpace(after)

	return before, after, nil
}

func GetStringBeforeSubstring(s, substr string) (string, error) {
	index := strings.Index(s, substr)
	if index == -1 {
		return "", fmt.Errorf("substring '%s' not found in string '%s'", substr, s)
	}
	return strings.TrimSpace(s[:index]), nil
}

func GetStringAfterSubstring(s, substr string) (string, error) {
	index := strings.Index(s, substr)
	if index == -1 {
		return "", fmt.Errorf("substring '%s' not found in string '%s'", substr, s)
	}
	return strings.TrimSpace(s[index+len(substr):]), nil
}

func SplitIfContainsAny(input string, chars string) (bool, rune, string, string) {
	for _, char := range chars {
		index := strings.IndexRune(input, char)
		if index != -1 {
			prefix := strings.TrimSpace(input[:index])
			suffix := strings.TrimSpace(input[index+1:])
			return true, char, prefix, suffix
		}

	}
	return false, 0, "", ""
}

var funcPattern = "^[a-zA-Z]+\\(.*\\)$"

func StartsWithFunctionCall(input string) bool {
	hasFunCall, _ := regexp.MatchString(funcPattern, strings.TrimSpace(input))
	return hasFunCall
}

var paranthesisPattern = "^\\(.*\\)$"

func EnclosedWithParanthesis(input string) bool {
	hasFunCall, _ := regexp.MatchString(paranthesisPattern, strings.TrimSpace(input))
	return hasFunCall
}

func ExtractTextInParanthesis(input string) (string, string) {
	input = strings.TrimSpace(input)
	index1 := strings.Index(input, "(")
	index2 := strings.LastIndex(input, ")")
	paranthesisText := input[index1+1 : index2]
	funcName := ""
	if index1 > 0 {
		funcName = input[:index1]
	}

	return paranthesisText, funcName
}
