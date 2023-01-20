package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"io/ioutil"
	"strings"
	"text/template"
)

type Poem struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Dynasty string `json:"dynasty"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Favor   bool   `json:"favor"`
}

func (p *Poem) Abstract() string {
	return fmt.Sprintf("%d. %s  (%s %s)", p.Id, p.Title, p.Dynasty, p.Author)
}

type PoemTemplateItem struct {
	*Poem
	MarkdownContent string
}

func NewPoemTemplateItem(poem *Poem, s *Search) *PoemTemplateItem {
	content := strings.ReplaceAll(poem.Content, "\n", "\n\n")

	for _, key := range s.Content {
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

func (p *Poem) Markdown(s *Search) string {
	var buf bytes.Buffer
	_ = poemMarkdownTpl.Execute(&buf, NewPoemTemplateItem(p, s))

	return buf.String()
}

func (p *Poem) PreviewMarkdown(s *Search) string {
	return strings.Join(strings.Split(p.Content, "\n")[:2], "")
}

func (p *Poem) Matched(s *Search) bool {
	if s.Id != 0 {
		if p.Id != s.Id {
			return false
		}
	}

	if s.FavorOnly {
		if !p.Favor {
			return false
		}
	}

	if len(s.Title) != 0 {
		if !strings.Contains(p.Title, s.Title) {
			return false
		}
	}

	if len(s.Dynasty) != 0 {
		if !strings.Contains(p.Dynasty, s.Dynasty) {
			return false
		}
	}

	if len(s.Author) != 0 {
		if !strings.Contains(p.Author, s.Author) {
			return false
		}
	}

	if len(s.Content) != 0 {
		for _, key := range s.Content {
			if !strings.Contains(p.Content, key) {
				return false
			}
		}
	}

	return true
}

type Poems struct {
	list []*Poem
}

func NewPoems() *Poems {
	return &Poems{
		list: make([]*Poem, 0),
	}
}

//go:embed poems.json
var _defaultPoems []byte

func (p *Poems) LoadDefault() {
	_ = json.Unmarshal(_defaultPoems, &p.list)
}

func (p *Poems) Load(reader fyne.URIReadCloser) (err error) {
	defer func(reader fyne.URIReadCloser) {
		err = reader.Close()
	}(reader)

	if data, err := ioutil.ReadAll(reader); err != nil {
		return err
	} else {
		return json.Unmarshal(data, &p.list)
	}
}

func (p *Poems) Store(writer fyne.URIWriteCloser) (err error) {
	defer func(writer fyne.URIWriteCloser) {
		err = writer.Close()
	}(writer)

	if data, err := json.Marshal(&p.list); err != nil {
		return err
	} else {
		block := data[:]
		for len(block) > 0 {
			if n, err := writer.Write(block); err != nil {
				return err
			} else {
				block = block[n:]
			}
		}
	}

	return nil
}

func (p *Poems) Filter(s *Search) []*Poem {
	filtered := make([]*Poem, 0, len(p.list))

	for _, poem := range p.list {
		if poem.Matched(s) {
			filtered = append(filtered, poem)
		}
	}

	return filtered
}

func (p *Poems) Remove(poem *Poem) (err error) {
	// TODO do remove
	return nil
}
