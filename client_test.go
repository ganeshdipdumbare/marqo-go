package marqo

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	MockServer *httptest.Server
	ReqClient  *req.Client
	Logger     *slog.Logger
}

func (suite *ClientTestSuite) SetupSuite() {
	suite.MockServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
			{
				"models": [
					{
						"model_name": "hf/all_datasets_v4_MiniLM-L6", 
						"model_device": "cpu"
					}
				]
			}`))
		}))
	suite.ReqClient = req.NewClient()
	suite.Logger = slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelError,
		}))
}

func (suite *ClientTestSuite) TearDownSuite() {
	suite.MockServer.Close()
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (suite *ClientTestSuite) TestNewClient() {
	t := suite.T()
	type args struct {
		url string
		opt []Options
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "test new client with logger option",
			args: args{
				url: suite.MockServer.URL,
				opt: []Options{
					WithLogger(suite.Logger),
				},
			},
			want: &Client{
				url:       suite.MockServer.URL,
				logger:    suite.Logger,
				reqClient: suite.ReqClient,
			},
			wantErr: false,
		},
		{
			name: "test new client with no options",
			args: args{
				url: suite.MockServer.URL,
			},
			want: &Client{
				url:       suite.MockServer.URL,
				reqClient: suite.ReqClient,
			},
			wantErr: false,
		},
		{
			name: "test new client with no url",
			args: args{
				url: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.url, tt.args.opt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// check url
			if got != nil && got.url != tt.want.url {
				t.Errorf("NewClient() url = %v, want %v", got.url, tt.want.url)
				return
			}
			// check logger as by default it should be set
			if got != nil && got.logger == nil {
				t.Errorf("NewClient() logger = %v, want %v", got.logger, tt.want.logger)
				return
			}
			// check req client
			if got != nil && got.reqClient.BaseURL != tt.args.url {
				t.Errorf("NewClient() reqClient.BaseURL = %v, want %v", got.reqClient.BaseURL, tt.args.url)
				return
			}
		})
	}
}
