package disgostutest

import (
	"os"
	"path"
	"fmt"
	"github.com/Oppodelldog/disgostu/capturing"
	"encoding/json"
	"io/ioutil"
)

const recordsFolderName = "records"

func getCaptureFilePath() (string, error) {
	recordsFolderPath, err := getRecordsFolderPath()
	if err != nil {
		return "", err
	}
	filePath := path.Join(recordsFolderPath, fmt.Sprintf("%s.json", captureFileName))

	return filePath, nil
}

func getRecordsFolderPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(wd, recordsFolderName), nil
}

func ensureRecordsFolder() error {
	recordsFolderPath, err := getRecordsFolderPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(recordsFolderPath, 0766)
}

func writeCapturedDataToFile(capture capturing.StoreCapture) error {
	err := ensureRecordsFolder()
	if err != nil {
		return err
	}

	capturedFilePath, err := getCaptureFilePath()
	if err != nil {
		return err
	}

	var encodedBytes []byte

	encodedBytes, err = json.Marshal(capture)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(capturedFilePath, encodedBytes, 0666)
}

func loadCapturedDataFromFile() (*capturing.StoreCapture, error) {

	capturedFilePath, err := getCaptureFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(capturedFilePath); err != nil {
		return nil, nil
	}

	fileBytes, err := ioutil.ReadFile(capturedFilePath)
	if err != nil {
		return nil, err
	}

	var capturedData *capturing.StoreCapture
	err = json.Unmarshal(fileBytes, &capturedData)
	if err != nil {
		return nil, err
	}
	return capturedData, nil
}
