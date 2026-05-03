package consoledto

type ColumnDTO struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	Sortable bool   `json:"sortable"`
	Visible  bool   `json:"visible"`
}

type ActionDTO struct {
	Key     string `json:"key"`
	Label   string `json:"label"`
	Scope   string `json:"scope"`
	Variant string `json:"variant"`
	Danger  bool   `json:"danger"`
}

type FilterDTO struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type ResourceDTO struct {
	Key     string      `json:"key"`
	Label   string      `json:"label"`
	Route   string      `json:"route"`
	Icon    string      `json:"icon"`
	Columns []ColumnDTO `json:"columns"`
	Actions []ActionDTO `json:"actions"`
	Filters []FilterDTO `json:"filters"`
}

type ServiceManifestDTO struct {
	Key         string `json:"key" example:"auth"`
	Name        string `json:"name" example:"Auth Service"`
	Description string `json:"description" example:"Auth Service"`
	Status      string `json:"status" example:"active"`
	Plan        string `json:"plan" example:"pro"`
	Icon        string `json:"icon" example:"auth"`
}

type ServiceManifestWithResourcesDTO struct {
	ServiceManifestDTO
	Resources []ResourceDTO `json:"resources"`
}
