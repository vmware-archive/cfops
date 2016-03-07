package cfbackup

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"

	"github.com/xchapter7x/lo"
)

//NewEncryptedStorageProvider - create a encrpyted wrapper for the given provider using the given encrpytion key.
//key lengths supported are 16, 24, 32 for AES-128, AES-192, or AES-256
func NewEncryptedStorageProvider(storageProvider StorageProvider, encryptionKey string) (encryptedStorageProvider *EncryptedStorageProvider, err error) {
	if encryptionKey != "" {
		encryptedStorageProvider = &EncryptedStorageProvider{
			EncryptionKey:          encryptionKey,
			wrappedStorageProvider: storageProvider,
		}

	} else {
		errMsg := "no encryption key provided"
		lo.G.Error(errMsg)
		err = errors.New(errMsg)
	}
	return
}

//Reader - returns the encrpyted reader for the given path
func (s *EncryptedStorageProvider) Reader(path ...string) (decryptReader io.ReadCloser, err error) {
	var unEncryptedReader io.ReadCloser
	stream := s.getStream()

	if unEncryptedReader, err = s.wrappedStorageProvider.Reader(path...); err == nil {

		decryptReader = &StreamReadCloser{
			StreamReader: cipher.StreamReader{S: stream, R: unEncryptedReader},
			Closer:       unEncryptedReader,
		}
	}
	return
}

//Writer - returns the encrpyted writer for the given path
func (s *EncryptedStorageProvider) Writer(path ...string) (cryptWriter io.WriteCloser, err error) {
	var unEncryptedWriter io.WriteCloser
	stream := s.getStream()

	if unEncryptedWriter, err = s.wrappedStorageProvider.Writer(path...); err == nil {
		cryptWriter = &cipher.StreamWriter{S: stream, W: unEncryptedWriter}
	}
	return
}

func (s *EncryptedStorageProvider) getStream() cipher.Stream {
	block, err := aes.NewCipher([]byte(s.EncryptionKey))
	if err != nil {
		lo.G.Error("there was an error in generating your cipher: ", err)
		panic(err)
	}

	var iv [aes.BlockSize]byte
	return cipher.NewOFB(block, iv[:])
}
