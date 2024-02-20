# Current Weather Service

A project which queries the current weather conditions.

The [OpenWeatherMap](https://openweathermap.org) service provides data on global weather and offers a variety of API endpoints. For the sake of this project only the current weather conditions are retrieved from Open WeatherMap.

## Objective

To demonstrate the use of a web server in the Go programming language and which interacts with one or more external services.

The source code was created to be simple yet graft in good programming practices which will allow it to be scaled out (if necessary) and easily tested. The code has also been anchored in with unit tests and mocked to prevent additional costs from interacting with a paid service like OpenWeatherMap.

## Dependencies

You will need to set an environment variable for the API key, named `OPENWEATHER_API_KEY` and would recommend creating a `.env` file in the project root to store that value.

## Usage

To run the application:

```sh
go run .
```

To run the tests:

```sh
go test
```


### Endpoints
The server provides the following endpoints:

- `/weather` (GET) - Retrieves the current weather conditions at the specified latitude and longitude
    - `lat` (required) - Query string parameter which specifies the latitude
    - `lon` (required) - Query string parameter which specifies the longitude

```
/weather?lat=44.34&lon=10.99
```
