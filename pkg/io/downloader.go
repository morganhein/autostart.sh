package io

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

	"golang.org/x/xerrors"
)

type Downloader interface {
	Download(ctx context.Context, from, to string) (string, error)
}

var _ Downloader = (*downloader)(nil)

type downloader struct{}

//Download copies a file 'from' the source online location and places
//it at the local 'to' location. If the 'to' location is a directly
//Content-Disposition: attachment; filename="filename.jpg"
func (d downloader) Download(ctx context.Context, from, to string) (string, error) {
	uri, err := url.Parse(from)
	if err != nil {
		return "", xerrors.Errorf("error parsing url: %v", err)
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(uri.String())
	if err != nil {
		return "", xerrors.Errorf("error getting the url: %v", err)
	}
	filename := determineFileName(ctx, to, resp.Header)
	//TODO (@morgan): probably need to make the folder if it doesn't exist
	//TODO (@morgan): this should also be an interface we control so we can mock it
	f, err := os.Create(filename)
	if err != nil {
		return "", xerrors.Errorf("error creating destination file: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", xerrors.Errorf("error copying file contents to destination file: %v", err)
	}
	return filename, nil
}

//determine filename places the file in targetDir/filename, where filename
//is extracted from the http header, or set to a default value
func determineFileName(ctx context.Context, targetDir string, header http.Header) string {
	folder := "/tmp/"
	name := "envy-file.tmp"
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
