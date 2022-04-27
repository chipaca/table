package table_test

import (
	"io"
	"os"

	"github.com/chipaca/escapes"
	"github.com/chipaca/table"
)

var countryData = [][]interface{}{
	{"004", "AF", "Afghanistan", "Kabul"},
	{"248", "AX", "Åland Islands", "Mariehamn"},
	{"008", "AL", "Albania", "Tirana"},
	{"012", "DZ", "Algeria", "Algiers"},
	{"016", "AS", "American Samoa", "Pago Pago"},
	{"020", "AD", "Andorra", "Andorra la Vella"},
	{"024", "AO", "Angola", "Luanda"},
}

func ExampleNew() {
	t := table.New(escapes.Bold("numeric"), escapes.Bold("alpha-2"), escapes.Bold("name"), escapes.Bold("capital"))
	t.Outdent = "\u200b" // this is to trick go's example tests' overly aggressive whitespace trimming
	t.Align[0] = table.ColumnAlignRight
	for _, row := range countryData {
		t.Append(row...)
	}
	t.Print(os.Stdout)
	// Output:
	//  [1mnumeric[0m  [1malpha-2[0m  [1mname[0m            [1mcapital[0m          ​
	//      004  AF       Afghanistan     Kabul            ​
	//      248  AX       Åland Islands   Mariehamn        ​
	//      008  AL       Albania         Tirana           ​
	//      012  DZ       Algeria         Algiers          ​
	//      016  AS       American Samoa  Pago Pago        ​
	//      020  AD       Andorra         Andorra la Vella ​
	//      024  AO       Angola          Luanda           ​
}

func ExampleNewGHFMD() {
	t := table.NewGHFMD("numeric", "alpha-2", "name", "capital")
	t.Outdent = "\u200b" // this is to trick go's example tests' overly aggressive whitespace trimming
	t.Align[0] = table.ColumnAlignRight
	for _, row := range countryData {
		t.Append(row...)
	}
	t.Print(os.Stdout)
	// Output:
	//  numeric | alpha-2 | name           | capital          ​
	// --------:|---------|----------------|------------------​
	//      004 | AF      | Afghanistan    | Kabul            ​
	//      248 | AX      | Åland Islands  | Mariehamn        ​
	//      008 | AL      | Albania        | Tirana           ​
	//      012 | DZ      | Algeria        | Algiers          ​
	//      016 | AS      | American Samoa | Pago Pago        ​
	//      020 | AD      | Andorra        | Andorra la Vella ​
	//      024 | AO      | Angola         | Luanda           ​
}

func ExampleNewGHFMD_BeginAndEnd() {
	t := table.NewGHFMD("numeric", "alpha-2", "name", "capital")
	t.Align[0] = table.ColumnAlignRight
	t.BeginRow = func(out io.Writer) {
		out.Write([]byte{'>'})
	}
	t.EndRow = func(out io.Writer) {
		out.Write([]byte{'<'})
	}
	for _, row := range countryData {
		t.Append(row...)
	}
	t.Print(os.Stdout)
	// Output:
	// > numeric | alpha-2 | name           | capital          <
	// >--------:|---------|----------------|------------------<
	// >     004 | AF      | Afghanistan    | Kabul            <
	// >     248 | AX      | Åland Islands  | Mariehamn        <
	// >     008 | AL      | Albania        | Tirana           <
	// >     012 | DZ      | Algeria        | Algiers          <
	// >     016 | AS      | American Samoa | Pago Pago        <
	// >     020 | AD      | Andorra        | Andorra la Vella <
	// >     024 | AO      | Angola         | Luanda           <
}
