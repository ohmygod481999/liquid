//go:generate ragel -Z scanner.rl
//go:generate goyacc expressions.y

package expression

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// var lexerTests = []struct{}{
// 	{"{{var}}", "value"},
// 	{"{{x}}", "1"},
// }

type testSymbol struct {
	tok int
	typ yySymType
}

func scanExpression(data string) ([]testSymbol, error) {
	var (
		lex     = newLexer([]byte(data))
		symbols []testSymbol
		s       yySymType
	)
	for {
		tok := lex.Lex(&s)
		if tok == 0 {
			break
		}
		symbols = append(symbols, testSymbol{tok, s})
	}
	return symbols, nil
}

func TestLex(t *testing.T) {
	ts, err := scanExpression("abc > 123")
	require.NoError(t, err)
	require.Len(t, ts, 3)
	require.Equal(t, IDENTIFIER, ts[0].tok)
	require.Equal(t, "abc", ts[0].typ.name)
	require.Equal(t, LITERAL, ts[2].tok)
	require.Equal(t, 123, ts[2].typ.val)

	// verify these don't match "for", "or", or "false"
	ts, _ = scanExpression("forage")
	require.Len(t, ts, 1)
	ts, _ = scanExpression("orange")
	require.Len(t, ts, 1)
	ts, _ = scanExpression("falsehood")
	require.Len(t, ts, 1)

	ts, _ = scanExpression("a.b-c")
	require.Len(t, ts, 2)
	require.Equal(t, PROPERTY, ts[1].tok)
	require.Equal(t, "b-c", ts[1].typ.name)

	// literals
	ts, _ = scanExpression(`true false nil 2 2.3 "abc" 'abc'`)
	require.Len(t, ts, 7)
	require.Equal(t, LITERAL, ts[0].tok)
	require.Equal(t, LITERAL, ts[1].tok)
	require.Equal(t, LITERAL, ts[2].tok)
	require.Equal(t, LITERAL, ts[3].tok)
	require.Equal(t, LITERAL, ts[4].tok)
	require.Equal(t, LITERAL, ts[5].tok)
	require.Equal(t, LITERAL, ts[6].tok)
	require.Equal(t, true, ts[0].typ.val)
	require.Equal(t, false, ts[1].typ.val)
	require.Equal(t, nil, ts[2].typ.val)
	require.Equal(t, 2, ts[3].typ.val)
	require.Equal(t, 2.3, ts[4].typ.val)
	require.Equal(t, "abc", ts[5].typ.val)
	require.Equal(t, "abc", ts[6].typ.val)
}