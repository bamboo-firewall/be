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
		return fmt.Errorf("unexpected status code when create globalnetworkset, status code: %d, response: %s", res.StatusCode, res.Body)
	}

	return nil
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
		return nil, fmt.Errorf("unexpected status code when get globalnetworkset by name, status code: %d, response: %s", res.StatusCode, res.Body)
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
		return fmt.Errorf("unexpected status code when delete globalnetworkset, status code: %d, response: %s", res.StatusCode, res.Body)
	}

	return nil
}
