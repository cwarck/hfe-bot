package currency

import (
	"encoding/json"
	"net/http"
)

// OxrResponse is the response from the Open Exchange Rates API.
type OxrResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

// GetRates gets the exchange rates from the Open Exchange Rates API.
func GetRates(appId string) (*OxrResponse, error) {
	var body OxrResponse

	resp, err := http.Get("https://openexchangerates.org/api/latest.json?app_id=" + appId)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}
