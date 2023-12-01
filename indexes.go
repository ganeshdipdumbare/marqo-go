package marqo

import (
	"fmt"
	"net/http"
)

// CreateIndexRequest is the request to create an index
type CreateIndexRequest struct {
	IndexName     string         `json:"-" validate:"required"`
	IndexDefaults *IndexDefaults `json:"index_defaults,omitempty"`
	// Number of shards for the index (default: 3)
	NumberOfShards *int `json:"number_of_shards,omitempty"`
	// Number of replicas for the index (default: 0)
	NumberOfReplicas *int `json:"number_of_replicas,omitempty"`
}

// IndexDefaults is the defaults for the index
type IndexDefaults struct {
	// Fetch images from points and URLs (default: false)
	TreatURLsAndPointersAsImages *bool `json:"treat_urls_and_pointers_as_images,omitempty"`
	// Model to vectorize doc content (default: hf/all_datasets_v4_MiniLM-L6)
	Model           *string          `json:"model,omitempty"`
	ModelProperties *ModelProperties `json:"model_properties,omitempty"`
	// TODO: add search model support in the future
	// SearchModel to vectorize queries (default: hf/all_datasets_v4_MiniLM-L6)
	// SearchModel           *string          `json:"search_model,omitempty"`
	// SearchModelProperties *ModelProperties `json:"search_model_properties,omitempty"`
	// Normalize embeddings to have unit length (default: true)
	NormalizeEmbeddings *bool               `json:"normalize_embeddings,omitempty"`
	TextPreprocessing   *TextPreprocessing  `json:"text_preprocessing,omitempty"`
	ImagePreprocessing  *ImagePreprocessing `json:"image_preprocessing,omitempty"`
	ANNParameters       *ANNParameters      `json:"ann_parameters,omitempty"`
}

// ModelProperties are the properties for the model
type ModelProperties struct {
	Name       *string `json:"name,omitempty"`
	Dimensions *int    `json:"dimensions,omitempty"`
	URL        *string `json:"url,omitempty"`
	Type       *string `json:"type,omitempty"`
}

// TextPreprocessing is the text preprocessing for the index
type TextPreprocessing struct {
	// SplitLength is length of chunks after splitting
	// by split method (default: 2)
	SplitLength *int `json:"split_length,omitempty"`
	// SplitOverlap is overlap between adjacent chunks (default: 0)
	SplitOverlap *int `json:"split_overlap,omitempty"`
	// SplitMethod method to split text into chunks (default: "sentence", options: "sentence", "word", "character" or "passage")
	SplitMethod *string `json:"split_method,omitempty"`
}

// ImagePreprocessing is the image preprocessing for the index
type ImagePreprocessing struct {
	// The method by which images are chunked (options: "simple" or "frcnn")
	PatchMethod *string `json:"patch_method,omitempty"`
}

// ANNParameters are the ANN parameters for the index
type ANNParameters struct {
	// The function used to measure the distance between two points in ANN (l1, l2, linf, or cosinesimil. default: cosinesimil)
	SpaceType *string `json:"space_type,omitempty"`
	// The hyperparameters for the ANN method (which is always hnsw for Marqo).
	Parameters *HSNWMethodParameters `json:"parameters,omitempty"`
}

// HSNWMethodParameters are the HSNW method parameters for the index
type HSNWMethodParameters struct {
	// The size of the dynamic list used during k-NN graph creation.
	// Higher values lead to a more accurate graph but slower indexing
	// speed. It is recommended to keep this between 2 and 800 (maximum is 4096)
	// (default: 128)
	EFConstruction *int `json:"ef_construction,omitempty"`
	// The number of bidirectional links that the plugin creates for each
	// new element. Increasing and decreasing this value can have a
	// large impact on memory consumption. Keep this value between 2 and 100.
	// (default: 16)
	M *int `json:"m,omitempty"`
}

// CreateIndexResponse is the response for creating an index
type CreateIndexResponse struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	Index              string `json:"index"`
}

// setDefaultCreateIndexRequest add default values to createIndexRequest if not set
func setDefaultCreateIndexRequest(createIndexReq *CreateIndexRequest) {
	if createIndexReq.NumberOfShards == nil {
		createIndexReq.NumberOfShards = new(int)
		*createIndexReq.NumberOfShards = 3
	}

	if createIndexReq.NumberOfReplicas == nil {
		createIndexReq.NumberOfReplicas = new(int)
		*createIndexReq.NumberOfReplicas = 0
	}

	if createIndexReq.IndexDefaults == nil {
		createIndexReq.IndexDefaults = new(IndexDefaults)
	}

	if createIndexReq.IndexDefaults.TreatURLsAndPointersAsImages == nil {
		createIndexReq.IndexDefaults.TreatURLsAndPointersAsImages = new(bool)
		*createIndexReq.IndexDefaults.TreatURLsAndPointersAsImages = false
	}

	if createIndexReq.IndexDefaults.Model == nil {
		createIndexReq.IndexDefaults.Model = new(string)
		*createIndexReq.IndexDefaults.Model = "hf/all_datasets_v4_MiniLM-L6"
	}

	// TODO: add search model support in the future
	// if createIndexReq.IndexDefaults.SearchModel == nil {
	// 	createIndexReq.IndexDefaults.SearchModel = new(string)
	// 	*createIndexReq.IndexDefaults.SearchModel =
	// 		"hf/all_datasets_v4_MiniLM-L6"
	// }

	if createIndexReq.IndexDefaults.NormalizeEmbeddings == nil {
		createIndexReq.IndexDefaults.NormalizeEmbeddings = new(bool)
		*createIndexReq.IndexDefaults.NormalizeEmbeddings = true
	}

	if createIndexReq.IndexDefaults.TextPreprocessing != nil {
		if createIndexReq.IndexDefaults.TextPreprocessing.SplitLength == nil {
			createIndexReq.IndexDefaults.TextPreprocessing.SplitLength =
				new(int)
			*createIndexReq.IndexDefaults.TextPreprocessing.SplitLength = 2
		}

		if createIndexReq.IndexDefaults.TextPreprocessing.SplitOverlap == nil {
			createIndexReq.IndexDefaults.TextPreprocessing.SplitOverlap =
				new(int)
			*createIndexReq.IndexDefaults.TextPreprocessing.SplitOverlap = 0
		}

		if createIndexReq.IndexDefaults.TextPreprocessing.SplitMethod == nil {
			createIndexReq.IndexDefaults.TextPreprocessing.SplitMethod =
				new(string)
			*createIndexReq.IndexDefaults.TextPreprocessing.SplitMethod =
				"sentence"
		}
	}

	if createIndexReq.IndexDefaults.ImagePreprocessing != nil {
		if createIndexReq.IndexDefaults.ImagePreprocessing.PatchMethod == nil {
			createIndexReq.IndexDefaults.ImagePreprocessing.PatchMethod =
				new(string)
			*createIndexReq.IndexDefaults.ImagePreprocessing.PatchMethod =
				"simple"
		}
	}

	if createIndexReq.IndexDefaults.ANNParameters != nil {
		if createIndexReq.IndexDefaults.ANNParameters.SpaceType == nil {
			createIndexReq.IndexDefaults.ANNParameters.SpaceType = new(string)
			*createIndexReq.IndexDefaults.ANNParameters.SpaceType =
				"cosinesimil"
		}

		if createIndexReq.IndexDefaults.ANNParameters.Parameters == nil {
			createIndexReq.IndexDefaults.ANNParameters.Parameters =
				new(HSNWMethodParameters)
		}

		if createIndexReq.IndexDefaults.ANNParameters.Parameters.EFConstruction ==
			nil {
			createIndexReq.IndexDefaults.ANNParameters.Parameters.EFConstruction =
				new(int)
			*createIndexReq.IndexDefaults.ANNParameters.Parameters.EFConstruction =
				128
		}

		if createIndexReq.IndexDefaults.ANNParameters.Parameters.M == nil {
			createIndexReq.IndexDefaults.ANNParameters.Parameters.M = new(int)
			*createIndexReq.IndexDefaults.ANNParameters.Parameters.M = 16
		}
	}

}

// CreateIndex creates an index
func (c *Client) CreateIndex(createIndexReq *CreateIndexRequest) (*CreateIndexResponse, error) {
	logger := c.logger.With("method", "CreateIndex")
	setDefaultCreateIndexRequest(createIndexReq)
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
