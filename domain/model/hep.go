package model

type CreateHostEndpointInput struct {
	Metadata    HostEndpointMetadataInput
	Spec        HostEndpointSpecInput
	Description string
}

type HostEndpointMetadataInput struct {
	Name   string
	Labels map[string]string
}

type HostEndpointSpecInput struct {
	InterfaceName string
	IPs           []string
	Ports         []HostEndpointSpecPortInput
}

type HostEndpointSpecPortInput struct {
	Name     string
	Port     int
	Protocol string
}
