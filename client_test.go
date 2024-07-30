package omgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestURLBuilder(t *testing.T) {
	c, err := NewClient()
	require.NoError(t, err)

	loc, err := NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	opts := ForecastOptions{
		TemperatureUnit:   "celsius",
		WindspeedUnit:     "kmh",
		PrecipitationUnit: "mm",
		Timezone:          "UTC",
		PastDays:          1,
		HourlyMetrics:     []string{"temperature_2m", "cloudcover", "direct_radiation", "diffuse_radiation", "precipitation", "windspeed_10m"},
		DailyMetrics:      []string{"temperature_2m_max"},
	}

	url := urlFromOptions(c.URL, loc, &opts)
	require.Equal(t, "https://api.open-meteo.com/v1/forecast?latitude=52.373800&longitude=4.891000&current_weather=true&temperature_unit=celsius&wind_speed_unit=kmh&precipitation_unit=mm&timezone=UTC&past_days=1&hourly=temperature_2m,cloudcover,direct_radiation,diffuse_radiation,precipitation,windspeed_10m&daily=temperature_2m_max", url)
}

func TestUrlFromHistoricalOptions(t *testing.T) {
	hc, err := NewClient()
	require.NoError(t, err)

	loc, err := NewLocation(52.5161, 13.4104) // Berlin
	require.NoError(t, err)

	hcOpts := HistoricalOptions{
		TemperatureUnit:   "celsius",
		WindspeedUnit:     "kmh",
		PrecipitationUnit: "mm",
		Timezone:          "UTC",
		StartDate:         "2023-07-25",
		EndDate:           "2023-08-08",
		HourlyMetrics:     []string{"temperature_2m", "rain"},
		DailyMetrics:      []string{"temperature_2m_max", "temperature_2m_mean"},
	}
	url, _ := urlFromHistoricalOptions(hc.HistoricalURL, loc, &hcOpts)
	require.Equal(t, "https://archive-api.open-meteo.com/v1/archive?latitude=52.516100&longitude=13.410400&start_date=2023-07-25&end_date=2023-08-08&temperature_unit=celsius&wind_speed_unit=kmh&precipitation_unit=mm&timezone=UTC&hourly=temperature_2m,rain&daily=temperature_2m_max,temperature_2m_mean", url)
}

func TestUrlFromHistoricalOptions2(t *testing.T) {
	hc, err := NewClient()
	require.NoError(t, err)

	loc, err := NewLocation(52.52, 13.41)
	require.NoError(t, err)

	hcOpts := HistoricalOptions{
		TemperatureUnit:   "fahrenheit",
		WindspeedUnit:     "mph",
		PrecipitationUnit: "inch",
		Timezone:          "America/Chicago",
		StartDate:         "2024-07-10",
		EndDate:           "2024-07-24",
		DailyMetrics:      []string{"weather_code", "temperature_2m_max", "temperature_2m_min", "temperature_2m_mean", "apparent_temperature_max", "apparent_temperature_min", "apparent_temperature_mean", "sunrise", "sunset", "daylight_duration", "sunshine_duration", "precipitation_sum", "rain_sum", "snowfall_sum", "precipitation_hours", "wind_speed_10m_max", "wind_gusts_10m_max", "wind_direction_10m_dominant", "shortwave_radiation_sum", "et0_fao_evapotranspiration"},
	}
	url, _ := urlFromHistoricalOptions(hc.HistoricalURL, loc, &hcOpts)
	require.Equal(t, "https://archive-api.open-meteo.com/v1/archive?latitude=52.520000&longitude=13.410000&start_date=2024-07-10&end_date=2024-07-24&temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch&timezone=America%2FChicago&daily=weather_code,temperature_2m_max,temperature_2m_min,temperature_2m_mean,apparent_temperature_max,apparent_temperature_min,apparent_temperature_mean,sunrise,sunset,daylight_duration,sunshine_duration,precipitation_sum,rain_sum,snowfall_sum,precipitation_hours,wind_speed_10m_max,wind_gusts_10m_max,wind_direction_10m_dominant,shortwave_radiation_sum,et0_fao_evapotranspiration", url)
}
