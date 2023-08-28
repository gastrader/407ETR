package client

import (
	"bytes"
	"context"
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

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	return &types.Invoice{
		OBUID: 123123,
		TotalDistance: 69,
		TotalAmount: 420,
	}, nil
}


//should be aggregate distance.
func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	
	b, err := json.Marshal(aggReq)
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