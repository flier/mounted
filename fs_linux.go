// +build linux

package mounted

import (
	"strings"
	"unsafe"
)

/*
#include <stdio.h>
#include <stdlib.h>
#include <mntent.h>
*/
import "C"

func FileSystems() ([]*FileSystem, error) {
	filename := C.CString("/proc/mounts")

	defer C.free(unsafe.Pointer(filename))

	mode := C.CString("r")

	defer C.free(unsafe.Pointer(mode))

	f := C.setmntent(filename, mode)

	if f == nil {
		return nil, Errno()
	}

	var result []*FileSystem

	for {
		entry := C.getmntent(f)

		if entry == nil {
			break
		}

		result = append(result, &FileSystem{
			Name:    C.GoString(entry.mnt_fsname),
			Path:    C.GoString(entry.mnt_dir),
			Type:    C.GoString(entry.mnt_type),
			Options: strings.Split(C.GoString(entry.mnt_opts), ","),
		})
	}

	return result, nil
}
