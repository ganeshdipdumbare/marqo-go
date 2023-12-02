package marqo

import (
	"fmt"
	"net/http"
)

// RefreshIndexRequest is the request to refresh the index
type RefreshIndexRequest struct {
	IndexName string `validate:"required" json:"-"`
}

// RefreshIndexResponse is the response from the server
type RefreshIndexResponse struct {
	Shards struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

// RefreshIndex refreshes the index on the server
func (c *Client) RefreshIndex(refreshIndexReq *RefreshIndexRequest) (*RefreshIndexResponse, error) {
	logger := c.logger.With("method", "RefreshIndex")
	err := validate.Struct(refreshIndexReq)
	if err != nil {
		logger.Error("error validating refresh index request",
			"error", err)
		return nil, err
	}

	var refreshIndexResp RefreshIndexResponse
	resp, err := c.reqClient.
		R().
		SetSuccessResult(&refreshIndexResp).
		Post(c.reqClient.BaseURL + "/indexes/" + refreshIndexReq.IndexName + "/refresh")
	if err != nil {
		logger.Error("error refreshing index", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error refreshing index", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error refreshing index: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info("index refreshed")
	return &refreshIndexResp, nil
}
