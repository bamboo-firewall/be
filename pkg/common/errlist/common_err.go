package errlist

import (
	"github.com/bamboo-firewall/be/pkg/httpbase/ierror"
)

var (
	ErrDatabase = ierror.NewCoreError("err_database", "")

	ErrNotFoundHostEndpoint         = ierror.NewCoreError("err_not_found_host_endpoint", "")
	ErrNotFoundGlobalNetworkPolicy  = ierror.NewCoreError("err_not_found_global_network_policy", "")
	ErrNotFoundGlobalNetworkSet     = ierror.NewCoreError("err_not_found_global_network_set", "")
	ErrDuplicateHostEndpoint        = ierror.NewCoreError("err_duplicate_host_endpoint", "")
	ErrDuplicateGlobalNetworkPolicy = ierror.NewCoreError("err_duplicate_global_network_policy", "")
	ErrDuplicateGlobalNetworkSet    = ierror.NewCoreError("err_duplicate_global_network_set", "")

	ErrUnmarshalFailed = ierror.NewCoreError("err_unmarshal_failed", "")

	ErrMalformedSelector = ierror.NewCoreError("err_malformed_selector", "")
)
