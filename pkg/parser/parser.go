package parser

import (
	"fmt"
	"io"
	"regexp"

	"github.com/kazukt/changelog/pkg/changelog"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Parse input into a Changelog struct
func Parse(r io.Reader) (*changelog.Changelog, error) {
	source, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	p := newParser()
	return p.parse(source)
}

type parser struct {
	markdown           goldmark.Markdown
	changelog          *changelog.Changelog
	currentVersion     *changelog.Version
	currentChangesType changelog.ChangeType
	currentChanges     *changelog.ChangeCollection

	versionRegexp *regexp.Regexp
}

func newParser() *parser {
	return &parser{
		markdown:      goldmark.New(),
		changelog:     changelog.NewChangelog(),
		versionRegexp: regexp.MustCompile(`\[(Unreleased|\d+.\d+.\d+)\](?: - (\d{4}-\d{2}-\d{2}))?`),
	}
}

func (p *parser) parse(source []byte) (*changelog.Changelog, error) {
	node := p.markdown.Parser().Parse(text.NewReader(source))
	err := ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch n.Kind() {
		case ast.KindHeading:
			return p.Heading(source, n.(*ast.Heading), entering)
		case ast.KindList:
			return p.List(source, n, entering)
		case ast.KindListItem:
			return p.ListItem(source, n.(*ast.ListItem), entering)
		case ast.KindParagraph:
			return p.Paragraph(source, n.(*ast.Paragraph), entering)
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	return p.changelog, nil
}

// Heading is called for each Heading node.
func (p *parser) Heading(src []byte, block *ast.Heading, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}

	switch block.Level {
	case 1: // Document title
		// We don't care about changelog title
		p.changelog.Title = string(block.Text(src))
		return ast.WalkSkipChildren, nil
	case 2: // It's version
		p.currentVersion = changelog.NewVersion()
		matches := p.versionRegexp.FindStringSubmatch(string(block.Text(src)))
		p.currentVersion.Name = matches[1]
		p.currentVersion.Date = matches[2]

		p.changelog.Versions = append(p.changelog.Versions, p.currentVersion)

		return ast.WalkSkipChildren, nil
	case 3: // It's type of changes
		text := block.Text(src)
		p.currentChangesType = changelog.ChangeTypeFromString(string(text))
		if p.currentChangesType == changelog.ChangeTypeUnknown {
			return ast.WalkStop, fmt.Errorf("unknown type of changes: %q", string(text))
		}

		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (p *parser) List(src []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	p.currentChanges = &changelog.ChangeCollection{
		Type: p.currentChangesType,
	}
	p.currentVersion.Changes = append(p.currentVersion.Changes, p.currentChanges)
	return ast.WalkContinue, nil
}

func (p *parser) ListItem(src []byte, block *ast.ListItem, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	p.currentChanges.Items = append(p.currentChanges.Items, string(block.Text(src)))
	return ast.WalkContinue, nil
}

func (p *parser) Paragraph(src []byte, block *ast.Paragraph, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	ps := block.PreviousSibling()
	if (ps.Kind() == ast.KindHeading) && (ps.(*ast.Heading).Level == 1) {
		p.changelog.Preamble = string(block.Text(src))
	}
	return ast.WalkContinue, nil
}
