package epubtomd

import "fmt"

type MarkdownGenerator interface {
	GenerateMarkdown(metadata *Metadata, chapterContents []string) (string, error)
}

type SimpleMarkdownGenerator struct{}

func (g *SimpleMarkdownGenerator) GenerateMarkdown(metadata *Metadata, chapterContents []string) (string, error) {
	result := fmt.Sprintf("# %s\n\n", metadata.Title)
	for _, content := range chapterContents {
		result += content + "\n\n"
	}
	return result, nil
}
