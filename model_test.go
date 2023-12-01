package marqo

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/imroc/req/v3"
)

func TestClient_GetModels(t *testing.T) {
	type fields struct {
		url    string
		logger *slog.Logger
		client *req.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    *GetModelsResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				url:       tt.fields.url,
				logger:    tt.fields.logger,
				reqClient: tt.fields.client,
			}
			got, err := c.GetModels()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetModels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_EjectModel(t *testing.T) {
	type fields struct {
		url    string
		logger *slog.Logger
		client *req.Client
	}
	type args struct {
		ejectModelReq *EjectModelRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.fields.url, WithLogger(tt.fields.logger))
			if err != nil {
				t.Errorf("Client.Connect() error = %v", err)
				return
			}
			if err = c.EjectModel(tt.args.ejectModelReq); (err != nil) != tt.wantErr {
				t.Errorf("Client.EjectModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
