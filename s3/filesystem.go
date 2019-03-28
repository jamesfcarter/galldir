package s3

import (
	"net/http"
	netURL "net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	as3 "github.com/aws/aws-sdk-go/service/s3"
)

// Filesystem implements the http.FileSystem interface for an S3 bucket.
type Filesystem struct {
	s3     s3Interface
	bucket string
}

// New builds and returns a new Filesystem.
// Credentials come from the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
// environment variables.
func New(url string) *Filesystem {
	host, region, bucket := configFromURL(url)
	creds := credentials.NewEnvCredentials()
	config := aws.NewConfig().
		WithCredentials(creds).
		WithEndpoint(host).
		WithRegion(region)
	sess := session.Must(session.NewSession(config))
	return &Filesystem{
		s3:     as3.New(sess),
		bucket: bucket,
	}
}

// Open returns an http.File interface for an object or simulated directory
// within the S3 bucket.
func (f *Filesystem) Open(path string) (http.File, error) {
	return &File{fs: f, path: prepPath(path)}, nil
}

func prepPath(path string) string {
	path = strings.TrimPrefix(path, "/")
	return path
}

func configFromURL(url string) (host, region, bucket string) {
	u, err := netURL.Parse(url)
	if err != nil {
		return
	}
	host = u.Host
	bucket = strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 2)[0]
	hostBits := strings.SplitN(host, ".", 3)
	if len(hostBits) > 1 {
		region = hostBits[1]
	}
	return
}
