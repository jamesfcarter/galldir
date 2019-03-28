package s3

import (
	"os"
	"syscall"
	"time"

	as3 "github.com/aws/aws-sdk-go/service/s3"
)

type fileInfo struct {
	o *as3.Object
}

type dirInfo struct {
	p *as3.CommonPrefix
}

func (fi fileInfo) Name() string {
	return *fi.o.Key
}

func (fi fileInfo) Size() int64 {
	return *fi.o.Size
}

func (fi fileInfo) Mode() os.FileMode {
	return os.FileMode(0644)
}

func (fi fileInfo) ModTime() time.Time {
	return *fi.o.LastModified
}

func (fi fileInfo) IsDir() bool {
	return false
}

func (fi fileInfo) Sys() interface{} {
	return &syscall.Stat_t{}
}

func (di dirInfo) Name() string {
	return *di.p.Prefix
}

func (di dirInfo) Size() int64 {
	return 0
}

func (di dirInfo) Mode() os.FileMode {
	return os.ModeDir
}

func (di dirInfo) ModTime() time.Time {
	return time.Now()
}

func (di dirInfo) IsDir() bool {
	return true
}

func (di dirInfo) Sys() interface{} {
	return &syscall.Stat_t{}
}
