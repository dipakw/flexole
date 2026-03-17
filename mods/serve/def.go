package serve

import "io/fs"

type Serve struct {
	cfg *Config
	fs  fs.FS
	vs  *virtualServer
}
