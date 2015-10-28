package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// lexer holds the state of the scanner.
type Lexer struct {
	name       string     // the name of the input; used only for error reports
	input      string     // the string being scanned
	state      stateFn    // the next lexing function to enter
	pos        Pos        // current position in the input
	start      Pos        // start position of this item
	width      Pos        // width of last rune read from input
	lastPos    Pos        // position of most recent item returned by nextItem
	items      chan Token // channel of scanned items
	parenDepth int        // nesting depth of ( ) exprs
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *Lexer) emit(t TokenType) {
	l.items <- Token{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by nextItem. Doing it this way
// means we don't have to worry about peek double counting.
func (l *Lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}

// colNumber reports which column on the current line we're on,
// based on the position of the current rune
func (l *Lexer) colNumber() int {
	ln := l.lineNumber()
	lines := strings.SplitN(l.input, "\n", ln)
	var total int
	for i := range lines[:ln - 1] {
		total += len(lines[i])
	}
	return int(l.pos) - total - (ln - 1)
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Token{ERROR, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *Lexer) NextItem() Token {
	token := <-l.items
	l.lastPos = token.Pos
	return token
}

// lex creates a new scanner for the input string.
func Lex(name, input string) *Lexer {
	l := &Lexer{
		name:  name,
		input: input,
		items: make(chan Token),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *Lexer) run() {
	for l.state = lexStart; l.state != nil; {
		l.state = l.state(l)
	}
}

// state functions

func lexStart(l *Lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(EOF)
		return nil
	case r == ';':
		return lexSemiColon
	case r == ':':
		return lexTag
	case r == ',':
		l.emit(COMMA)
		return lexStart
	case r == '"':
		return lexString
	case r == '(':
		l.emit(LEFTPAREN)
		l.parenDepth++
		return lexStart
	case r == ')':
		l.emit(RIGHTPAREN)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren at line %d:%d with %#U", l.lineNumber(), l.colNumber(), r)
		}
		return lexStart
	case r == '[':
		l.emit(LEFTBRACKET)
		return lexStart
	case r == ']':
		l.emit(RIGHTBRACKET)
		return lexStart
	case r == '{':
		l.emit(LEFTCURLY)
		return lexStart
	case r == '}':
		l.emit(RIGHTCURLY)
		return lexStart
	case isSpace(r):
		return lexSpace
	case isEndOfLine(r):
		return lexEndOfLine
	case ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
	case r == '#':
		return lexLineComment
	case r == '/' && l.peek() == '*':
		return lexBlockComment
	default:
		return l.errorf("unknown syntax: %q", l.input[l.start:l.pos])
	}
}

// lexSemiColon scans a semicolon
func lexSemiColon(l *Lexer) stateFn {
	l.emit(SEMICOLON)
	return lexStart
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *Lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.ignore()
	return lexStart
}

// lexEndOfLine scans a end of line character.
func lexEndOfLine(l *Lexer) stateFn {
	for isEndOfLine(l.peek()) {
		l.next()
	}
	l.ignore()
	return lexStart
}

// lexNumber scans a number: decimal with optional KMG
//
func lexNumber(l *Lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax at %d:%d with %q", l.lineNumber(), l.colNumber(), l.input[l.start:l.pos])
	}
	l.emit(NUMBER)
	return lexStart
}

func (l *Lexer) scanNumber() bool {
	digits := "0123456789"
	l.acceptRun(digits)
	l.accept("kKmMgG")
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

func lexString(l *Lexer) stateFn {
	Loop:
	for {
		switch r := l.next(); {
		case r != '"':
		// absorb.
		case r == eof:
			return l.errorf("Non-terminating string literal at %#U", r)
		default:
			l.emit(STRING)
			break Loop
		}
	}
	return lexStart
}

// lexTag scans an tag.
func lexTag(l *Lexer) stateFn {
	Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
		// absorb.
		default:
			l.backup()
			l.emit(TAG)
			break Loop
		}
	}
	return lexStart
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *Lexer) stateFn {
	Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
		// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			// if !l.atTerminator() {
			// 	return l.errorf("bad character %#U", r)
			// }
			switch {
			case key[word] > COMMAND:
				l.emit(key[word])
			default:
				l.emit(IDENTIFIER)
			}
			break Loop
		}
	}
	return lexStart
}

// atTerminator reports whether the input is at valid termination character to
// appear after an identifier. Breaks .X.Y into two pieces. Also catches cases
// like "$x+2" not being acceptable without a space, in case we decide one
// day to implement arithmetic.
func (l *Lexer) atTerminator() bool {
	r := l.peek()
	if isSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case eof, '.', ',', ':', ')', '(':
		return true
	}
	return false
}

func lexLineComment(l *Lexer) stateFn {
	Loop:
	for {
		switch r := l.next(); {
		case !isEndOfLine(r):
		// absorb.
		default:
			l.backup()
			if !l.atTerminator() {
				return l.errorf("bad character %#U at %d:%d", r, l.lineNumber(), l.colNumber())
			}
			l.emit(LINECOMMENT)
			break Loop
		}
	}
	return lexStart
}

func lexBlockComment(l *Lexer) stateFn {
	Loop:
	for {
		// if we find '*' and the next is  '/'
		switch r := l.next(); {
		case !l.atEndBlockComment():
		// absorb.
		case r == eof:
			return l.errorf("Non-terminating block comment at %#U", r)
		default:
			// l.backup()
			// l.next()
			word := l.input[l.start:l.pos]
			switch {
			case strings.Index(word, "*/") == len(word) - len("*/"):
				l.emit(BLOCKCOMMENT)
			default:
				return l.errorf("error in  block comment at %#U", r)
			}
			break Loop
		}
	}
	return lexStart
}

func (l *Lexer) atEndBlockComment() bool {
	word := l.input[l.pos - 2 : l.pos]
	if strings.Index(word, "*/") == len(word) - len("*/") {
		return true
	}
	return false
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}