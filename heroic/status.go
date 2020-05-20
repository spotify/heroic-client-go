package heroic

import "context"

func (c *Client) Status(ctx context.Context) (*StatusResponse, error) {
	req, err := c.NewRequest("GET", "status", nil)
	if err != nil {
		return nil, err
	}
	var sr StatusResponse
	_, err = c.Do(ctx, req, &sr)
	return &sr, err
}

type StatusResponse struct {
	Ok               bool            `json:"ok"`
	Service          ServiceInfo     `json:"service"`
	Consumers        Consumer        `json:"consumers"`
	Backends         Backend         `json:"backends"`
	MetadataBackends MetadataBackend `json:"metadataBackends"`
	Cluster          Cluster         `json:"cluster"`
}

type ServiceInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Id      string `json:"id"`
}

type Consumer struct {
	*Status
	Errors           int64 `json:"errors"`
	ConsumingThreads int64 `json:"consumingThreads"`
	TotalThreads     int64 `json:"totalThreads"`
}

type Backend struct {
	*Status
}

type MetadataBackend struct {
	*Status
}

type Cluster struct {
	*Status
}

type Status struct {
	Ok        bool `json:"ok"`
	Available int  `json:"available"`
	Ready     int  `json:"ready"`
}
