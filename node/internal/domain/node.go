package domain

type Node struct {
	Name          string `json:"name,omitempty"`
	DriverVersion string `json:"driverversion,omitempty"`
	Labels        *map[string]string
}
