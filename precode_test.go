package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := len(cafeList["moscow"])
	reqURL := fmt.Sprintf("/cafe?count=%d&city=moscow", totalCount+1)
	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("expected status code: %d, got %d", http.StatusOK, status)
	}

	got := strings.Split(responseRecorder.Body.String(), ",")
	expected := cafeList["moscow"]

	assert.Len(t, got, totalCount)
	require.Equal(t, expected, got)
}

func TestMainHandlerWhitoutCount(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/cafe?city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code: %d, got %d", http.StatusBadRequest, status)
	}

	got := responseRecorder.Body.String()

	require.Equal(t, errorCountMissing, got)
}

func TestMainHandlerWhenCountIsAString(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/cafe?count=five&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code: %d, got %d", http.StatusBadRequest, status)
	}

	got := responseRecorder.Body.String()

	require.Equal(t, errorCountWrongValue, got)
}

func TestMainHandlerWhenCountLessThanHave(t *testing.T) {
	totalCount := 2
	reqURL := fmt.Sprintf("/cafe?count=%d&city=moscow", totalCount)

	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("expected status code: %d, got %d", http.StatusOK, status)
	}

	expected := []string{"Мир кофе", "Сладкоежка"}
	got := strings.Split(responseRecorder.Body.String(), ",")

	require.Len(t, got, totalCount)
	require.Equal(t, expected, got)
}

func TestMainHandlerWhenCityIsNotInList(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/cafe?count=4&city=cairo", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code: %d, got %d", http.StatusBadRequest, status)
	}

	got := responseRecorder.Body.String()

	require.Equal(t, errorCityWrongValue, got)
}
