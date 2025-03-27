package epubtomd

import (
	"fmt"
	"io"
	"path/filepath"
)

// SaleClose 安全的关闭
func SaleClose(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}

// ResolvePath 解析包含的路径，将相对路径转换为绝对路径
func ResolvePath(baseDir string, path string) (string, error) {
	if filepath.IsAbs(path) {
		// 绝对路径，直接使用
		return path, nil
	}
	// 相对路径，基于 baseDir 转换为绝对路径
	absPath, err := filepath.Abs(filepath.Join(baseDir, path))
	if err != nil {
		return "", fmt.Errorf("failed to resolve path %s: %w", path, err)
	}
	return absPath, nil
}
