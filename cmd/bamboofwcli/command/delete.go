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

var (
	deleteHEPByTenantID uint64
	deleteHEPByIP       string
	fileDeletes         []string
)

var deleteCMD = &cobra.Command{
	Use:   "delete [resourceType]",
	Short: "Delete resources",
	Long: `The delete command is used to delete resources by name(Global Network Policy, Global Network Set),
by tenantID,IP(Host Endpoint) or filename. 

  Resource type available:
    * HostEndpoint(or hep)
    * GlobalNetworkSet(or gns)
    * GlobalNetworkPolicy(or gnp)`,
	Example: `  # Delete a policy with name
  bbfw delete gnp allow_ssh

  # Delete many policy with name
  bbfw delete gnp allow_ssh allow_ping

  # Delete many policy with filename
  bbfw delete gnp -f allow_ssh.yaml allow_ping.yaml

  # Delete a set with name
  bbfw delete gns server

  # Delete many policy with name
  bbfw delete gns server vm

  # Delete many policy with filename
  bbfw delete gns -f server.yaml vm.yaml

  # Delete a hep with tenantID and ip
  bbfw delete hep --tenantID=1 --ip=192.168.1.1

  # Delete many hep with filename
  bbfw delete hep -f server.yaml vm.yaml
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := deleteResources(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	deleteCMD.Flags().Uint64Var(&deleteHEPByTenantID, "tenantID", 0, "HEP: get by tenantID")
	deleteCMD.Flags().StringVar(&deleteHEPByIP, "ip", "", "HEP: get by ip")
	deleteCMD.Flags().StringArrayVarP(&fileDeletes, "file", "f", []string{}, "file to read")
}

func deleteResources(cmd *cobra.Command, args []string) error {
	resourceType := args[0]
	resourceMgr, err := common.GetResourceMgrByType(resourceType)
	if err != nil {
		return err
	}
	var resourcesName []string
	if resourceMgr.GetResourceType() == resouremanager.ResourceTypeHEP {
		if deleteHEPByTenantID > 0 && deleteHEPByIP != "" && len(fileDeletes) > 0 {
			return fmt.Errorf("cannot use tenantID, ip with file param together")
		} else if (deleteHEPByTenantID == 0 || deleteHEPByIP == "") && len(fileDeletes) == 0 {
			return fmt.Errorf("must specify tenantID, IP or file to delete")
		}
	} else {
		if len(args) > 1 {
			resourcesName = args[1:]
		}
		if len(resourcesName) > 0 && len(fileDeletes) > 0 {
			return fmt.Errorf("cannot use name resource with file param together")
		} else if len(resourcesName) == 0 && len(fileDeletes) == 0 {
			return fmt.Errorf("must specify at least one resource to delete")
		}
	}

	var resources []*common.ResourceFile
	if len(fileDeletes) > 0 {
		switch resourceMgr.GetResourceType() {
		case resouremanager.ResourceTypeHEP:
			resources, err = common.GetResourceFilesByFileNames[dto.DeleteHostEndpointInput](fileDeletes)
		case resouremanager.ResourceTypeGNS:
			resources, err = common.GetResourceFilesByFileNames[dto.DeleteGlobalNetworkSetInput](fileDeletes)
		case resouremanager.ResourceTypeGNP:
			resources, err = common.GetResourceFilesByFileNames[dto.DeleteGlobalNetworkPolicyInput](fileDeletes)
		default:
			return fmt.Errorf("unsupported resource type: %s", resourceType)
		}
		if err != nil {
			return err
		}
	} else {
		for _, name := range resourcesName {
			switch resourceMgr.GetResourceType() {
			case resouremanager.ResourceTypeHEP:
			case resouremanager.ResourceTypeGNS:
				resources = append(resources, &common.ResourceFile{
					Name: name,
					Content: &dto.DeleteGlobalNetworkSetInput{
						Metadata: dto.GNSMetadataInput{
							Name: name,
						},
					},
				})
			case resouremanager.ResourceTypeGNP:
				resources = append(resources, &common.ResourceFile{
					Name: name,
					Content: &dto.DeleteGlobalNetworkPolicyInput{
						Metadata: dto.GNPMetadataInput{
							Name: name,
						},
					},
				})
			default:
				return fmt.Errorf("unsupported resource type: %s", resourceType)
			}
		}

		if resourceMgr.GetResourceType() == resouremanager.ResourceTypeHEP {
			resources = append(resources, &common.ResourceFile{
				Name: fmt.Sprintf("%d_%s", deleteHEPByTenantID, deleteHEPByIP),
				Content: &dto.DeleteHostEndpointInput{
					Spec: dto.HostEndpointSpecInput{
						TenantID: deleteHEPByTenantID,
						IP:       deleteHEPByIP,
						IPs:      []string{deleteHEPByIP},
					},
				},
			})
		}
	}

	apiServer := client.NewAPIServer(os.Getenv(common.APIServerENV))
	var numHandled int
	for _, r := range resources {
		err = resourceMgr.Delete(context.Background(), apiServer, r.Content)
		if err != nil {
			fmt.Printf("fail to delete resource from: %v\n", err)
		} else {
			fmt.Printf("successsfully deleted resource from %s\n", r.Name)
			numHandled++
		}
	}

	fmt.Printf("Total: %d resources. Success: %d. Fail: %d.\n", len(resources), numHandled, len(resources)-numHandled)
	return nil
}
