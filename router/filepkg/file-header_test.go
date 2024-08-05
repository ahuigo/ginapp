package filepkg

import (
	"bytes"
	"errors"
	"fmt"
	_ "ginapp/conf"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFileHeader(t *testing.T) {
	fh, err := createFileHeader("./go.mod")
	fmt.Println(fh, err)
}

func createFileHeader(filePath string) (*multipart.FileHeader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file1", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	_, params, err := mime.ParseMediaType(writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	boundary, ok := params["boundary"]
	if !ok {
		return nil, errors.New("no boundary")
	}

	reader := multipart.NewReader(body, boundary)
	mf, _ := reader.ReadForm(1 << 8)
	fileHeader := mf.File["file1"][0]

	return fileHeader, nil
}
