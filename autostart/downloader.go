package autostart

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//Download copies a file 'from' the source online location and places
//it at the local 'to' location. If the 'to' location is a directly
//Content-Disposition: attachment; filename="filename.jpg"
func Download(ctx context.Context, from, to string) error {
	uri, err := url.Parse(from)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout:       10 * time.Second,
	}
	resp, err := client.Get(uri.String())
	if err != nil {
		return err
	}
	filename := determineFileName(ctx, to, resp.Header)


}

func determineFileName(ctx context.Context, target string, headers http.Header) string {
	filename := "/tmp/autostart.tmp"
	if target == "" {
		return filename
	}
	//is the target dir a folder?
	f, err := os.Stat(target)
	if err != nil {
		return filename
	}
	if f.IsDir() {

	}
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
