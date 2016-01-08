package cfbackup

import (
	"io"
	"strings"

	"github.com/pivotalservices/gtils/storage"
)

// S3Provider is a storage provider that allows backups
// to be stored to an S3 compatible blobstore
type S3Provider struct {
	S3Domain        string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
}

// NewS3Provider creates a new instance of the S3 storage provider
func NewS3Provider(domain, key, secret, bucket string) StorageProvider {
	return &S3Provider{
		AccessKeyID:     key,
		SecretAccessKey: secret,
		BucketName:      bucket,
		S3Domain:        domain,
	}
}

// Writer for writing to an S3 bucket
func (s *S3Provider) Writer(path ...string) (io.WriteCloser, error) {
	s3FilePath := strings.Join(path, "/")
	s3, err := storage.SafeCreateS3Bucket(s.S3Domain, s.BucketName, s.AccessKeyID, s.SecretAccessKey)

	if err != nil {
		return nil, err
	}
	return s3.NewWriter(s3FilePath)
}

// Reader for reading from an S3 bucket
func (s *S3Provider) Reader(path ...string) (io.ReadCloser, error) {
	s3FilePath := strings.Join(path, "/")
	s3, err := storage.SafeCreateS3Bucket(s.S3Domain, s.BucketName, s.AccessKeyID, s.SecretAccessKey)

	if err != nil {
		return nil, err
	}
	return s3.NewReader(s3FilePath)
}
