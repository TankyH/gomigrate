package indexstruct

type Index struct {
	ModuleName string `json:"module_name"`
	FullName   string `json:"full_name"`
	Seq        int    `json:"seq"`
}

type IndexConfig []Index
