package transit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/merlinfuchs/blimp/internal/config"
)

type DeparturesData struct {
	Boards []DeparturesBoardData `json:"boards"`
}

type DeparturesBoardData struct {
	Place struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Location struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lng"`
		} `json:"location"`
	} `json:"place"`
	Departures []struct {
		Time      time.Time `json:"time"`
		Platform  string    `json:"platform"`
		Transport struct {
			Mode      string `json:"mode"`
			Name      string `json:"name"`
			Category  string `json:"category"`
			Color     string `json:"color"`
			TextColor string `json:"textColor"`
			HeadSign  string `json:"headsign"`
			ShortName string `json:"shortName"`
		} `json:"transport"`
		Agency struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"agency"`
	} `json:"departures"`
}

func getNextDeparturesData() (DeparturesData, error) {
	url := fmt.Sprintf(
		"https://transit.hereapi.com/v8/departures?in=%f,%f;r=%d&apiKey=%s",
		config.K.Float64("widgets.transit.here_lat"),
		config.K.Float64("widgets.transit.here_lon"),
		config.K.Int("widgets.transit.here_radius"),
		config.K.String("widgets.transit.here_api_key"),
	)

	resp, err := http.Get(url)
	if err != nil {
		return DeparturesData{}, fmt.Errorf("failed to get next departures data: %w", err)
	}

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return DeparturesData{}, fmt.Errorf("Failed to get next departures aata, status code %d", resp.StatusCode)
	}

	var data DeparturesData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return DeparturesData{}, fmt.Errorf("Failed to decode next departures data, %w", err)
	}

	return data, nil
}
