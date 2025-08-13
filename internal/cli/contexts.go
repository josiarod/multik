package cli

import (
	"fmt"

	"github.com/josiarod/multik/internal/kube"
	"github.com/spf13/cobra"
)

func init() {
	NewRootCmd().AddCommand(contextsCmd())
}

func contextsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "contexts",
		Short: "List kubeconfig contexts detected and client availability",
		RunE: func(cmd *cobra.Command, args []string) error {
			clis, err := kube.LoadAllContexts()
			if err != nil {
				return err
			}
			if len(clis) == 0 {
				fmt.Println("No contexts found. Check KUBECONFIG or ~/.kube/config.")
				return nil
			}
			for _, c := range clis {
				fmt.Println("-", c.Name)
			}
			return nil
		},
	}
}
