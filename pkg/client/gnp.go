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

func (c *apiServer) CreateGNP(ctx context.Context, input *dto.CreateGlobalNetworkPolicyInput) error {
	inputBytes, _ := json.Marshal(input)
	res := c.client.NewRequest().
		SetSubURL("/api/v1/globalNetworkPolicies").
		SetHeader(httpbase.HeaderContentType, httpbase.MIMEApplicationJSON).
		SetBody(bytes.NewReader(inputBytes)).
		SetMethod(http.MethodPost).
		DoRequest(ctx)

	if res.Err != nil {
		return fmt.Errorf("failed to create globalnetworkpolicy: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code when create globalnetworkpolicy, status code: %d, response: %s", res.StatusCode, res.Body)
	}

	return nil
}

func (c *apiServer) GetGNP(ctx context.Context, input *dto.GetGNPInput) (*dto.GlobalNetworkPolicy, error) {
	res := c.client.NewRequest().
		SetSubURL(fmt.Sprintf("/api/v1/globalNetworkPolicies/byName/%s", input.Name)).
		SetMethod(http.MethodGet).
		DoRequest(ctx)

	if res.Err != nil {
		return nil, fmt.Errorf("failed to get globalnetworkpolicy by name: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code when get globalnetworkpolicy by name, status code: %d, response: %s", res.StatusCode, res.Body)
	}

	var gnp *dto.GlobalNetworkPolicy
	if err := json.Unmarshal(res.Body, &gnp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal when get globalnetworkpolicy by name, response: %s, err: %w", string(res.Body), err)
	}
	return gnp, nil
}

func (c *apiServer) DeleteGNP(ctx context.Context, input *dto.DeleteGlobalNetworkPolicyInput) error {
	inputBytes, _ := json.Marshal(input)
	res := c.client.NewRequest().
		SetSubURL("/api/v1/globalNetworkPolicies").
		SetHeader(httpbase.HeaderContentType, httpbase.MIMEApplicationJSON).
		SetBody(bytes.NewReader(inputBytes)).
		SetMethod(http.MethodDelete).
		DoRequest(ctx)

	if res.Err != nil {
		return fmt.Errorf("failed to delete globalnetworkpolicy: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code when delete globalnetworkpolicy, status code: %d, response: %s", res.StatusCode, res.Body)
	}

	return nil
}
