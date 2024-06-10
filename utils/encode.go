package utils

import "encoding/base64"

func Encode(code string) string {
	data := base64.StdEncoding.EncodeToString([]byte(code))
	return string(data)
}

func Decode(s string) string {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(data)
}
