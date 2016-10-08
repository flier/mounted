package mounted

import (
	"fmt"
	"strings"
)

// the mounted file system description
type FileSystem struct {
	Name    string   // name of mounted file system
	Path    string   // file system path prefix
	Type    string   // mount type
	Options []string // mount options
}

func (fs *FileSystem) String() string {
	return fmt.Sprintf("%s on %s type %s (%s)", fs.Name, fs.Path, fs.Type, strings.Join(fs.Options, ","))
}
