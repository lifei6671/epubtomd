package epubtomd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/strikethrough"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"

	"github.com/PuerkitoBio/goquery"
)

// XHTMLConverter html 文件转换接口
type XHTMLConverter interface {
	// Convert 将指定文件系统中的指定文件转换为 Markdown 文件
	Convert(filename string) (string, string, error)
}

type BasicXHTMLConverter struct {
	f fs.FS
}

func NewBasicXHTMLConverter(f fs.FS) XHTMLConverter {
	return &BasicXHTMLConverter{
		f: f,
	}
}
func (c *BasicXHTMLConverter) Convert(filename string) (string, string, error) {
	// 用来修正畸形的文件路径
	filename = filepath.Clean(filename)
	// 使用 goquery 解析 XHTML，转换为 Markdown
	file, err := c.f.Open(filename)
	if err != nil {
		return "", "", fmt.Errorf("can't open file %s: %w", filename, err)
	}
	defer SaleClose(file)
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return "", "", fmt.Errorf("can't parse file %s: %w", filename, err)
	}
	title := strings.TrimSpace(doc.Find("title").Text())
	data, err := doc.Find("body").Html()
	if err != nil {
		return "", "", fmt.Errorf("can't parse file %s: %w", filename, err)
	}
	body, err := HtmlToMarkdown(data)
	if err != nil {
		return "", "", fmt.Errorf("can't convert file %s: %w", filename, err)
	}
	return title, fmt.Sprintf("# %s\n\n%s", title, body), nil
}

// HtmlToMarkdown 将 html 格式内容转换为 markdown 格式
func HtmlToMarkdown(htmlContent string) (string, error) {
	convert := converter.NewConverter(converter.WithPlugins(
		base.NewBasePlugin(),
		commonmark.NewCommonmarkPlugin(
			commonmark.WithStrongDelimiter("__"),
		),
		table.NewTablePlugin(),
		strikethrough.NewStrikethroughPlugin(),
	))

	markdownContent, err := convert.ConvertString(htmlContent)
	if err != nil {
		return "", fmt.Errorf("convert markdown content: %w", err)
	}
	return markdownContent, nil
}
