package marqo

import (
	"fmt"
	"net/http"
)

// GetModelsResponse is the response from the server
type GetModelsResponse struct {
	Models []Model `json:"models"`
}

// Model is the model from the server
type Model struct {
	ModelName   string `json:"model_name"`
	ModelDevice string `json:"model_device"`
}

// Model returns the loaded models from the server
func (c *Client) GetModels() (*GetModelsResponse, error) {
	logger := c.logger.With("method", "Model")
	var result GetModelsResponse

	resp, err := c.reqClient.
		R().
		SetSuccessResult(&result).
		Get(c.reqClient.BaseURL + "/models")
	if err != nil {
		logger.Error("error getting models", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error getting models", "status_code", resp.
			Response.StatusCode)
		return nil, fmt.Errorf("error getting models: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response models: %+v",
		result))
	return &result, nil
}

// EjectModelRequest is the request to eject a model
type EjectModelRequest struct {
	ModelName   string `validate:"required" json:"model_name"`
	ModelDevice string `validate:"oneof=cpu cuda" json:"model_device"`
}

// EjectModel ejects the model from the server cache
func (c *Client) EjectModel(ejectModelReq *EjectModelRequest) error {
	logger := c.logger.With("method", "EjectModel")
	err := validate.Struct(ejectModelReq)
	if err != nil {
		logger.Error("error validating eject model request",
			"error", err)
		return err
	}

	resp, err := c.reqClient.
		R().
		SetQueryParams(
			map[string]string{
				"model_name":   ejectModelReq.ModelName,
				"model_device": ejectModelReq.ModelDevice,
			},
		).
		Delete(c.reqClient.BaseURL + "/models")

	if err != nil {
		logger.Error("error ejecting model", "error", err)
		return err
	}
	if resp.Response.StatusCode != http.StatusOK {
		logger.Error("error ejecting model", "status_code", resp.
			Response.StatusCode)
		return fmt.Errorf("error ejecting model: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info("ejected model successfully",
		"model_name", ejectModelReq.ModelName,
		"model_device", ejectModelReq.ModelDevice)
	return nil
}
