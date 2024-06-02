package extractor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
)

func MCZ(r io.Reader, targetExt string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, fmt.Errorf("copying mcz file: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(len(buf.Bytes())))
	if err != nil {
		return nil, fmt.Errorf("initing zip reader: %w", err)
	}

	for _, file := range zipReader.File {
		if filepath.Ext(file.Name) == targetExt {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("opening target file: %w", err)
			}
			defer rc.Close()

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, rc); err != nil {
				return nil, fmt.Errorf("copying target file: %w", err)
			}

			return buf.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("no file with extension %s found in mcz", targetExt)
}
