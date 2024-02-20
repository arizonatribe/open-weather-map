package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Using a light dependency-injection pattern, modeled here:
// - https://stackoverflow.com/a/60324692/2022208

type WeatherController struct {
	Service *ApiService
}

func (controller *WeatherController) GetCurrentWeatherByLocation(c *gin.Context) {
	q := c.Request.URL.Query()

	coords, err := ParseCoordinates(q.Get("lat"), q.Get("lon"))
	if err != nil {
		c.String(http.StatusBadRequest, "Lat/lon value is invalid")
		return
	}

	res, apiErr := controller.Service.GetCurrentWeather(coords)
	if apiErr != nil {
		c.String(apiErr.Code, apiErr.Message)
		return
	}

	c.JSON(http.StatusOK, res.GetWeatherCondition())
}

func SetupRoutes(baseUrl string, apiKey string) *gin.Engine {
	r := gin.New()
	c := WeatherController{Service: NewApiService(baseUrl, apiKey)}

	r.GET("/weather", c.GetCurrentWeatherByLocation)

	return r
}
