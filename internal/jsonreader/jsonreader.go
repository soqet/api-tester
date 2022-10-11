package jsonreader

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"os"
	"path/filepath"
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
	if err != nil {
		return nil, err
	}
	err = res.checkPaths(filePath)
	return res, err
}

func (task *MainSchema) checkPaths(taskPath string) error {
	absTaskPath, err := filepath.Abs(taskPath)
	if err != nil {
		return err
	}
	taskDir := filepath.Dir(absTaskPath)
	for _, r := range task.Requests {
		r.processPaths(taskDir)
	}
	for _, g := range task.Groups {
		for _, r := range g.Requests {
			r.processPaths(taskDir)
		}
	}
	return nil
}

func (r *RequestSchema) processPaths(taskDir string) {
	if r.BodyFile != "" && !filepath.IsAbs(r.BodyFile) {
		r.BodyFile = filepath.Join(taskDir, r.BodyFile)
	}
	if r.ResponseFile != "" && !filepath.IsAbs(r.ResponseFile) {
		r.ResponseFile = filepath.Join(taskDir, r.ResponseFile)
	}
}
