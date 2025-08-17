package output

import (
	"fmt"
	"sort"
	"strings"

	"github.com/josiarod/multik/internal/types"
)

func PodDescribeText(items []types.PodDescribe, errs []ErrItem) string {
	var b strings.Builder

	sort.Slice(items, func(i, j int) bool {
		if items[i].Cluster != items[j].Cluster {
			return items[i].Cluster < items[j].Cluster
		}
		if items[i].Namespace != items[j].Namespace {
			return items[i].Namespace < items[j].Namespace
		}
		return items[i].Name < items[j].Name
	})

	for idx, d := range items {
		if idx > 0 {
			b.WriteString("\n---\n")
		}
		fmt.Fprintf(&b, "Cluster: %s\n", d.Cluster)
		fmt.Fprintf(&b, "Namespace: %s\n", d.Namespace)
		fmt.Fprintf(&b, "Name: %s\n", d.Name)
		fmt.Fprintf(&b, "Node: %s\n", d.Node)
		fmt.Fprintf(&b, "Status: %s (Ready %s, Restarts %d, Age%s)\n",
			d.Status, d.Ready, d.Restarts, humanAge(d.Age))
		if d.IP != "" {
			fmt.Fprintf(&b, "Pod IP:     %s\n", d.IP)
		}
		if d.QoSClass != "" {
			fmt.Fprintf(&b, "QoS:        %\n", d.QoSClass)
		}
		b.WriteString("Containers:\n")
		for _, c := range d.Containers {
			line := fmt.Sprintf("  -%s (image: %s", c.Name, c.Image)
			if len(c.Ports) > 0 {
				line += ", ports: " + strings.Join(c.Ports, ",")
			}
			line += ")"
			if c.EnvLen > 0 {
				line += fmt.Sprintf(" [env: %d]", c.EnvLen)
			}
			b.WriteString(line + "\n")
		}
	}

	if len(errs) > 0 {
		b.WriteString("\nErrors:\n")
		for _, e := range errs {
			fmt.Fprintf(&b, "- %s: %s\n", e.Cluster, e.Error)
		}
	}

	return b.String()
}
