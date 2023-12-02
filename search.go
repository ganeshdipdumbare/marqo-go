package marqo

import (
	"fmt"
	"strconv"
	"strings"
)

// SearchRequest is the request to search
type SearchRequest struct {
	// IndexName is the name of the index
	// Same is used in bulk search request as well so the json tag is "index"
	IndexName string `validate:"required" json:"index"`

	// Body params
	// Q is the query string
	Q *string `json:"q,omitempty"`
	// Limit is the number of results to return (default: 20)
	Limit *int `json:"limit,omitempty"`
	// Offset is the number of results to skip (default: 0)
	Offset *int    `json:"offset,omitempty"`
	Filter *string `json:"filter,omitempty"`
	// SearchableAttributes is the list of
	// attributes to search in (default: ["*"])
	SearchableAttributes []string `json:"searchableAttributes,omitempty"`
	// ShowHighlights return highlights for the document match.
	// Only applicable for TENSOR search. With LEXICAL search,
	// highlights will always be []. (default: true)
	ShowHighlights *bool `json:"showHighlights,omitempty"`
	// SearchMethod is the search method to use,
	// can be LEXICAL or TENSOR (default: LEXICAL)
	SearchMethod *string `json:"searchMethod,omitempty"`
	// AttributesToRetrieve is the list of attributes to retrieve
	// (default: ["*"])
	AttributesToRetrieve []string `json:"attributesToRetrieve,omitempty"`
	// ReRanker method to use for re-ranking results
	// (default: null)
	// options: null, "owl/ViT-B/32", "owl/ViT-B/16", "owl/ViT-L/14"
	ReRanker *string `json:"reRanker,omitempty"`
	// Boost is the map of attribute
	// (string): 2-Array [weight (float), bias (float)]
	// e.g. {"price": [0.5, 0.5]}
	Boost map[string][2]float64 `json:"boost,omitempty"`
	// ImageDownloadHeaders is for the image download. Can be used to authenticate the images for download.
	ImageDownloadHeaders map[string]interface{} `json:"image_download_headers,omitempty"`
	// Context allows you to use your own vectors as context for your queries.
	// Your vectors will be incorporated into the query using a weighted sum
	// approach, allowing you to reduce the number of inference requests for
	// duplicated content. The dimension of the provided vectors should be
	// consistent with the index dimension.
	Context *Context `json:"context,omitempty"`
	// ScoreModifiers is an object with two optional keys: multiply_score_by
	// and add_to_score. The value of each of these keys is an array of objects
	// that each contain the name of a numeric field in the document as the
	// field_name key and the weighting that should be applied to the numeric
	// value, as the weight key.
	ScoreModifiers *ScoreModifiers `json:"scoreModifiers,omitempty"`
	// ModelAuth is an authorisation details used by Marqo to download
	// non-publicly available models.
	ModelAuth map[string]interface{} `json:"modelAuth,omitempty"`
	// TextChunkPrefix is a string to be added to the start of all text queries
	// before vectorisation. Text itself will not be returned or used for
	// lexical search. Only affects vectors generated.
	TextQueryPrefix *string `json:"textQueryPrefix,omitempty"`

	// Query params
	// Device used to search. If device is not specified and CUDA devices are
	// available to Marqo (see here for more info), Marqo will speed up search
	// by using an available CUDA device. Otherwise, the CPU will be used.
	Device *string `json:"device,omitempty"`
	// Telemetry if true, the telemtry object is returned in the search
	// response body. This includes information like latency metrics.
	Telemetry *bool `json:"telemetry,omitempty"`
}

// Context is the tensor context for the search
type Context struct {
	Tensor []Tensor `json:"tensor"`
}

// Tensor is the tensor for the search
type Tensor struct {
	Vector []float64 `json:"vector"`
	Weight float64   `json:"weight"`
}

// ScoreModifiers is the score modifiers for the search
// options for keys: "multiply_score_by", "add_to_score"
type ScoreModifiers map[string][]ScoreModifier

// ScoreModifier is the score modifier for the search
type ScoreModifier struct {
	FieldName string `json:"field_name"`
	// Weight is the weighting that should be applied to the numeric value
	// default: 1 for multiply_score_by, 0 for add_to_score
	Weight float64 `json:"weight"`
}

// SearchResponse is the response from the server
type SearchResponse struct {
	// Hits is the list of hits
	Hits []map[string]interface{} `json:"hits"`
	// Limit is the number of results to return (default: 20)
	Limit int `json:"limit"`
	// Offset is the number of results to skip (default: 0)
	Offset int `json:"offset"`
	// ProcessingTimeMS is the processing time in milliseconds
	ProcessingTimeMS float64 `json:"processingTimeMs"`
	// Query is the query string
	Query string `json:"query"`
}

// Search searches the index
func (c *Client) Search(searchReq *SearchRequest) (*SearchResponse, error) {
	logger := c.logger.With("method", "Search")
	err := validate.Struct(searchReq)
	if err != nil {
		logger.Error("error validating search request",
			"error", err)
		return nil, err
	}

	var searchResp SearchResponse
	queryParams := map[string]string{}
	if searchReq.Limit != nil {
		queryParams["limit"] = strconv.Itoa(*searchReq.Limit)
	}
	if searchReq.Offset != nil {
		queryParams["offset"] = strconv.Itoa(*searchReq.Offset)
	}
	if searchReq.Filter != nil {
		queryParams["filter"] = *searchReq.Filter
	}
	if searchReq.SearchableAttributes != nil {
		queryParams["searchableAttributes"] = strings.Join(searchReq.SearchableAttributes, ",")
	}
	if searchReq.ShowHighlights != nil {
		queryParams["showHighlights"] = strconv.FormatBool(*searchReq.ShowHighlights)
	}
	if searchReq.SearchMethod != nil {
		queryParams["searchMethod"] = *searchReq.SearchMethod
	}
	if searchReq.AttributesToRetrieve != nil {
		queryParams["attributesToRetrieve"] = strings.Join(searchReq.AttributesToRetrieve, ",")
	}
	if searchReq.ReRanker != nil {
		queryParams["reRanker"] = *searchReq.ReRanker
	}
	if searchReq.TextQueryPrefix != nil {
		queryParams["textQueryPrefix"] = *searchReq.TextQueryPrefix
	}
	if searchReq.Device != nil {
		queryParams["device"] = *searchReq.Device
	}
	if searchReq.Telemetry != nil {
		queryParams["telemetry"] = strconv.FormatBool(*searchReq.Telemetry)
	}

	resp, err := c.reqClient.
		R().
		SetQueryParams(queryParams).
		SetBody(searchReq).
		SetSuccessResult(&searchResp).
		Post(c.reqClient.BaseURL + "/indexes/" + searchReq.IndexName + "/search")
	if err != nil {
		logger.Error("error searching", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != 200 {
		logger.Error("error searching", "status_code",
			resp.Response.StatusCode)
		return nil, fmt.Errorf(
			"error searching: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response search: %+v\n",
		searchResp))
	return &searchResp, nil
}
