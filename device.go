package marqo

import "fmt"

// GetCPUInfoResponse is the response from the server
type GetCPUInfoResponse struct {
	CPUUsagePercent   string `json:"cpu_usage_percent"`
	MemoryUsedPercent string `json:"memory_used_percent"`
	MemoryUsedGB      string `json:"memory_used_gb"`
}

// GetCPUInfo returns the cpu info from the server
func (c *Client) GetCPUInfo() (*GetCPUInfoResponse, error) {
	logger := c.logger.With("method", "GetCPUInfo")
	var result GetCPUInfoResponse

	resp, err := c.reqClient.
		R().
		SetSuccessResult(&result).
		Get(c.reqClient.BaseURL + "/device/cuda")
	if err != nil {
		logger.Error("error getting cpu info", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != 200 {
		logger.Error("error getting cpu info", "status_code",
			resp.Response.StatusCode)
		return nil, fmt.Errorf(
			"error getting cpu info: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf(
		"response cpu info: %+v\n", result))
	return &result, nil
}

type GetCUDAInfoResponse struct {
	CUDADevices []CUDADevice `json:"cuda_devices"`
}

type CUDADevice struct {
	DeviceName  string `json:"device_name"`
	MemoryUsed  string `json:"memory_used"`
	TotalMemory string `json:"total_memory"`
}

// GetCUDAInfo returns the cuda info from the server
func (c *Client) GetCUDAInfo() (*GetCUDAInfoResponse, error) {
	logger := c.logger.With("method", "GetCUDAInfo")
	var result GetCUDAInfoResponse

	resp, err := c.reqClient.
		R().
		SetSuccessResult(&result).
		Get(c.reqClient.BaseURL + "/device/cuda")
	if err != nil {
		logger.Error("error getting cuda info", "error", err)
		return nil, err
	}
	if resp.Response.StatusCode != 200 {
		logger.Error("error getting cuda info", "status_code",
			resp.Response.StatusCode)
		return nil, fmt.Errorf(
			"error getting cuda info: status code: %v",
			resp.Response.StatusCode)
	}

	logger.Info(fmt.Sprintf("response cuda info: %+v\n",
		result))
	return &result, nil
}
