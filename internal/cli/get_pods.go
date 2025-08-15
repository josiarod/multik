package cli

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/josiarod/multik/internal/engine"
	"github.com/josiarod/multik/internal/kube"
	"github.com/josiarod/multik/internal/output"
	"github.com/josiarod/multik/internal/types"
	"github.com/josiarod/multik/internal/util"
	"github.com/spf13/cobra"
)

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resources across clusters",
	}
	cmd.AddCommand(getPodsCmd())
	return cmd
}

func getPodsCmd() *cobra.Command {
	var (
		contextsCSV  string
		ns           string
		allNS        bool
		labels       string
		outFmt       string
		maxParallel  int
		timeout      time.Duration
		failOnErrors bool
	)

	cmd := &cobra.Command{
		Use:   "pods",
		Short: "List pods across selected clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := kube.LoadAllContexts()
			if err != nil {
				return err
			}
			if len(all) == 0 {
				fmt.Println("No contexts found. Check KUBECONFIG or ~/.kube/config.")
				return nil
			}

			var patterns []string
			if contextsCSV != "" {
				patterns = strings.Split(contextsCSV, ",")
			}

			// Select contexts
			selected := make(map[string]kube.ClusterClient)
			for _, c := range all {
				if util.MatchAny(c.Name, patterns) {
					selected[c.Name] = c
				}
			}
			if len(selected) == 0 {
				fmt.Println("No contexts matched. Use --context or --context-csv to specify contexts.")
				return nil
			}

			// Stable order
			var keys []string
			for k := range selected {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			opts := types.ListOpts{
				Namespace:     ns,
				AllNamespaces: allNS,
				LabelSelector: labels,
			}

			// Fan out per cluster
			type rowsT = []types.PodRow
			results := engine.FanOut[rowsT](
				context.Background(),
				keys,
				maxParallel,
				timeout,
				func(ctx context.Context, name string) ([]types.PodRow, error) {
					list, err := kube.ListPods(ctx, selected[name], opts)
					if err != nil {
						return nil, err
					}
					return kube.ToPodRows(name, list), nil
				},
			)

			rows := []types.PodRow{}
			errs := []output.ErrItem{}
			for _, r := range results {
				if r.Err != nil {
					errs = append(errs, output.ErrItem{Cluster: r.Key, Error: r.Err.Error()})
				}
				rows = append(rows, r.Item...)
			}

			switch outFmt {
			case "json":
				data, _ := output.JSON("Podlist", rows, errs)
				fmt.Println(string(data))
			default:
				fmt.Println(output.PodsTable(rows))
				if len(errs) > 0 {
					fmt.Println("\nErrors:")
					for _, e := range errs {
						fmt.Printf("- %s: %s\n", e.Cluster, e.Error)
					}
				}
			}

			if failOnErrors && len(errs) > 0 {
				return fmt.Errorf("some errors occurred")
			}
			return nil
		},
	}

	// Flags
	cmd.Flags().StringVar(&contextsCSV, "contexts", "", "CSV or glob patterns (e.g., prod-*, staging-* )")
	cmd.Flags().StringVarP(&ns, "namespace", "n", "", "Namespace (omit or use -A for all)")
	cmd.Flags().BoolVarP(&allNS, "all-namespaces", "A", false, "List all namespaces")
	cmd.Flags().StringVarP(&labels, "labels", "l", "", "Label selector (key=value[,key2=value2])")
	cmd.Flags().StringVarP(&outFmt, "output", "o", "table", "Output format (table|json)")
	cmd.Flags().IntVar(&maxParallel, "max-parallel", 8, "Max parallel cluster queries")
	cmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "Per-cluster timeout")
	cmd.Flags().BoolVar(&failOnErrors, "fail-on-error", false, "Exit with non-zero exit code if any errors occur")

	return cmd
}
