package httpUtils

import (
	"compress/gzip"
	"net/http"
	"sync"
)

type GzipResponseWriter struct {
	http.ResponseWriter

	GW            *gzip.Writer
	StatusCode    int
	HeaderWritten bool
}

var (
	Pool = sync.Pool{
		New: func() interface{} {
			GW, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
			return &GzipResponseWriter{
				GW: GW,
			}
		},
	}
)

func (gzr *GzipResponseWriter) WriteHeader(statusCode int) {
	gzr.StatusCode = statusCode
	gzr.HeaderWritten = true

	if gzr.StatusCode != http.StatusNotModified && gzr.StatusCode != http.StatusNoContent {
		gzr.ResponseWriter.Header().Del("Content-Length")
		gzr.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}

	gzr.ResponseWriter.WriteHeader(statusCode)
}

func (gzr *GzipResponseWriter) Write(b []byte) (int, error) {
	if _, ok := gzr.Header()["Content-Type"]; !ok {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		gzr.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(b))
	}

	if !gzr.HeaderWritten {
		// This is exactly what Go would also do if it hasn't been written yet.
		gzr.WriteHeader(http.StatusOK)
	}

	return gzr.GW.Write(b)
}

func (gzr *GzipResponseWriter) Flush() {
	if gzr.GW != nil {
		gzr.GW.Flush()
	}

	if fw, ok := gzr.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}
