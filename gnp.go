package be

type MetaData struct {
	Annotations map[string]string
}

type GlobalNetworkPolicy struct {
	ApiVersion string
	Kind       string
	MetaData   MetaData
	Name       string
	UUID       string
}
