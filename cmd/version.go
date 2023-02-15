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

// versionCmd represents 