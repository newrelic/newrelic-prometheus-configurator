package promapi

type Metadata struct {
	Status string                      `json:"status"`
	Data   map[string][]MetricMetadata `json:"data"`
}

type MetricMetadata struct {
	Type string `json:"type"`
	Help string `json:"help"`
	Unit string `json:"unit"`
}
