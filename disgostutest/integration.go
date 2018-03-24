package disgostutest

import (
	"github.com/Oppodelldog/disgostu/capturing"
	"testing"
	"github.com/kr/pretty"
	"strings"
)

var session *capturing.Session
var previousCapturedData *capturing.StoreCapture
var testingT *testing.T
var captureFileName string
var diffFunc captureDiffFunc

type captureDiffFunc func(capturing.StoreCapture, capturing.StoreCapture) []string

func StartRecording(t *testing.T, config capturing.CapturingConfig) {
	testingT = t

	captureFileName = config.RecordName

	var err error
	previousCapturedData, err = loadCapturedDataFromFile()
	if err != nil {
		t.Fatalf("error while loading previously captured file: %v", err)
	}
	session = capturing.NewSession(config)
	session.Start()
}

func Assert() {
	session.Stop()
	if capturedData, ok := session.GetCapturedData(); ok {
		if previousCapturedData != nil {
			var df captureDiffFunc
			if diffFunc != nil {
				df = diffFunc
			} else {
				df = ignoreTimeOffsetDiffFunc
			}

			differences := df(*previousCapturedData, capturedData)
			if len(differences) > 0 {
				testingT.Fatalf("captured data does not match the recorded data:\n%v", strings.Join(differences, "\n"))
			}

		} else {
			err := writeCapturedDataToFile(capturedData)
			if err != nil {
				testingT.Fatalf("error while writing captured data to file: %v", err)
			}
		}
	} else {
		testingT.Fatalf("captured data not found")

	}
}

func defaultDiffFunc(previousCapture, currentCapture capturing.StoreCapture) []string {
	return pretty.Diff(previousCapture, currentCapture)
}

func ignoreTimeOffsetDiffFunc(previousCapture, currentCapture capturing.StoreCapture) []string {
	res := []string{}
	diffs := pretty.Diff(previousCapture, currentCapture)
	for _, diff := range diffs {
		if strings.Contains(diff, "TimeOffset") {
			continue
		}

		res = append(res, diff)
	}

	return res
}
