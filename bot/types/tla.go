package types

type Queries map[string]string

type TLAOP struct {
	OperationName string `json:"operationName"`
	Query         string `json:"query"`
	Variables     any    `json:"variables"`
}

type TLAInputExtensions struct {
	PersistedQuery struct {
		Version    int 	  `json:"version"`
		Sha256Hash string `json:"sha256Hash"`
	} `json:"persistedQuery"`
}

type TLAGenericRes struct {
	Data       map[string]interface{} `json:"data"`
	Errors     []TLAErrors            `json:"errors"`
	Extensions TLAExtensions          `json:"extensions"`
}

type TLAErrors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type TLAExtensions struct {
	DurationMilliseconds int    `json:"durationMilliseconds"`
	OperationName        string `json:"operationName"`
	RequestID            string `json:"requestID"`
}

type TLAUserRes struct {
	Data struct {
		User TwitchUser `json:"user"`
	} `json:"data"`
	*TLAGenericRes
}

type TLAUserOrErrorRes struct {
	Data struct {
		User TwitchUser `json:"userResultByLogin"`
	} `json:"data"`
	*TLAGenericRes
}