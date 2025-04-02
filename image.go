package epubtomd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type ImageHandler interface {
	// CopyImage 返回 Markdown 兼容路径
	CopyImage(srcImagePath string, dstImagePath string) (string, error)
	// CopyWithRename 复制并重命名指定文件
	CopyWithRename(srcImagePath string, namePathFn func(b []byte) string) (string, error)
}

// LocalImageHandler 本地文件处理，将指定文件系统中的文件复制到本地文件，并返回相对路径
type LocalImageHandler struct {
	f    fs.FS
	path string
}

func NewLocalImageHandler(f fs.FS, dirPath string) ImageHandler {
	return &LocalImageHandler{
		f:    f,
		path: dirPath,
	}
}

// CopyImage 将指定文件系统中的文件复制到指定目录中
func (h *LocalImageHandler) CopyImage(srcImagePath string, dstImagePath string) (string, error) {
	return h.CopyWithRename(srcImagePath, func(b []byte) string {
		return dstImagePath
	})
}

func (h *LocalImageHandler) CopyWithRename(srcImagePath string, namePathFn func(b []byte) string) (string, error) {
	file, err := h.f.Open(srcImagePath)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %w", srcImagePath, err)
	}
	defer SaleClose(file)

	body, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", srcImagePath, err)
	}
	dstImagePath := namePathFn(body)

	fullLocalImagePath := filepath.Join(h.path, dstImagePath)

	fullLocalImageDir := filepath.Dir(fullLocalImagePath)

	if _, err := os.Stat(fullLocalImageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(fullLocalImageDir, 0766); err != nil {
			return "", fmt.Errorf("error creating directory %s: %w", fullLocalImageDir, err)
		}
	}
	localFile, err := os.Create(fullLocalImagePath)
	if err != nil {
		return "", fmt.Errorf("error creating file %s: %w", fullLocalImagePath, err)
	}
	defer SaleClose(localFile)
	n, err := localFile.Write(body)
	if err != nil {
		return "", fmt.Errorf("error copying file %s: %w", fullLocalImagePath, err)
	}
	fmt.Println("copy file", srcImagePath, "to", fullLocalImagePath, "size", n)
	return dstImagePath, nil
}
