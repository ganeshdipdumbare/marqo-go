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

// GetModels returns the loaded models from the server.
//
// This method sends a GET request to the server to retrieve the list of loaded models.
//
// Returns:
//
//	*GetModelsResponse: The response containing the list of models.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Sends a GET request to the server to retrieve the models.
// 2. Checks the response status code and logs any errors.
// 3. Returns the list of models if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	modelsResponse, err := client.GetModels()
//	if err != nil {
//	    log.Fatalf("Failed to get models: %v", err)
//	}
//	fmt.Printf("Loaded models: %+v\n", modelsResponse.Models)
func (c *Client) GetModels() (*GetModelsResponse, error) {
	logger := c.logger.With("method", "GetModels")
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
		logger.Error("error getting models", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error getting models: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response models: %+v", result))
	return &result, nil
}

// EjectModelRequest is the request to eject a model
type EjectModelRequest struct {
	ModelName   string `validate:"required" json:"model_name"`
	ModelDevice string `validate:"oneof=cpu cuda" json:"model_device"`
}

// EjectModel ejects the model from the server cache.
//
// This method sends a DELETE request to the server to remove the specified model
// from the server's cache. The model to be ejected is specified in the
// EjectModelRequest parameter.
//
// Parameters:
//
//	ejectModelReq (*EjectModelRequest): The request containing the model name and device.
//
// Returns:
//
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Validates the ejectModelReq parameter.
// 2. Sends a DELETE request to the server with the model name and device as query parameters.
// 3. Checks the response status code and logs any errors.
// 4. Returns an error if the operation fails, otherwise returns nil.
//
// Example usage:
//
//	ejectModelReq := &EjectModelRequest{
//	    ModelName:   "example_model",
//	    ModelDevice: "cpu",
//	}
//	err := client.EjectModel(ejectModelReq)
//	if err != nil {
//	    log.Fatalf("Failed to eject model: %v", err)
//	}
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
