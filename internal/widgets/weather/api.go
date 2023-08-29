package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/merlinfuchs/blimp/internal/config"
)

type weatherMain struct {
	Temp      float64 `json:"temp"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  float64 `json:"pressure"`
	SeaLevel  float64 `json:"sea_level"`
	GrndLevel float64 `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
}

type weatherWeather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type CurrentWeatherData struct {
	Name string `json:"name"`
	Sys  struct {
		Sunrise int `json:"sunrise"`
		Sunset  int `json:"sunset"`
	} `json:"sys"`
	Weather []weatherWeather `json:"weather"`
	Main    weatherMain      `json:"main"`
}

type ForecastWeatherData struct {
	List []struct {
		Dt      int                `json:"dt"`
		Main    weatherMain        `json:"main"`
		Weather []weatherWeather   `json:"weather"`
		Rain    map[string]float32 `json:"rain"`
	}
}

func getCurrentWeatherData() (CurrentWeatherData, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=%s&lang=%s",
		config.K.Float64("widgets.weather.owm_lat"),
		config.K.Float64("widgets.weather.owm_lon"),
		config.K.String("widgets.weather.owm_api_key"),
		config.K.String("widgets.weather.owm_unit"),
		config.K.String("widgets.weather.owm_language"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return CurrentWeatherData{}, fmt.Errorf("Failed to get weather data, %w", err)
	}

	var data CurrentWeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return CurrentWeatherData{}, fmt.Errorf("Failed to decode weather data, %w", err)
	}

	return data, nil
}

func getWeatherForecast() (ForecastWeatherData, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&appid=%s&units=%s&lang=%s",
		config.K.Float64("widgets.weather.owm_lat"),
		config.K.Float64("widgets.weather.owm_lon"),
		config.K.String("widgets.weather.owm_api_key"),
		config.K.String("widgets.weather.owm_unit"),
		config.K.String("widgets.weather.owm_language"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return ForecastWeatherData{}, fmt.Errorf("Failed to get weather data, %w", err)
	}

	var data ForecastWeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ForecastWeatherData{}, fmt.Errorf("Failed to decode weather data, %w", err)
	}

	return data, nil
}
