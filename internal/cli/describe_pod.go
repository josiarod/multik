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

func describeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a resource across slected clusters.",
	}
	cmd.AddCommand(describePodCmd())
	return cmd
}

func describePodCmd() *cobra.Command {
	var (
		contextsCSV   string
		ns            string
		outFmt        string
		maxParallel   int
		timeout       time.Duration
		faileOnErrors bool
	)

	cmd := &cobra.Command{
		Use:   "pod Name",
		Short: "Describe a pod by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]

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
			selected := make(map[string]kube.ClusterClient)
			for _, c := range all {
				if util.MatchAny(c.Name, patterns) {
					selected[c.Name] = c
				}
			}

			if len(selected) == 0 {
				fmt.Println("No contexts matched. Use --contexts like 'prod-*, staging-*', or a CSV list.")
				return nil
			}

			var keys []string
			for k := range selected {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			// Fan-out: per cluster -> Get pod -> map to view
			results := engine.FanOut[types.PodDescribe](
				context.Background(),
				keys,
				maxParallel,
				timeout,
				func(ctx context.Context, cluster string) (types.PodDescribe, error) {
					p, err := kube.GetPod(ctx, selected[cluster], ns, name)
					if err != nil {
						return types.PodDescribe{}, err
					}
					return kube.ToPodDescribe(cluster, p), nil
				},
			)

			descs := make([]types.PodDescribe, 0, len(results))
			errs := []output.ErrItem{}
			for _, r := range results {
				if r.Err != nil {
					errs = append(errs, output.ErrItem{Cluster: r.Key, Error: r.Err.Error()})
					continue
				}
				descs = append(descs, r.Item)
			}

			switch outFmt {
			case "json":
				b, _ := output.JSON("PodDescribeList", descs, errs)
				fmt.Println(string(b))
			default:
				fmt.Println(output.PodDescribeText(descs, errs))
			}

			if faileOnErrors && len(errs) > 0 {
				return fmt.Errorf("one or more clusters failed")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&contextsCSV, "contexts", "", "CSV or glob patterns (e.g., prod-*, staging-*)")
	cmd.Flags().StringVarP(&ns, "namespace", "n", "default", "Namespace")
	cmd.Flags().StringVarP(&outFmt, "output", "o", "text", "Output format text|json")
	cmd.Flags().IntVar(&maxParallel, "max-parallel", 8, "Max parallel cluster queries")
	cmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "Per-cluster timeout")
	cmd.Flags().BoolVar(&faileOnErrors, "fail-on-error", false, "Exit non-zero if any cluster errors occur")

	return cmd

}
