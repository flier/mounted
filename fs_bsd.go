// +build darwin dragonfly freebsd netbsd openbsd

package mounted

import (
	"unsafe"
)

/*
#include <sys/param.h>
#include <sys/ucred.h>
#include <sys/mount.h>
*/
import "C"

// the mount flags
type MountFlags int

var mountFlags = []struct {
	value MountFlags
	name  string
}{
	{C.MNT_RDONLY, "ro"},               /* read only filesystem */
	{C.MNT_SYNCHRONOUS, "sync"},        /* file system written synchronously */
	{C.MNT_NOEXEC, "noexec"},           /* can't exec from filesystem */
	{C.MNT_NOSUID, "nosuid"},           /* don't honor setuid bits on fs */
	{C.MNT_NODEV, "nodev"},             /* don't interpret special files */
	{C.MNT_UNION, "union"},             /* union with underlying filesystem */
	{C.MNT_ASYNC, "async"},             /* file system written asynchronously */
	{C.MNT_CPROTECT, "cprotect"},       /* file system supports content protection */
	{C.MNT_EXPORTED, "exported"},       /* file system is exported */
	{C.MNT_QUARANTINE, "quarantined"},  /* file system is quarantined */
	{C.MNT_LOCAL, "local"},             /* filesystem is stored locally */
	{C.MNT_QUOTA, "quota"},             /* quotas are enabled on filesystem */
	{C.MNT_ROOTFS, "rootfs"},           /* identifies the root filesystem */
	{C.MNT_DOVOLFS, "volfs"},           /* FS supports volfs (deprecated flag in Mac OS X 10.5) */
	{C.MNT_DONTBROWSE, "nobrowse"},     /* file system is not appropriate path to user data */
	{C.MNT_IGNORE_OWNERSHIP, ""},       /* VFS will ignore ownership information on filesystem objects */
	{C.MNT_AUTOMOUNTED, "automounted"}, /* filesystem was mounted by automounter */
	{C.MNT_JOURNALED, "journaled"},     /* filesystem is journaled */
	{C.MNT_NOUSERXATTR, "nouserxattr"}, /* Don't allow user extended attributes */
	{C.MNT_DEFWRITE, "deferwrite"},     /* filesystem should defer writes */
	{C.MNT_MULTILABEL, "multilabel"},   /* MAC support for individual labels */
	{C.MNT_NOATIME, "noatime"},         /* disable update of file access time */
}

func (flags MountFlags) Options() (opts []string) {
	for _, flag := range mountFlags {
		if (flags & flag.value) == flag.value {
			opts = append(opts, flag.name)
		}
	}

	return
}

// get information about mounted file systems
func FileSystems() ([]*FileSystem, error) {
	var p *C.struct_statfs

	ret := C.getmntinfo(&p, C.MNT_NOWAIT)

	if ret == 0 {
		return nil, Errno()
	}

	entries := (*[1 << 30]C.struct_statfs)(unsafe.Pointer(p))[:ret:ret]

	var result []*FileSystem

	for _, entry := range entries {
		result = append(result, &FileSystem{
			Name:    C.GoString(&entry.f_mntfromname[0]),
			Path:    C.GoString(&entry.f_mntonname[0]),
			Type:    C.GoString(&entry.f_fstypename[0]),
			Options: MountFlags(entry.f_flags).Options(),
		})
	}

	return result, nil
}
