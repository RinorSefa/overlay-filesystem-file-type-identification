# Instructions for installation

Tested environment requirement
- Raspberry PI 4 B Rev 1.4
- Debian GNU/Linux 11
- go version go1.18.4 linux/arm64
- $Home/Desktop directory should exist (Mount Point)

Navigate to the filesystem application for the file type and start the application.

```
$ cd cmd/file-type-filesystem
$ go run main.go
```

## Author of code in files 

- cmd/file-type-filesystem/main.go (author)
- fs/loopback.go (modified rename and create method)
- fuse/opcode.go (modified doRead method)
- fuse/read.go (modified Bytes method)