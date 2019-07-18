package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestStatusHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/healthz", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(healthHandler)

    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := fmt.Sprintf(`{"status":"ok","running_since":"%v"}`, RUNNING_SINCE)
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
