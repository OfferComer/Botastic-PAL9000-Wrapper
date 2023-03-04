package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	versionRevision = "unknown"
	versionTime     = "unknown"
	versionModified = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version info",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]string{
			"revision": versionRevision,
			"time":     versionTime,
			"modified": versionModified,
		}
		for k, v := range m {
			fmt.Printf("%s: %s\n", k, v)
		}
	},
}

func init() {
	rootCmd.