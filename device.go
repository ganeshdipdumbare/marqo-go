package marqo

import "fmt"

// GetCPUInfoResponse is the response from the server containing CPU information.
type GetCPUInfoResponse struct {
	CPUUsagePercent   string `json:"cpu_usage_percent"`
	MemoryUsedPercent string `json:"memory_used_percent"`
	MemoryUsedGB      string `json:"memory_used_gb"`
}

// GetCPUInfo returns the CPU info from the server.
//
// This method sends a GET request to the server to retrieve the CPU information.
//
// Returns:
//
//	*GetCPUInfoResponse: The response containing the CPU information.
//	error: An error if the operation fails, otherwise nil.
//
// The function performs the following steps:
// 1. Sends a GET request to the server to retrieve the CPU information.
// 2. Checks the response status code and logs any errors.
// 3. Returns the response from the server if the operation is successful, otherwise returns an error.
//
// Example usage:
//
//	resp, err := client.GetCPUInfo()
//	if err != nil {
//	    log.Fatalf("Failed to get CPU info: %v", err)
//	}
//	fmt.Printf("GetCPUInfoResponse: %+v\n", resp)
func (c *Client) GetCPUInfo() (*GetCPUInfoResponse, error) {
	logger := c.logger.With("method", "GetCPUInfo")
	var result GetCPUInfoResponse

	resp, err := c.reqClient.
		R().
		SetSuccessResult(&result).
		Get(c.reqClient.BaseURL + "/device/cpu")
	if err != nil {
		logger.Error("error getting CPU info", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != 200 {
		logger.Error("error getting CPU info", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error getting CPU info: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response CPU info: %+v", result))
	return &result, nil
}

// GetCUDAInfoResponse is the response from the server containing CUDA information.
type GetCUDAInfoResponse struct {
	CUDADevices []CUDADevice `json:"cuda_devices"`
}

// GetCUDAInfo returns the CUDA info from the server.
// Returns the response from the server and an error if the operation fails.
func (c *Client) GetCUDAInfo() (*GetCUDAInfoResponse, error) {
	logger := c.logger.With("method", "GetCUDAInfo")
	var result GetCUDAInfoResponse

	resp, err := c.reqClient.
		R().
		SetSuccessResult(&result).
		Get(c.reqClient.BaseURL + "/device/cuda")
	if err != nil {
		logger.Error("error getting CUDA info", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != 200 {
		logger.Error("error getting CUDA info", "status_code", resp.Response.StatusCode)
		return nil, fmt.Errorf("error getting CUDA info: status code: %v", resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response CUDA info: %+v", result))
	return &result, nil
}

type CUDADevice struct {
	DeviceName  string `json:"device_name"`
	MemoryUsed  string `json:"memory_used"`
	TotalMemory string `json:"total_memory"`
}
