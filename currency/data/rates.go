// This package retrieves exchange rate data from european central bank

package data

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
)

// ExchangeRates struct is used to hold decoded data from ecb api
type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

// Returns a new instance of ExhangeRates
func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: l, rates: map[string]float64{}}

	err := er.getRates()

	return er, err
}

// the ratio between two currencies = dest / base
func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency base %s", base)
	}

	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency dest %s", dest)
	}

	return dr / br, nil
}

// MonitorRates checks the rates in the ECB API every given interval and sends a messasge
// to the returned channel when there are changes

// Note that the ECB API returns data once a day.
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				// adding a random differenc to the rate to simulate the fluctuattion in currency rates
				for k, v := range e.rates {
					change := (rand.Float64() / 10)
					direction := rand.Intn(1)
					if direction == 0 {
						change = 1 - change
					} else {
						change = 1 + change
					}

					e.rates[k] = v * change
				}
				// send back the change through channel by handing an empty struct
				ret <- struct{}{}
			}
		}
	}()

	return ret
}

// Fetches the xml data from ecb api
func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200, but got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}

		e.rates[c.Currency] = r
	}

	e.rates["EUR"] = 1

	return nil
}

// Struct to decode the xml data returned from ecb api
type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

// Struct to decode each entries of the xml data
type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
