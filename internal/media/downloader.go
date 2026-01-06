package media

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader struct {
	ImagesDir       string
	ImagePathPrefix string
	PDSHost         string
}

func NewDownloader(imagesDir, imagePathPrefix, pdsHost string) *Downloader {
	return &Downloader{
		ImagesDir:       imagesDir,
		ImagePathPrefix: imagePathPrefix,
		PDSHost:         pdsHost,
	}
}

func (d *Downloader) DownloadBlob(ctx context.Context, did string, cid string) (string, error) {
	// https://bsky.social/xrpc/com.atproto.sync.getBlob?did=did:plc:xxx&cid=bafyxxx
	url := fmt.Sprintf("%s/xrpc/com.atproto.sync.getBlob?did=%s&cid=%s", d.PDSHost, did, cid)

	if err := os.MkdirAll(d.ImagesDir, 0755); err != nil {
		return "", err
	}

	// Check if file already exists (try common extensions)
	for _, ext := range []string{".jpg", ".png", ".webp", ".gif", ".bin"} {
		fileName := cid + ext
		filePath := filepath.Join(d.ImagesDir, fileName)
		if _, err := os.Stat(filePath); err == nil {
			return filepath.Join(d.ImagePathPrefix, fileName), nil
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download blob: %s", resp.Status)
	}

	// Determine extension from Content-Type header
	ext := ".bin"
	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	case "image/gif":
		ext = ".gif"
	}
	fileName := cid + ext
	filePath := filepath.Join(d.ImagesDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filepath.Join(d.ImagePathPrefix, fileName), nil
}
