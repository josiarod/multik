package output

import "encoding/json"

type ErrItem struct {
	Cluster string `json:"cluster"`
	Error   string `json:"error"`
}

type Envelope struct {
	APIVersion string      `json:"apiVersion"` //multik.io/v1
	Kind       string      `json:"kind"`
	Items      interface{} `json:"items"`
	Errors     []ErrItem   `json:"errors"`
}

func JSON(kind string, items interface{}, errs []ErrItem) ([]byte, error) {
	return json.MarshalIndent(Envelope{
		APIVersion: "multik.io/v1",
		Kind:       kind,
		Items:      items,
		Errors:     errs,
	}, "", "  ")
}
