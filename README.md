# Open-Meteo-Go

A simple go client for the [open meteo](https://open-meteo.com) API. It supports all options of the API as of Sept 20 2021.

This fork will provide and implementation for Historical API, so client can get the past data of Open Meteo.

The implementation will separate one client for the Forecast and another for Historical.

## Usage

Simple example of how it'll look like:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lum8rjack/omgo"
)

func main() {
	f, err := omgo.NewClient()
	if err != nil {
		fmt.Printf("error creating client: %v", err)
		os.Exit(0)
	}

	// Get the current weather for Chicago
	loc, err := omgo.NewLocation(41.8781, 87.6298)
	if err != nil {
		fmt.Printf("error creating location: %v", err)
		os.Exit(0)
	}
	res, err := f.CurrentWeather(context.Background(), loc, nil)
	if err != nil {
		fmt.Printf("error getting current weather: %v", err)
		os.Exit(0)
	}
	fmt.Printf("The temperature in Chicago is: %.2f degrees C\n", res.Temperature)

	// Get the humidity and cloud cover forecast for Chicago,
	// including the last 2 days and non-metric units
	opts := omgo.ForecastOptions{
		TemperatureUnit:   "fahrenheit",
		WindspeedUnit:     "mph",
		PrecipitationUnit: "inch",
		Timezone:          "US/Central",
		HourlyMetrics:     []string{"cloudcover"},
		DailyMetrics:      []string{"temperature_2m_max"},
	}

	fres, err := f.Forecast(context.Background(), loc, &opts)
	if err != nil {
		fmt.Printf("error getting forecast: %v", err)
		os.Exit(0)
	}
	fmt.Printf("Current Temperature: %.2f degrees F\n", fres.CurrentWeather.Temperature)

	// Loop over the cloud cover every hour for 7 days
	for x, h := range fres.HourlyMetrics["cloudcover"] {
		fmt.Printf("%d - Cloud cover: %.2f\n", x, h)
	}
	// fres.HourlyTimes contains the timestamps for each prediction
	// fres.DailyMetrics["temperature_2m_max"] contains daily maximum values for the temperature_2m metric
	// fres.DailyTimes contains the timestamps for all daily predictions

	hopts := omgo.HistoricalOptions{
		TemperatureUnit:   "fahrenheit",
		WindspeedUnit:     "mph",
		PrecipitationUnit: "inch",
		Timezone:          "US/Central",
		StartDate:         "2024-06-01",
		EndDate:           "2024-06-10",
		DailyMetrics:      []string{"temperature_2m_max", "temperature_2m_mean"},
	}

	hres, err := f.Historical(context.Background(), loc, &hopts)
	if err != nil {
		fmt.Printf("error getting historical data: %v", err)
		os.Exit(0)
	}

	// Convert the start date from string to time.Time
	startDate, err := time.Parse("2006-01-02", hopts.StartDate)
	if err != nil {
		fmt.Printf("error converting start date: %v", err)
		os.Exit(0)
	}

	// Loop over each day
	for x, d := range hres.DailyMetrics["temperature_2m_max"] {
		newdate := startDate.Add(time.Hour * time.Duration(24*x))
		fmt.Printf("%s - Mean Temp: %.2f\n", newdate.Format("2006-01-02"), d)
	}
}
```
