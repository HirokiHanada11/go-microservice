package server

import (
	"context"
	"io"
	"time"

	"github.com/HirokiHanada11/go-microservices/currency/data"
	protos "github.com/HirokiHanada11/go-microservices/currency/protos"
	"github.com/hashicorp/go-hclog"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
	protos.UnimplementedCurrencyServer
}

// NewCurrency creates a new Currency server
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{r, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest), protos.UnimplementedCurrencyServer{}}
	go c.handleUpdates()

	return c
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got Update rates")

		// loop over subscribed clients
		for k, v := range c.subscriptions {

			// loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}

				err = k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r})
				if err != nil {
					c.log.Error("Unable to get send rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}
		}
	}

}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

func (c *Currency) SubscribeRates(stream protos.Currency_SubscribeRatesServer) error {
	// Listen for client input in a goroutine
	for {
		rr, err := stream.Recv()
		if err == io.EOF { // EOF in the stdi indicates that the client has close connection
			c.log.Info("Client has closed connection")
			break
		}
		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			break
		}

		c.log.Info("Handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())

		rrs, ok := c.subscriptions[stream]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		rrs = append(rrs, rr)
		c.subscriptions[stream] = rrs
	}

	return nil
}
