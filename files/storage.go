package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	BasePath string
}

func (s *Storage) Save(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	dest := filepath.Join(s.BasePath, filename)
	dst, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("could not create file: %w", err)
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("could not write file: %w", err)
	}

	return dest, nil
}

func (s *Storage) Dlt(filename string) error {

	dest := filepath.Join(s.BasePath, filename)
	return os.Remove(dest)

}
