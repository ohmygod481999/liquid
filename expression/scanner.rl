package expression

import "strconv"

%%{
	machine expression;
	write data;
	access lex.;
	variable p lex.p;
	variable pe lex.pe;
}%%

type lexer struct {
    data []byte
    p, pe, cs int
    ts, te, act int
		val func(Context) interface{}
}

func (l* lexer) token() string {
	return string(l.data[l.ts:l.te])
}

func newLexer(data []byte) *lexer {
	lex := &lexer{
			data: data,
			pe: len(data),
	}
	%% write init;
	return lex
}

func (lex *lexer) Lex(out *yySymType) int {
	eof := lex.pe
	tok := 0

	%%{
		action Bool {
			tok = LITERAL
			out.val = lex.token() == "true"
			fbreak;
		}
		action Identifier {
			tok = IDENTIFIER
			out.name = lex.token()
			fbreak;
		}
		action Int {
			tok = LITERAL
			n, err := strconv.ParseInt(lex.token(), 10, 64)
			if err != nil {
				panic(err)
			}
			out.val = int(n)
			fbreak;
		}
		action Float {
			tok = LITERAL
			n, err := strconv.ParseFloat(lex.token(), 64)
			if err != nil {
				panic(err)
			}
			out.val = n
			fbreak;
		}
		action String {
			tok = LITERAL
			// TODO unescape \x
			out.val = string(lex.data[lex.ts+1:lex.te-1])
			fbreak;
		}
		action Relation { tok = RELATION; out.name = lex.token(); fbreak; }

		identifier = (alpha | '_') + ;
		# TODO what can a property name contain?
		property = '.' (alnum | '_' | '-')+ ;
		int = '-'? digit+ ;
		float = '-'? (digit+ '.' digit* | '.' digit+) ;
		string = '"' (any - '"')* '"' | "'" (any - "'")* "'" ; # TODO escapes

		main := |*
			"%assign" => { tok = ASSIGN; fbreak; };
			"%loop" => { tok = LOOP; fbreak; };
			int => Int;
			float => Float;
			string => String;
			("true" | "false") => Bool;
			"nil" => { tok = LITERAL; out.val = nil; fbreak; };
			"==" => { tok = EQ; fbreak; };
			"!=" => { tok = NEQ; fbreak; };
			">=" => { tok = GE; fbreak; };
			"<=" => { tok = LE; fbreak; };
			"and" => { tok = AND; fbreak; };
			"or" => { tok = OR; fbreak; };
			"contains" => { tok = CONTAINS; fbreak; };
			"for" => { tok = FOR; fbreak; };
			"in" => { tok = IN; fbreak; };
			identifier ':' => { tok = KEYWORD; out.name = string(lex.data[lex.ts:lex.te-1]); fbreak; };
			identifier => Identifier;
			property => { tok = PROPERTY; out.name = string(lex.data[lex.ts+1:lex.te]); fbreak; };
			space+;
			any => { tok = int(lex.data[lex.ts]); fbreak; };
		*|;

		write exec;
	}%%

	return tok
}

func (lex *lexer) Error(e string) {
    // fmt.Println("scan error:", e)
}