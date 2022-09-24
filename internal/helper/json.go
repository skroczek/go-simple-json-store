package helper

import (
	"encoding/json"
	"log"
)

// FromJSON is a helper function to convert a byte array to an interface direct usable with functions like io.ReadAll
func FromJSON(data []byte, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}
	var o interface{}
	err = json.Unmarshal(data, &o)
	return o, err
}

func ToJSON(object interface{}) []byte {
	data, err := json.Marshal(object)
	if err != nil {
		log.Panicf("Error: %+v", err)
	}
	return data
}
