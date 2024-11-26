package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase"
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
		return fmt.Errorf("unexpected status code when create hostendpoint, status code: %d, response: %s", res.StatusCode, res.Body)
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
		return nil, fmt.Errorf("unexpected status code when list hostendpoint, status code: %d, response: %s", res.StatusCode, res.Body)
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
		return nil, fmt.Errorf("unexpected status code when get hostendpoint, status code: %d, response: %s", res.StatusCode, res.Body)
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
		return fmt.Errorf("unexpected status code when delete hostendpoint, status code: %d, response: %s", res.StatusCode, res.Body)
	}

	return nil
}
