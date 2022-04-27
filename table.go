// Package table lets you create simple text-based tables that are
// width-aware, meaning your cells can have escape sequences or emoji
// in them without everything breaking.
//
// It has only been tested with left-to-right text.
package table

import (
	"fmt"
	"io"
	"strings"

	"github.com/chipaca/width"
	"github.com/mattn/go-runewidth"
)

// ColumnAlign is used to let columns align in different ways.
type ColumnAlign int8

const (
	// ColumnAlignRight means to align things to the right of a
	// cell, like for whole numbers.
	ColumnAlignRight = ColumnAlign(-1)
	// ColumnAlignCenter would align things in the middle of a
	// cell. Not currently implemented.
	ColumnAlignCenter = ColumnAlign(0)
	// ColumnAlignLeft is the default.
	ColumnAlignLeft = ColumnAlign(1)
)

// A Rule contains information for drawing a horizontal rule between
// the header and the data in the table
type Rule struct {
	// the rule Gutter, if set, must have the same width as the table gutter.
	// If not set, falls back to the table gutter.
	Gutter string
	// all these runes should have width 1
	// Rule is the rule itself; if not set, no horizontal rule is drawn
	Rule rune
	// if not set these fall back to Rule
	RightAlignedLeftPad  rune
	RightAlignedRightPad rune
	LeftAlignedLeftPad   rune
	LeftAlignedRightPad  rune
}

func (r *Rule) isSet() bool {
	return r.Rule != 0
}

// setDefaults defaults all the pads to the rule
func (r *Rule) setDefaults() {
	if r.RightAlignedLeftPad == 0 {
		r.RightAlignedLeftPad = r.Rule
	}
	if r.RightAlignedRightPad == 0 {
		r.RightAlignedRightPad = r.Rule
	}
	if r.LeftAlignedLeftPad == 0 {
		r.LeftAlignedLeftPad = r.Rule
	}
	if r.LeftAlignedRightPad == 0 {
		r.LeftAlignedRightPad = r.Rule
	}
}

// A Table is a way to print tabular data to the terminal, where each column of
// data has a header, there is a minimum fixed gutter between columns (that can
// be empty), and all columns except the last are reasonably short such that no
// truncation is necessary to fit them, and a truncated last column, in the
// terminal.
type Table struct {
	// Gutter is used to separate cells
	Gutter string
	// Indent is added to the beginning of every line
	Indent string
	// Outdent is added to the end of every line
	Outdent string
	// Rule can be set for drawing a line between headers and data
	// (see details in Rule itself)
	Rule Rule
	// Alignment of each column (including headers)
	Align []ColumnAlign
	// TermWidth is the width of the terminal
	TermWidth int
	// BeginRow, is called before printing each row
	BeginRow func(out io.Writer)
	// EndRow is called after printing each row (before printing the newline)
	EndRow func(out io.Writer)
	// Space is the rune used to fill a cell to a uniform width. Must have width 1.
	Space rune
	// Pad is the rune used as padding around cells. Must be width 1.
	Pad rune

	maxw    []int
	headers []string
	headerW []int
	data    [][]string
	widths  [][]int
}

func stringAndWidth(a interface{}) (string, int) {
	switch b := a.(type) {
	case width.StringWidther:
		return b.String(), b.Width()
	case string:
		return b, runewidth.StringWidth(b)
	default:
		panic(fmt.Sprintf("can't determine width from %T", b))
	}
}

func nopLiner(io.Writer) {}

// New initializes a Table with reasonable defaults for the given headers.
//
// The headers can be a string, or a StringAndWidther; anything else will panic.
func New(headers ...interface{}) *Table {
	maxw := make([]int, len(headers))
	headerW := make([]int, len(headers))
	align := make([]ColumnAlign, len(headers))
	strheaders := make([]string, len(headers))
	for i := range headers {
		strheaders[i], headerW[i] = stringAndWidth(headers[i])
		maxw[i] = headerW[i]
		align[i] = ColumnAlignLeft
	}
	return &Table{
		headers:   strheaders,
		headerW:   headerW,
		maxw:      maxw,
		Align:     align,
		TermWidth: 80,
		BeginRow:  nopLiner,
		EndRow:    nopLiner,
		Space:     ' ',
		Pad:       ' ',
	}
}

// NewGHFMD initializes a Table with reasonable defaults for the given headers
// to produce a table formatted for GitHub-flavoured markdown.
func NewGHFMD(headers ...string) *Table {
	maxw := make([]int, len(headers))
	headerW := make([]int, len(headers))
	align := make([]ColumnAlign, len(headers))
	for i := range headers {
		headerW[i] = runewidth.StringWidth(headers[i])
		maxw[i] = headerW[i]
		align[i] = ColumnAlignLeft
	}
	return &Table{
		headers:   headers,
		headerW:   headerW,
		maxw:      maxw,
		Align:     align,
		TermWidth: 200,
		Rule: Rule{
			Rule:                 '-',
			RightAlignedRightPad: ':',
		},
		Gutter:   "|",
		BeginRow: nopLiner,
		EndRow:   nopLiner,
		Space:    '\u00a0',
		Pad:      '\u00a0',
	}
}

// Len returns the number of rows of data in the table.
func (t *Table) Len() int {
	return len(t.data)
}

// Set lets you assign a value to a particular cell of the data in the table.
func (t *Table) Set(i, j int, v string) {
	t.data[i][j] = v
	t.incMaxW(i, j, runewidth.StringWidth(v))
}

// incMaxW will adjust the max width for the given data.
func (t *Table) incMaxW(i, j int, w int) {
	// TODO: shrink
	t.widths[i][j] = w
	if t.maxw[j] < w {
		t.maxw[j] = w
	}
}

// Append adds a row to the table's data.
//
// Note it must have the same number of entries as there were headers
// in the constructor.
//
// Also note the entries can be a string, or a StringAndWidther.
func (t *Table) Append(row ...interface{}) {
	if len(row) != len(t.headers) {
		panic("row has wrong number of items")
	}
	t.widths = append(t.widths, make([]int, len(row)))
	strrow := make([]string, len(row))
	var w int
	for i := range row {
		strrow[i], w = stringAndWidth(row[i])
		t.incMaxW(len(t.data), i, w)
	}
	t.data = append(t.data, strrow)
}

func (t *Table) printLine(out io.Writer, tpl string, row []interface{}) {
	t.BeginRow(out)
	fmt.Fprintf(out, tpl, row...)
	t.EndRow(out)
	out.Write([]byte{'\n'})
}

// Print the table.
//
// It attempts to determine the terminal width, falling back to 80 columns if it
// fails. The last column of data will be truncated so the whole table fits
// inside this width; note this might not be possible if the non-truncated
// columns already overflow.
func (t *Table) Print(out io.Writer) {
	// 'last' is two things:
	// 1. the index of the last item in each row (and t.headers)
	// 2. the number of items in max (and the number of untrucated columns)
	last := len(t.maxw) - 1
	// each cell has a left pad char, data, right pad char. Then there's the gutter between cells.
	rowTemplate := t.Indent + strings.Repeat("%c%s%c%s", last) + "%c%s%c" + t.Outdent

	width := t.TermWidth
	if width < 80 {
		// if the terminal is too narrow, give up on making things
		// pretty and just aim for 'standard' 80 columns
		// NOTE when not connected to a terminal (and other errors), on
		// ~unix width will be zero, whereas on windows it'll be 80.
		width = 80
	}

	// see how much width we have left for the last column
	width -= runewidth.StringWidth(t.Indent) + last*(2+runewidth.StringWidth(t.Gutter)) + 2 + runewidth.StringWidth(t.Outdent)
	for _, w := range t.maxw[:len(t.maxw)-1] {
		width -= w
	}
	overflow := make([]bool, len(t.data))
	if width < t.maxw[last] {
		// need to shorten the last column
		// 1. reset column width to header width
		t.maxw[last] = t.headerW[last]
		// 2. shorten
		for i := range t.data {
			o := runewidth.Truncate(t.data[i][last], width, "")
			if o != t.data[i][last] {
				overflow[i] = true
				// note this is what increments maxw[last]
				t.Set(i, last, o)
			}
		}
	}
	row := make([]interface{}, 4*len(t.maxw)-1)
	for j := 0; j <= last; j++ {
		space := strings.Repeat(string(t.Space), t.maxw[j]-t.headerW[j])
		row[4*j] = t.Pad
		switch t.Align[j] {
		case ColumnAlignRight:
			row[4*j+1] = space + t.headers[j]
		case ColumnAlignLeft:
			row[4*j+1] = t.headers[j] + space
		default:
			panic("only left- and right-align implemented yet")
		}
		row[4*j+2] = t.Pad
		if j != last {
			row[4*j+3] = t.Gutter
		}
	}
	t.printLine(out, rowTemplate, row)
	if t.Rule.isSet() {
		t.Rule.setDefaults()
		for j := 0; j <= last; j++ {
			row[4*j+1] = strings.Repeat(string(t.Rule.Rule), t.maxw[j])
			switch t.Align[j] {
			case ColumnAlignRight:
				row[4*j] = t.Rule.RightAlignedLeftPad
				row[4*j+2] = t.Rule.RightAlignedRightPad
			case ColumnAlignLeft:
				row[4*j] = t.Rule.LeftAlignedLeftPad
				row[4*j+2] = t.Rule.LeftAlignedRightPad
			}
			if j != last && t.Rule.Gutter != "" {
				row[4*j+3] = t.Rule.Gutter
			}
		}
		t.printLine(out, rowTemplate, row)
		for j := 0; j <= last; j++ {
			row[4*j] = t.Pad
			row[4*j+2] = t.Pad
			if j != last {
				row[4*j+3] = t.Gutter
			}
		}
	}

	for i := range t.data {
		for j := 0; j <= last; j++ {
			space := strings.Repeat(string(t.Space), t.maxw[j]-t.widths[i][j])
			switch t.Align[j] {
			case ColumnAlignRight:
				row[4*j+1] = space + t.data[i][j]
			case ColumnAlignLeft:
				row[4*j+1] = t.data[i][j] + space
			}
		}
		if overflow[i] {
			row[4*last+2] = '…'
		} else {
			row[4*last+2] = t.Pad
		}
		t.printLine(out, rowTemplate, row)
	}
}

func (t *Table) String() string {
	var b strings.Builder
	t.Print(&b)
	return b.String()
}
