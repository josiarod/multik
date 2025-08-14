package cli

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/josiarod/multik/internal/engine"
	"github.com/josiarod/multik/internal/kube"
	"github.com/josiarod/multik/internal/util"
	"github.com/spf13/cobra"
)

func pingCmd() *cobra.Command {
	var (
		contextCSV  string
		maxParallel int
		timeout     time.Duration
	)

	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Query selected clusters in parallel and print theur Kubernetes server version",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load all contexts
			all, err := kube.LoadAllContexts()
			if err != nil {
				return err
			}
			if len(all) == 0 {
				fmt.Println("No contexts found. Check KUBECONFIG or ~/.kube/config.")
				return nil
			}

			// Parse filter patterns
			var patterns []string
			if contextCSV != "" {
				patterns = strings.Split(contextCSV, ",")
			}

			// Filter selected cluster clients
			selected := make(map[string]kube.ClusterClient)
			for _, c := range all {
				if util.MatchAny(c.Name, patterns) {
					selected[c.Name] = c
				}
			}
			if len(selected) == 0 {
				fmt.Println("Not contexts matched. Use --contexts like 'prod-*,dev-*', or a CSV list.")
				return nil
			}

			// Stable order (nice UX)
			var keys []string
			for k := range selected {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			// Fan out: per selected context, get server version
			results := engine.FanOut[string](
				context.Background(),
				keys,
				maxParallel,
				timeout,
				func(ctx context.Context, name string) (string, error) {
					return kube.ServerVersion(ctx, selected[name])
				},
			)

			// Print results
			ok := 0
			fail := 0
			for _, r := range results {
				if r.Err != nil {
					fmt.Printf("%-20s ERROR: %v\n", r.Key, r.Err)
					fail++
					continue
				}
				fmt.Printf("%-20s %s\n", r.Key, r.Item)
				ok++
			}
			// Footer
			fmt.Printf("\n%d OK, %d ERRORS\n", ok, fail)
			return nil
		},
	}

	cmd.Flags().StringVar(&contextCSV, "contexts", "", "CSV or glob patterns (e.g. prod-*,dev-*)")
	cmd.Flags().IntVar(&maxParallel, "max-parallel", 8, "Maximum number of parallel queries")
	cmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "Timeout per cluster")
	return cmd
}
