package marqo

import (
	"fmt"
	"net/http"
)

// BulkSearchRequest is the request to bulk search
type BulkSearchRequest struct {
	// Body params
	// Queries is the list of search requests
	Queries []SearchRequest `validate:"required" json:"queries"`

	// Query params
	// Device is the device to run the search on
	Device *string `json:"device,omitempty"`
	// Telemetry is whether to send telemetry
	Telemetry *bool `json:"telemetry,omitempty"`
}

// BulkSearchResponse is the response from the server
type BulkSearchResponse struct {
	// Result is the list of search responses
	Result           []SearchResponse `json:"result"`
	ProcessingTimeMS float64          `json:"processingTimeMs"`
}

// BulkSearch performs a bulk search on the server.
//
// This method sends a POST request to the server to perform a bulk search with the specified queries.
//
// Parameters:
//
//	bulkSearchReq (*BulkSearchRequest): The request containing the search queries.
//
// Returns:
//
//	*BulkSearchResponse: The response from the server.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the bulkSearchReq parameter.
// 2. Sends a POST request to the server with the search queries in the request body.
// 3. Checks the response status code and logs any errors.
// 4. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	bulkSearchReq := &BulkSearchRequest{
//	    Queries: []SearchRequest{...},
//	}
//	resp, err := client.BulkSearch(bulkSearchReq)
//	if err != nil {
//	    log.Fatalf("Failed to perform bulk search: %v", err)
//	}
//	fmt.Printf("BulkSearchResponse: %+v\n", resp)
func (c *Client) BulkSearch(bulkSearchReq *BulkSearchRequest) (*BulkSearchResponse, error) {
	logger := c.logger.With("method", "BulkSearch")
	for i := range bulkSearchReq.Queries {
		setDefaultSearchRequest(&bulkSearchReq.Queries[i])
	}
	err := validate.Struct(bulkSearchReq)
	if err != nil {
		logger.Error("error validating bulk search request", "error", err)
		return nil, err
	}

	var bulkSearchResp BulkSearchResponse
	resp, err := c.reqClient.
		R().
		SetSuccessResult(&bulkSearchResp).
		Post(c.reqClient.BaseURL + "/indexes/_bulk_search")
	if err != nil {
		logger.Error("error bulk searching", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error bulk searching", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error bulk searching: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response bulk search: %+v", bulkSearchResp))
	return &bulkSearchResp, nil
}
