package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/bamboo-firewall/be/buildinfo"
)

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		version := fmt.Sprintf("Version: \t %s \nBranch: \t %s\nBuild: \t %s\nOrganization: \t %s", buildinfo.Version, buildinfo.GitBranch, buildinfo.BuildDate, buildinfo.Organization)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, version)
		w.Flush()
		os.Exit(1)
	},
}
