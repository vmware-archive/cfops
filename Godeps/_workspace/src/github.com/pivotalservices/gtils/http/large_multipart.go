package http

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/cheggaaa/pb"
	"github.com/xchapter7x/lo"
)

func getFileSize(filename string) (fsize int64, err error) {
	var (
		input *os.File
		stat  os.FileInfo
	)
	input, err = os.Open(filename)
	defer input.Close()
	stat, err = input.Stat()
	fsize = stat.Size()
	return
}

var LargeMultiPartUpload = func(conn ConnAuth, paramName, filename string, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	var fsize int64
	fsize, err = getFileSize(filename)
	pipeOut, pipeIn := io.Pipe()
	bar := pb.New(int(fsize)).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	writer := multipart.NewWriter(pipeIn)
	done := make(chan AsyncResponse)
	go asyncRequest(done, conn, filename, pipeOut, fsize, writer.FormDataContentType(), bar)
	part, err := writer.CreateFormFile(paramName, filename)
	out := io.MultiWriter(part, bar)
	_, err = io.Copy(out, fileRef)
	writer.Close()
	pipeIn.Close()
	asyncResponse := <-done
	bar.FinishPrint("Upload done!")
	return asyncResponse.Res, asyncResponse.Err
}

type AsyncResponse struct {
	Res *http.Response
	Err error
}

func asyncRequest(done chan AsyncResponse, conn ConnAuth, filename string, pipeOut io.Reader, fsize int64, contentType string, bar *pb.ProgressBar) {
	var (
		req           *http.Request
		asyncResponse = AsyncResponse{}
	)
	if req, asyncResponse.Err = http.NewRequest("POST", conn.Url, pipeOut); asyncResponse.Err == nil {

		if conn.Username != "" && conn.Password != "" {
			req.SetBasicAuth(conn.Username, conn.Password)
		}
		req.ContentLength = fsize
		req.ContentLength += 227
		req.ContentLength += int64(len(filename))
		req.ContentLength += 19
		req.Header.Set("Content-Type", contentType)
		lo.G.Debug("Created Request")
		bar.Start()
		asyncResponse.Res, asyncResponse.Err = http.DefaultClient.Do(req)
	}
	done <- asyncResponse
}
