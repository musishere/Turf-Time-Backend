package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// LocationIQResponse matches the LocationIQ API JSON response.
type LocationIQResponse struct {
	PlaceID     string   `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmID       string   `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Government  string   `json:"government,omitempty"`
	HouseNumber string   `json:"house_number,omitempty"`
	Road        string   `json:"road,omitempty"`
	Quarter     string   `json:"quarter,omitempty"`
	Suburb      string   `json:"suburb,omitempty"`
	City        string   `json:"city,omitempty"`
	State       string   `json:"state,omitempty"`
	Postcode    string   `json:"postcode,omitempty"`
	Country     string   `json:"country,omitempty"`
	CountryCode string   `json:"country_code,omitempty"`
	Boundingbox []string `json:"boundingbox,omitempty"`
}

// GeocodeAddress calls LocationIQ Geocoding API and returns the full location response for the given address.
func GeocodeAddress(address string) (*LocationIQResponse, error) {
	apiKey := os.Getenv("LOCATION_IQ_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("LOCATION_IQ_KEY environment variable is not set")
	}

	u := fmt.Sprintf(
		"https://us1.locationiq.com/v1/search?key=%s&q=%s&format=json",
		apiKey,
		url.QueryEscape(address),
	)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geocoding request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status %d: check if LOCATION_IQ_KEY is valid", res.StatusCode)
	}

	// LocationIQ returns an array, so we need to parse it as such
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var results []LocationIQResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("geocoding response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	fmt.Printf("Geocoding result for '%s': lat=%s, lon=%s\n", address, results[0].Lat, results[0].Lon)

	return &results[0], nil
}
