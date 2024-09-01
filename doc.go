/*
Package marqo provides a client for interacting with the Marqo server.

The simplest way to use Marqo is to create a new client and use its methods to interact with the server:

	package main

	import (
	  "fmt"
	  "log"

	  marqo "github.com/ganeshdipdumbare/marqo-go"
	)

	func main() {
	  client, err := marqo.NewClient("http://localhost:8882")
	  if err != nil {
	    log.Fatalf("Failed to create client: %v", err)
	  }

	  modelsResponse, err := client.GetModels()
	  if err != nil {
	    log.Fatalf("Failed to get models: %v", err)
	  }
	  fmt.Printf("Loaded models: %+v\n", modelsResponse.Models)

	  createIndexReq := &marqo.CreateIndexRequest{
	    IndexName: "example_index",
	  }
	  createIndexResp, err := client.CreateIndex(createIndexReq)
	  if err != nil {
	    log.Fatalf("Failed to create index: %v", err)
	  }
	  fmt.Printf("CreateIndexResponse: %+v\n", createIndexResp)

	  upsertDocumentsReq := &marqo.UpsertDocumentsRequest{
	    IndexName: "example_index",
	    Documents: []interface{}{
	      map[string]interface{}{"_id": "1", "title": "Document 1"},
	      map[string]interface{}{"_id": "2", "title": "Document 2"},
	    },
	  }
	  upsertDocumentsResp, err := client.UpsertDocuments(upsertDocumentsReq)
	  if err != nil {
	    log.Fatalf("Failed to upsert documents: %v", err)
	  }
	  fmt.Printf("UpsertDocumentsResponse: %+v\n", upsertDocumentsResp)

	  searchReq := &marqo.SearchRequest{
	    IndexName: "example_index",
	    Q:         "Document",
	  }
	  searchResp, err := client.Search(searchReq)
	  if err != nil {
	    log.Fatalf("Failed to perform search: %v", err)
	  }
	  fmt.Printf("SearchResponse: %+v\n", searchResp)
	}

For a full guide visit https://github.com/ganeshdipdumbare/marqo-go

# License

This project is licensed under the MIT License - see the LICENSE file for details.
*/
package marqo
