# Overview
Get information about the mounted file systems.

# Features
 - [x] list the mounted file system
 - [x] support Linux
 - [x] support MacOS/BSD
 - [x] support Windows
 - [x] command line tools (fstab)

# Install
```bash
$ go get github.com/flier/mounted
```

# Usage
```go
fstab, err := mounted.FileSystems()

if err != nil {
    fmt.Printf("fail to get mounted file systems, %s", err)

    os.Exit(-1)
}

for _, fs := range fstab {
    fmt.Printf("%s\n", fs)
}
```

# Tools
```bash
$ go get github.com/flier/mounted/cmd/fstab
$ fstab
...
sysfs on /sys type sysfs (rw,nosuid,nodev,noexec,relatime)
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
...
```