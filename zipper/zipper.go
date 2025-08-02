package zipper

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func ZipStream(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")

			gzWriter := gzip.NewWriter(w)
			defer gzWriter.Close()

			wrapper := &WrappedWriter{
				RespWriter: w,
				GzipWriter: gzWriter,
			}

			next.ServeHTTP(wrapper, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type WrappedWriter struct {
	RespWriter http.ResponseWriter
	GzipWriter *gzip.Writer
}

func (ww *WrappedWriter) Header() http.Header {
	return ww.RespWriter.Header()
}

func (ww *WrappedWriter) Write(d []byte) (int, error) {
	return ww.GzipWriter.Write(d)
}

func (ww *WrappedWriter) WriteHeader(statusCode int) {
	ww.RespWriter.WriteHeader(statusCode)
}

func (ww *WrappedWriter) Flush() {
	ww.GzipWriter.Flush()
}
