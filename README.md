# epubtomd

`epubtomd` 是一个用 Golang 实现的工具，它能够将 EPUB 文件转换为 Markdown 文件。该工具支持提取 EPUB 文件的元数据，如标题、作者等，同时可以处理
EPUB 中的图片，并将其复制到指定目录或上传到 S3 存储桶。此外，它还能将 EPUB 中的 XHTML 章节文件转换为 Markdown
格式，并生成一个包含所有章节链接的 `README.md` 文件。

## 功能特性

- **元数据提取**：从 EPUB 文件中提取标题、作者、语言、版本等元数据。
- **图片处理**：支持将 EPUB 中的图片复制到本地目录或上传到 S3 存储桶。
- **XHTML 转换**：将 EPUB 中的 XHTML 章节文件转换为 Markdown 格式。
- **Markdown 生成**：生成一个包含所有章节链接的 `README.md` 文件。

## 安装

确保你已经安装了 Go 1.23 或更高版本。然后，使用以下命令下载并安装 `epubtomd`：

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
	err := epubtomd.Convert(epubPath, outputDir)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Println("转换成功")
	}
}
```

### 测试代码

项目中包含了多个测试文件，用于验证各个功能的正确性。你可以使用以下命令运行所有测试：

```sh
go test ./...
```

## 代码结构

- **`epub.go`**：负责解压 EPUB 文件，解析 OPF 文件，提取元数据和章节、图片文件列表。
- **`converter.go`**：实现了将 XHTML 文件转换为 Markdown 文件的功能。
- **`image.go`**：定义了图片处理接口和本地图片处理实现。
- **`image_s3.go`**：实现了将图片上传到 S3 存储桶的功能。
- **`generator.go`**：定义了 Markdown 生成接口和简单的 Markdown 生成器实现。
- **`epubtomd.go`**：提供了将 EPUB 文件转换为 Markdown 文件的主要功能。
- **`util.go`**：包含了一些辅助函数，如安全关闭文件和路径解析。

## 依赖项

项目依赖了多个第三方库，具体可以查看 `go.mod` 文件。以下是一些主要的依赖项：

- `github.com/JohannesKaufmann/html-to-markdown`：用于将 HTML 内容转换为 Markdown 格式。
- `github.com/aws/aws-sdk-go-v2`：用于与 AWS S3 存储桶进行交互。
- `github.com/PuerkitoBio/goquery`：用于解析 XHTML 文件。

## 贡献

如果你发现了 bug 或者有新的功能建议，欢迎提交 issue 或 pull request。在提交 pull request 之前，请确保你的代码通过了所有测试，并且遵循了项目的代码风格。

## 许可证

本项目采用 [MIT 许可证](LICENSE)。