package mounted

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

/*
#include <stdlib.h>
#include <Windows.h>
*/
import "C"

const FALSE = 0

func FileSystems() ([]*FileSystem, error) {
	volumeNames, err := findVolumes()

	if err != nil {
		return nil, err
	}

	var result []*FileSystem

	for _, volumeName := range volumeNames {
		deviceName, err := getDeviceName(volumeName[4 : len(volumeName)-1])

		if err != nil {
			return nil, fmt.Errorf("fail to get device name, %s", err)
		}

		pathNames, err := getVolumePathNames(volumeName)

		if err != nil {
			return nil, fmt.Errorf("fail to get volume path names, %s", err)
		}

		if len(pathNames) > 0 {
			volumeInfo, err := getVolumeInfo(pathNames[0])

			if err != nil {
				return nil, fmt.Errorf("fail to get volume information, %s", err)
			}

			result = append(result, &FileSystem{
				Name:    deviceName,
				Path:    pathNames[0],
				Type:    volumeInfo.FileSystemName,
				Options: volumeInfo.FileSystemFlags.Options(),
			})
		}
	}

	return result, nil
}

func findVolumes() ([]string, error) {
	buf := make([]byte, C.MAX_PATH)

	p := (*C.CHAR)(unsafe.Pointer(&buf[0]))
	l := C.DWORD(len(buf))

	h := C.FindFirstVolumeA(p, l)

	if h == C.HANDLE(syscall.InvalidHandle) {
		return nil, Errno()
	}

	defer C.FindVolumeClose(h)

	var volumeNames []string

	for {
		volumeName := C.GoString((*C.char)(unsafe.Pointer(&buf[0])))

		if !strings.HasPrefix(volumeName, `\\?\`) || !strings.HasSuffix(volumeName, `\`) {
			return nil, fmt.Errorf("bad volume name, %s", volumeName)
		}

		volumeNames = append(volumeNames, volumeName)

		if FALSE == C.FindNextVolume(h, p, l) {
			break
		}
	}

	return volumeNames, nil
}

func getDeviceName(volumeName string) (string, error) {
	name := C.CString(volumeName)
	defer C.free(unsafe.Pointer(name))

	buf := make([]byte, C.MAX_PATH)

	if 0 == C.QueryDosDeviceA((*C.CHAR)(name), (*C.CHAR)(unsafe.Pointer(&buf[0])), C.DWORD(len(buf))) {

		if err := Errno(); err != syscall.ERROR_NO_MORE_FILES {
			return "", fmt.Errorf("fail to get device name of %s, %s", volumeName, err)
		}
	}

	return C.GoString((*C.char)(unsafe.Pointer(&buf[0]))), nil
}

func getVolumePathNames(volumeName string) ([]string, error) {
	name := C.CString(volumeName)
	defer C.free(unsafe.Pointer(name))

	buf := make([]byte, C.MAX_PATH)
	var size C.DWORD

	for {
		if FALSE == C.GetVolumePathNamesForVolumeNameA((*C.CHAR)(name),
			(*C.CHAR)(unsafe.Pointer(&buf[0])), C.DWORD(len(buf)), &size) {

			if err := Errno(); err != syscall.ERROR_MORE_DATA {
				return nil, err
			}

			buf = make([]byte, len(buf)*2)
		} else {
			break
		}
	}

	var pathNames []string

	for {
		if buf[0] == 0 {
			break
		}

		name := C.GoString((*C.char)(unsafe.Pointer(&buf[0])))

		pathNames = append(pathNames, name)

		buf = buf[len(name)+1:]
	}

	return pathNames, nil
}

type FileSystemFlag uint

var fileSystemFlags = []struct {
	value FileSystemFlag
	name  string
}{
	{C.FILE_CASE_PRESERVED_NAMES, "preserved-case"},      // The specified volume supports preserved case of file names when it places a name on disk.
	{C.FILE_CASE_SENSITIVE_SEARCH, "case-sensitive"},     // The specified volume supports case-sensitive file names.
	{C.FILE_FILE_COMPRESSION, "file-compression"},        // The specified volume supports file-based compression.
	{C.FILE_NAMED_STREAMS, "named-streams"},              // The specified volume supports named streams.
	{C.FILE_PERSISTENT_ACLS, "acl"},                      // The specified volume preserves and enforces access control lists (ACL).
	{C.FILE_READ_ONLY_VOLUME, "ro"},                      // The specified volume is read-only.
	{C.FILE_SEQUENTIAL_WRITE_ONCE, "seq-write"},          // The specified volume supports a single sequential write.
	{C.FILE_SUPPORTS_ENCRYPTION, "encryption"},           // The specified volume supports the Encrypted File System (EFS).
	{C.FILE_SUPPORTS_EXTENDED_ATTRIBUTES, "ext-attrs"},   // The specified volume supports extended attributes.
	{C.FILE_SUPPORTS_HARD_LINKS, "hard-links"},           // The specified volume supports hard links.
	{C.FILE_SUPPORTS_OBJECT_IDS, "obj-ids"},              // The specified volume supports object identifiers.
	{C.FILE_SUPPORTS_OPEN_BY_FILE_ID, "open-by-file-id"}, // The file system supports open by FileID.
	{C.FILE_SUPPORTS_REPARSE_POINTS, "reparse-points"},   // The specified volume supports reparse points.
	{C.FILE_SUPPORTS_SPARSE_FILES, "sparse"},             // The specified volume supports sparse files.
	{C.FILE_SUPPORTS_TRANSACTIONS, "trans"},              // The specified volume supports transactions.
	{C.FILE_SUPPORTS_USN_JOURNAL, "usn"},                 // The specified volume supports update sequence number (USN) journals.
	{C.FILE_UNICODE_ON_DISK, "unicode"},                  // The specified volume supports Unicode in file names as they appear on disk.
	{C.FILE_VOLUME_IS_COMPRESSED, "compressed"},          // The specified volume is a compressed volume
	{C.FILE_VOLUME_QUOTAS, "quota"},                      //The specified volume supports disk quotas.
}

func (flags FileSystemFlag) Options() (opts []string) {
	for _, flag := range fileSystemFlags {
		if (flags & flag.value) == flag.value {
			opts = append(opts, flag.name)
		}
	}

	return
}

type volumeInfo struct {
	VolumeName             string
	VolumeSerialNumber     uint
	MaximumComponentLength uint
	FileSystemFlags        FileSystemFlag
	FileSystemName         string
}

func getVolumeInfo(volumePath string) (*volumeInfo, error) {
	path := C.CString(volumePath)
	defer C.free(unsafe.Pointer(path))

	volumeName := make([]byte, C.MAX_PATH)
	fileSystemName := make([]byte, C.MAX_PATH)

	var volumeSerialNumber, maximumComponentLength, fileSystemFlags C.DWORD

	if FALSE == C.GetVolumeInformationA((*C.CHAR)(path),
		(*C.CHAR)(unsafe.Pointer(&([]byte)(volumeName)[0])), C.DWORD(len(volumeName)),
		&volumeSerialNumber, &maximumComponentLength, &fileSystemFlags,
		(*C.CHAR)(unsafe.Pointer(&([]byte)(fileSystemName)[0])), C.DWORD(len(fileSystemName))) {
		return nil, Errno()
	}

	return &volumeInfo{
		VolumeName:             C.GoString((*C.char)(unsafe.Pointer(&volumeName[0]))),
		VolumeSerialNumber:     uint(volumeSerialNumber),
		MaximumComponentLength: uint(maximumComponentLength),
		FileSystemFlags:        FileSystemFlag(fileSystemFlags),
		FileSystemName:         C.GoString((*C.char)(unsafe.Pointer(&fileSystemName[0]))),
	}, nil
}
