package backup

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (scp *SecureCopier) processDir(procWriter io.Writer, srcFilePath string, srcFileInfo os.FileInfo, outPipe io.Writer, errPipe io.Writer) error {
	err := scp.sendDir(procWriter, srcFilePath, srcFileInfo, errPipe)
	if err != nil {
		return err
	}
	dir, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if fi.IsDir() {
			err = scp.processDir(procWriter, filepath.Join(srcFilePath, fi.Name()), fi, outPipe, errPipe)
			if err != nil {
				return err
			}
		} else {
			err = scp.sendFile(procWriter, filepath.Join(srcFilePath, fi.Name()), fi, outPipe, errPipe)
			if err != nil {
				return err
			}
		}
	}
	//TODO process errors
	err = scp.sendEndDir(procWriter, errPipe)
	return err
}

func (scp *SecureCopier) sendEndDir(procWriter io.Writer, errPipe io.Writer) error {
	header := fmt.Sprintf("E\n")
	if scp.IsVerbose {
		fmt.Fprintf(errPipe, "Sending end dir: %s", header)
	}
	_, err := procWriter.Write([]byte(header))
	return err
}

func (scp *SecureCopier) sendDir(procWriter io.Writer, srcPath string, srcFileInfo os.FileInfo, errPipe io.Writer) error {
	mode := uint32(srcFileInfo.Mode().Perm())
	header := fmt.Sprintf("D%04o 0 %s\n", mode, filepath.Base(srcPath))
	if scp.IsVerbose {
		fmt.Fprintf(errPipe, "Sending Dir header : %s", header)
	}
	_, err := procWriter.Write([]byte(header))
	return err
}

func (scp *SecureCopier) sendFile(procWriter io.Writer, srcPath string, srcFileInfo os.FileInfo, outPipe io.Writer, errPipe io.Writer) error {
	//single file
	mode := uint32(srcFileInfo.Mode().Perm())
	fileReader, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer fileReader.Close()
	size := srcFileInfo.Size()
	header := fmt.Sprintf("C%04o %d %s\n", mode, size, filepath.Base(srcPath))
	if scp.IsVerbose {
		fmt.Fprintf(errPipe, "Sending File header: %s", header)
	}
	pb := NewProgressBarTo(srcPath, size, outPipe)
	pb.Update(0)
	_, err = procWriter.Write([]byte(header))
	if err != nil {
		return err
	}
	//TODO buffering
	_, err = io.Copy(procWriter, fileReader)
	if err != nil {
		return err
	}
	// terminate with null byte
	err = sendByte(procWriter, 0)
	if err != nil {
		return err
	}

	err = fileReader.Close()
	if scp.IsVerbose {
		fmt.Fprintln(errPipe, "Sent file plus null-byte.")
	}
	pb.Update(size)
	fmt.Fprintln(errPipe)

	if err != nil {
		fmt.Fprintln(errPipe, err.Error())
	}
	return err
}

//to-scp
func (scp *SecureCopier) scpToRemote(srcFile, password, dstUser, dstHost, dstFile string, outPipe io.Writer, errPipe io.Writer) error {

	srcFileInfo, err := os.Stat(srcFile)
	if err != nil {
		fmt.Fprintln(errPipe, "Could not stat source file "+srcFile)
		return err
	}
	session, err := Connect(dstUser, password, dstHost, scp.Port, scp.IsVerbose, errPipe)
	if err != nil {
		return err
	} else if scp.IsVerbose {
		fmt.Fprintln(errPipe, "Got session")
	}
	defer session.Close()
	ce := make(chan error)
	if dstFile == "" {
		dstFile = filepath.Base(srcFile)
		//dstFile = "."
	}
	go func() {
		procWriter, err := session.StdinPipe()
		if err != nil {
			fmt.Fprintln(errPipe, err.Error())
			ce <- err
			return
		}
		defer procWriter.Close()
		if scp.IsRecursive {
			if srcFileInfo.IsDir() {
				err = scp.processDir(procWriter, srcFile, srcFileInfo, outPipe, errPipe)
				if err != nil {
					fmt.Fprintln(errPipe, err.Error())
					ce <- err
				}
			} else {
				err = scp.sendFile(procWriter, srcFile, srcFileInfo, outPipe, errPipe)
				if err != nil {
					fmt.Fprintln(errPipe, err.Error())
					ce <- err
				}
			}
		} else {
			if srcFileInfo.IsDir() {
				ce <- errors.New("Error: Not a regular file")
				return
			} else {
				err = scp.sendFile(procWriter, srcFile, srcFileInfo, outPipe, errPipe)
				if err != nil {
					fmt.Fprintln(errPipe, err.Error())
					ce <- err
				}
			}
		}
		err = procWriter.Close()
		if err != nil {
			fmt.Fprintln(errPipe, err.Error())
			ce <- err
			return
		}
	}()
	go func() {
		select {
		case err, ok := <-ce:
			fmt.Fprintln(errPipe, "Error:", err, ok)
			os.Exit(1)
		}
	}()

	remoteOpts := "-t"
	if scp.IsQuiet {
		remoteOpts += "q"
	}
	if scp.IsRecursive {
		remoteOpts += "r"
	}
	err = session.Run("/usr/bin/scp " + remoteOpts + " " + dstFile)
	if err != nil {
		fmt.Fprintln(errPipe, "Failed to run remote scp: "+err.Error())
	}
	return err
}
