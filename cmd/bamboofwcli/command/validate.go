package command

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/wI2L/jsondiff"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command/common"
	"github.com/bamboo-firewall/be/cmd/bamboofwcli/command/resourcemanager"
	"github.com/bamboo-firewall/be/pkg/client"
	"github.com/bamboo-firewall/be/pkg/httpbase"
	"github.com/bamboo-firewall/be/pkg/httpbase/ierror"
)

var fileValidates []string

var validateCommand = &cobra.Command{
	Use:   "validate [resourceType]",
	Short: "validate resource by filename",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := validate(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	validateCommand.Flags().StringSliceVarP(&fileValidates, "file", "f", []string{}, "resource file path")
	validateCommand.MarkFlagRequired("file")
}

func validate(cmd *cobra.Command, args []string) error {
	resourceType := args[0]
	resourceMgr, err := common.GetResourceMgrByType(resourceType)
	if err != nil {
		return fmt.Errorf("get resource by type: %s", resourceType)
	}

	var resources []*common.ResourceFile
	switch resourceMgr.GetResourceType() {
	case resourcemanager.ResourceTypeHEP:
		resources, err = common.GetResourceFilesByFileNames[dto.CreateHostEndpointInput](fileValidates)
	case resourcemanager.ResourceTypeGNS:
		resources, err = common.GetResourceFilesByFileNames[dto.CreateGlobalNetworkSetInput](fileValidates)
	case resourcemanager.ResourceTypeGNP:
		resources, err = common.GetResourceFilesByFileNames[dto.CreateGlobalNetworkPolicyInput](fileValidates)
	default:
		return fmt.Errorf("invalid resource type: %s", resourceType)
	}
	if err != nil {
		return err
	}

	apiServer := client.NewAPIServer(os.Getenv(common.APIServerENV))
	for _, r := range resources {
		fmt.Printf("Validate for resource %s\n", r.Name)
		validateOutput, errValidate := resourceMgr.Validate(context.Background(), apiServer, r.FilePath, r.Content)
		if errValidate != nil {
			var ierr *ierror.Error
			if errors.As(errValidate, &ierr) {
				if ierr.Code == httpbase.ErrorCodeBadRequest {
					fmt.Printf("Resoure invalid. Detail:\n")
					var buf bytes.Buffer
					encoder := json.NewEncoder(&buf)
					encoder.SetEscapeHTML(false)
					encoder.SetIndent("", "  ")
					if errValidate = encoder.Encode(ierr.Detail); errValidate != nil {
						fmt.Printf("Error encoding detail: %v. Error: %v\n", ierr.Detail, errValidate)
					} else {
						fmt.Printf("%s\n", buf.String())
					}
					continue
				}
			}
			fmt.Printf("Fail to validate resource: %s. Error: %v\n", r.Name, errValidate)
			continue
		} else {
			fmt.Printf("Resource is valid.\n")
		}

		var jsondiffOpts = []jsondiff.Option{
			jsondiff.Ignores(
				"/id",
				"/uuid",
				"/version",
				"/createdAt",
				"/updatedAt",
			),
		}
		switch resourceMgr.GetResourceType() {
		case resourcemanager.ResourceTypeHEP:
			validateHEPOutput, ok := validateOutput.(*dto.ValidateHostEndpointOutput)
			if !ok {
				fmt.Printf("invalid validate output. Raw: %v", validateHEPOutput)
				break
			}
			if validateHEPOutput.HEPExisted != nil {
				patch, errDiff := jsondiff.Compare(
					validateHEPOutput.HEPExisted,
					validateHEPOutput.HEP,
					jsondiffOpts...,
				)
				if errDiff != nil {
					fmt.Printf("Fail to compare HEP. Error: %v\n", errDiff)
					break
				}
				if patch != nil {
					fmt.Printf("Resource will change:\n")
					if errDiff = printDiff(patch); errDiff != nil {
						fmt.Printf("Fail to print diff. Error: %v\n", errDiff)
					}
				} else {
					fmt.Printf("Resouce willn't change.\n")
				}
			} else {
				fmt.Printf("Resource doesn't exist and will be create new one.\n")
			}

			if len(validateHEPOutput.ParsedGNPs) > 0 {
				fmt.Printf("Resource will have %d global network policies:\n", len(validateHEPOutput.ParsedGNPs))
				for _, policy := range validateHEPOutput.ParsedGNPs {
					fmt.Printf("%s\n", policy.Name)
				}
			} else {
				fmt.Printf("Resource willn't have any global network policies.\n")
			}
		case resourcemanager.ResourceTypeGNP:
			validateGNPOutput, ok := validateOutput.(*dto.ValidateGlobalNetworkPolicyOutput)
			if !ok {
				fmt.Printf("invalid validate output. Raw: %v", validateGNPOutput)
				break
			}
			if validateGNPOutput.GNPExisted != nil {
				patch, errDiff := jsondiff.Compare(validateGNPOutput.GNPExisted, validateGNPOutput.GNP, jsondiffOpts...)
				if errDiff != nil {
					fmt.Printf("Fail to compare GNP. Error: %v\n", errDiff)
					break
				}
				if patch != nil {
					fmt.Printf("Resource will change:\n")
					if errDiff = printDiff(patch); errDiff != nil {
						fmt.Printf("Fail to print diff. Error: %v\n", errDiff)
					}
				} else {
					fmt.Printf("Resouce willn't change.\n")
				}
			} else {
				fmt.Printf("Resource doesn't exist and will be create new one.\n")
			}

			if len(validateGNPOutput.ParsedHEPs) > 0 {
				fmt.Printf("Resource will be match %d host endpoints:\n", len(validateGNPOutput.ParsedHEPs))
				if errValidate = printParsedHEPs(validateGNPOutput.ParsedHEPs); errValidate != nil {
					fmt.Printf("Fail to print related host endpoint. Error: %v\n", errValidate)
				}
			} else {
				fmt.Printf("Resouce willn't be match with any host endpoint.\n")
			}
		case resourcemanager.ResourceTypeGNS:
			validateGNSOutput, ok := validateOutput.(*dto.ValidateGlobalNetworkSetOutput)
			if !ok {
				fmt.Printf("invalid validate output. Raw: %v", validateGNSOutput)
				break
			}

			if validateGNSOutput.GNSExisted != nil {
				patch, errDiff := jsondiff.Compare(validateGNSOutput.GNSExisted, validateGNSOutput.GNS, jsondiffOpts...)
				if errDiff != nil {
					fmt.Printf("Fail to compare GNS. Error: %v\n", errDiff)
					break
				}
				if patch != nil {
					fmt.Printf("Resource will change:\n")
					if errDiff = printDiff(patch); errDiff != nil {
						fmt.Printf("Fail to print diff. Error: %v\n", errDiff)
					}
				} else {
					fmt.Printf("Resouce willn't change.\n")
				}
			} else {
				fmt.Printf("Resource doesn't exist and will be create new one.\n")
			}
		default:
			return fmt.Errorf("invalid resource type: %s", resourceType)
		}
		fmt.Println("--------------------------------------------------------------------")
	}
	return nil
}

func replaceSlashToDot(s string) string {
	return strings.ReplaceAll(s, "/", ".")
}

func printDiff(patch jsondiff.Patch) error {
	header := []string{"FIELD", "CURRENT_VALUE", "NEW_VALUE"}
	headerFieldValue := []string{"{{replaceSlashToDot .Path}}", "{{.OldValue}}", "{{.Value}}"}
	buf := new(bytes.Buffer)
	for _, h := range header {
		buf.WriteString(h)
		buf.WriteByte('\t')
	}
	buf.WriteByte('\n')

	buf.WriteString("{{range .}}")
	for _, h := range headerFieldValue {
		buf.WriteString(h)
		buf.WriteByte('\t')
	}
	buf.WriteByte('\n')
	buf.WriteString("{{end}}")

	tmpl, err := template.New("validate").Funcs(map[string]any{
		"replaceSlashToDot": replaceSlashToDot,
	}).Parse(buf.String())
	if err != nil {
		return fmt.Errorf("parse validate template: %w", err)
	}
	writer := tabwriter.NewWriter(os.Stdout, 5, 1, 3, ' ', 0)
	if err = tmpl.Execute(writer, patch); err != nil {
		return fmt.Errorf("execute validate template: %w", err)
	}
	writer.Flush()
	fmt.Printf("\n")
	return nil
}

func printParsedHEPs(parsedHEPs []*dto.ParsedHEP) error {
	header := []string{"TENANT_ID", "IP", "NAME"}
	headerFieldValue := []string{"{{.TenantID}}", "{{.IP}}", "{{.Name}}"}

	buf := new(bytes.Buffer)
	for _, h := range header {
		buf.WriteString(h)
		buf.WriteByte('\t')
	}
	buf.WriteByte('\n')

	buf.WriteString("{{range .}}")
	for _, h := range headerFieldValue {
		buf.WriteString(h)
		buf.WriteByte('\t')
	}
	buf.WriteByte('\n')
	buf.WriteString("{{end}}")

	tmpl, err := template.New("validate").Parse(buf.String())
	if err != nil {
		return fmt.Errorf("parse validate template: %w", err)
	}
	writer := tabwriter.NewWriter(os.Stdout, 5, 1, 3, ' ', 0)
	if err = tmpl.Execute(writer, parsedHEPs); err != nil {
		return fmt.Errorf("execute validate template: %w", err)
	}
	writer.Flush()
	fmt.Printf("\n")
	return nil
}
