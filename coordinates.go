package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// --------------------------------------------------
// Regular expression used to validate lat/lon values
//
// References:
// - https://ihateregex.io/expr/lat-long/
// - https://stackoverflow.com/a/37976143/2022208
// - https://stackoverflow.com/q/30784608/2022208
var /* const */ patternLatLon = regexp.MustCompile(`^(\-?|\+?)?\d+(\.\d+)?$`)

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func ParseLatLon(val string) (float64, error) {
	if patternLatLon.MatchString(strings.TrimSpace(val)) {
		latlon, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		return latlon, err
	}
	return 0, fmt.Errorf("invalid lat/lon value: '%s'", val)
}

func ParseCoordinates(strLat string, strLon string) (Coordinates, error) {
	lat, err := ParseLatLon(strLat)
	if err != nil {
		return Coordinates{}, err
	}

	lon, err := ParseLatLon(strLon)
	if err != nil {
		return Coordinates{}, err
	}

	return Coordinates{
		Lat: lat,
		Lon: lon,
	}, nil
}
