package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	baseUrl = "https://api.openweathermap.org/data/2.5"
	apiKey  = "abc123"
)

func Test_CreateUrl(t *testing.T) {
	urlPattern := regexp.MustCompile(`^(https?:\/\/)(api.openweathermap.org\/data)\/(\d+\.?\d*)\/(weather)\?(lat=` + strLatDenver + `)&(lon=` + strLonDenver + `)&(apikey=` + apiKey + `)(&units=(imperial|metric|standard))?$`)
	coords, _ := ParseCoordinates(strLatDenver, strLonDenver)
	url := createUrl(baseUrl, coords, apiKey)
	assert.Regexp(t, urlPattern, url)
}

func Test_GetTemperaturDescription(t *testing.T) {
	main := Main{
		Temp:        73.2,
		FeelsLike:   78.0,
		TempMin:     68.7,
		TempMax:     82.0,
		Pressure:    12,
		Humidity:    32,
		SeaLevel:    5984,
		GroundLevel: 0,
	}
	assert.Equal(t, main.GetTemperatureDescription(), "moderate")
	main.FeelsLike = 58.4
	assert.Equal(t, main.GetTemperatureDescription(), "cold")
	main.FeelsLike = 89.2
	assert.Equal(t, main.GetTemperatureDescription(), "hot")
}

func Test_NewOpenWeathermapService(t *testing.T) {
	service := NewApiService(baseUrl, apiKey)
	baseUrlPattern := regexp.MustCompile(`^(https?:\/\/)(api.openweathermap.org\/data)\/(\d+\.?\d*)$`)
	assert.Equal(t, service.ApiKey, apiKey)
	assert.Regexp(t, baseUrlPattern, service.BaseUrl)
}

var successMarshaledResponse = ApiResponse{
	ID:       3163858,
	Name:     "Zocca",
	Code:     200,
	Timezone: 7200,
	Datetime: 1661870592,
	Rain: Rain{
		OneHour: 3.16,
	},
	Clouds: Clouds{
		All: 100,
	},
	Visibility: 10000,
	Wind: Wind{
		Speed:  0.62,
		Degree: 349,
		Gust:   1.18,
	},
	System: System{
		ID:      2075663,
		Type:    2,
		Country: "IT",
		Sunrise: 1661834187,
		Sunset:  1661882248,
	},
	Base: "stations",
	Coords: Coordinates{
		Lon: 10.99,
		Lat: 44.34,
	},
	Main: Main{
		Temp:        298.48,
		FeelsLike:   298.74,
		TempMin:     297.56,
		TempMax:     300.05,
		Pressure:    1015,
		Humidity:    64,
		SeaLevel:    1015,
		GroundLevel: 933,
	},
	Weather: []Weather{
		{
			ID:          501,
			Main:        "Rain",
			Description: "moderate rain",
			Icon:        "10d",
		},
	},
}

func Test_GetWeatherCondition(t *testing.T) {
	condition := WeatherCondition{
		Activity:  []string{"Rain"},
		Condition: "hot",
	}
	assert.Equal(t, successMarshaledResponse.GetWeatherCondition(), condition)
}

const successExampleResponse = `{
  "id": 3163858,
  "cod": 200,
  "name": "Zocca",
  "dt": 1661870592,
  "timezone": 7200,
  "base": "stations",
  "visibility": 10000,
  "coord": {
    "lon": 10.99,
    "lat": 44.34
  },
  "weather": [
    {
      "id": 501,
      "main": "Rain",
      "description": "moderate rain",
      "icon": "10d"
    }
  ],
  "main": {
    "temp": 298.48,
    "feels_like": 298.74,
    "temp_min": 297.56,
    "temp_max": 300.05,
    "pressure": 1015,
    "humidity": 64,
    "sea_level": 1015,
    "grnd_level": 933
  },
  "wind": {
    "speed": 0.62,
    "deg": 349,
    "gust": 1.18
  },
  "rain": {
    "1h": 3.16
  },
  "clouds": {
    "all": 100
  },
  "sys": {
    "id": 2075663,
    "type": 2,
    "country": "IT",
    "sunrise": 1661834187,
    "sunset": 1661882248
  }
}`

const invalidExampleResponse = `{
  "cod": 401,
  "message": "Invalid API key. Please see https://openweathermap.org/faq#error401 for more info."
}`

func handleInvalidResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(invalidExampleResponse))
}

func handleSuccessResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(successExampleResponse))
}

func Test_GetCurrentWeather_Success(t *testing.T) {
	successResponseServer := httptest.NewServer(http.HandlerFunc(handleSuccessResponse))
	service := NewApiService(successResponseServer.URL, apiKey)
	apiResponse, apiError := service.GetCurrentWeather(successMarshaledResponse.Coords)
	assert.Nil(t, apiError)
	assert.Equal(t, apiResponse, &successMarshaledResponse)
}

func TestApiService_GetUserDetails_InvalidResponse(t *testing.T) {
	invalidResponseServer := httptest.NewServer(http.HandlerFunc(handleInvalidResponse))
	service := NewApiService(invalidResponseServer.URL, apiKey)
	_, apiError := service.GetCurrentWeather(successMarshaledResponse.Coords)
	assert.Equal(t, http.StatusUnauthorized, apiError.Code)
	assert.Equal(t, apiError.Message, "Invalid API key. Please see https://openweathermap.org/faq#error401 for more info.")
}
