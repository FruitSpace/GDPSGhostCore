package connectors

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

// OutputConnector used to write anything to anywhere (stdio/http/files)
type OutputConnector interface {
	Write()
}

type HttpWriter struct {
	Endpoint string
}

func (hw HttpWriter) Write(dat []byte ) (int, error) {
	resp, err:= http.Post(hw.Endpoint,"text/plain",bytes.NewBuffer(dat))
	return resp.StatusCode, err
}

type FileWriter struct {
	Endpoint string
}
func (fw FileWriter) Write(dat []byte) (int, error) {
	f,_:=os.OpenFile(fw.Endpoint,os.O_APPEND,644)
	defer f.Close()
	return f.Write(dat)
}

func GetWriter(t string, endpoint string) io.Writer {
	switch t {
	case "http":
		return HttpWriter{endpoint}
	case "file":
		return FileWriter{endpoint}
	default:
		return os.Stderr
	}
}