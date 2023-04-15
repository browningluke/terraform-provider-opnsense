package opnsense

// Data structs

type UnboundHostAlias struct {
	Enabled     string `json:"enabled"`
	Host        string `json:"host"`
	Hostname    string `json:"hostname"`
	Domain      string `json:"domain"`
	Description string `json:"description"`
}

// Response structs

type unboundHostAliasGetResp struct {
	Alias struct {
		Enabled     string `json:"enabled"`
		Hostname    string `json:"hostname"`
		Domain      string `json:"domain"`
		Description string `json:"description"`
	} `json:"alias"`
}
