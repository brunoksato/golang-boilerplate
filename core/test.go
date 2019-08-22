package core

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/brunoksato/golang-boilerplate/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
)

func OpenTestConnection() (db *gorm.DB, err error) {
	testDB := os.Getenv("SERVER_DB")
	if testDB == "" {
		testDB = "user=postgres dbname=server_test sslmode=disable"
	}
	db, err = gorm.Open("postgres", testDB)
	return
}

// Rest Test
func typeTestHandler(c echo.Context) {
	parent := expectedTypeName(c.Get("ParentType").(reflect.Type))
	child := expectedTypeName(c.Get("Type").(reflect.Type))

	fmt.Fprintf(c.Response(), "Parent Type: %s, Type: %s", parent, child)
}

func expectedTypeName(t reflect.Type) (expected string) {
	if t == nil {
		expected = "nil"
	} else {
		expected = t.Name()
	}
	return
}

func NewTestRequest(method, path string) (*httptest.ResponseRecorder, *http.Request) {
	request, err := http.NewRequest(method, path, nil)
	request.Header.Set("X-Company", "Office")
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error configuring test request:", err)
	}
	recorder := httptest.NewRecorder()

	return recorder, request
}

func NewTestPost(method, path string, body interface{}) (*httptest.ResponseRecorder, *http.Request) {
	jsonStr := ModelToJson(body)
	bodyStr := bytes.NewBufferString(jsonStr)
	request, err := http.NewRequest(method, path, bodyStr)
	request.Header.Set("X-Company", "Office")
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error configuring test request:", err)
	}
	recorder := httptest.NewRecorder()

	return recorder, request
}

func NewTestForm(path string, params map[string]interface{}) (*httptest.ResponseRecorder, *http.Request) {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, fmt.Sprintf("%v", v))
	}

	// Can't seem to get this working. So, for now...
	querypath := fmt.Sprintf("%s?%s", path, values.Encode())
	recorder, request := NewTestRequest("POST", querypath)

	return recorder, request
}

// Assert methods

func AssertResponse(t *testing.T, rr *httptest.ResponseRecorder, body string, code int) {
	if gotBody := strings.TrimSpace(string(rr.Body.Bytes())); body != gotBody {
		t.Errorf("assertResponse: expected body to be %s but got %s. (caller: %s)", body, gotBody, util.CallerInfo())
	}
	AssertResponseCode(t, rr, code)
}

func AssertResponseCode(t *testing.T, rr *httptest.ResponseRecorder, code int) {
	if code != rr.Code {
		t.Errorf("assertResponse: expected code to be %d but got %d. (caller: %s)", code, rr.Code, util.CallerInfo())
	}
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		kind := reflect.ValueOf(expected).Kind()
		if kind == reflect.Map {
			if mExp, ok := expected.(map[string]interface{}); ok {
				mAct := actual.(map[string]interface{})
				for key, value := range mExp {
					AssertEqual(t, value, mAct[key])
				}
			} else if mExp, ok := expected.(map[float32][]uint); ok {
				mAct := actual.(map[float32][]uint)
				for key, value := range mExp {
					AssertEqual(t, value, mAct[key])
				}
			} else {
				mExp := expected.(map[uint]uint)
				mAct := actual.(map[uint]uint)
				for key, value := range mExp {
					AssertEqual(t, value, mAct[key])
				}
			}
		} else if kind == reflect.Array || kind == reflect.Slice {
			eVal := reflect.ValueOf(expected)
			aVal := reflect.ValueOf(actual)
			AssertEqual(t, eVal.Len(), aVal.Len())
			for i := 0; i < aVal.Len(); i++ {
				AssertEqual(t, eVal.Index(i).Interface(), aVal.Index(i).Interface())
			}
		}
	}
}

func AssertNil(t *testing.T, actual interface{}) {
	if actual != nil && !reflect.ValueOf(actual).IsNil() {
		t.Errorf("assertNil: Expected nil value, but got:\n%v\n (caller: %s)", actual, CallerInfo())
	}
}
func AssertNotNil(t *testing.T, actual interface{}) {
	if actual == nil || reflect.ValueOf(actual).IsNil() {
		t.Errorf("assertNotNil: Expected object to not be nil (caller: %s)", CallerInfo())
	}
}
func AssertZeroStruct(t *testing.T, actual interface{}) {
	if !util.IsZeroStruct(actual) {
		t.Errorf("assertZeroStruct: Expected zero struct, but got:\n%v\n (caller: %s)", actual, CallerInfo())
	}
}
func AssertFalse(t *testing.T, result bool) {
	if result {
		t.Errorf("assertFalse: Assertion is true (caller: %s)", CallerInfo())
	}
}
func AssertTrue(t *testing.T, result bool) {
	if !result {
		t.Errorf("assertTrue: Assertion is false (caller: %s)", CallerInfo())
	}
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("assertNoError: Expected no error, but received message: %s (caller: %s)", err, CallerInfo())
	}
}

func AssertWarning(t *testing.T, msg string, err DefaultError) {
	AssertDefaultError(t, "Warning", msg, err, CallerInfo())
}

func AssertBusinessError(t *testing.T, msg string, err DefaultError) {
	AssertDefaultError(t, "Business Error", msg, err, CallerInfo())
}

func AssertPermissionError(t *testing.T, msg string, err DefaultError) {
	AssertDefaultError(t, "Permission Error", msg, err, CallerInfo())
}

func AssertNotFoundError(t *testing.T, msg string, err DefaultError) {
	AssertDefaultError(t, "Not Found Error", msg, err, CallerInfo())
}

func AssertServerError(t *testing.T, msg string, err DefaultError) {
	AssertDefaultError(t, "Server Error", msg, err, CallerInfo())
}

func AssertDefaultError(t *testing.T, prefix, msg string, err DefaultError, callerInfo string) {
	if err == nil {
		t.Errorf("AssertDefaultError: Expected an error, but it was nil (caller: %s)", callerInfo)
		return
	}

	if _, ok := err.(DefaultError); !ok {
		t.Errorf("AssertDefaultError: Expected a petmondo error, but got: %s (caller: %s)", err, callerInfo)
		return
	}

	errMsg := err.Error()

	idx := strings.Index(errMsg, "(caller:")
	if idx != -1 {
		errMsg = errMsg[:idx-1]
	}

	if msg != errMsg {
		t.Errorf("AssertDefaultError: Expected error \"%s\", but got \"%s\" (caller: %s)", msg, errMsg, callerInfo)
		return
	}
}

func LineInfo() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func CallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func ParentCallerInfo() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}
