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

// GetIndexStats gets the index stats from the server.
//
// This method sends a GET request to the server to retrieve the statistics of the specified index.
//
// Parameters:
//
//	getIndexStatsReq (*GetIndexStatsRequest): The request containing the index name.
//
// Returns:
//
//	*GetIndexStatsResponse: The response containing the index statistics.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the getIndexStatsReq parameter.
// 2. Sends a GET request to the server with the index name as a query parameter.
// 3. Checks the response status code and logs any errors.
// 4. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	getIndexStatsReq := &GetIndexStatsRequest{
//	    IndexName: "example_index",
//	}
//	resp, err := client.GetIndexStats(getIndexStatsReq)
//	if err != nil {
//	    log.Fatalf("Failed to get index stats: %v", err)
//	}
//	fmt.Printf("GetIndexStatsResponse: %+v\n", resp)
func (c *Client) GetIndexStats(getIndexStatsReq *GetIndexStatsRequest) (*GetIndexStatsResponse, error) {
	logger := c.logger.With("method", "GetIndexStats")
	err := validate.Struct(getIndexStatsReq)
	if err != nil {
		logger.Error("error validating get index stats request", "error", err)
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
		logger.Error("error getting index stats", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error getting index stats: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response index stats: %+v", getIndexStatsResp))
	return &getIndexStatsResp, nil
}
