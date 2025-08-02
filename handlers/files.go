package handlers

import (
	"corenethttp/files"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Files struct {
	Lg *log.Logger
	St *files.Storage
}

func FilesController(l *log.Logger, s *files.Storage) *Files {
	return &Files{l, s}
}

func (fs *Files) StoreFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error : "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, "Error reading file: "+err.Error(), http.StatusBadRequest)
		return
	}

	path, err := fs.St.Save(file, header)
	if err != nil {
		http.Error(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fs.Lg.Println("File saved to:", path)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func (fs *Files) DeleteFile(w http.ResponseWriter, r *http.Request) {

	filename := chi.URLParam(r, "file")
	if filename == "" {
		http.Error(w, "Error reading file: No file provided", http.StatusBadRequest)
		return
	}

	err := fs.St.Dlt(filename)
	if err != nil {
		http.Error(w, "Error deleting file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fs.Lg.Println("File delete")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
}
