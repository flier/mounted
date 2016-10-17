package mounted

import (
	"syscall"
)

/* 
#include <Windows.h>
*/
import "C"

func Errno() error {
	return syscall.Errno(C.GetLastError())
}