package marqo

import (
	"fmt"
	"net/http"
)

// GetIndexStatsRequest is the request to get the index stats
type GetIndexStatsRequest struct {
	IndexName string `validate:"required" json:"-"`
}

// GetIndexStatsResponse is the response from the server
type GetIndexStatsResponse struct {
	NumberOfDocuments int `json:"numberOfDocuments"`
	NumberOfVectors   int `json:"numberOfVectors"`
}

// GetIndexStats gets the index stats from the server
func (c *Client) GetIndexStats(getIndexStatsReq *GetIndexStatsRequest) (*GetIndexStatsResponse, error) {
	logger := c.logger.With("method", "GetIndexStats")
	err := validate.Struct(getIndexStatsReq)
	if err != nil {
		logger.Error("error validating get index stats request",
			"error", err)
		return nil, err
	}

	var getIndexStatsResp GetIndexStatsResponse
	resp, err := c.reqClient.
		R().
		SetSuccessResult(&getIndexStatsResp).
		Get(c.reqClient.BaseURL + "/indexes/" + getIndexStatsReq.IndexName + "/stats")

	if err != nil {
		logger.Error("error getting index stats", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error getting index stats", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error getting index stats: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info("index stats retrieved")
	return &getIndexStatsResp, nil
}
