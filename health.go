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

// GetIndexHealth gets the index health from the server.
//
// This method sends a GET request to the server to retrieve the health status of the specified index.
//
// Parameters:
//
//	getIndexHealthReq (*GetIndexHealthRequest): The request containing the index name.
//
// Returns:
//
//	*GetIndexHealthResponse: The response containing the index health status.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the getIndexHealthReq parameter.
// 2. Sends a GET request to the server with the index name as a query parameter.
// 3. Checks the response status code and logs any errors.
// 4. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	getIndexHealthReq := &GetIndexHealthRequest{
//	    IndexName: "example_index",
//	}
//	resp, err := client.GetIndexHealth(getIndexHealthReq)
//	if err != nil {
//	    log.Fatalf("Failed to get index health: %v", err)
//	}
//	fmt.Printf("GetIndexHealthResponse: %+v\n", resp)
func (c *Client) GetIndexHealth(getIndexHealthReq *GetIndexHealthRequest) (*GetIndexHealthResponse, error) {
	logger := c.logger.With("method", "GetIndexHealth")
	err := validate.Struct(getIndexHealthReq)
	if err != nil {
		logger.Error("error validating get index health request", "error", err)
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
		logger.Error("error getting index health", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error getting index health: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response index health: %+v", getIndexHealthResp))
	return &getIndexHealthResp, nil
}
