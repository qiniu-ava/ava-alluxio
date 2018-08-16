package common

type ListOptions struct {
	CommonOptions struct {
		SyncIntervalMs int `json:"syncIntervalMs"`
	} `json:"commonOptions,omitempty"`
	LoadMetadataType string `json:"loadMetadataType"`
}

var DefaultListOptions ListOptions = ListOptions{
	LoadMetadataType: "Always",
}
