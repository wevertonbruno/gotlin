package scanner

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"gotlin/frontend/token"
)

const (
	Eof = '\u0000'
	Eol = '\n'
)

type Scanner struct {
	current byte
	peek    byte
	line    uint
	col     uint
	tokens  []token.Token
	reader  *bufio.Reader
}

func NewScanner(reader io.Reader) *Scanner {
	s := &Scanner{
		reader: bufio.NewReader(reader),
		line:   1,
		col:    1,
	}
	s.advance()
	s.advance()
	return s
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.scan()
		s.advance()
	}

	if !s.lastMatch(token.NEWLINE) {
		s.addToken(token.NEWLINE)
	}

	s.addToken(token.EOF)
	return s.tokens
}

func (s *Scanner) scan() {
	switch s.current {
	case '(':
		s.addToken(token.LParen)
		break
	case ')':
		s.addToken(token.RParen)
		break
	case '{':
		s.addToken(token.LBrace)
		break
	case '}':
		s.addToken(token.RBrace)
		break
	case ';':
		s.addToken(token.SEMICOLON)
		break
	case ':':
		s.addToken(token.COLON)
		break
	case ',':
		s.addToken(token.COLON)
		break
	case '.':
		s.addToken(token.DOT)
		break
	case '-':
		if s.match('=') {
			s.addToken(token.MINUS_ASSIGN)
		} else {
			s.addToken(token.OpMinus)
		}
		break
	case '+':
		if s.match('=') {
			s.addToken(token.PLUS_ASSIGN)
		} else {
			s.addToken(token.OpPlus)
		}
		break
	case '/':
		if s.match('/') {
			for ; s.current != '\n' && !s.isAtEnd(); s.advance() {
			}
		} else {
			s.addToken(token.OpDivide)
		}
		break
	case '*':
		s.addToken(token.OpMulti)
		break
	case '"':
		s.addTokenString()
		break
	case '!':
		if s.match('=') {
			s.addToken(token.OpNotEq)
		} else if s.match('!') {
			s.addToken(token.OpNotNull)
		} else {
			s.addToken(token.NOT)
		}
		break
	case '=':
		if s.match('=') {
			s.addToken(token.OpEq)
		} else {
			s.addToken(token.ASSIGN)
		}
		break
	case '<':
		if s.match('=') {
			s.addToken(token.OpLte)
		} else {
			s.addToken(token.OpLt)
		}
		break
	case '>':
		if s.match('=') {
			s.addToken(token.OpGte)
		} else {
			s.addToken(token.OpGt)
		}
		break
	case '&':
		if s.match('&') {
			s.addToken(token.AND)
		} else {
			s.addToken(token.ERROR)
		}
		break
	case '|':
		if s.match('|') {
			s.addToken(token.OR)
		} else {
			s.addToken(token.ERROR)
		}
		break
	case '?':
		if s.match(':') {
			s.addToken(token.ELVIS)
		} else if s.match('.') {
			s.addToken(token.OpSafeCall)
		} else {
			s.addToken(token.QUESTION)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		// whitespace
		break
	case '\n':
		if !s.lastMatch(token.NEWLINE) {
			s.addToken(token.NEWLINE)
		}
		break
	default:
		if isAlpha(s.current) {
			s.addTokenIdentifier()
			return
		}

		if isDigit(s.current) {
			s.addTokenNumber()
			return
		}

		s.addToken(token.ERROR)
		break
	}
}

func (s *Scanner) lastMatch(kind token.Kind) bool {
	return len(s.tokens) > 0 && s.tokens[len(s.tokens)-1].Kind == kind
}

func (s *Scanner) match(expected byte) bool {
	if s.peek == expected {
		s.advance()
		return true
	}
	return false
}

func (s *Scanner) read() byte {
	b, err := s.reader.ReadByte()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return Eof
		}
		panic(err)
	}
	return b
}

func (s *Scanner) advance() byte {
	b := s.read()
	if b == Eol {
		s.line++
		s.col = 1
	} else {
		s.col++
	}
	s.current = s.peek
	s.peek = b
	return s.current
}

//TODO REMOVE
//func (s *Scanner) peekNext() byte {
//	b, err := s.reader.Peek(1)
//	if err != nil {
//		if err == io.EOF {
//			return Eof
//		}
//		panic(err)
//	}
//	return b[0]
//}

func (s *Scanner) isAtEnd() bool {
	return s.current == Eof
}

func (s *Scanner) addToken(kind token.Kind) {
	if kind == token.ERROR {
		s.tokens = append(s.tokens, token.Token{
			Kind:     kind,
			Spelling: string(s.current),
			Position: token.Pos{Line: s.line, Col: s.col},
		})
		return
	}
	s.tokens = append(s.tokens, token.NewToken(kind, s.line, s.col))
}

func (s *Scanner) addTokenString() {
	var sb strings.Builder
	s.advance()

	for s.current != '"' && !s.isAtEnd() {
		if s.current == '\n' {
			s.line++
			s.col = 1
		}
		sb.WriteByte(s.current)
		s.advance()
	}

	if s.isAtEnd() {
		panic("Unterminated string") // TODO REFACTOR
	}

	s.advance()
	s.tokens = append(s.tokens, token.NewTokenLiteral(token.StringLit, sb.String(), s.line, s.col))
}

func (s *Scanner) addTokenNumber() {
	var sb strings.Builder
	sb.WriteByte(s.current)
	hasDot := false
	for {
		if s.peek == '.' {
			hasDot = true
		} else if !isDigit(s.peek) {
			break
		}
		sb.WriteByte(s.peek)
		s.advance()
	}

	if hasDot {
		s.tokens = append(s.tokens, token.NewTokenLiteral(token.DoubleLit, sb.String(), s.line, s.col))
	} else {
		s.tokens = append(s.tokens, token.NewTokenLiteral(token.IntLit, sb.String(), s.line, s.col))
	}
}

func (s *Scanner) addTokenIdentifier() {
	var sb strings.Builder
	sb.WriteByte(s.current)

	for {
		if !isAlpha(s.peek) && !isDigit(s.peek) {
			break
		}
		sb.WriteByte(s.peek)
		s.advance()
	}

	s.tokens = append(s.tokens, token.NewTokenLiteral(token.IDENTIFIER, sb.String(), s.line, s.col))
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isAlpha(b byte) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_'
}
