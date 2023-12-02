package marqo

import (
	"fmt"
	"net/http"
)

// GetIndexHealthRequest is the request to get the index health
type GetIndexHealthRequest struct {
	IndexName string `validate:"required" json:"-"`
}

// GetIndexHealthResponse is the response from the server
type GetIndexHealthResponse struct {
	Status  string `json:"status"`
	Backend struct {
		Status             string `json:"status"`
		StorageIsAvailable bool   `json:"storage_is_available"`
	} `json:"backend"`
}

// GetIndexHealth gets the index health from the server
func (c *Client) GetIndexHealth(getIndexHealthReq *GetIndexHealthRequest) (*GetIndexHealthResponse, error) {
	logger := c.logger.With("method", "GetIndexHealth")
	err := validate.Struct(getIndexHealthReq)
	if err != nil {
		logger.Error("error validating get index health request",
			"error", err)
		return nil, err
	}

	var getIndexHealthResp GetIndexHealthResponse
	resp, err := c.reqClient.
		R().
		SetSuccessResult(&getIndexHealthResp).
		Get(c.reqClient.BaseURL + "/indexes/" + getIndexHealthReq.IndexName + "/health")

	if err != nil {
		logger.Error("error getting index health", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error getting index health", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error getting index health: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info("index health retrieved")
	return &getIndexHealthResp, nil
}
