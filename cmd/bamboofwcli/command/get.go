package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command/common"
	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command/resouremanager"
	"github.com/bamboo-firewall/be/pkg/client"
)

var (
	getHEPByTenantID uint64
	getHEPByIP       string
	outputFormat     string
)

var getCMD = &cobra.Command{
	Use:   "get",
	Short: "Get resource",
	Example: `  # Get a global network policy by name
  bbfw get gnp allow_ssh

  # Get a global network policy by name with json output format
  bbfw get gnp allow_ssh -o json

 # Get a host endpoint
  bbfw get hep --tenantID=1 --ip=192.168.123.0

  # Get a global network set by name
  bbfw get gns allow_ssh

  # Get a global network set by name with json output format
  bbfw get gns my_set -o json
`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := get(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	getCMD.Flags().Uint64Var(&getHEPByTenantID, "tenantID", 0, "HEP: get by tenantID")
	getCMD.Flags().StringVar(&getHEPByIP, "ip", "", "HEP: get by ip")
	getCMD.Flags().StringVarP(&outputFormat, "output", "o", "", "output format(yaml|json). Default: yaml")
}

func get(cmd *cobra.Command, args []string) error {
	resourceType := args[0]
	resourceMgr, err := common.GetResourceMgrByType(resourceType)
	if err != nil {
		return err
	}

	var resourceName string
	if len(args) > 1 {
		resourceName = args[1]
	}

	var input interface{}
	switch resourceMgr.GetResourceType() {
	case resouremanager.ResourceTypeHEP:
		if getHEPByTenantID == 0 || getHEPByIP == "" {
			return fmt.Errorf("get HEP by tenantID or ip is required")
		}
		input = &dto.GetHostEndpointInput{
			TenantID: getHEPByTenantID,
			IP:       getHEPByIP,
		}
	case resouremanager.ResourceTypeGNS:
		if resourceName == "" {
			return fmt.Errorf("no resource name provided")
		}
		input = &dto.GetGNSInput{Name: resourceName}
	case resouremanager.ResourceTypeGNP:
		if resourceName == "" {
			return fmt.Errorf("no resource name provided")
		}
		input = &dto.GetGNPInput{Name: resourceName}
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	apiServer := client.NewAPIServer(os.Getenv(common.APIServerENV))

	resource, err := resourceMgr.Get(context.Background(), apiServer, input)
	if err != nil {
		return fmt.Errorf("get resource by name %s failed: %v", resourceName, err)
	}

	var output []byte
	switch common.FileExtension(outputFormat) {
	case common.FileExtensionJSON:
		output, err = json.MarshalIndent(resource, "", "  ")
	default:
		var buf bytes.Buffer
		yamlEncoder := yaml.NewEncoder(&buf)
		yamlEncoder.SetIndent(2)
		err = yamlEncoder.Encode(resource)
		output = buf.Bytes()
	}
	if err != nil {
		return fmt.Errorf("fail to marshal resource. Error: %v\n", err)
	}
	fmt.Printf("%s\n", string(output))
	return nil
}
