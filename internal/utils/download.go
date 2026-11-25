package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DownloadFile(url, savePath, username, password string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status: %s", resp.Status)
	}

	dir := filepath.Dir(savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %w", err)
	}

	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("create file failed: %w", err)
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		file.Close()
		os.Remove(savePath)
		return fmt.Errorf("download failed: %w", err)
	}

	return nil
}

func ReadTxt(txtfile string, action func(line string)) {
	f, err := os.Open(txtfile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, "\r")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		action(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func UntarGz(src string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue

		case tar.TypeReg, tar.TypeRegA:
			if !strings.HasSuffix(header.Name, ".mmdb") {
				continue
			}
			baseName := filepath.Base(header.Name)
			out, err := os.OpenFile(
				baseName,
				os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
				os.FileMode(header.Mode),
			)
			if err != nil {
				return err
			}

			_, copyErr := io.Copy(out, tarReader)
			closeErr := out.Close()

			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		}
	}

	return nil
}
