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

func (c *apiServer) CreateGNS(ctx context.Context, input *dto.CreateGlobalNetworkSetInput) error {
	inputBytes, _ := json.Marshal(input)
	res := c.client.NewRequest().
		SetSubURL("/api/v1/globalNetworkSets").
		SetHeader(httpbase.HeaderContentType, httpbase.MIMEApplicationJSON).
		SetBody(bytes.NewReader(inputBytes)).
		SetMethod(http.MethodPost).
		DoRequest(ctx)

	if res.Err != nil {
		return fmt.Errorf("failed to create globalnetworkset: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return responseBodyToIError(ctx, res)
	}

	return nil
}

func (c *apiServer) ListGNSs(ctx context.Context) ([]*dto.GlobalNetworkSet, error) {
	res := c.client.NewRequest().
		SetSubURL("/api/v1/globalNetworkSets").
		SetMethod(http.MethodGet).
		DoRequest(ctx)

	if res.Err != nil {
		return nil, fmt.Errorf("failed to list gnss by name: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, responseBodyToIError(ctx, res)
	}

	var gnss []*dto.GlobalNetworkSet
	if err := json.Unmarshal(res.Body, &gnss); err != nil {
		return nil, fmt.Errorf("failed to unmarshal when list gnss, response: %s, err: %w", string(res.Body), err)
	}
	return gnss, nil
}

func (c *apiServer) GetGNS(ctx context.Context, input *dto.GetGNSInput) (*dto.GlobalNetworkSet, error) {
	res := c.client.NewRequest().
		SetSubURL(fmt.Sprintf("/api/v1/globalNetworkSets/byName/%s", input.Name)).
		SetMethod(http.MethodGet).
		DoRequest(ctx)

	if res.Err != nil {
		return nil, fmt.Errorf("failed to get globalnetworkset by name: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, responseBodyToIError(ctx, res)
	}

	var gns *dto.GlobalNetworkSet
	if err := json.Unmarshal(res.Body, &gns); err != nil {
		return nil, fmt.Errorf("failed to unmarshal when get globalnetworkset by name, response: %s, err: %w", string(res.Body), err)
	}
	return gns, nil
}

func (c *apiServer) DeleteGNS(ctx context.Context, input *dto.DeleteGlobalNetworkSetInput) error {
	inputBytes, _ := json.Marshal(input)
	res := c.client.NewRequest().
		SetSubURL("/api/v1/globalNetworkSets").
		SetHeader(httpbase.HeaderContentType, httpbase.MIMEApplicationJSON).
		SetBody(bytes.NewReader(inputBytes)).
		SetMethod(http.MethodDelete).
		DoRequest(ctx)

	if res.Err != nil {
		return fmt.Errorf("failed to delete globalnetworkset: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return responseBodyToIError(ctx, res)
	}

	return nil
}

func (c *apiServer) ValidateGlobalNetworkSet(ctx context.Context, input *dto.CreateGlobalNetworkSetInput) (*dto.ValidateGlobalNetworkSetOutput, error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input to validate global network set: %w", err)
	}

	res := c.client.NewRequest().
		SetSubURL("/api/v1/globalNetworkSets/validate").
		SetHeader(httpbase.HeaderContentType, httpbase.MIMEApplicationJSON).
		SetBody(bytes.NewReader(inputBytes)).
		SetMethod(http.MethodPost).
		DoRequest(ctx)

	if res.Err != nil {
		return nil, fmt.Errorf("failed to validate global network set: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, responseBodyToIError(ctx, res)
	}

	var validateGlobalNetworkSetOutput *dto.ValidateGlobalNetworkSetOutput
	if err = json.Unmarshal(res.Body, &validateGlobalNetworkSetOutput); err != nil {
		return nil, fmt.Errorf("failed to unmarshal when validate global network set response: %s, err: %w", string(res.Body), err)
	}

	return validateGlobalNetworkSetOutput, nil
}
