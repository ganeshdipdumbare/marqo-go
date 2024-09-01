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
//
// This method sends a POST request to the server to refresh the specified index.
//
// Parameters:
//
//	refreshIndexReq (*RefreshIndexRequest): The request containing the index name.
//
// Returns:
//
//	*RefreshIndexResponse: The response from the server.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the refreshIndexReq parameter.
// 2. Sends a POST request to the server with the index name as a query parameter.
// 3. Checks the response status code and logs any errors.
// 4. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	refreshIndexReq := &RefreshIndexRequest{
//	    IndexName: "example_index",
//	}
//	resp, err := client.RefreshIndex(refreshIndexReq)
//	if err != nil {
//	    log.Fatalf("Failed to refresh index: %v", err)
//	}
//	fmt.Printf("RefreshIndexResponse: %+v\n", resp)
func (c *Client) RefreshIndex(refreshIndexReq *RefreshIndexRequest) (*RefreshIndexResponse, error) {
	logger := c.logger.With("method", "RefreshIndex")
	err := validate.Struct(refreshIndexReq)
	if err != nil {
		logger.Error("error validating refresh index request", "error", err)
		return nil, err
	}

	var refreshIndexResp RefreshIndexResponse
	resp, err := c.reqClient.
		R().
		SetSuccessResult(&refreshIndexResp).
		Post(c.reqClient.BaseURL + "/indexes/" + refreshIndexReq.IndexName + "/_refresh")
	if err != nil {
		logger.Error("error refreshing index", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error refreshing index", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error refreshing index: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response refresh index: %+v", refreshIndexResp))
	return &refreshIndexResp, nil
}
