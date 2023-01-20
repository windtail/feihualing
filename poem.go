package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"
)

type Poem struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Dynasty  string `json:"dynasty"`
	Author   string `json:"author"`
	Content  string `json:"content"`
	Favor    bool   `json:"favor"`
	segments []string
}

func (p *Poem) Abstract() string {
	return fmt.Sprintf("%d. %s  (%s %s)", p.Id, p.Title, p.Dynasty, p.Author)
}

func (p *Poem) MakeSegments() {
	content := strings.ReplaceAll(p.Content, "\n", "")
	r := regexp.MustCompile(`.*?[，。：？！,.:?!]`)
	p.segments = r.FindAllString(content, -1)
}

func highlight(s, key string) string {
	return strings.ReplaceAll(s, key, fmt.Sprintf(" **%s** ", key))
}

type PoemDetailTemplateContext struct {
	*Poem
	MarkdownContent string
}

func NewPoemDetailTemplateContext(poem *Poem, s *Search) *PoemDetailTemplateContext {
	content := strings.ReplaceAll(poem.Content, "\n", "\n\n")

	for _, key := range s.Content {
		content = highlight(content, key)
	}

	return &PoemDetailTemplateContext{Poem: poem, MarkdownContent: content}
}

var poemDetailMarkdownTpl *template.Template

type PoemPreviewTemplateContext struct {
	*Poem
	MarkdownContent string
}

func NewPoemPreviewTemplateContext(poem *Poem, s *Search) *PoemPreviewTemplateContext {
	const MaxSegment = 2

	segments := make([]string, 0, MaxSegment)

	matched := func(seg string) bool {
		for _, key := range s.Content {
			if strings.Contains(seg, key) {
				return true
			}
		}
		return false
	}

	highlighted := func(seg string) string {
		for _, key := range s.Content {
			seg = highlight(seg, key)
		}
		return seg
	}

	for _, seg := range poem.segments {
		if matched(seg) {
			segments = append(segments, highlighted(seg))
			if len(segments) == MaxSegment {
				break
			}
		}
	}

	return &PoemPreviewTemplateContext{Poem: poem, MarkdownContent: strings.Join(segments, "")}
}

var poemPreviewMarkdownTpl *template.Template

func init() {
	poemDetailMarkdownTpl, _ = template.New("detail").Parse(`# {{.Title}}

{{.Dynasty}} {{.Author}}

{{.MarkdownContent}}`)

	poemPreviewMarkdownTpl, _ = template.New("preview").Parse(`	{{.Abstract}}

{{.MarkdownContent}}`)
}

func (p *Poem) DetailMarkdown(s *Search) string {
	var buf bytes.Buffer
	_ = poemDetailMarkdownTpl.Execute(&buf, NewPoemDetailTemplateContext(p, s))

	return buf.String()
}

func (p *Poem) PreviewMarkdown(s *Search) string {
	var buf bytes.Buffer
	_ = poemPreviewMarkdownTpl.Execute(&buf, NewPoemPreviewTemplateContext(p, s))

	return buf.String()
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

func (p *Poems) MakeSegments() {
	for _, poem := range p.list {
		poem.MakeSegments()
	}
}

func (p *Poems) LoadDefault() {
	_ = json.Unmarshal(_defaultPoems, &p.list)
	p.MakeSegments()
}

func (p *Poems) Load(reader fyne.URIReadCloser) (err error) {
	defer func(reader fyne.URIReadCloser) {
		err = reader.Close()
	}(reader)

	if data, err := ioutil.ReadAll(reader); err != nil {
		return err
	} else {
		if err := json.Unmarshal(data, &p.list); err != nil {
			return err
		}
		p.MakeSegments()
		return nil
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
