package marqo

import (
	"fmt"
	"net/http"
)

// GetIndexSettingsRequest is the request to get the index settings
type GetIndexSettingsRequest struct {
	IndexName string `validate:"required" json:"-"`
}

// GetIndexSettingsResponse is the response from the server
type GetIndexSettingsResponse struct {
	IndexDefaults *IndexDefaults `json:"index_defaults"`
}

// GetIndexSettings gets the index settings from the server
func (c *Client) GetIndexSettings(getIndexSettingsReq *GetIndexSettingsRequest) (*GetIndexSettingsResponse, error) {
	logger := c.logger.With("method", "GetIndexSettings")
	err := validate.Struct(getIndexSettingsReq)
	if err != nil {
		logger.Error("error validating get index settings request",
			"error", err)
		return nil, err
	}

	var getIndexSettingsResp GetIndexSettingsResponse
	resp, err := c.reqClient.
		R().
		SetSuccessResult(&getIndexSettingsResp).
		Get(c.reqClient.BaseURL + "/indexes/" + getIndexSettingsReq.IndexName + "/settings")

	if err != nil {
		logger.Error("error getting index settings", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error getting index settings", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error getting index settings: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info("index settings retrieved")
	return &getIndexSettingsResp, nil
}
