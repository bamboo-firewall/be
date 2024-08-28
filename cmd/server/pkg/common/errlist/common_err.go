package errlist

import "github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"

var (
	ErrDatabase = ierror.NewCoreError("err_database", "")

	ErrNotFoundHostEndpoint = ierror.NewCoreError("err_not_found_host_endpoint", "")
)
