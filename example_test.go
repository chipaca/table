package table_test

import (
	"fmt"
	"os"

	"github.com/chipaca/table"
)

func Example() {
	t := table.New("numeric", "alpha-2", "name", "capital")
	t.Rule = &table.Rule{
		Rule:                 '-',
		RightAlignedLeftPad:  '-',
		RightAlignedRightPad: ':',
		LeftAlignedLeftPad:   '-',
		LeftAlignedRightPad:  '-',
	}
	t.Gutter = "|"
	t.Outdent = "\u200b" // this is to trick go's example tests' overly aggressive whitespace trimming
	t.Align[0] = table.ColumnAlignRight
	for _, row := range [][]interface{}{
		{"004", "AF", "Afghanistan", "Kabul"},
		{"248", "AX", "Åland Islands", "Mariehamn"},
		{"008", "AL", "Albania", "Tirana"},
		{"012", "DZ", "Algeria", "Algiers"},
		{"016", "AS", "American Samoa", "Pago Pago"},
		{"020", "AD", "Andorra", "Andorra la Vella"},
		{"024", "AO", "Angola", "Luanda"},
	} {
		t.Append(row...)
	}
	fmt.Println("github-flavoured markdown table:")
	t.Print(os.Stdout)
	fmt.Println(".")
	// Output:
	// github-flavoured markdown table:
	//  numeric | alpha-2 | name           | capital          ​
	// --------:|---------|----------------|------------------​
	//      004 | AF      | Afghanistan    | Kabul            ​
	//      248 | AX      | Åland Islands  | Mariehamn        ​
	//      008 | AL      | Albania        | Tirana           ​
	//      012 | DZ      | Algeria        | Algiers          ​
	//      016 | AS      | American Samoa | Pago Pago        ​
	//      020 | AD      | Andorra        | Andorra la Vella ​
	//      024 | AO      | Angola         | Luanda           ​
	// .
}
