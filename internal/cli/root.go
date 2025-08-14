package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev" // will be overridden by ldflags

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multik",
		Short: "Multi-cluster Kubernetes CLI",
		Long:  "multik: query multiple kubernetes clusters in parallel and agregate results.",
	}
	cmd.AddCommand(versionCmd())
	cmd.AddCommand(contextsCmd())
	cmd.AddCommand(pingCmd())

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		Run: func(*cobra.Command, []string) {
			fmt.Println("multik", version)
		},
	}
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
