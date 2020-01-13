package utils

import (
	"encoding/json"
)

func PrintJson(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	println(string(jsonData))
}
