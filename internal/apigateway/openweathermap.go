package apigateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"sync"
	"time"
)

const openWeatherURI = "https://api.openweathermap.org/data/2.5/"
const openWeatherTodayUrlTmpl = "%sweather?APPID=%s&units=%s&lat=%f&lon=%f"
const openWeatherForecastUrlTmpl = "%sforecast?APPID=%s&units=%s&lat=%f&lon=%f"

type blockMain struct {
	Temp    float64 `json:"temp"`
	TempMax float64 `json:"temp_max"`
	TempMin float64 `json:"temp_min"`
}

type blockWeather struct {
	Icon string `json:"icon"`
}

type OpenWeatherMapTodayResponse struct {
	Date    int64          `json:"dt"`
	Main    blockMain      `json:"main"`
	Weather []blockWeather `json:"weather"`
}

type OpenWeatherMapForecastResponse struct {
	List []OpenWeatherMapTodayResponse `json:"list"`
}

func OpenWeatherMap(ctx context.Context, r WeatherRequest) ([]WeatherReport, error) {
	reports := make(chan WeatherReport, 10)
	done := make(chan bool, 1)

	go func() {
		var wg sync.WaitGroup
		wg.Add(2)
		go openWeatherMapGetToday(reports, r, &wg)
		go openWeatherMapGetForecast(reports, r, &wg)

		wg.Wait()
		close(reports)
		done <- true
	}()

	weatherReports := []WeatherReport{}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		for report := range reports {
			weatherReports = append(weatherReports, report)
		}
		sort.Sort(WeatherReportByDate(weatherReports))
	}

	return weatherReports, nil
}

func openWeatherMapGetToday(ch chan WeatherReport, r WeatherRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf(openWeatherTodayUrlTmpl, openWeatherURI, r.APIKey, r.Units, r.Latitude, r.Longitude)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	var weatherRec OpenWeatherMapTodayResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherRec); err != nil {
		log.Println(err)
		return
	}

	ch <- WeatherReport{
		Date:       time.Unix(weatherRec.Date, 0).Format(DateFormat),
		TempMax:    weatherRec.Main.TempMax,
		TempMin:    weatherRec.Main.TempMin,
		TempDayAvg: weatherRec.Main.Temp,
		Icon:       openWeatherMapNormalizeIcon(weatherRec.Weather[0].Icon),
	}
}

func openWeatherMapGetForecast(ch chan WeatherReport, r WeatherRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf(openWeatherForecastUrlTmpl, openWeatherURI, r.APIKey, r.Units, r.Latitude, r.Longitude)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	var data OpenWeatherMapForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
		return
	}

	report := WeatherReport{}
	var tempSum float64
	var tempCount float64
	today := time.Now().Format(DateFormat)
	for _, weatherRec := range data.List {
		datetime := time.Unix(weatherRec.Date, 0)
		date := datetime.Format(DateFormat)

		if date == today {
			continue
		}

		if date != report.Date {
			if report.Date != "" {
				report.TempDayAvg = tempSum / tempCount

				ch <- report
			}
			report = WeatherReport{
				Date:    date,
				TempMax: weatherRec.Main.TempMax,
				TempMin: weatherRec.Main.TempMin,
			}
			tempSum = tempSum + weatherRec.Main.Temp
			tempCount++
		} else {
			report.TempMax = math.Max(report.TempMax, weatherRec.Main.TempMax)
			report.TempMin = math.Max(report.TempMin, weatherRec.Main.TempMin)
			tempSum = tempSum + weatherRec.Main.Temp
			tempCount++
		}
		if datetime.Hour() > 12 && datetime.Hour() < 16 {
			report.Icon = openWeatherMapNormalizeIcon(weatherRec.Weather[0].Icon)
		}

	}
}

func openWeatherMapNormalizeIcon(icon string) string {
	switch icon[:2] {
	case "01":
		return "clear-sky"
	case "02":
		return "few-clouds"
	case "03":
		return "scattered-clouds"
	case "04":
		return "broken-clouds"
	case "09":
		return "shower-rain"
	case "10":
		return "rain"
	case "11":
		return "thunderstorm"
	case "13":
		return "snow"
	case "50":
		return "mist"
	default:
		log.Printf("unknown weather icon '%s' from OpenWeatherMap")
	}

	return ""
}
