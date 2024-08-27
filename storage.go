package be

import "context"

type Storage interface {
	AddHostEndpoint(ctx context.Context, hostEndpointName string, ip []string, networkInterface []string) error
	DelHostEndpoint(ctx context.Context, hostEndpointName string) error
	AddGroupPolicy(ctx context.Context, groupPolicyName string, groupPolicy string) error
	DeleteGroupPolicy(ctx context.Context, groupPolicyName string)
}
