package utils

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wDRxxx/avito-shop/internal/models"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()

	body := &models.AuthResponse{
		Token: "test",
	}

	err := WriteJSON(body, w)
	if err != nil {
		t.Fatal(err)
	}

	expectedBody := `{"token":"test"}`
	require.Equal(t, expectedBody, w.Body.String())
}

func TestWriteJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	resultErr := errors.New("error text")
	err := WriteJSONError(resultErr, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedBody := `{"errors":"error text"}`
	require.Equal(t, expectedBody, w.Body.String())
}

func TestReadReqJSON(t *testing.T) {
	type Request struct {
		Name string `json:"name"`
	}

	jsonData := `{"name":"test"}`
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(jsonData)))
	w := httptest.NewRecorder()

	var req Request
	err := ReadReqJSON(w, r, &req)

	require.NoError(t, err)
	require.Equal(t, "test", req.Name)
}

func TestReadJSON_InvalidJSON(t *testing.T) {
	invalidJSON := `{"name":}`
	reader := bytes.NewReader([]byte(invalidJSON))

	type Request struct {
		Name string `json:"name"`
	}
	var req Request

	err := ReadJSON(reader, &req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
