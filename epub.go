package epubtomd

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

type Metadata struct {
	Title      string
	Author     string
	Language   string
	Version    string
	BasePath   string
	TextFiles  []string // XHTML 章节文件列表
	ImageFiles []string // 图片文件列表
}

type EpubReader interface {
	// Extract 返回解压目录和文件列表
	Extract(epubPath string) (fs.FS, error)
	ParseMetadata(r fs.FS) (*Metadata, error)
	Close() error
}

type zipEpubReader struct {
	zipReader *zip.ReadCloser
}

func NewZipEpubReader() EpubReader {
	return &zipEpubReader{}
}

func (z *zipEpubReader) Extract(epubPath string) (fs.FS, error) {
	// 解压 EPUB，返回解压后的路径
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return nil, fmt.Errorf("zip: failed to open epub: %w", err)
	}
	z.zipReader = r
	return z.zipReader, nil
}

func (z *zipEpubReader) Close() error {
	if z.zipReader != nil {
		return z.zipReader.Close()
	}
	return nil
}

// Package EPUB 2.0 和 3.0 的通用 OPF 解析结构
type Package struct {
	XMLName  xml.Name `xml:"package"`
	Version  string   `xml:"version,attr"`
	Metadata struct {
		Title    string   `xml:"title"`
		Creator  []string `xml:"creator"`
		Language string   `xml:"language"`
	} `xml:"metadata"`
	Manifest struct {
		Items []struct {
			ID        string `xml:"id,attr"`
			Href      string `xml:"href,attr"`
			MediaType string `xml:"media-type,attr"`
		} `xml:"item"`
	} `xml:"manifest"`
}

func (z *zipEpubReader) ParseMetadata(r fs.FS) (*Metadata, error) {
	// 1. 读取 META-INF/container.xml 找到 content.opf
	opfPath, err := findOpfFile(r)
	if err != nil {
		return nil, err
	}

	// 2. 解析 content.opf
	f, err := r.Open(opfPath)
	if err != nil {
		return nil, fmt.Errorf("zip: failed to open opf file: %w", err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("zip: failed to read opf file: %w", err)
	}
	var pkg Package
	if err := xml.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	// 3. 提取元数据
	metadata := &Metadata{
		Title:    pkg.Metadata.Title,
		Author:   strings.Join(pkg.Metadata.Creator, ","),
		Language: pkg.Metadata.Language,
		Version:  pkg.Version,
	}

	// 4. 提取章节和图片文件
	basePath := filepath.Dir(opfPath)
	metadata.BasePath = basePath
	for _, item := range pkg.Manifest.Items {
		fullPath := filepath.Join(basePath, item.Href)
		if strings.HasPrefix(item.MediaType, "application/xhtml+xml") {
			metadata.TextFiles = append(metadata.TextFiles, fullPath)
		} else if strings.HasPrefix(item.MediaType, "image/") {
			metadata.ImageFiles = append(metadata.ImageFiles, fullPath)
		}
	}

	return metadata, nil
}

type RootFile struct {
	FullPath string `xml:"full-path,attr"`
}

type Container struct {
	RootFiles []RootFile `xml:"rootfiles>rootfile"`
}

// findOpfFile 找到opf文件
func findOpfFile(r fs.FS) (string, error) {
	containerPath := "META-INF/container.xml"
	f, err := r.Open(containerPath)
	if err != nil {
		return "", fmt.Errorf("open container: %w", err)
	}
	defer SaleClose(f)

	data, err := io.ReadAll(f)

	if err != nil {
		return "", fmt.Errorf("failed to read container file: %w", err)
	}
	var container Container
	if err := xml.Unmarshal(data, &container); err != nil {
		return "", err
	}

	if len(container.RootFiles) > 0 {
		return container.RootFiles[0].FullPath, nil
	}
	return "", fmt.Errorf("content.opf not found")
}
