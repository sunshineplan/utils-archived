package service

var defaultName = "Service"

// Service represents a windows service.
type Service struct {
	Name    string
	Desc    string
	Exec    func()
	Options Options
}

// Options is Service options
type Options struct {
	Dependencies []string
	Arguments    string
	Others       []string
}

// New creates a new service name.
func New() *Service {
	return &Service{Name: defaultName}
}
