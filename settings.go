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

// GetIndexSettings gets the index settings from the server.
//
// This method sends a GET request to the server to retrieve the settings of the specified index.
//
// Parameters:
//
//	getIndexSettingsReq (*GetIndexSettingsRequest): The request containing the index name.
//
// Returns:
//
//	*GetIndexSettingsResponse: The response containing the index settings.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the getIndexSettingsReq parameter.
// 2. Sends a GET request to the server with the index name as a query parameter.
// 3. Checks the response status code and logs any errors.
// 4. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	getIndexSettingsReq := &GetIndexSettingsRequest{
//	    IndexName: "example_index",
//	}
//	resp, err := client.GetIndexSettings(getIndexSettingsReq)
//	if err != nil {
//	    log.Fatalf("Failed to get index settings: %v", err)
//	}
//	fmt.Printf("GetIndexSettingsResponse: %+v\n", resp)
func (c *Client) GetIndexSettings(getIndexSettingsReq *GetIndexSettingsRequest) (*GetIndexSettingsResponse, error) {
	logger := c.logger.With("method", "GetIndexSettings")
	err := validate.Struct(getIndexSettingsReq)
	if err != nil {
		logger.Error("error validating get index settings request", "error", err)
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
		logger.Error("error getting index settings", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error getting index settings: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response index settings: %+v", getIndexSettingsResp))
	return &getIndexSettingsResp, nil
}
