package backend

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func init() {
	log = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%] %time% - %msg%",
		},
	}
}

func TestHandleGetTime(t *testing.T) {

	// Create a ResponseRecorder to record the response
	responseRecorder := httptest.NewRecorder()

	// Create a request to pass to the handler
	request, err := http.NewRequest("GET", "/get-time", nil)
	if err != nil {
		t.Fatal(err)
	}

	initialTime := time.Now().Format(time.RFC3339)

	// Call the handler
	handleGetTime()(responseRecorder, request)

	// Check the status code is what we expect
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedTime string
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &returnedTime)
	if err != nil {
		t.Errorf("unmarshal failed, for %v", responseRecorder.Body.String())
	}

	if initialTime != returnedTime {
		t.Errorf("time is not matching, for %v and %v", initialTime, returnedTime)
	}
}
