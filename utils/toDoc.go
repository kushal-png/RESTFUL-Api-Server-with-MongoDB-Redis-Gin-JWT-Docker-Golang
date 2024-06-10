package utils

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func ToDoc(data interface{}) (bson.D, error) {
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON to BSON
	var doc bson.D
	err = bson.UnmarshalExtJSON(jsonData, true, &doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
