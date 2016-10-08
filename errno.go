package mounted

import (
	"syscall"
)

/*
#include <errno.h>

int _errno() {
    return errno;
}
*/
import "C"

func Errno() error {
	return syscall.Errno(C._errno())
}
