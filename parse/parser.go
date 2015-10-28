package parse
import (
	"strings"
	"fmt"
	"os"
	"strconv"
	"github.com/qingshan/sieve/ast"
)

// parser holds the state of the scanner.
type Parser struct {
	name      string        // the name of the input; used only for error reports
	input     string        // the string being scanned
	pos       int           // the position of token in Items; pos == -1 when Items is nil
	Items     []Token       // the unreduced items received from the lexer
	lastToken Token         // Used for error and debug messages
	Lexer     *Lexer        // the lexer
	List      []ast.Command // the file being parsed
}

// ------------------------------------------------------------------------------
// parsing support

// next return the next token in the input.
// calls getItem if at end of Items
func (p *Parser) next() Token {
	if p.pos != -1 && p.Items[p.pos].Typ == EOF {
		return p.Items[p.pos]
	}
	p.pos += 1
	if p.pos >= len(p.Items) {
		p.errorf("Internal error in next(): parser.pos moving out of bounds of lexed tokens\n")
	}
	// call p.errorf if lexing error
	if p.Items[p.pos].Typ == ERROR {
		p.errorf(p.Items[p.pos].String())
	}
	p.lastToken = p.Items[p.pos]
	return p.Items[p.pos]
}

// peek returns the k forward token in items but does not move the pos.
func (p *Parser) peek(k int) Token {
	if p.pos + k >= len(p.Items) {
		return p.Items[p.pos]
	}
	return p.Items[p.pos + k]
}

// backup steps back one token.
// Can only be called as many times as there are unreduced tokens in Items
// return error if there aren't enough tokens in Items
func (p *Parser) backup() error {
	if p.pos <= -1 {
		p.errorf("Internal error in backup: Cannot backup anymore pos is at start of Items\n")
	}
	p.pos -= 1
	// if p.pos != -1 {
	// 	p.lastToken = p.Items[p.pos]
	// }
	return nil
}

// accept consumes the next token if it's from the valid set.
func (p *Parser) accept(valid []TokenType) bool {
	item := p.next()
	for _, tokTyp := range valid {
		if item.Typ == tokTyp {
			return true
		}
	}
	p.backup()
	return false
}

// acceptRun consumes a run of tokens from the valid set.
func (p *Parser) acceptRun(valid []TokenType) {
	for p.accept(valid) {
	}
}

func (p *Parser) expect(valid TokenType) bool {
	if p.Items[p.pos].Typ == valid {
		return true
	}
	return false
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by next.
func (p *Parser) lineNumber() int {
	item := p.Items[p.pos]
	return 1 + strings.Count(p.input[:item.Pos], "\n")
}

// colNumber reports which column on the current line we're on,
// based on the position of the previous item returned by next
func (p *Parser) colNumber() int {
	ln := p.lineNumber()
	lines := strings.SplitN(p.input, "\n", ln)
	var total int
	for i := range lines[:ln - 1] {
		total += len(lines[i])
	}
	return int(p.pos) - total - (ln - 1)
}

// lineNumber reports which line we're on, based a lex.Pos
func (p *Parser) lineNumberAt(pos Pos) int {
	return 1 + strings.Count(p.input[:pos], "\n")
}

// errorf prints an error and terminates the scan
func (p *Parser) errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

// Parse creates a new parser for the input string.
// It uses lex to tokenize the input
func Parse(name, input string) *Parser {
	l := Lex(name, input)
	p := &Parser{
		name:     name,
		input:    input,
		pos:      -1,
		Lexer:    l,
	}
	p.run()
	return p
}

// runs the parser
func (p *Parser) run() {
	// lex everything
	t := p.Lexer.NextItem()
	for ; t.Typ != EOF; t = p.Lexer.NextItem() {
		p.Items = append(p.Items, t)
	}
	p.Items = append(p.Items, t)

	parseFile(p)
	return
}

// --------------------------------------------------------------------------------------------
// Recursive descent parser
// Mutually recursive functions

func parseFile(p *Parser) {
	switch t := p.next(); {
	case t.IsCommand() || t.Typ == LINECOMMENT || t.Typ == BLOCKCOMMENT || t.Typ == IDENTIFIER:
		p.backup()
		command := parseCommand(p)
		p.List = append(p.List, command)
	case t.Typ == EOF:
		return
	default:
		p.errorf("Invalid statement at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
}

func parseCommand(p *Parser) ast.Command {
	switch t := p.next(); {
	case t.Typ == IF || t.Typ == ELSE || t.Typ == ELSIF:
		name := t.Val
		var test ast.Test
		if t.Typ == IF || t.Typ == ELSIF {
			test = parseTest(p)
		}
		block := parseBlock(p)
		return &ast.ControlCommand{name, test, block}
	case t.Typ == LINECOMMENT:
		text := t.Val
		return &ast.CommentCommand{"line", text}
	case t.Typ == BLOCKCOMMENT:
		text := t.Val
		return &ast.CommentCommand{"block", text}
	case t.Typ == STOP:
		if !p.accept([]TokenType{SEMICOLON}) {
			p.errorf("Invalid command at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			return nil
		}
		return &ast.StopCommand{}
	case t.Typ == IDENTIFIER:
		name := t.Val
		al := parseArguments(p)
		if !p.accept([]TokenType{SEMICOLON}) {
			p.errorf("Invalid command at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			return nil
		}
		return &ast.GenericCommand{name, al}
	default:
		p.errorf("Invalid command at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		return nil
	}
}

func parseBlock(p *Parser) []ast.Command {
	var cl []ast.Command
	switch t := p.next(); {
	case t.Typ == LEFTCURLY:
	default:
		p.errorf("Invalid block statement at line %d:%d with token '%s', in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}

	Loop:
	for {
		cl = append(cl, parseCommand(p))
		switch t := p.next(); {
		case t.Typ == RIGHTCURLY:
			break Loop
		default:
			p.backup()
		}
	}
	return cl;
}

func parseTest(p *Parser) ast.Test {
	switch t := p.next(); {
	case t.Typ == NOT:
		return &ast.NotTest{parseTest(p)}
	case t.Typ == ANYOF:
		return &ast.AnyofTest{parseTests(p)}
	case t.Typ == ALLOF:
		return &ast.AllofTest{parseTests(p)}
	case t.Typ == TRUE:
		return &ast.TrueTest{}
	case t.Typ == FALSE:
		return &ast.FalseTest{}
	case t.Typ == IDENTIFIER:
		name := t.Val
		al := parseArguments(p)
		return &ast.GenericTest{name, al}
	default:
		p.errorf("Invalid Test statement at line %d:%d with token '%s', in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		return nil
	}
}

func parseTests(p *Parser) []ast.Test {
	var tl []ast.Test
	switch t := p.next(); {
	case t.Typ == LEFTPAREN:
	default:
		p.errorf("Invalid Tests statement at line %d:%d with token '%s', in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
	Loop:
	for {
		tl = append(tl, parseTest(p))
		switch t := p.next(); {
		case t.Typ == COMMA:
		//absorb
		case t.Typ == RIGHTPAREN:
			break Loop
		default:
			p.errorf("Invalid Tests statement at line %d:%d with token '%s', in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	}
	return tl;
}

func parseArguments(p *Parser) []ast.Argument {
	var al []ast.Argument
	Loop:
	for {
		switch t := p.next(); {
		case t.Typ == NUMBER:
			al = append(al, ast.NumberArgument(t.Val))
		case t.Typ == TAG:
			al = append(al, ast.TagArgument(t.Val))
		case t.Typ == STRING:
			s, err := strconv.Unquote(t.Val)
			if err != nil {
				p.errorf(err.Error())
				return al;
			}
			al = append(al, ast.StringArgument([]string{s}))
		case t.Typ == LEFTBRACKET:
			p.backup()
			al = append(al, ast.StringArgument(parseStrings(p)))
		case atTerminator(t):
			p.backup()
			break Loop
		default:
			p.errorf("Invalid Arguments statement at line %d:%d with token '%s', in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	}
	return al;
}


func parseStrings(p *Parser) []string {
	var sl []string
	switch t := p.next(); {
	case t.Typ == LEFTBRACKET:
	//absorb
	default:
		p.errorf("Invalid Strings statement at line %d:%d with token '%s', in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
	Loop:
	for {
		switch t := p.next(); {
		case t.Typ == STRING:
			s, err := strconv.Unquote(t.Val)
			if err != nil {
				p.errorf(err.Error())
				return sl;
			}
			sl = append(sl, s)
		default:
			p.errorf("Invalid Strings statement at line %d:%d with token '%s', in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
		switch t := p.next(); {
		case t.Typ == COMMA:
		//absorb
		case t.Typ == RIGHTBRACKET:
			break Loop
		default:
			p.errorf("Invalid Strings statement at line %d:%d with token '%s', in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	}
	return sl
}

func atTerminator(t Token) bool {
	if t.Typ == SEMICOLON || t.Typ == COMMA || t.Typ == LEFTCURLY || t.Typ == EOF {
		return true
	}
	return false
}