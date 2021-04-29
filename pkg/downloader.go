package autostart

import (
	"context"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

//Download copies a file 'from' the source online location and places
//it at the local 'to' location. If the 'to' location is a directly
//Content-Disposition: attachment; filename="filename.jpg"
func Download(ctx context.Context, from, to string) (string, error) {
	uri, err := url.Parse(from)
	if err != nil {
		return "", err
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(uri.String())
	if err != nil {
		return "", err
	}
	filename := determineFileName(ctx, to, resp.Header)
	//TODO (@morgan): probably need to make the folder if it doesn't exist
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}
	return filename, nil
}

//determine filename places the file in targetDir/filename, where filename
//is extracted from the http header, or set to a default value
func determineFileName(ctx context.Context, targetDir string, header http.Header) string {
	folder := "/tmp/"
	name := "autostart-file.tmp"
	if targetDir != "" {
		folder = targetDir
	}
	headerName := extractHeaderFilename(header)
	if headerName != "" {
		name = headerName
	}
	return path.Join(folder, name)
}

func extractHeaderFilename(header http.Header) string {
	h := header.Get("Content-Disposition")
	if h == "" {
		return ""
	}
	if !strings.Contains(h, "filename=") {
		return ""
	}
	_, params, err := mime.ParseMediaType(h)
	if err != nil {
		return ""
	}
	if filename, ok := params["filename"]; ok {
		return filename
	}
	return ""
}
