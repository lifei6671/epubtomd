# epubtomd

`epubtomd` 是一个用 Golang 编写的工具，可以将 EPUB 电子书转换为 Markdown 格式。它不仅支持提取元数据（如标题、作者等），还能处理书中的图片，并将其保存到本地或上传到
S3。此外，它还能将 EPUB 的 XHTML 章节转换为 Markdown，并自动生成包含所有章节链接的 `README.md` 文件。

## 主要功能

- **元数据提取**：获取 EPUB 书籍的标题、作者、语言、版本等信息。
- **图片处理**：支持将书中的图片复制到指定目录或上传至 S3。
- **XHTML 转换**：将 EPUB 章节转换为 Markdown 格式，保留原有结构。
- **自动生成目录**：创建 `README.md`，方便查看和导航章节内容。

## 安装

请确保你的 Go 版本为 1.23 或更高，然后运行以下命令安装 `epubtomd`：

```sh
go get github.com/lifei6671/epubtomd
```

## 使用方法

### 转换 EPUB 文件

```go
package main

import (
	"fmt"
	"github.com/lifei6671/epubtomd"
)

func main() {
	epubPath := "./testdata/history.epub"
	outputDir := "./history"
	if err := epubtomd.Convert(epubPath, outputDir); err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Println("转换成功")
	}
}
```

### 运行测试

项目包含多个测试文件，可以使用以下命令执行所有测试：

```sh
go test ./...
```

## 代码结构

- **`epub.go`**：解压 EPUB，解析 OPF 文件，提取元数据、章节和图片列表。
- **`converter.go`**：负责将 XHTML 章节转换为 Markdown。
- **`image.go`**：定义图片处理逻辑，包括本地存储。
- **`image_s3.go`**：实现将图片上传至 S3 的功能。
- **`generator.go`**：用于生成 Markdown 目录和章节内容。
- **`epubtomd.go`**：核心逻辑，负责 EPUB 到 Markdown 的转换。
- **`util.go`**：提供辅助函数，如文件操作和路径解析。

## 依赖

`epubtomd` 依赖多个 Golang 库，主要包括：

- `github.com/JohannesKaufmann/html-to-markdown`（HTML 转 Markdown）
- `github.com/aws/aws-sdk-go-v2`（S3 交互）
- `github.com/PuerkitoBio/goquery`（XHTML 解析）

详细依赖可查看 `go.mod`。

## 贡献

欢迎提交 issue 或 pull request。如果有 bug 或改进建议，请在提交 PR 之前确保代码通过所有测试，并符合项目代码风格。

## 许可证

本项目采用 [MIT 许可证](LICENSE)。
