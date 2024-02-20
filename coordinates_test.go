package main

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	strLatDenver = "39.7392"
	strLonDenver = "-105.1704"
)

func Test_ParseLatLon_Fail(t *testing.T) {
	lat, err := ParseLatLon("abc")
	assert.Equal(t, lat, float64(0))
	assert.Equal(t, err, errors.New("invalid lat/lon value: 'abc'"))
}

func Test_ParseLatLon_Succeed(t *testing.T) {
	lat, _ := ParseLatLon(strLatDenver)
	assert.Equal(t, lat, 39.7392)
	lon, _ := ParseLatLon(strLonDenver)
	assert.Equal(t, lon, -105.1704)
	lat2, _ := ParseLatLon(" " + strLatDenver + " ")
	assert.Equal(t, lat2, 39.7392)
}

func Test_ParseCoordinates_Fail(t *testing.T) {
	latlon, err := ParseCoordinates("", strLonDenver)
	assert.Equal(t, latlon, Coordinates{})
	assert.Equal(t, err, errors.New("invalid lat/lon value: ''"))
}

func Test_ParseCoordinates_Succeed(t *testing.T) {
	latlon, _ := ParseCoordinates(strLatDenver, strLonDenver)
	latDenver, _ := strconv.ParseFloat(strLatDenver, 64)
	lonDenver, _ := strconv.ParseFloat(strLonDenver, 64)
	assert.Equal(t, latlon, Coordinates{Lat: latDenver, Lon: lonDenver})
}
