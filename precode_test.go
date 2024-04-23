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

const testingMapKey = "TestCity#42"

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {

	_, ok := cafeList[testingMapKey]
	require.Falsef(t, ok, "Testing data already appering with key \"%s\"", testingMapKey)
	cafeList[testingMapKey] = []string{"test1", "test2", "test3"}
	defer delete(cafeList, testingMapKey)

	totalCount := 3
	reqURL := fmt.Sprintf("/cafe?count=%d&city=%s", totalCount+1, testingMapKey)
	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	got := strings.Split(responseRecorder.Body.String(), ",")
	expected := cafeList[testingMapKey]

	assert.Len(t, got, totalCount)
	require.ElementsMatch(t, expected, got)
}

func TestMainHandlerWhitoutCount(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/cafe?city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	got := responseRecorder.Body.String()

	require.Equal(t, errorCountMissing, got)
}

func TestMainHandlerWhenCountIsAString(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/cafe?count=five&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	got := responseRecorder.Body.String()

	require.Equal(t, errorCountWrongValue, got)
}

func TestMainHandlerWhenCountLessThanHave(t *testing.T) {
	_, ok := cafeList[testingMapKey]
	require.Falsef(t, ok, "Testing data already appering with key \"%s\"", testingMapKey)
	cafeList[testingMapKey] = []string{"test1", "test2", "test3"}
	defer delete(cafeList, testingMapKey)

	totalCount := 2
	reqURL := fmt.Sprintf("/cafe?count=%d&city=%s", totalCount, testingMapKey)

	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	expected := []string{"test1", "test2"}
	got := strings.Split(responseRecorder.Body.String(), ",")

	require.Len(t, got, totalCount)
	require.ElementsMatch(t, expected, got)
}

func TestMainHandlerWhenCityIsNotInList(t *testing.T) {

	_, ok := cafeList[testingMapKey]
	require.Falsef(t, ok, "Testing data already appering with key \"%s\"", testingMapKey)

	reqURL := fmt.Sprintf("/cafe?count=3&city=%s", testingMapKey)
	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	got := responseRecorder.Body.String()

	require.Equal(t, errorCityWrongValue, got)
}
