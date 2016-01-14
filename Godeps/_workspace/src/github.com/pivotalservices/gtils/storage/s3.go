package storage

import (
	"errors"
	"io"

	"github.com/rlmcpherson/s3gof3r"
)

//S3Bucket - Represents an s3 bucket
type S3Bucket struct {
	Name      string
	Domain    string
	Bucket    string
	AccessKey string
	SecretKey string
	bucket    *s3gof3r.Bucket
}

//SafeCreateS3Bucket creates an s3 bucket for storing files to an s3-compatible blobstore
func SafeCreateS3Bucket(domain, bucket, accessKey, secretKey string) (*S3Bucket, error) {
	s := &S3Bucket{
		Bucket:    bucket,
		Name:      "s3",
		Domain:    domain,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
	if s.Bucket == "" {
		return nil, errors.New("bucket name is undefined")
	}
	var k s3gof3r.Keys
	var err error

	if s.AccessKey == "" || s.SecretKey == "" {
		k, err = s3gof3r.EnvKeys() // get S3 keys from environment
		if err != nil {
			return nil, err
		}
	} else {
		k = s3gof3r.Keys{
			AccessKey: s.AccessKey,
			SecretKey: s.SecretKey,
		}
	}
	s3 := s3gof3r.New(s.Domain, k)
	s.bucket = s3.Bucket(s.Bucket)
	return s, nil
}

//NewWriter - get a new s3 writer
func (s *S3Bucket) NewWriter(path string) (io.WriteCloser, error) {
	return s.bucket.PutWriter(path, nil, nil)
}

//NewReader - get a new s3 reader
func (s *S3Bucket) NewReader(path string) (io.ReadCloser, error) {
	r, _, err := s.bucket.GetReader(path, nil)
	return r, err
}

//Delete - delete an s3 bucket
func (s *S3Bucket) Delete(path string) error {
	return s.bucket.Delete(path)
}
