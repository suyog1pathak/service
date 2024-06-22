package util

import (
	"encoding/json"
	"strconv"
)

func StringToInt(data string) (int, error) {
	o, err := strconv.Atoi(data)
	if err != nil {
		return 0, err
	}
	return o, nil
}

func StructToJson(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}
