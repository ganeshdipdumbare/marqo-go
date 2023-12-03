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

type StatsTestSuite struct {
	suite.Suite
	MockServer *httptest.Server
	ReqClient  *req.Client
	Logger     *slog.Logger
}

func (suite *StatsTestSuite) SetupSuite() {
	suite.MockServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// disabeling lint as this is a mock server
			// nolint
			w.Write([]byte(`
			{
				"numberOfDocuments": 2,
				"numberOfVectors": 3
			}
			  `))
		}))
	suite.ReqClient = req.NewClient()
	suite.Logger = slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelError,
		}))
}

func (suite *StatsTestSuite) TearDownSuite() {
	suite.MockServer.Close()
}

func TestStatsTestSuite(t *testing.T) {
	suite.Run(t, new(StatsTestSuite))
}

func (suite *StatsTestSuite) TestClient_GetIndexStats() {
	t := suite.T()
	type fields struct {
		url    string
		logger *slog.Logger
	}
	type args struct {
		getIndexStatsReq *GetIndexStatsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GetIndexStatsResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				url:    suite.MockServer.URL,
				logger: suite.Logger,
			},
			args: args{
				getIndexStatsReq: &GetIndexStatsRequest{
					IndexName: "test",
				},
			},
			want: &GetIndexStatsResponse{
				NumberOfDocuments: 2,
				NumberOfVectors:   3,
			},
			wantErr: false,
		},
		{
			name: "empty index name",
			fields: fields{
				url:    suite.MockServer.URL,
				logger: suite.Logger,
			},
			args: args{
				getIndexStatsReq: &GetIndexStatsRequest{
					IndexName: "",
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
			got, err := c.GetIndexStats(tt.args.getIndexStatsReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetIndexStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetIndexStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
