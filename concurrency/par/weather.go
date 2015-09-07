package main

import (
	"code.google.com/p/go.net/context"
	"fmt"
	"github.com/savaki/openweathermap"
	"github.com/savaki/par"
	// "golang.org/x/net/context"
	"log"
	"time"
)

func ok(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func findById(cityId int, responses chan *openweathermap.Forecast) par.RequestFunc {
	return func(ctx context.Context) error {
		fmt.Printf("cityId=%v\n", cityId)
		forecast, err := openweathermap.New().ByCityId(cityId)
		fmt.Printf("forecast=%v\n", forecast)
		ok(err)
		responses <- forecast
		return nil
	}
}

func find(city string, responses chan *openweathermap.Forecast) par.RequestFunc {
	return func(ctx context.Context) error {
		fmt.Printf("city=%v\n", city)
		forecast, err := openweathermap.New().ByCityName(city)
		fmt.Printf("forecast=%v\n", forecast)
		ok(err)
		responses <- forecast
		return nil
	}
}

func main() {
	const numRequests = 3
	// create a channel to capture our results
	forecasts := make(chan *openweathermap.Forecast, numRequests)

	// create our channel of requests
	requests := make(chan par.RequestFunc, numRequests)
	requests <- findById(4288809, forecasts) // Covington, VA
	requests <- findById(4288809, forecasts)
	requests <- findById(4140963, forecasts) // DC
	close(requests)                          // important to remember to close the channel

	// resolver := par.Requests(requests).WithRedundancy(1)
	// resolver := par.Requests(requests).WithConcurrency(numRequests)
	resolver := par.Requests(requests).WithConcurrency(0)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	err := resolver.DoWithContext(ctx)
	cancel()
	ok(err)

	// the forecasts channel now contains all our forecasts
	close(forecasts)
	cities := map[string]*openweathermap.Forecast{}
	for forecast := range forecasts {
		cities[forecast.Name] = forecast
	}
}
