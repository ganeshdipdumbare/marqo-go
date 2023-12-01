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
	Documents            []interface{}          `json:"documents" validate:"required"`
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

// DeleteDocumentsRequest is the request to delete documents
type DeleteDocumentsRequest struct {
	IndexName   string   `json:"-" validate:"required"`
	DocumentIDs []string `json:"document_ids" validate:"required"`
}

// DeleteDocumentsResponse is the response from the server
type DeleteDocumentsResponse struct {
	IndexName string `json:"index_name"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Details   struct {
		ReceivedDocumentIds int `json:"receivedDocumentIds"`
		DeletedDocuments    int `json:"deletedDocuments"`
	} `json:"details"`
	Items      []DeletedItem `json:"items"`
	Duration   string        `json:"duration"`
	StartedAt  string        `json:"startedAt"`
	FinishedAt string        `json:"finishedAt"`
}

// DeletedItem is the item which was deleted
type DeletedItem struct {
	ID     string `json:"_id"`
	Shards struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Status int    `json:"status"`
	Result string `json:"result"`
}

// DeleteDocuments deletes documents from the server
func (c *Client) DeleteDocuments(deleteDocumentsReq *DeleteDocumentsRequest) (*DeleteDocumentsResponse, error) {
	logger := c.logger.With("method", "DeleteDocuments")
	err := validate.Struct(deleteDocumentsReq)
	if err != nil {
		logger.Error("error validating delete documents request",
			"error", err)
		return nil, err
	}

	var deleteDocumentsResp DeleteDocumentsResponse
	resp, err := c.reqClient.
		R().
		SetBody(deleteDocumentsReq.DocumentIDs).
		SetSuccessResult(&deleteDocumentsResp).
		Post(c.reqClient.BaseURL + "/indexes/" + deleteDocumentsReq.IndexName + "/documents/delete-batch")

	if err != nil {
		logger.Error("error deleting documents", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error deleting documents", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error deleting documents: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response delete documents: %+v",
		deleteDocumentsResp))
	return &deleteDocumentsResp, nil
}

// GetDocumentRequest is the request to get a document
type GetDocumentRequest struct {
	IndexName    string `json:"-" validate:"required"`
	DocumentID   string `json:"document_id" validate:"required"`
	ExposeFacets bool   `json:"expose_facets,omitempty"`
}

// GetDocumentResponse is the response from the server
type GetDocumentResponse map[string]interface{}

// GetDocument gets a document from the server
func (c *Client) GetDocument(getDocumentReq *GetDocumentRequest) (*GetDocumentResponse, error) {
	logger := c.logger.With("method", "GetDocument")
	err := validate.Struct(getDocumentReq)
	if err != nil {
		logger.Error("error validating get document request",
			"error", err)
		return nil, err
	}

	var getDocumentResp GetDocumentResponse
	queryParams := map[string]string{}
	if getDocumentReq.ExposeFacets {
		queryParams["expose_facets"] = strconv.FormatBool(getDocumentReq.ExposeFacets)
	}

	resp, err := c.reqClient.
		R().
		SetQueryParams(queryParams).
		SetSuccessResult(&getDocumentResp).
		Get(c.reqClient.BaseURL + "/indexes/" + getDocumentReq.IndexName + "/documents/" + getDocumentReq.DocumentID)

	if err != nil {
		logger.Error("error getting document", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error getting document", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error getting document: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response get document: %+v",
		getDocumentResp))
	return &getDocumentResp, nil
}
