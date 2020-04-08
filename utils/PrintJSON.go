package utils

import (
	"encoding/json"
)

func PrintJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	println(string(jsonData))
}
