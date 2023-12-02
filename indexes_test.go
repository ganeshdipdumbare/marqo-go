package marqo

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/suite"
)

type IndexesTestSuite struct {
	suite.Suite
	Logger *slog.Logger
}

func (suite *IndexesTestSuite) SetupSuite() {
	suite.Logger = slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelError,
		}))
}

func (suite *IndexesTestSuite) TearDownSuite() {

}

func TestIndexesTestSuite(t *testing.T) {
	suite.Run(t, new(IndexesTestSuite))
}

func getMockServerForCreateIndex() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
			{
				"acknowledged": true,
				"shards_acknowledged": true,
				"index": "test"
			}`))
		}))
}
func (suite *IndexesTestSuite) TestClient_CreateIndex() {
	t := suite.T()
	modelName := "hf/all_datasets_v4_MiniLM-L6"
	mockServer := getMockServerForCreateIndex()
	type fields struct {
		url    string
		logger *slog.Logger
		client *req.Client
	}
	type args struct {
		createIndexReq *CreateIndexRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CreateIndexResponse
		wantErr bool
	}{
		{
			name: "create index successfully",
			fields: fields{
				url:    mockServer.URL,
				logger: suite.Logger,
			},
			args: args{
				createIndexReq: &CreateIndexRequest{
					IndexName: "test",
					IndexDefaults: &IndexDefaults{
						Model: &modelName,
					},
				},
			},
			want: &CreateIndexResponse{
				Acknowledged:       true,
				ShardsAcknowledged: true,
				Index:              "test",
			},
			wantErr: false,
		},
		{
			name: "create index with empty name fails",
			fields: fields{
				url:    mockServer.URL,
				logger: suite.Logger,
			},
			args: args{
				createIndexReq: &CreateIndexRequest{
					IndexDefaults: &IndexDefaults{
						Model: &modelName,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.fields.url, WithLogger(tt.fields.logger))
			if err != nil {
				t.Errorf("Client.Connect() error = %v", err)
				return
			}
			got, err := c.CreateIndex(tt.args.createIndexReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.CreateIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getMockServerForDeleteIndex() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
			{"acknowledged": true}`))
		}),
	)
}

func getMockServerForListIndexes() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
			{
				"results": [
				  {
					"index_name": "Book Collection"
				  },
				  {
					"index_name": "Animal facts"
				  }
				]
			}`))
		}),
	)
}

func (suite *IndexesTestSuite) TestClient_ListIndexes() {
	t := suite.T()
	type fields struct {
		url    string
		logger *slog.Logger
		client *req.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ListIndexesResponse
		wantErr bool
	}{
		{
			name: "list indexes successfully",
			fields: fields{
				url:    getMockServerForListIndexes().URL,
				logger: suite.Logger,
			},
			want: &ListIndexesResponse{
				Results: []Result{
					{
						IndexName: "Book Collection",
					},
					{
						IndexName: "Animal facts",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.fields.url, WithLogger(tt.fields.logger))
			if err != nil {
				t.Errorf("Client.Connect() error = %v", err)
				return
			}
			got, err := c.ListIndexes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListIndexes() = %v, want %v", got, tt.want)
			}
		})
	}
}
