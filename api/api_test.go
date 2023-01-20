package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPut(t *testing.T) {
	handler := NewSdHandler()

	// Test bad request
	request_str := []byte(`{"asdfgajnk":"adsfhadsgf","afdadfs":"asfds"}`)
	req, err := http.NewRequest("PUT", "/", bytes.NewBuffer(request_str))

	if err != nil {
		t.Error("Error creating request")
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	expected_body := "Invalid Data Passed"
	if rr.Code != http.StatusBadRequest {
		t.Errorf("On bad data expected code %v, received %v",
			http.StatusBadRequest, rr.Code)
	}

	if strings.TrimSpace(rr.Body.String()) != expected_body {
		t.Errorf(`On bad data expected message "%s", received "%s"`,
			expected_body, strings.TrimSpace(rr.Body.String()))
	}

	// test proper request
	request_str = []byte(`{"target":"test1", "labels" : [{"x":"y"}]}`)
	req, err = http.NewRequest("PUT", "/", bytes.NewBuffer(request_str))

	if err != nil {
		t.Error("Error creating request")
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	// check that data was stored in handler
	_, target_ok := handler.Targets["test1"]
	if !target_ok {
		t.Error("Failed to create target")
	} else {
		label, label_ok := handler.Targets["test1"][0]["x"]
		if !label_ok {
			t.Error("Labels not created")
		}
		if label != "y" {
			t.Error("Incorrect label registered")
		}
	}
}

func TestGet(t *testing.T) {
	handler := NewSdHandler()

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Body.String() != "[]" {
		t.Error("Incorrect response from empty service")
	}

	handler.Targets["test1"] = []map[string]string{{"x": "y"}, {"a": "b"}}
	rr = httptest.NewRecorder()
	expected := `[{"target":"test1","labels":[{"x":"y"},{"a":"b"}]}]`
	handler.ServeHTTP(rr, req)
	if rr.Body.String() != expected {
		t.Errorf("Response Error:  \n  expected: %s \n received: %s", expected, rr.Body.String())
	}
}

func TestDelete(t *testing.T) {
	handler := NewSdHandler()

	// test request missing form

	req, err := http.NewRequest("DELETE", "/?x=y", nil)

	if err != nil {
		t.Error("Unable to create request")
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected_body := "No deletion target specified"
	if rr.Code != http.StatusBadRequest {
		t.Errorf("On bad data expected code %v, received %v",
			http.StatusBadRequest, rr.Code)
	}

	if strings.TrimSpace(rr.Body.String()) != expected_body {
		t.Errorf(`On bad delete request expected message "%s", received "%s"`,
			expected_body, strings.TrimSpace(rr.Body.String()))
	}

	// test delete of non-existent tarvet
	req, err = http.NewRequest("DELETE", "/?target=test1", nil)

	if err != nil {
		t.Error("Unable to create request")
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected_body = "Invalid target specified"
	if rr.Code != http.StatusBadRequest {
		t.Errorf("On bad data expected code %v, received %v",
			http.StatusBadRequest, rr.Code)
	}

	if strings.TrimSpace(rr.Body.String()) != expected_body {
		t.Errorf(`On bad delete request expected message "%s", received "%s"`,
			expected_body, strings.TrimSpace(rr.Body.String()))
	}

	// test happy path
	handler.Targets["test1"] = []map[string]string{{"x": "y"}, {"a": "b"}}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	_, ok := handler.Targets["test1"]

	if ok {
		t.Error("Failed to delete target")
	}

}
