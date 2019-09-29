package apigateway

import (
	"context"
	"errors"
	"time"
)

const DateFormat = "2006-01-02"

type WeatherRequest struct {
	Provider  string  `json:"provider"`
	APIKey    string  `json:"api_key"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Units     string  `json:"units"`
}

type WeatherReport struct {
	Date       string  `json:"date"`
	TempMax    float64 `json:"temp_max"`
	TempMin    float64 `json:"temp_min"`
	TempDayAvg float64 `json:"temp_day_avg"`
	Icon       string  `json:"icon"`
}

type WeatherReportByDate []WeatherReport

func HandleRequest(ctx context.Context, r WeatherRequest) ([]WeatherReport, error) {
	switch r.Provider {
	case "openweathermap":
		return OpenWeatherMap(ctx, r)
	default:
		return nil, errors.New("not supported weather provider")
	}
}

func (r WeatherReportByDate) Len() int {
	return len(r)
}

func (r WeatherReportByDate) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r WeatherReportByDate) Less(i, j int) bool {
	date1, _ := time.Parse(DateFormat, r[i].Date)
	date2, _ := time.Parse(DateFormat, r[j].Date)
	return date1.Before(date2)
}
