package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CurrentWeather_Route_BadRequest(t *testing.T) {
	successResponseServer := httptest.NewServer(http.HandlerFunc(handleSuccessResponse))
	router := SetupRoutes(successResponseServer.URL, apiKey)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/weather", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "Lat/lon value is invalid", w.Body.String())
}

func Test_CurrentWeather_Route(t *testing.T) {
	requestUrl := fmt.Sprintf(
		"/weather?lat=%s&lon=%s",
		strconv.FormatFloat(successMarshaledResponse.Coords.Lat, 'f', -1, 64),
		strconv.FormatFloat(successMarshaledResponse.Coords.Lon, 'f', -1, 64),
	)

	successResponseServer := httptest.NewServer(http.HandlerFunc(handleSuccessResponse))
	router := SetupRoutes(successResponseServer.URL, apiKey)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", requestUrl, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"activity":["Rain"],"condition":"hot"}`, w.Body.String())
}
