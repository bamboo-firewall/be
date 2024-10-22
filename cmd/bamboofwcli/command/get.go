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

var outputFormat string

var getCMD = &cobra.Command{
	Use:   "get",
	Short: "Get resource by name",
	Example: `  # Get a global network policy by name
  bbfwcli get gnp allow_ssh

  # Get a global network policy by name with json output format
  bbfwcli get gnp allow_ssh -o json
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := get(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	getCMD.Flags().StringVarP(&outputFormat, "output", "o", "", "output format(yaml|json). Default: yaml")
}

func get(cmd *cobra.Command, args []string) error {
	resourceType := args[0]
	resourceMgr, err := common.GetResourceMgrByType(resourceType)
	if err != nil {
		return err
	}

	resourceName := args[1]

	var input interface{}
	switch resourceMgr.GetResourceType() {
	case resouremanager.ResourceTypeHEP:
		input = &dto.GetHostEndpointInput{Name: resourceName}
	case resouremanager.ResourceTypeGNS:
		input = &dto.GetGNSInput{Name: resourceName}
	case resouremanager.ResourceTypeGNP:
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
