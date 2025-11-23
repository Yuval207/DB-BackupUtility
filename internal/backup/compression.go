package backup

import (
	"compress/gzip"
	"io"
	"os"
)

func CompressFile(srcPath string) (string, error) {
	destPath := srcPath + ".gz"
	
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	gzipWriter := gzip.NewWriter(destFile)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, srcFile)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

func DecompressFile(srcPath string) (string, error) {
	// Remove .gz extension
	destPath := srcPath[:len(srcPath)-3]
	
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	gzipReader, err := gzip.NewReader(srcFile)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, gzipReader)
	if err != nil {
		return "", err
	}

	return destPath, nil
}
