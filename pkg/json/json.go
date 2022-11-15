package json

import json "github.com/goccy/go-json"

type FastJSON struct{}

func (self *FastJSON) EncodeToByte(obj any) ([]byte, error) {
	return json.Marshal(obj)
}

func (self *FastJSON) EncodeToString(obj any) (*string, error) {
	jsonByte, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	jsonStr := string(jsonByte)
	return &jsonStr, nil
}

func (self *FastJSON) DecodeFromString(jsonStr *string, obj any) error {
	return json.Unmarshal([]byte(*jsonStr), &obj)
}

func (self *FastJSON) DecodeFromBytes(jsonByte []byte, obj any) error {
	return json.Unmarshal(jsonByte, &obj)
}
