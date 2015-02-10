package mock

import "errors"

var (
	READ_FAIL_ERROR  error = errors.New("copy failed on read")
	WRITE_FAIL_ERROR error = errors.New("copy failed on write")
	CLOSE_FAIL_ERROR error = errors.New("close file failed")
)

func NewReadWriteCloser(readErr, writeErr, closeErr error) *MockReadWriteCloser {
	return &MockReadWriteCloser{
		ReadErr:  readErr,
		WriteErr: writeErr,
		CloseErr: closeErr,
	}
}

type MockReadWriteCloser struct {
	BytesRead    []byte
	BytesWritten []byte
	ReadErr      error
	WriteErr     error
	CloseErr     error
}

func (r *MockReadWriteCloser) Read(p []byte) (n int, err error) {

	if err = r.ReadErr; err == nil {
		r.BytesRead = p
		n = len(p)
	}
	return
}

func (r *MockReadWriteCloser) Close() (err error) {
	err = r.CloseErr
	return
}

func (r *MockReadWriteCloser) Write(p []byte) (n int, err error) {

	if err = r.WriteErr; err != nil {
		r.BytesWritten = p
		n = len(p)
	}
	return
}
