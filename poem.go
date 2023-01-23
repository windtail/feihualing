package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type Poem struct {
	ID       uint64     `json:"id" gorm:"primarykey"`
	Title    string     `json:"title"`
	Dynasty  string     `json:"dynasty"`
	Author   string     `json:"author"`
	Content  string     `json:"content"`
	Favor    bool       `json:"favor"`
	Segments []*Segment `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Segment struct {
	ID      uint64 `gorm:"primarykey"`
	Content string
	PoemID  uint64
}

func (p *Poem) Abstract() string {
	return fmt.Sprintf("%d. %s  (%s %s)", p.ID, p.Title, p.Dynasty, p.Author)
}

func (p *Poem) MakeSegments() {
	content := strings.ReplaceAll(p.Content, "\n", "")
	r := regexp.MustCompile(`.*?[，。：？！,.:?!]`)
	segments := r.FindAllString(content, -1)
	p.Segments = make([]*Segment, 0, len(segments))
	for _, seg := range segments {
		p.Segments = append(p.Segments, &Segment{
			Content: seg,
			PoemID:  p.ID,
		})
	}
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

	for _, seg := range poem.Segments {
		if matched(seg.Content) {
			segments = append(segments, highlighted(seg.Content))
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

var db *gorm.DB

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
		if p.ID != s.Id {
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

func toFilePath(uri string) (string, error) {
	if !strings.HasPrefix(uri, "file://") {
		return "", errors.New("unexpected uri")
	} else {
		return uri[len("file://"):], nil
	}
}

func transaction(f func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = errors.New("unexpected error")
		}
	}()

	err = f(tx)
	if err == nil {
		err = tx.Commit().Error
	}

	return
}

func (p *Poems) Init(uri string) error {
	path, err := toFilePath(uri)
	if err != nil {
		return err
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		dir := filepath.Dir(path)
		_, err = os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
			if err != nil {
				return err
			}
		}

		f, err := os.Create(path)
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Poem{}, &Segment{})
	if err != nil {
		return err
	}
	err = db.Exec("PRAGMA foreign_keys=ON").Error
	if err != nil {
		return err
	}

	err = db.Model(&Poem{}).Preload("Segments").Find(&p.list).Error
	if err != nil {
		return err
	}

	return nil
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

func (p *Poems) Clear() error {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Poem{}).Error; err != nil {
		return err
	}

	p.list = make([]*Poem, 0)
	return nil
}

func (p *Poems) createAll() error {
	return transaction(func(tx *gorm.DB) error {
		for _, poem := range p.list {
			if err := db.Create(poem).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	})
}

func (p *Poems) Import(reader fyne.URIReadCloser) error {
	if err := p.Clear(); err != nil {
		return nil
	}

	if err := p.Load(reader); err != nil {
		return nil
	}

	return p.createAll()
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

func (p *Poems) Export(writer fyne.URIWriteCloser) error {
	return p.Store(writer)
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

func (p *Poems) Remove(poem *Poem) error {
	if err := db.Delete(poem).Error; err != nil {
		return err
	}

	for i, pm := range p.list {
		if pm == poem {
			p.list = append(p.list[:i], p.list[i+1:]...)
			break
		}
	}

	return nil
}
