package marqo

import (
	"fmt"
	"net/http"
	"strconv"
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
	// Results is the list of search responses
	Results          []SearchResponse `json:"results"`
	ProcessingTimeMS float64          `json:"processingTimeMs"`
}

// BulkSearch searches the index
func (c *Client) BulkSearch(bulkSearchReq *BulkSearchRequest) (*BulkSearchResponse, error) {
	logger := c.logger.With("method", "BulkSearch")
	err := validate.Struct(bulkSearchReq)
	if err != nil {
		logger.Error("error validating bulk search request",
			"error", err)
		return nil, err
	}

	var bulkSearchResp BulkSearchResponse
	queryParams := map[string]string{}
	if bulkSearchReq.Device != nil {
		queryParams["device"] = *bulkSearchReq.Device
	}
	if bulkSearchReq.Telemetry != nil {
		queryParams["telemetry"] = strconv.FormatBool(*bulkSearchReq.Telemetry)
	}

	resp, err := c.reqClient.
		R().
		SetQueryParams(queryParams).
		SetBody(bulkSearchReq).
		SetSuccessResult(&bulkSearchResp).
		Post(c.reqClient.BaseURL + "/indexes/bulk/search")

	if err != nil {
		logger.Error("error bulk searching", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error bulk searching", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error bulk searching: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info("bulk search completed")
	return &bulkSearchResp, nil
}
