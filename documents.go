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

// UpsertDocuments upserts documents to the server.
//
// This method sends a POST request to the server to upsert the specified documents.
//
// Parameters:
//
//	upsertDocumentsReq (*UpsertDocumentsRequest): The request containing the documents to be upserted.
//
// Returns:
//
//	*UpsertDocumentsResponse: The response from the server.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the upsertDocumentsReq parameter.
// 2. Sends a POST request to the server with the documents in the request body.
// 3. Checks the response status code and logs any errors.
// 4. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	upsertDocumentsReq := &UpsertDocumentsRequest{
//	    IndexName: "example_index",
//	    Documents: []interface{}{...},
//	}
//	resp, err := client.UpsertDocuments(upsertDocumentsReq)
//	if err != nil {
//	    log.Fatalf("Failed to upsert documents: %v", err)
//	}
//	fmt.Printf("UpsertDocumentsResponse: %+v\n", resp)
func (c *Client) UpsertDocuments(upsertDocumentsReq *UpsertDocumentsRequest) (*UpsertDocumentsResponse, error) {
	logger := c.logger.With("method", "UpsertDocuments")
	err := validate.Struct(upsertDocumentsReq)
	if err != nil {
		logger.Error("error validating upsert documents request", "error", err)
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
		logger.Error("error upserting documents", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error upserting documents: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response upsert documents: %+v", upsertDocumentsResp))
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

// DeleteDocuments deletes documents from the server.
//
// This method sends a POST request to the server to delete the specified documents.
//
// Parameters:
//
//	deleteDocumentsReq (*DeleteDocumentsRequest): The request containing the document IDs to be deleted.
//
// Returns:
//
//	*DeleteDocumentsResponse: The response from the server.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the deleteDocumentsReq parameter.
// 2. Sends a POST request to the server with the document IDs in the request body.
// 3. Checks the response status code and logs any errors.
// 4. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	deleteDocumentsReq := &DeleteDocumentsRequest{
//	    IndexName:   "example_index",
//	    DocumentIDs: []string{"doc1", "doc2"},
//	}
//	resp, err := client.DeleteDocuments(deleteDocumentsReq)
//	if err != nil {
//	    log.Fatalf("Failed to delete documents: %v", err)
//	}
//	fmt.Printf("DeleteDocumentsResponse: %+v\n", resp)
func (c *Client) DeleteDocuments(deleteDocumentsReq *DeleteDocumentsRequest) (*DeleteDocumentsResponse, error) {
	logger := c.logger.With("method", "DeleteDocuments")
	err := validate.Struct(deleteDocumentsReq)
	if err != nil {
		logger.Error("error validating delete documents request", "error", err)
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
		logger.Error("error deleting documents", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error deleting documents: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response delete documents: %+v", deleteDocumentsResp))
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

// GetDocument gets a document from the server.
//
// This method sends a GET request to the server to retrieve the specified document.
//
// Parameters:
//
//	getDocumentReq (*GetDocumentRequest): The request containing the document ID to be retrieved.
//
// Returns:
//
//	*GetDocumentResponse: The response from the server.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the getDocumentReq parameter.
// 2. Sends a GET request to the server with the document ID as a query parameter.
// 3. Checks the response status code and logs any errors.
// 4. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	getDocumentReq := &GetDocumentRequest{
//	    IndexName:  "example_index",
//	    DocumentID: "doc1",
//	}
//	resp, err := client.GetDocument(getDocumentReq)
//	if err != nil {
//	    log.Fatalf("Failed to get document: %v", err)
//	}
//	fmt.Printf("GetDocumentResponse: %+v\n", resp)
func (c *Client) GetDocument(getDocumentReq *GetDocumentRequest) (*GetDocumentResponse, error) {
	logger := c.logger.With("method", "GetDocument")
	err := validate.Struct(getDocumentReq)
	if err != nil {
		logger.Error("error validating get document request", "error", err)
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
		logger.Error("error getting document", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error getting document: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response get document: %+v", getDocumentResp))
	return &getDocumentResp, nil
}

// GetDocumentsRequest is the request to get documents
type GetDocumentsRequest struct {
	IndexName    string   `json:"-" validate:"required"`
	DocumentIDs  []string `json:"document_ids" validate:"required"`
	ExposeFacets bool     `json:"expose_facets,omitempty"`
}

// GetDocumentsResponse is the response from the server
type GetDocumentsResponse struct {
	Results []GetDocumentResponse `json:"results"`
}

// GetDocuments gets documents from the server.
//
// This method sends a GET request to the server to retrieve the specified documents.
//
// Parameters:
//
//	getDocumentsReq (*GetDocumentsRequest): The request containing the document IDs to be retrieved.
//
// Returns:
//
//	*GetDocumentsResponse: The response from the server.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the getDocumentsReq parameter.
// 2. Sends a GET request to the server with the document IDs in the request body.
// 3. Checks the response status code and logs any errors.
// 4. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	getDocumentsReq := &GetDocumentsRequest{
//	    IndexName:   "example_index",
//	    DocumentIDs: []string{"doc1", "doc2"},
//	}
//	resp, err := client.GetDocuments(getDocumentsReq)
//	if err != nil {
//	    log.Fatalf("Failed to get documents: %v", err)
//	}
//	fmt.Printf("GetDocumentsResponse: %+v\n", resp)
func (c *Client) GetDocuments(getDocumentsReq *GetDocumentsRequest) (*GetDocumentsResponse, error) {
	logger := c.logger.With("method", "GetDocuments")
	err := validate.Struct(getDocumentsReq)
	if err != nil {
		logger.Error("error validating get documents request", "error", err)
		return nil, err
	}

	var getDocumentsResp GetDocumentsResponse
	queryParams := map[string]string{}
	if getDocumentsReq.ExposeFacets {
		queryParams["expose_facets"] = strconv.FormatBool(getDocumentsReq.ExposeFacets)
	}

	resp, err := c.reqClient.
		R().
		SetQueryParams(queryParams).
		SetBody(getDocumentsReq.DocumentIDs).
		SetSuccessResult(&getDocumentsResp).
		Get(c.reqClient.BaseURL + "/indexes/" + getDocumentsReq.IndexName + "/documents")
	if err != nil {
		logger.Error("error getting documents", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error getting documents", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error getting documents: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response get documents: %+v", getDocumentsResp))
	return &getDocumentsResp, nil
}
