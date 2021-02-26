package database

import (
	"encoding/json"
)

func marshal(model Model) (string, error) {
	b, err := json.Marshal(model)
	if nil != err {
		return "", err
	}
	return string(b), err
}

func unmarshal(model Model, data string) error {
	return json.Unmarshal([]byte(data), model)
}
