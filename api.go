package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// ------------------------------------
// Open Weather Map API Docs:
// - https://openweathermap.org/current
// ------------------------------------

type WeatherCondition struct {
	Activity  []string `json:"activity"`
	Condition string   `json:"condition"`
}

type Wind struct {
	Speed  float32 `json:"speed"`
	Degree int     `json:"deg"`
	Gust   float32 `json:"gust"`
}

type Clouds struct {
	All int `json:"all"`
}

type Rain struct {
	OneHour float32 `json:"1h"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type System struct {
	ID      int32  `json:"id"`
	Type    int    `json:"type"`
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

type Main struct {
	Temp        float32 `json:"temp"`
	FeelsLike   float32 `json:"feels_like"`
	TempMin     float32 `json:"temp_min"`
	TempMax     float32 `json:"temp_max"`
	Pressure    int     `json:"pressure"`
	Humidity    int     `json:"humidity"`
	SeaLevel    int     `json:"sea_level"`
	GroundLevel int     `json:"grnd_level"`
}

type ApiResponse struct {
	ID         int32       `json:"id"`
	Code       int         `json:"cod"`
	Name       string      `json:"name"`
	Base       string      `json:"base"`
	Timezone   int         `json:"timezone"`
	Visibility int         `json:"visibility"`
	Coords     Coordinates `json:"coord"`
	Datetime   int64       `json:"dt"`
	Wind       Wind        `json:"wind"`
	System     System      `json:"sys"`
	Main       Main        `json:"main"`
	Weather    []Weather   `json:"weather"`
	Clouds     Clouds      `json:"clouds"`
	Rain       Rain        `json:"rain"`
}

type ApiError struct {
	Code    int    `json:"cod"`
	Message string `json:"message"`
}

type ApiService struct {
	BaseUrl string `json:"base_url"`
	ApiKey  string `json:"api_key"`
}

func createUrl(baseUrl string, coords Coordinates, apiKey string) string {
	return fmt.Sprintf(
		"%s/weather?lat=%s&lon=%s&apikey=%s&units=imperial",
		baseUrl,
		strconv.FormatFloat(coords.Lat, 'f', -1, 64),
		strconv.FormatFloat(coords.Lon, 'f', -1, 64),
		apiKey,
	)
}

func NewApiService(baseUrl string, apiKey string) *ApiService {
	return &ApiService{
		BaseUrl: baseUrl,
		ApiKey:  apiKey,
	}
}

func (res *ApiResponse) GetWeatherCondition() WeatherCondition {
	weatherCondition := WeatherCondition{
		Condition: res.Main.GetTemperatureDescription(),
	}

	for _, weather := range res.Weather {
		if weather.Main != "" {
			weatherCondition.Activity = append(weatherCondition.Activity, weather.Main)
		}
	}

	return weatherCondition
}

func (main *Main) GetTemperatureDescription() string {
	if main.FeelsLike > 85 {
		return "hot"
	} else if main.FeelsLike < 60 {
		return "cold"
	}
	return "moderate"
}

// Certain errors aren't intended to be sent back to the calling client
// Instead they are logged and a more generalized error is returned back to the caller
func handleServerError(err error, apiErr *ApiError) {
	log.Println(err)
	apiErr.Message = "An unknown error occurred. Please contact support."
	apiErr.Code = http.StatusInternalServerError
}

func (service *ApiService) GetCurrentWeather(coords Coordinates) (*ApiResponse, *ApiError) {
	url := createUrl(service.BaseUrl, coords, service.ApiKey)
	res, err := http.Get(url)
	if err != nil {
		var apiErr ApiError
		handleServerError(err, &apiErr)
		return nil, &apiErr
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		var apiErr ApiError
		handleServerError(err, &apiErr)
		return nil, &apiErr
	}

	if res.StatusCode != http.StatusOK {
		var apiErr ApiError
		err = json.Unmarshal(body, &apiErr)
		if err != nil {
			handleServerError(err, &apiErr)
			return nil, &apiErr
		}
		return nil, &apiErr
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		var apiErr ApiError
		handleServerError(err, &apiErr)
		return nil, &apiErr
	}

	return &apiResponse, nil
}
