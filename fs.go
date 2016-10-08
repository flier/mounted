package mounted

import (
	"fmt"
	"strings"
)

type FileSystem struct {
	Name    string
	Path    string
	Type    string
	Options []string
}

func (fs *FileSystem) String() string {
	return fmt.Sprintf("%s %s %s %s", fs.Name, fs.Path, fs.Type, strings.Join(fs.Options, ","))
}
