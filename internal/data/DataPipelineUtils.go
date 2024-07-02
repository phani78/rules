package data

import (
	"fmt"
	"strconv"
	"strings"
)

func ComparePipelineValue(jsonPath string, dataPipeline map[string]interface{}, value interface{}, operator string) bool {
	pipelineValue, _ := GetMapValueByJSONPath(dataPipeline, jsonPath)
	pipelineFloat := GetFloat(pipelineValue)
	valueFloat := GetFloat(value)
	switch operator {
	case TOKEN_LTE:
		return pipelineFloat <= valueFloat
	case TOKEN_LT:
		return pipelineFloat < valueFloat
	case TOKEN_GTE:
		return pipelineFloat >= valueFloat
	case TOKEN_GT:
		return pipelineFloat > valueFloat
	}
	return false
}

func GetPipelineOrStaticValue(dataPipeline map[string]interface{}, prefix string) interface{} {
	value, _ := GetMapValueByJSONPath(dataPipeline, prefix)
	if value == nil {
		value = DetectAndConvert(prefix)
	}
	return value
}

func GetMapValueByJSONPath(data map[string]interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".")
	var result interface{} = data

	for _, key := range keys {
		hasArrayVal, arrayName, indexStr := ExtractArrayParts(key)
		if hasArrayVal {
			key = arrayName
		}
		// Type assert the current level to a map
		if currentMap, ok := result.(map[string]interface{}); ok {
			// Fetch the value associated with the key
			if val, exists := currentMap[key]; exists {
				if hasArrayVal {
					valArray := val.([]interface{})
					if len(indexStr) > 0 {
						index, _ := strconv.Atoi(indexStr)
						result = valArray[index]
					} else {
						result = valArray
					}
				} else {
					result = val
				}
			} else {
				return nil, fmt.Errorf("key '%s' not found in path '%s'", key, path)
			}
		} else {
			return nil, fmt.Errorf("path '%s' does not resolve to a nested map structure", path)
		}

		/*
			// Type assert the current level to a map
			if currentMap, ok := result.(map[string]interface{}); ok {
				// Fetch the value associated with the key
				if val, exists := currentMap[key]; exists {
					result = val
				} else {
					return nil, fmt.Errorf("key '%s' not found in path '%s'", key, path)
				}
			} else {
				return nil, fmt.Errorf("path '%s' does not resolve to a nested map structure", path)
			}
		*/

	}
	return result, nil
}

func SetMapValueByJSONPath(m map[string]interface{}, jsonPath string, value interface{}) {
	// Split the JSON path into keys
	keys := strings.Split(jsonPath, ".")

	// Initialize current map as the input map
	currentMap := m

	// Iterate through keys except the last one
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		var index int
		hasArrayVal, arrayName, indexStr := ExtractArrayParts(key)
		if hasArrayVal {
			key = arrayName
			index, _ = strconv.Atoi(indexStr)
		}

		// Check if the key exists and is a map, otherwise create a new map
		/*--------------- array -----------*/
		if _, ok := currentMap[key]; !ok {
			if hasArrayVal {
				arrayObj := make([]interface{}, 1)
				currObj := make(map[string]interface{})
				arrayObj[index] = currObj
				currentMap[key] = arrayObj

				currentMap = currObj
			} else {
				currentMap[key] = make(map[string]interface{})

				if nextMap, ok := currentMap[key].(map[string]interface{}); ok {
					currentMap = nextMap
				}
			}
		} else { //if key exists get the map
			if hasArrayVal {
				if arrayObj, ok := currentMap[key].([]interface{}); ok {
					var currObj map[string]interface{}

					fmt.Printf("\nwhile setting array at index %d, len of arry is %d ", index, len(arrayObj))
					if index >= len(arrayObj) {
						fmt.Printf(", so resizing it to %d", (index + 1))
						newArray := make([]interface{}, index+1)
						copy(newArray, arrayObj)
						currentMap[key] = newArray
						arrayObj = newArray
					}

					ObjAtIndex := arrayObj[index]
					if ObjAtIndex == nil {
						currObj = make(map[string]interface{})
						arrayObj[index] = currObj
					} else {
						currObj = ObjAtIndex.(map[string]interface{})
					}

					currentMap = currObj
				}

				/*
					if arrayObj, ok := currentMap[key].([]interface{}); ok {
						currObj := arrayObj[index].(map[string]interface{})
						currentMap = currObj
					} else {

					}
				*/

			} else {
				if nextMap, ok := currentMap[key].(map[string]interface{}); ok {
					currentMap = nextMap
				}
			}
		}
		//if hasArrayVal {
		//
		//	if nextMap, ok := currentMap[key].(map[string]interface{}); ok {
		//		currentMap = nextMap
		//	}
		//} else {
		//	if nextMap, ok := currentMap[key].(map[string]interface{}); ok {
		//		currentMap = nextMap
		//	}
		//
		//}

		/* ------------- single ---------
		if _, ok := currentMap[key]; !ok {
			currentMap[key] = make(map[string]interface{})
		}
		if nextMap, ok := currentMap[key].(map[string]interface{}); ok {
			currentMap = nextMap
		}
		*/
	}

	// Set the value at the last key
	lastKey := keys[len(keys)-1]
	currentMap[lastKey] = value
	/* ------------- single ---------
	lastKey := keys[len(keys)-1]
	currentMap[lastKey] = value
	*/
}
