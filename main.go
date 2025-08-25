package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const version = "v1.0.0"

const (
	currentWeatherAPI = "https://api.openweathermap.org/data/2.5/weather"
	forecastAPI       = "https://api.openweathermap.org/data/2.5/forecast"
)

type WeatherResponse struct {
	Name    string        `json:"name"`
	Main    MainData      `json:"main"`
	Weather []WeatherInfo `json:"weather"`
}

type MainData struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMax   float64 `json:"temp_max"`
	Humidity  float64 `json:"humidity"`
}

type WeatherInfo struct {
	Description string `json:"description"`
}

type ForecastItem struct {
	Weather []WeatherInfo `json:"weather"`
	Main    MainData      `json:"main"`
	DtTxt   string        `json:"dt_txt"`
}

type CityInfo struct {
	Name string `json:"name"`
}

type ForecastResponse struct {
	List []ForecastItem `json:"list"`
	City CityInfo       `json:"city"`
}

func main() {
	cityflag := flag.String("city", "", "The city for which to fetch the weather forecast.")
	unitflag := flag.String("unit", "metric", "The unit for the temperature.")

	forecastFlag := flag.Bool("forecast", false, "Get the forecast for the next 5 days.")

	versionFlag := flag.Bool("version", false, "Get the version of the program.")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Weather Fetcher CLI version is: %s\n", version)
		os.Exit(0)
	}

	// fmt.Printf("city: %s, unit: %s\n", *cityflag, *unitflag)

	if *cityflag == "" {
		fmt.Fprintln(os.Stderr, "Please provide a city name using the --city flag.")
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	weather_api := os.Getenv("API_Key")
	if weather_api == "" {
		log.Fatalf("Weather API Key is not set")
	}

	if *forecastFlag {
		forecastData, err := getForecast(*cityflag, *unitflag, weather_api)
		if err != nil {
			log.Fatal("error getting forecast data:", err)
		}
		displayForecast(forecastData, *unitflag)
	} else {
		weatherData, errors := getWeather(*cityflag, *unitflag, weather_api)
		if errors != nil {
			log.Fatal(errors)
		}
		displayWeather(weatherData, *unitflag)
	}
}

func getWeather(location string, unit string, weather_api string) (WeatherResponse, error) {
	// fmt.Printf("fetching weather for %s using this api key: %s\n", location, weather_api)
	var weatherData WeatherResponse

	fullUrl := fmt.Sprintf("%s?q=%s&appid=%s&units=%s", currentWeatherAPI, location, weather_api, unit)
	// fmt.Println(fullUrl)

	res, err := http.Get(fullUrl)
	if err != nil {
		return weatherData, fmt.Errorf("error while getting response from server: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return weatherData, fmt.Errorf("error while reading response body: %w", err)
	}
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return weatherData, fmt.Errorf("error while unmarshalling response body: %w", err)
	}
	return weatherData, nil
}

func getForecast(location string, unit string, weather_api string) (ForecastResponse, error) {
	var forecastData ForecastResponse

	forecastURL := fmt.Sprintf("%s?q=%s&appid=%s&units=%s", forecastAPI, location, weather_api, unit)

	forecastResponse, err := http.Get(forecastURL)
	if err != nil {
		return forecastData, fmt.Errorf("error while getting response from server: %w", err)
	}
	defer forecastResponse.Body.Close()

	body, err := io.ReadAll(forecastResponse.Body)
	if err != nil {
		return forecastData, fmt.Errorf("error while reading response body: %w", err)
	}
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &forecastData)
	if err != nil {
		return forecastData, fmt.Errorf("error while unmarshalling response body: %w", err)
	}
	return forecastData, nil

}

func displayWeather(data WeatherResponse, unit string) {
	fmt.Printf("Current Weather for %s:\n", data.Name)

	var unitSymbol string
	if unit == "imperial" {
		unitSymbol = "°F"
	} else {
		unitSymbol = "°C"
	}

	fmt.Printf("Temperature: %.2f%s\n", data.Main.Temp, unitSymbol)
	fmt.Printf("Maximum Temperature: %.2f%s\n", data.Main.TempMax, unitSymbol)
	fmt.Printf("Condition: %s\n", data.Weather[0].Description)
	fmt.Printf("Humidity: %.2f%%\n", data.Main.Humidity)
	fmt.Printf("Feels like: %.2f%s\n", data.Main.FeelsLike, unitSymbol)
}

func displayForecast(data ForecastResponse, unit string) {
	fmt.Printf("5 day forecast for %s:\n", data.City.Name)
	for _, item := range data.List {
		var unitSymbol string
		if unit == "imperial" {
			unitSymbol = "°F"
		} else {
			unitSymbol = "°C"
		}
		if strings.Contains(item.DtTxt, "12:00:00") {
			date := strings.Split(item.DtTxt, " ")[0]

			fmt.Printf("%s: Temp: %.1f%s, Max Temp: %.1f%s, Feels Like: %.1f%s, Humidity: %.1f%% Condition: %s\n", date,
				item.Main.Temp, unitSymbol,
				item.Main.TempMax, unitSymbol,
				item.Main.FeelsLike, unitSymbol,
				item.Main.Humidity,
				item.Weather[0].Description)
		}
	}
}
