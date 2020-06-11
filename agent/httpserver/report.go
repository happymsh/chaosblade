package httpserver

import (
	"archive/zip"
	"fmt"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade/agent/metrics"
	"io"
	"net/http"
	"os"
)

func reportHandler(w http.ResponseWriter, r *http.Request) {
	//reader
	file, err := os.Open(metrics.ReportFileName)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("Open node.report file error, %s", err.Error())).ToString()))
		return
	}
	defer file.Close()
	//writer
	writer := zip.NewWriter(w)
	defer writer.Close()
	f, err := writer.Create(metrics.ReportFile)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("create node.report.zip file error, %s", err.Error())).ToString()))
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(spec.ReturnFail(spec.Code[spec.ServerError],
			fmt.Sprintf("zip node.report file error, %s", err.Error())).ToString()))
		return
	}
}
