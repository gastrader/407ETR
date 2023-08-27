package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gastrader/407ETR/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}
//should be aggregate distance.
func (c *HTTPClient) AggregateInvoice(distance types.Distance) error {
	
	b, err := json.Marshal(distance)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("srevice responded with %d instead of 200 status code", resp.StatusCode)	
	}
	return nil
}