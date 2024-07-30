package omgo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	URL           string
	HistoricalURL string
	UserAgent     string
	Client        *http.Client
}

const DefaultUserAgent = "Open-Meteo_Go_Client"

// Returns a new
func NewClient() (Client, error) {
	return Client{
		URL:           "https://api.open-meteo.com/v1/forecast",
		HistoricalURL: "https://archive-api.open-meteo.com/v1/archive",
		UserAgent:     DefaultUserAgent,
		Client:        http.DefaultClient,
	}, nil
}

type Location struct {
	lat, lon float64
}

// Converts the provided latitude and longitude floats into a Location struct
func NewLocation(lat, lon float64) (Location, error) {
	return Location{lat: lat, lon: lon}, nil
}

type ForecastOptions struct {
	TemperatureUnit   string   // Default "celsius"
	WindspeedUnit     string   // Default "kmh",
	PrecipitationUnit string   // Default "mm"
	Timezone          string   // Default "UTC"
	PastDays          int      // Default 0
	HourlyMetrics     []string // Lists required hourly metrics, see https://open-meteo.com/en/docs for valid metrics
	DailyMetrics      []string // Lists required daily metrics, see https://open-meteo.com/en/docs for valid metrics
}

type HistoricalOptions struct {
	TemperatureUnit   string   // Default "celsius"
	WindspeedUnit     string   // Default "kmh",
	PrecipitationUnit string   // Default "mm"
	Timezone          string   // Default "UTC"
	StartDate         string   // Format yyyy-mm-dd ISO8601 date
	EndDate           string   // Format yyyy-mm-dd ISO8601 date
	HourlyMetrics     []string // Lists required hourly metrics, see https://open-meteo.com/en/docs for valid metrics
	DailyMetrics      []string // Lists required daily metrics, see https://open-meteo.com/en/docs for valid metrics
}

// Returns the Forecast URL used in the web request
func urlFromOptions(baseURL string, loc Location, opts *ForecastOptions) string {
	// TODO: Validate the Options
	newurl := fmt.Sprintf(`%s?latitude=%f&longitude=%f&current_weather=true`, baseURL, loc.lat, loc.lon)
	if opts == nil {
		return newurl
	}

	if opts.TemperatureUnit != "" {
		newurl = fmt.Sprintf(`%s&temperature_unit=%s`, newurl, opts.TemperatureUnit)
	}
	if opts.WindspeedUnit != "" {
		newurl = fmt.Sprintf(`%s&wind_speed_unit=%s`, newurl, opts.WindspeedUnit)
	}
	if opts.PrecipitationUnit != "" {
		newurl = fmt.Sprintf(`%s&precipitation_unit=%s`, newurl, opts.PrecipitationUnit)
	}
	if opts.Timezone != "" {
		encodedTimezone := url.QueryEscape(opts.Timezone)
		newurl = fmt.Sprintf(`%s&timezone=%s`, newurl, encodedTimezone)
	}

	newurl = fmt.Sprintf("%s&past_days=%d", newurl, opts.PastDays)
	if opts.HourlyMetrics != nil && len(opts.HourlyMetrics) > 0 {
		metrics := strings.Join(opts.HourlyMetrics, ",")
		newurl = fmt.Sprintf(`%s&hourly=%s`, newurl, metrics)
	}

	if opts.DailyMetrics != nil && len(opts.DailyMetrics) > 0 {
		metrics := strings.Join(opts.DailyMetrics, ",")
		newurl = fmt.Sprintf(`%s&daily=%s`, newurl, metrics)
	}

	return newurl
}

// Returns the Historical URL used in the web request
func urlFromHistoricalOptions(baseURL string, loc Location, opts *HistoricalOptions) (string, error) {
	newurl := fmt.Sprintf(`%s?latitude=%f&longitude=%f`, baseURL, loc.lat, loc.lon)

	if opts.StartDate == "" || opts.EndDate == "" {
		return "", fmt.Errorf("please provide a start date and end date for historical option")
	}

	if opts == nil {
		return newurl, nil
	}

	// TODO: if end_date is equal or less than end_date return an error
	newurl = fmt.Sprintf(`%s&start_date=%s&end_date=%s`, newurl, opts.StartDate, opts.EndDate)

	if opts.TemperatureUnit != "" {
		newurl = fmt.Sprintf(`%s&temperature_unit=%s`, newurl, opts.TemperatureUnit)
	}
	if opts.WindspeedUnit != "" {
		newurl = fmt.Sprintf(`%s&wind_speed_unit=%s`, newurl, opts.WindspeedUnit)
	}
	if opts.PrecipitationUnit != "" {
		newurl = fmt.Sprintf(`%s&precipitation_unit=%s`, newurl, opts.PrecipitationUnit)
	}
	if opts.Timezone != "" {
		encodedTimezone := url.QueryEscape(opts.Timezone)
		newurl = fmt.Sprintf(`%s&timezone=%s`, newurl, encodedTimezone)
	}

	if opts.HourlyMetrics != nil && len(opts.HourlyMetrics) > 0 {
		metrics := strings.Join(opts.HourlyMetrics, ",")
		newurl = fmt.Sprintf(`%s&hourly=%s`, newurl, metrics)
	}

	if opts.DailyMetrics != nil && len(opts.DailyMetrics) > 0 {
		metrics := strings.Join(opts.DailyMetrics, ",")
		newurl = fmt.Sprintf(`%s&daily=%s`, newurl, metrics)
	}

	return newurl, nil
}

func (c Client) GetForecast(ctx context.Context, loc Location, opts *ForecastOptions) ([]byte, error) {
	newurl := urlFromOptions(c.URL, loc, opts)
	req, err := http.NewRequestWithContext(ctx, "GET", newurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%s - %s", res.Status, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) GetHistorical(ctx context.Context, loc Location, opts *HistoricalOptions) ([]byte, error) {
	newurl, err := urlFromHistoricalOptions(c.HistoricalURL, loc, opts)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to form url from historical options %q", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", newurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("error status %s, response %s", res.Status, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
