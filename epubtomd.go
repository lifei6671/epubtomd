package epubtomd

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 正则匹配本地图片 ![alt text](image.png) 或 ![](image.png)
var localImageRegex = regexp.MustCompile(`!\[(.*?)?\]\(([^:\)]+)\)`) // 兼容没有 alt 文本的情况

// Convert 将指定路径的epub文件转换为markdown文件，并报错到指定路径
func Convert(filename string, output string) error {
	reader := NewZipEpubReader()
	r, err := reader.Extract(filename)
	if err != nil {
		return err
	}
	defer SaleClose(reader)

	metadata, err := reader.ParseMetadata(r)
	if err != nil {
		return err
	}
	htmlConvert := NewBasicXHTMLConverter(r)
	imageHandler := NewLocalImageHandler(r, output)

	if _, err := os.Stat(output); os.IsNotExist(err) {
		if err := os.MkdirAll(output, 0777); err != nil {
			return fmt.Errorf("unable to create output directory: %w", err)
		}
	}

	readme, err := os.Create(filepath.Join(output, "README.md"))
	if err != nil {
		return fmt.Errorf("unable to create README.md: %w", err)
	}
	_, _ = readme.WriteString("# " + metadata.Title + "\n")
	for _, p := range metadata.TextFiles {
		title, mdContent, err := htmlConvert.Convert(p)
		if err != nil {
			return err
		}

		relativePath := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p)) + ".md"
		mdPath := filepath.Join(output, relativePath)

		mdContent = localImageRegex.ReplaceAllStringFunc(mdContent, func(s string) string {
			matches := localImageRegex.FindStringSubmatch(s)
			if len(matches) > 2 {
				altText := matches[1]
				localPath := matches[2]
				if !strings.HasPrefix(localPath, "http") && !strings.HasPrefix(localPath, "https") { // 确保不是外部链接
					if altText == "" {
						altText = "image" // 默认 alt 文本
					}
					imagePath := filepath.Join(metadata.BasePath, localPath)
					imagePath = strings.ReplaceAll(imagePath, "\\", "/") // 替换为正斜杠以兼容 Windows 和 Unix 系统
					ext := filepath.Ext(localPath)

					rename, err := imageHandler.CopyWithRename(imagePath, func(b []byte) string {
						hash := md5.New()
						hash.Write(b)
						return fmt.Sprintf("images/%s/%s", time.Now().Format("20060102"), fmt.Sprintf("%x", hash.Sum(nil))+ext)
					})
					if err != nil {
						log.Printf("unable to rename image: %s", localPath)
						//如果失败则使用原地址
						return fmt.Sprintf("![%s](%s)", altText, localPath)
					}
					return fmt.Sprintf("![%s](%s)", altText, rename)
				}
			}
			return s
		})

		wErr := os.WriteFile(mdPath, []byte(mdContent), 0777)
		if wErr != nil {
			return fmt.Errorf("write file error: %w", wErr)
		}
		_, _ = readme.WriteString(fmt.Sprintf("- [%s](%s)\n", title, relativePath))
	}
	return err
}
