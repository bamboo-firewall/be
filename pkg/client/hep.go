package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/pkg/httpbase"
)

func (c *apiServer) CreateHEP(ctx context.Context, input *dto.CreateHostEndpointInput) error {
	inputBytes, _ := json.Marshal(input)
	res := c.client.NewRequest().
		SetSubURL("/api/v1/hostEndpoints").
		SetHeader(httpbase.HeaderContentType, httpbase.MIMEApplicationJSON).
		SetBody(bytes.NewReader(inputBytes)).
		SetMethod(http.MethodPost).
		DoRequest(ctx)

	if res.Err != nil {
		return fmt.Errorf("failed to create hostendpoint: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return responseBodyToIError(ctx, res)
	}

	return nil
}

func (c *apiServer) ListHEPs(ctx context.Context, input *dto.ListHostEndpointsInput) ([]*dto.HostEndpoint, error) {
	params := make(map[string]string)
	if input != nil {
		if input.TenantID != nil {
			params["tenantID"] = fmt.Sprint(*input.TenantID)
		}
		if input.IP != nil {
			params["ip"] = *input.IP
		}
	}
	res := c.client.NewRequest().
		SetSubURL("/api/v1/hostEndpoints").
		SetParams(params).
		SetMethod(http.MethodGet).
		DoRequest(ctx)

	if res.Err != nil {
		return nil, fmt.Errorf("failed to list hostendpoint: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, responseBodyToIError(ctx, res)
	}

	var heps []*dto.HostEndpoint
	if err := json.Unmarshal(res.Body, &heps); err != nil {
		return nil, fmt.Errorf("failed to unmarshal when list hostendpoint, response: %s, err: %w", string(res.Body), err)
	}
	return heps, nil
}

func (c *apiServer) GetHEP(ctx context.Context, input *dto.GetHostEndpointInput) (*dto.HostEndpoint, error) {
	res := c.client.NewRequest().
		SetSubURL(fmt.Sprintf("/api/v1/hostEndpoints/byTenantID/%d/byIP/%s", input.TenantID, input.IP)).
		SetMethod(http.MethodGet).
		DoRequest(ctx)

	if res.Err != nil {
		return nil, fmt.Errorf("failed to get hostendpoint: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, responseBodyToIError(ctx, res)
	}

	var hep *dto.HostEndpoint
	if err := json.Unmarshal(res.Body, &hep); err != nil {
		return nil, fmt.Errorf("failed to unmarshal when get hostendpoint, response: %s, err: %w", string(res.Body), err)
	}
	return hep, nil
}

func (c *apiServer) DeleteHEP(ctx context.Context, input *dto.DeleteHostEndpointInput) error {
	inputBytes, _ := json.Marshal(input)
	res := c.client.NewRequest().
		SetSubURL("/api/v1/hostEndpoints").
		SetHeader(httpbase.HeaderContentType, httpbase.MIMEApplicationJSON).
		SetBody(bytes.NewReader(inputBytes)).
		SetMethod(http.MethodDelete).
		DoRequest(ctx)

	if res.Err != nil {
		return fmt.Errorf("failed to delete hostendpoint: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return responseBodyToIError(ctx, res)
	}

	return nil
}

func (c *apiServer) ValidateHostEndpoint(ctx context.Context, input *dto.CreateHostEndpointInput) (*dto.ValidateHostEndpointOutput, error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input to validate hostendpoint: %w", err)
	}

	res := c.client.NewRequest().
		SetSubURL("/api/v1/hostEndpoints/validate").
		SetHeader(httpbase.HeaderContentType, httpbase.MIMEApplicationJSON).
		SetBody(bytes.NewReader(inputBytes)).
		SetMethod(http.MethodPost).
		DoRequest(ctx)

	if res.Err != nil {
		return nil, fmt.Errorf("failed to validate hostendpoint: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, responseBodyToIError(ctx, res)
	}

	var validateHostEndpointOutput *dto.ValidateHostEndpointOutput
	if err = json.Unmarshal(res.Body, &validateHostEndpointOutput); err != nil {
		return nil, fmt.Errorf("failed to unmarshal when validate hostendpoint response: %s, err: %w", string(res.Body), err)
	}

	return validateHostEndpointOutput, nil
}
