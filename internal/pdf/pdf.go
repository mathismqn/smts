package pdf

type PDF struct {
	path   string
	buffer []byte
}

func New(path string) *PDF {
	return &PDF{path: path}
}
