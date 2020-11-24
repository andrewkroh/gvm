package common

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Extract(sourceFile, destinationDir string) error {
	switch {
	case strings.HasSuffix(sourceFile, ".tar.gz"), strings.HasSuffix(sourceFile, ".tgz"):
		return untar(sourceFile, destinationDir)
	case strings.HasSuffix(sourceFile, ".zip"):
		return unzip(sourceFile, destinationDir)
	default:
		return fmt.Errorf("failed to extract %v, unhandled file type", sourceFile)
	}
}

func unzip(sourceFile, destinationDir string) error {
	r, err := zip.OpenReader(sourceFile)
	if err != nil {
		return err
	}
	defer r.Close()

	if err = os.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("failed to mkdir %v: %w", destinationDir, err)
	}

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(destinationDir, f.Name)

		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(path, f.Mode()); err != nil {
				return fmt.Errorf("failed to mkdir %v: %w", path, err)
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
				return fmt.Errorf("failed to mkdir %v: %w", path, err)
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return fmt.Errorf("failed extracting %q from %q: %w", f.Name, sourceFile, err)
		}
	}

	return nil
}

func untar(sourceFile, destinationDir string) error {
	file, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var fileReader io.ReadCloser = file

	if strings.HasSuffix(sourceFile, ".gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {
			return err
		}
		defer fileReader.Close()
	}

	tarBallReader := tar.NewReader(fileReader)

	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		filename := filepath.Join(destinationDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer
			if err != nil {
				return err
			}

		case tar.TypeReg:
			writer, err := os.Create(filename)
			if err != nil {
				return err
			}

			if _, err = io.Copy(writer, tarBallReader); err != nil {
				return err
			}

			if err = os.Chmod(filename, os.FileMode(header.Mode)); err != nil {
				return err
			}

			writer.Close()
		default:
			return fmt.Errorf("unable to untar type: %c in file %s", header.Typeflag, filename)
		}
	}
	return nil
}
