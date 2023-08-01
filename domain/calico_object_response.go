package domain

type CalicoObjectResponse struct {
	Kind       string      `json:"kind"`
	ApiVersion string      `json:"apiVersion"`
	Metadata   interface{} `json:"metadata,omitempty"`
	Spec       interface{} `json:"spec,omitempty"`
}
