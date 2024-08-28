package dto

import "time"

type HostEndpoint struct {
	ID          string               `json:"id"`
	UUID        string               `json:"uuid"`
	Version     uint                 `json:"version"`
	Metadata    HostEndpointMetadata `json:"metadata" yaml:"metadata"`
	Spec        HostEndpointSpec     `json:"spec" yaml:"spec"`
	Description string               `json:"description" yaml:"description"`
	CreatedAt   time.Time            `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" yaml:"updated_at"`
}

type HostEndpointMetadata struct {
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type HostEndpointSpec struct {
	InterfaceName string                 `json:"interfaceName" yaml:"interfaceName"`
	IPs           []string               `json:"ips" yaml:"ips"`
	Ports         []HostEndpointSpecPort `json:"ports" yaml:"ports"`
}

type HostEndpointSpecPort struct {
	Name     string `json:"name" yaml:"name"`
	Port     int    `json:"port" yaml:"port"`
	Protocol string `json:"protocol" yaml:"protocol"`
}

type CreateHostEndpointInput struct {
	Metadata    HostEndpointMetadataInput `json:"metadata" yaml:"metadata" validate:"required"`
	Spec        HostEndpointSpecInput     `json:"spec" yaml:"spec"`
	Description string                    `json:"description" yaml:"description"`
}

type HostEndpointMetadataInput struct {
	Name   string            `json:"name" yaml:"name" validate:"required"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type HostEndpointSpecInput struct {
	InterfaceName string                      `json:"interfaceName" yaml:"interfaceName"`
	IPs           []string                    `json:"ips" yaml:"ips"`
	Ports         []HostEndpointSpecPortInput `json:"ports" yaml:"ports"`
}

type HostEndpointSpecPortInput struct {
	Name     string `json:"name" yaml:"name"`
	Port     int    `json:"port" yaml:"port"`
	Protocol string `json:"protocol" yaml:"protocol"`
}

type DeleteHostEndpointInput struct {
	Metadata HostEndpointMetadataInput `json:"metadata" yaml:"metadata" validate:"required"`
}
