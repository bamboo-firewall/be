package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/pkg/httpbase"
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
		return responseBodyToIError(ctx, res)
	}

	return nil
}

func (c *apiServer) ListGNPs(ctx context.Context, input *dto.ListGNPsInput) ([]*dto.GlobalNetworkPolicy, error) {
	res := c.client.NewRequest().
		SetSubURL("/api/v1/globalNetworkPolicies").
		SetParam("isOrder", strconv.FormatBool(input.IsOrder)).
		SetMethod(http.MethodGet).
		DoRequest(ctx)

	if res.Err != nil {
		return nil, fmt.Errorf("failed to list gnp by name: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, responseBodyToIError(ctx, res)
	}

	var gnps []*dto.GlobalNetworkPolicy
	if err := json.Unmarshal(res.Body, &gnps); err != nil {
		return nil, fmt.Errorf("failed to unmarshal when list gnp, response: %s, err: %w", string(res.Body), err)
	}
	return gnps, nil
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
		return nil, responseBodyToIError(ctx, res)
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
		return responseBodyToIError(ctx, res)
	}

	return nil
}

func (c *apiServer) ValidateGlobalNetworkPolicy(ctx context.Context, input *dto.CreateGlobalNetworkPolicyInput) (*dto.ValidateGlobalNetworkPolicyOutput, error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input to validate global network policy: %w", err)
	}

	res := c.client.NewRequest().
		SetSubURL("/api/v1/globalNetworkPolicies/validate").
		SetHeader(httpbase.HeaderContentType, httpbase.MIMEApplicationJSON).
		SetBody(bytes.NewReader(inputBytes)).
		SetMethod(http.MethodPost).
		DoRequest(ctx)

	if res.Err != nil {
		return nil, fmt.Errorf("failed to validate global network policy: %w", res.Err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, responseBodyToIError(ctx, res)
	}

	var validateGlobalNetworkPolicyOutput *dto.ValidateGlobalNetworkPolicyOutput
	if err = json.Unmarshal(res.Body, &validateGlobalNetworkPolicyOutput); err != nil {
		return nil, fmt.Errorf("failed to unmarshal when validate global network policy response: %s, err: %w", string(res.Body), err)
	}

	return validateGlobalNetworkPolicyOutput, nil
}
