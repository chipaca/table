package table_test

import (
	"testing"

	"github.com/chipaca/table"
)

func TestTableRules(t *testing.T) {
	tbl := table.New("foo", "bar", "baz", "meh")
	tbl.Gutter = "|"
	tbl.Rule.Rule = '-'
	tbl.Rule.Gutter = "+"
	tbl.Indent = ">"
	tbl.Outdent = "<"
	tbl.Space = '_'
	tbl.Pad = '#'
	tbl.TermWidth = 80
	tbl.Align = []table.ColumnAlign{
		table.ColumnAlignRight,
		table.ColumnAlignCenter,
		table.ColumnAlignLeft,
		table.ColumnAlignRight,
	}

	tbl.Append("a", "b", "c", "d")
	tbl.Append("1", "2", "3", "4")

	expected := `
>#foo#|#bar#|#baz#|#meh#<
>-----+-----+-----+-----<
>#__a#|#_b_#|#c__#|#__d#<
>#__1#|#_2_#|#3__#|#__4#<
`[1:]
	got := tbl.String()
	if got != expected {
		t.Errorf("wanted:\n%+q\ngot:\n%#q", expected, got)
	}
}
