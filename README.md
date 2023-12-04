[![Go](https://github.com/ganeshdipdumbare/marqo-go/actions/workflows/go.yml/badge.svg)](https://github.com/ganeshdipdumbare/marqo-go/actions/workflows/go.yml)

<p align="center">
<img src="https://uploads-ssl.webflow.com/62dfa8e3960a6e2b47dc7fae/62fdf9cef684e6f16158b094_MARQO%20LOGO-UPDATED-GREEN.svg" width="25%" height="25%">
</p>

## Marqo (Unofficial Go client)

Marqo is more than a vector database, it's an end-to-end vector search engine for both text and images. Vector generation, storage and retrieval are handled out of the box through a single API. No need to bring your own embeddings.

The Go client(Unofficial) allows you to connect to Marqo in less than 3 lines.

## Getting started

1. Marqo requires docker. To install Docker go to the [Docker Official website](https://docs.docker.com/get-docker/). Ensure that docker has at least 8GB memory and 50GB storage.

2. Use docker to run Marqo for Mac users:

```bash
## will start marqo-os and marqo service mentioned in docker compose.

docker compose up -d
```

3. Get the Marqo Go client and start using it in your program:

```bash
go get github.com/ganeshdipdumbare/marqo-go@latest
```

4. Start indexing and searching.

```Golang
package main

import (
	"fmt"

	marqo "github.com/ganeshdipdumbare/marqo-go"
)

// Document represents a document

type Document struct {
	ID          string `json:"_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
}

func main() {

	marqoClient, err := marqo.NewClient("http://localhost:8882")
	if err != nil {
		panic(err)
	}
	model, err := marqoClient.GetModels()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Model: %+v\n", model)

	indexes, err := marqoClient.ListIndexes()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Indexes: %+v\n", indexes)

	// create index
	resp, err := marqoClient.CreateIndex(&marqo.CreateIndexRequest{
		IndexName: "test1",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("CreateIndexResponse: %+v\n", resp)

	// create document
	documents := []Document{
		{
			ID:          "1",
			Title:       "The Great Gatsby",
			Description: "The Great Gatsby is a 1925 novel by American writer F. Scott Fitzgerald.",
			Genre:       "Novel",
		},
		{
			ID:          "2",
			Title:       "The Catcher in the Rye",
			Description: "The Catcher in the Rye is a novel by J. D. Salinger, partially published in serial form in 1945–1946 and as a novel in 1951.",
			Genre:       "Novel",
		},
	}
	docInterface := make([]interface{}, len(documents))
	for i, v := range documents {
		docInterface[i] = v
	}

	upsertResp, err := marqoClient.UpsertDocuments(&marqo.UpsertDocumentsRequest{
		IndexName: "test1",
		Documents: docInterface,
		TensorFields: []string{
			"title",
			"description",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("UpsertDocumentsResponse: %+v\n", upsertResp)

	// delete documents
	deleteResp, err := marqoClient.DeleteDocuments(&marqo.DeleteDocumentsRequest{
		IndexName:   "test1",
		DocumentIDs: []string{"1", "2"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("DeleteDocumentsResponse: %+v\n", deleteResp)

	// Refresh index
	refreshResp, err := marqoClient.RefreshIndex(&marqo.RefreshIndexRequest{
		IndexName: "test1",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("RefreshIndexResponse: %+v\n", refreshResp)

	// Get index stats
	statsResp, err := marqoClient.GetIndexStats(&marqo.GetIndexStatsRequest{
		IndexName: "test1",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetIndexStatsResponse: %+v\n", statsResp)

	// Get index settings
	settingsResp, err := marqoClient.GetIndexSettings(&marqo.GetIndexSettingsRequest{
		IndexName: "test1",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetIndexSettingsResponse: %+v\n", *settingsResp.IndexDefaults.Model)

	// Get index health
	healthResp, err := marqoClient.GetIndexHealth(&marqo.GetIndexHealthRequest{
		IndexName: "test1",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetIndexHealthResponse: %+v\n", healthResp)

	// Search
	searchQuery := "The Great Gatsby"
	searchResp, err := marqoClient.Search(&marqo.SearchRequest{
		IndexName: "test1",
		Q:         &searchQuery,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("SearchResponse: %+v\n", searchResp)

	// Bulk search
	//searchQuery := "The Great Gatsby"
	indexName := "test1"
	bulkSearchResp, err := marqoClient.BulkSearch(&marqo.BulkSearchRequest{
		Queries: []marqo.SearchRequest{
			{
				Index: &indexName,
				Q:     &searchQuery,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("BulkSearchResponse: %+v\n", bulkSearchResp)

}
```

## Improvements

- Add UT for all the APIs.
- Add end to end examples for search.
