package epubtomd

import (
	"fmt"
	"io/fs"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
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
	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.Table())
	converter.Use(plugin.TaskListItems())
	converter.Use(plugin.YoutubeEmbed())
	converter.Use(plugin.EXPERIMENTALMoveFrontMatter())

	markdownContent, err := converter.ConvertString(htmlContent)
	if err != nil {
		return "", fmt.Errorf("convert markdown content: %w", err)
	}
	return markdownContent, nil
}
