package config

// Database represents information for connecting DB
type Database struct {
	Driver   string
	Hostname string
	Port     int
	Username string
	Password string
	Name     string
}
