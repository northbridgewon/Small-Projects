package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Weather struct {
	CurrentCondition []struct {
		FeelsLikeC  string `json:"FeelsLikeC"`
		FeelsLikeF  string `json:"FeelsLikeF"`
		TempC       string `json:"temp_C"`
		TempF       string `json:"temp_F"`
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
	} `json:"current_condition"`
}

func main() {
	fahrenheit := flag.Bool("f", false, "Display temperature in Fahrenheit")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: go run main.go [-f] <city>")
		return
	}
	city := flag.Arg(0)

	resp, err := http.Get(fmt.Sprintf("https://wttr.in/%s?format=j1", city))
	if err != nil {
		fmt.Printf("Error fetching weather: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	var weather Weather
	if err := json.Unmarshal(body, &weather); err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return
	}

	if len(weather.CurrentCondition) > 0 {
		current := weather.CurrentCondition[0]
		fmt.Printf("Weather in %s:\n", city)

		if *fahrenheit {
			fmt.Printf("  Temperature: %s°F (Feels like %s°F)\n", current.TempF, current.FeelsLikeF)
		} else {
			fmt.Printf("  Temperature: %s°C (Feels like %s°C)\n", current.TempC, current.FeelsLikeC)
		}

		if len(current.WeatherDesc) > 0 {
			fmt.Printf("  Description: %s\n", current.WeatherDesc[0].Value)
		}
	} else {
		fmt.Println("Could not retrieve weather information.")
	}
}
