package command

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command/common"
	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command/resourcemanager"
	"github.com/bamboo-firewall/be/pkg/client"
)

var (
	ListHEPsByTenantID uint64
	ListHEPsByIP       string

	ListGNPsByIsOrder bool
)

var listCMD = &cobra.Command{
	Use:   "list",
	Short: "List resource",
	Example: `  # List global network sets
  bbfw list gns

  # List global network policy
  bbfw list gnp

  # List global network policy with order
  bbfw list gnp --isOrder

  # List host endpoint
  bbfw list hep

  # List host endpoint with tenantID
  bbfw list hep --tenantID=1

  # List host endpoint with IP
  bbfw list hep --ip=192.168.0.1

  # List host endpoint with tenantID and IP
  bbfw list hep --tenantID=1 --ip=192.168.0.1,
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := list(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	listCMD.Flags().Uint64Var(&ListHEPsByTenantID, "tenantID", 0, "Host Endpoint: filter by TenantID")
	listCMD.Flags().StringVar(&ListHEPsByIP, "ip", "", "Host Endpoint: filter by IP")
	listCMD.Flags().BoolVar(&ListGNPsByIsOrder, "isOrder", false, "Global Network Policy: filter by Order")
}

func list(cmd *cobra.Command, args []string) error {
	resourceType := args[0]
	resourceMgr, err := common.GetResourceMgrByType(resourceType)
	if err != nil {
		return err
	}

	var input interface{}
	switch resourceMgr.GetResourceType() {
	case resourcemanager.ResourceTypeHEP:
		listHEPsInput := &dto.ListHostEndpointsInput{}
		if ListHEPsByTenantID > 0 {
			listHEPsInput.TenantID = &ListHEPsByTenantID
		}
		if ListHEPsByIP != "" {
			listHEPsInput.IP = &ListHEPsByIP
		}
		input = listHEPsInput
	case resourcemanager.ResourceTypeGNS:
	case resourcemanager.ResourceTypeGNP:
		input = &dto.ListGNPsInput{IsOrder: ListGNPsByIsOrder}
	default:
		return fmt.Errorf("unsupported resources type: %s", resourceType)
	}

	apiServer := client.NewAPIServer(os.Getenv(common.APIServerENV))

	resources, err := resourceMgr.List(context.Background(), apiServer, input)
	if err != nil {
		return fmt.Errorf("list resources failed: %w", err)
	}

	if err = printResources(resourceMgr, resources); err != nil {
		return err
	}
	return nil
}

func printResources(resourceMgr resourcemanager.Resource, resources interface{}) error {
	header := resourceMgr.GetHeader()
	headerMap := resourceMgr.GetHeaderMap()

	buf := new(bytes.Buffer)
	for _, h := range header {
		buf.WriteString(h)
		buf.WriteByte('\t')
	}
	buf.WriteByte('\n')

	buf.WriteString("{{range .}}")

	for _, h := range header {
		value, ok := headerMap[h]
		if !ok {
			continue
		}
		buf.WriteString(value)
		buf.WriteByte('\t')
	}
	buf.WriteByte('\n')

	buf.WriteString("{{end}}")

	tmpl, err := template.New("list").Parse(buf.String())
	if err != nil {
		return fmt.Errorf("parse template failed: %w", err)
	}
	writer := tabwriter.NewWriter(os.Stdout, 5, 1, 3, ' ', 0)
	err = tmpl.Execute(writer, resources)
	if err != nil {
		return fmt.Errorf("execute template failed: %w", err)
	}
	writer.Flush()
	fmt.Printf("\n")
	return nil
}
