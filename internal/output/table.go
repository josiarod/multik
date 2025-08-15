package output

import (
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/josiarod/multik/internal/types"
)

func humanAge(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	h := int(d.Hours())
	if h < 24 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dd", h/24)
}

func PodsTable(rows []types.PodRow) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"CLUSTER", "NAMESPACE", "NAME", "READY", "STATUS", "RESTARTS", "AGE", "NODE", "IMAGE"})
	for _, r := range rows {
		t.AppendRow(table.Row{
			r.Cluster, r.Namespace, r.Name, r.Ready, r.Status, r.Restarts, humanAge(r.Age), r.Node, r.Image,
		})
	}
	return t.Render()
}
