package marqo

import (
	"fmt"
	"net/http"
	"strconv"
)

// UpsertDocumentsRequest is the request to upsert documents
type UpsertDocumentsRequest struct {
	IndexName string `json:"-" validate:"required"`
	// Query params
	Refresh   *bool   `json:"-"`
	Device    *string `json:"-"`
	Telemetry *bool   `json:"-"`
	// Body params
	Documents            [][]byte               `json:"documents" validate:"required"`
	TensorFields         []string               `json:"tensorFields,omitempty"`
	UseExistingTensors   *bool                  `json:"useExistingTensors,omitempty"`
	ImageDownloadHeaders map[string]string      `json:"imageDownloadHeaders,omitempty"`
	Mappings             map[string]interface{} `json:"mappings,omitempty"`
	ModelAuth            map[string]interface{} `json:"modelAuth,omitempty"`
	TextChunkPrefix      *string                `json:"textChunkPrefix,omitempty"`
	ClientBatchSize      *int                   `json:"client_batch_size,omitempty"`
}

// UpsertDocumentsResponse is the response from the server
type UpsertDocumentsResponse struct {
	Errors           bool    `json:"errors"`
	Items            []Item  `json:"items"`
	ProcessingTimeMS float64 `json:"processingTimeMs"`
	IndexName        string  `json:"index_name"`
}

// Item is the item from the server
type Item struct {
	ID     string `json:"_id"`
	Result string `json:"result"`
	Status int    `json:"status"`
}

// UpsertDocuments upserts documents to the server
func (c *Client) UpsertDocuments(upsertDocumentsReq *UpsertDocumentsRequest) (*UpsertDocumentsResponse, error) {
	logger := c.logger.With("method", "UpsertDocuments")
	err := validate.Struct(upsertDocumentsReq)
	if err != nil {
		logger.Error("error validating upsert documents request",
			"error", err)
		return nil, err
	}

	var upsertDocumentsResp UpsertDocumentsResponse
	queryParams := map[string]string{}
	if upsertDocumentsReq.Refresh != nil {
		queryParams["refresh"] = strconv.FormatBool(*upsertDocumentsReq.Refresh)
	}
	if upsertDocumentsReq.Device != nil {
		queryParams["device"] = *upsertDocumentsReq.Device
	}
	if upsertDocumentsReq.Telemetry != nil {
		queryParams["telemetry"] = strconv.FormatBool(*upsertDocumentsReq.Telemetry)
	}

	resp, err := c.reqClient.
		R().
		SetQueryParams(queryParams).
		SetBody(upsertDocumentsReq).
		SetSuccessResult(&upsertDocumentsResp).
		Post(c.reqClient.BaseURL + "/indexes/" + upsertDocumentsReq.IndexName + "/documents")

	if err != nil {
		logger.Error("error upserting documents", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error upserting documents", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error upserting documents: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response upsert documents: %+v",
		upsertDocumentsResp))
	return &upsertDocumentsResp, nil
}
