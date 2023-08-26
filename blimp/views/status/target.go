package status

type StatusTarget struct {
	Name string `koanf:"name"`
	Type string `koanf:"type"`
	Host string `koanf:"host"`
}

type StatusEntry struct {
	Target     StatusTarget
	Online     bool
	HTTPStatus int
}
