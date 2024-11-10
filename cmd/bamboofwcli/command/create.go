package command

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command/common"
	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command/resouremanager"
	"github.com/bamboo-firewall/be/pkg/client"
)

var fileCreates []string

var createCMD = &cobra.Command{
	Use:   "create [resourceType]",
	Short: "Create resources by filename",
	Long: `The create command is used to create resources by filename.

  Resource type available:
    * HostEndpoint(or hep)
    * GlobalNetworkSet(or gns)
    * GlobalNetworkPolicy(or gnp)`,
	Example: `  # Create a global network policy
  bbfwcli create gnp -f policy.yaml

  # Create many global network policy
  bbfwcli create gnp -f policy1.yaml policy2.yaml`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := create(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	createCMD.Flags().StringArrayVarP(&fileCreates, "file", "f", []string{}, "file to read")
	createCMD.MarkFlagRequired("file")
}

func create(cmd *cobra.Command, args []string) error {
	resourceType := args[0]
	resourceMgr, err := common.GetResourceMgrByType(resourceType)
	if err != nil {
		return err
	}

	var resources []*common.ResourceFile
	switch resourceMgr.GetResourceType() {
	case resouremanager.ResourceTypeHEP:
		resources, err = common.GetResourceFilesByFileNames[dto.CreateHostEndpointInput](fileCreates)
	case resouremanager.ResourceTypeGNS:
		resources, err = common.GetResourceFilesByFileNames[dto.CreateGlobalNetworkSetInput](fileCreates)
	case resouremanager.ResourceTypeGNP:
		resources, err = common.GetResourceFilesByFileNames[dto.CreateGlobalNetworkPolicyInput](fileCreates)
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	if err != nil {
		return err
	}

	apiServer := client.NewAPIServer(os.Getenv(common.APIServerENV))
	var numHandled int
	for _, r := range resources {
		err = resourceMgr.Create(context.Background(), apiServer, r.Content)
		if err != nil {
			fmt.Printf("Fail to create resource. Error: %v\n", err)
		} else {
			fmt.Printf("Successsfully created resource from %s\n", r.Name)
			numHandled++
		}
	}

	fmt.Printf("Total: %d resources. Success: %d. Fail: %d.\n", len(resources), numHandled, len(resources)-numHandled)
	return nil
}
