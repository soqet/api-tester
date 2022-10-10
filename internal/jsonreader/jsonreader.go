package jsonreader

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"os"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ReadFile(filePath string) (*MainSchema, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := bytes.Buffer{}
	buf.ReadFrom(file)
	res := new(MainSchema)
	err = json.Unmarshal(buf.Bytes(), res)
	return res, err
}
