package marqo

import (
	"fmt"
	"net/http"
)

// CreateIndexRequest is the request to create an index
type CreateIndexRequest struct {
	IndexName     string         `json:"-" validate:"required"`
	IndexDefaults *IndexDefaults `json:"index_defaults"`
	// Number of shards for the index (default: 3)
	NumberOfShards *int `json:"number_of_shards"`
	// Number of replicas for the index (default: 0)
	NumberOfReplicas *int `json:"number_of_replicas"`
}

// IndexDefaults is the defaults for the index
type IndexDefaults struct {
	TreatURLAndPointersAsImages *bool               `json:"treat_url_and_pointers_as_images"`
	Model                       *string             `json:"model"`
	ModelProperties             *ModelProperties    `json:"model_properties"`
	NormalizeEmbeddings         *bool               `json:"normalize_embeddings"`
	TextPreprocessing           *TextPreprocessing  `json:"text_preprocessing"`
	ImagePreprocessing          *ImagePreprocessing `json:"image_preprocessing"`
	ANNParameters               *ANNParameters      `json:"ann_parameters"`
}

// ModelProperties are the properties for the model
type ModelProperties struct {
	Name       *string `json:"name"`
	Dimensions *int    `json:"dimensions"`
	URL        *string `json:"url"`
	Type       *string `json:"type"`
}

// TextPreprocessing is the text preprocessing for the index
type TextPreprocessing struct {
	SplitLength             *int    `json:"split_length"`
	SplitOverlap            *int    `json:"split_overlap"`
	SplitMethod             *string `json:"split_method"`
	OverrideTextChunkPrefix *string `json:"override_text_chunk_prefix"`
	OverrideTextQueryPrefix *string `json:"override_text_query_prefix"`
}

// ImagePreprocessing is the image preprocessing for the index
type ImagePreprocessing struct {
	PatchMethod *string `json:"patch_method"`
}

// ANNParameters are the ANN parameters for the index
type ANNParameters struct {
	SpaceType  *string               `json:"space_type"`
	Parameters *HSNWMethodParameters `json:"parameters"`
}

// HSNWMethodParameters are the HSNW method parameters for the index
type HSNWMethodParameters struct {
	EFConstruction *int `json:"ef_construction"`
	M              *int `json:"m"`
}

// CreateIndexResponse is the response for creating an index
type CreateIndexResponse struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	Index              string `json:"index"`
}

// CreateIndex creates an index
func (c *Client) CreateIndex(createIndexReq *CreateIndexRequest) (*CreateIndexResponse, error) {
	logger := c.logger.With("method", "CreateIndex")
	err := validate.Struct(createIndexReq)
	if err != nil {
		logger.Error("error validating create index request",
			"error", err)
		return nil, err
	}

	var createIndexResp CreateIndexResponse
	resp, err := c.reqClient.
		R().
		SetBody(createIndexReq).
		SetSuccessResult(&createIndexResp).
		Post(c.reqClient.BaseURL + "/indexes/" + createIndexReq.IndexName)
	if err != nil {
		logger.Error("error creating index", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error creating index", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error creating index: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info("index created")
	return &createIndexResp, nil
}

// DeleteIndex deletes an index
func (c *Client) DeleteIndex(indexName string) error {
	logger := c.logger.With("method", "DeleteIndex")
	resp, err := c.reqClient.
		R().
		Delete(c.reqClient.BaseURL + "/indexes/" + indexName)
	if err != nil {
		logger.Error("error deleting index", "error", err)
		return err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error deleting index", "status_code", resp.
			Response.StatusCode)
		return fmt.Errorf("error deleting index: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info("index deleted")
	return nil
}

// ListIndexesResponse is the response for listing indexes
type ListIndexesResponse struct {
	Results []Result `json:"results"`
}

// Result is the result for listing one index
type Result struct {
	IndexName string `json:"index_name"`
}

// ListIndexes lists the indexes
func (c *Client) ListIndexes() (*ListIndexesResponse, error) {
	logger := c.logger.With("method", "ListIndexes")
	var result ListIndexesResponse

	resp, err := c.reqClient.
		R().
		SetSuccessResult(&result).
		Get(c.reqClient.BaseURL + "/indexes")
	if err != nil {
		logger.Error("error listing indexes", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error listing indexes", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error listing indexes: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response indexes: %+v", result))
	return &result, nil
}
