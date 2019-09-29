package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/Quard/weatherclock-api/internal/apigateway"
)

func main() {
	providerPtr := flag.String("provider", "openweathermap", "weather api provider: openweathermap")
	apiKeyPtr := flag.String("api-key", "", "API key for provider")
	unitsPtr := flag.String("units", "metric", "display units")
	latitudePtr := flag.Float64("latitude", 0, "latitude")
	longitudePtr := flag.Float64("longitude", 0, "longitude")
	timeoutPtr := flag.String("timeout", "5s", "API fetch timeout")

	flag.Parse()

	timeout, errDurationParse := time.ParseDuration(*timeoutPtr)
	if errDurationParse != nil {
		panic(errDurationParse)
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	reqData := apigateway.WeatherRequest{
		Provider:  *providerPtr,
		APIKey:    *apiKeyPtr,
		Units:     *unitsPtr,
		Latitude:  *latitudePtr,
		Longitude: *longitudePtr,
	}
	reports, err := apigateway.HandleRequest(ctx, reqData)
	if err != nil {
		panic(err)
	}

	for _, report := range reports {
		fmt.Printf(
			"%s\t%0.2f (%0.2f - %0.2f)\t%s\n",
			report.Date,
			report.TempDayAvg,
			report.TempMin,
			report.TempMax,
			report.Icon,
		)
	}
}
