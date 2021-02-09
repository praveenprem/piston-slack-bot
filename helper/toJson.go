package helper

import (
	"encoding/json"
	"log"
)

func ToJson(data interface{}) (*string, error) {
	jsonObject, err := json.Marshal(data)
	if err != nil {
		log.Printf("%#v", err.Error())
		return nil, err
	}
	jsonString := string(jsonObject)
	log.Printf("%s", jsonString)
	return &jsonString, nil
}
