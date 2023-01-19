package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type Poem struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Dynasty string `json:"dynasty"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

func (p *Poem) Abstract() string {
	return fmt.Sprintf("%d. %s  (%s %s)", p.Id, p.Title, p.Dynasty, p.Author)
}

type PoemTemplateItem struct {
	*Poem
	MarkdownContent string
}

func NewPoemTemplateItem(poem *Poem, key string) *PoemTemplateItem {
	content := strings.ReplaceAll(poem.Content, "\n", "\n\n")
	if len(key) != 0 {
		content = strings.ReplaceAll(content, key, fmt.Sprintf(" **%s** ", key))
	}

	return &PoemTemplateItem{Poem: poem, MarkdownContent: content}
}

var poemMarkdownTpl *template.Template

func init() {
	poemMarkdownTpl, _ = template.New("poem").Parse(`# {{.Title}}

{{.Dynasty}} {{.Author}}

{{.MarkdownContent}}`)
}

func (p *Poem) Markdown(key string) string {
	var buf bytes.Buffer
	_ = poemMarkdownTpl.Execute(&buf, NewPoemTemplateItem(p, key))

	return buf.String()
}
