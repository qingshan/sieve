package parse

import (
	"testing"
)

func TestNumbers(t *testing.T) {
	input := `1024 789k 512M 2G`
	lexer := Lex("TestNumbers", input)
	var output []Token
	expected := []Token{
		Token{Typ: NUMBER, Val: "1024"},
		Token{Typ: NUMBER, Val: "789k"},
		Token{Typ: NUMBER, Val: "512M"},
		Token{Typ: NUMBER, Val: "2G"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}


func TestComments(t *testing.T) {
	input := `
#aoeu
#/*test*/
/*test*/`
	lexer := Lex("TestComments", input)
	var output []Token
	expected := []Token{
		Token{Typ: LINECOMMENT, Val: "#aoeu"},
		Token{Typ: LINECOMMENT, Val: "#/*test*/"},
		Token{Typ: BLOCKCOMMENT, Val: "/*test*/"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF || item.Typ == ERROR {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestStringlists(t *testing.T) {
	input := `["Subject", "From", "To"]`

	lexer := Lex("TestStringlists", input)
	var output []Token
	expected := []Token{
		Token{Typ: LEFTBRACKET, Val: "["},
		Token{Typ: STRING, Val: "\"Subject\""},
		Token{Typ: COMMA, Val: ","},
		Token{Typ: STRING, Val: "\"From\""},
		Token{Typ: COMMA, Val: ","},
		Token{Typ: STRING, Val: "\"To\""},
		Token{Typ: RIGHTBRACKET, Val: "]"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF || item.Typ == ERROR {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}


func TestArguments(t *testing.T) {
	input := `"subject" :under 1024k`

	lexer := Lex("TestArguments", input)
	var output []Token
	expected := []Token{
		Token{Typ: STRING, Val: "\"subject\""},
		Token{Typ: TAG, Val: ":under"},
		Token{Typ: NUMBER, Val: "1024k"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF || item.Typ == ERROR {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}


func TestIfElse(t *testing.T) {
	input := `if not allof (true, false) { discard; stop; } else { stop; }`
	lexer := Lex("TestIfElse", input)
	var output []Token
	expected := []Token{
		Token{Typ: IF, Val: "if"},
		Token{Typ: NOT, Val: "not"},
		Token{Typ: ALLOF, Val: "allof"},
		Token{Typ: LEFTPAREN, Val: "("},
		Token{Typ: TRUE, Val: "true"},
		Token{Typ: COMMA, Val: ","},
		Token{Typ: FALSE, Val: "false"},
		Token{Typ: RIGHTPAREN, Val: ")"},
		Token{Typ: LEFTCURLY, Val: "{"},
		Token{Typ: IDENTIFIER, Val: "discard"},
		Token{Typ: SEMICOLON, Val: ";"},
		Token{Typ: STOP, Val: "stop"},
		Token{Typ: SEMICOLON, Val: ";"},
		Token{Typ: RIGHTCURLY, Val: "}"},
		Token{Typ: ELSE, Val: "else"},
		Token{Typ: LEFTCURLY, Val: "{"},
		Token{Typ: STOP, Val: "stop"},
		Token{Typ: SEMICOLON, Val: ";"},
		Token{Typ: RIGHTCURLY, Val: "}"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			//t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestAction(t *testing.T) {
	input := `size :under "Subject" 10M;`
	lexer := Lex("TestGenericCommand", input)
	var output []Token
	expected := []Token{
		Token{Typ: IDENTIFIER, Val: "size"},
		Token{Typ: TAG, Val: ":under"},
		Token{Typ: STRING, Val: "\"Subject\""},
		Token{Typ: NUMBER, Val: "10M"},
		Token{Typ: SEMICOLON, Val: ";"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			//t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}
