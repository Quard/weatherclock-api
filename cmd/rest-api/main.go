package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/Quard/weatherclock-api/internal/apigateway"
)

func weatherRestAPIHandler(w http.ResponseWriter, r *http.Request) {
	var requestData apigateway.WeatherRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	reports, err := apigateway.HandleRequest(ctx, requestData)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)
}

func main() {
	bindPtr := flag.String("bind", ":8000", "host:port for listen")
	flag.Parse()

	http.HandleFunc("/api/v1/weather/", weatherRestAPIHandler)

	log.Printf("listen: %s", *bindPtr)
	log.Fatal(http.ListenAndServe(*bindPtr, nil))
}
