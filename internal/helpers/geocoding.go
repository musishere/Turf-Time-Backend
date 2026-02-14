package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// GeocodeResponse matches the Google Geocoding API JSON response (subset we need).
type GeocodeResponse struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status  string `json:"status"`
	Message string `json:"error_message,omitempty"`
}

// GeocodeAddress calls Google Maps Geocoding API and returns latitude and longitude for the given address.
// Requires GOOGLE_MAPS_API_KEY in the environment.
func GeocodeAddress(address string) (latitude, longitude float64, err error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return 0, 0, fmt.Errorf("GOOGLE_MAPS_API_KEY is not set")
	}

	u := "https://maps.googleapis.com/maps/api/geocode/json?address=" + url.QueryEscape(address) + "&key=" + url.QueryEscape(apiKey)
	resp, err := http.Get(u)
	if err != nil {
		return 0, 0, fmt.Errorf("geocoding request: %w", err)
	}
	defer resp.Body.Close()

	var data GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, fmt.Errorf("geocoding response: %w", err)
	}

	if data.Status != "OK" {
		msg := data.Message
		if msg == "" {
			msg = data.Status
		}
		return 0, 0, fmt.Errorf("geocoding API: %s", msg)
	}
	if len(data.Results) == 0 {
		return 0, 0, fmt.Errorf("no results for address: %s", address)
	}

	loc := data.Results[0].Geometry.Location
	return loc.Lat, loc.Lng, nil
}
