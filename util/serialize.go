package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// MustJsonToBytes converts JSON object to bytes and if error occurs, it prints error to Stderr
func MustJsonToBytes(data interface{}) []byte {
	bytes, err := JsonToBytes(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	return bytes
}

// JsonToBytes converts JSON object to bytes and if error occurs, returns it
func JsonToBytes(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %v", err)
	}

	return bytes, nil
}
