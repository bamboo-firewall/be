package errlist

import "github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"

var (
	ErrDatabase = ierror.NewCoreError("err_database", "")

	ErrNotFoundHostEndpoint        = ierror.NewCoreError("err_not_found_host_endpoint", "")
	ErrNotFoundGlobalNetworkPolicy = ierror.NewCoreError("err_not_found_global_network_policy", "")
	ErrNotFoundGlobalNetworkSet    = ierror.NewCoreError("err_not_found_global_network_set", "")

	ErrUnmarshalFailed = ierror.NewCoreError("err_unmarshal_failed", "")
)
