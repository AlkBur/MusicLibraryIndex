package MusicLibraryIndex

import "testing"

func TestParseQuery(t *testing.T) {
	library := New()

	//strictEqual(library.parseQuery("").toString(), "()")
	//strictEqual(library.parseQuery("a").toString(), `(fuzzy "a")`)
	//strictEqual(library.parseQuery(" ab   cd ").toString(), `((fuzzy "ab") AND (fuzzy "cd"))`)
	//strictEqual(library.parseQuery(`"a  b"`).toString(), `(exact "a  b")`)
	strictEqual(library.parseQuery(`"a  b\\" c"`).toString(), `(exact "a  b\\" c")`)
	//strictEqual(library.parseQuery(`\\"a  b"`).toString(), `((fuzzy "\\\"a") AND (fuzzy "b\""))`)
	//strictEqual(library.parseQuery(`"`).toString(), `(fuzzy "\"")`)
	//strictEqual(library.parseQuery(`\\`).toString(), `(fuzzy "\\\\")`)
	//strictEqual(library.parseQuery(`""`).toString(), `()`)
	//strictEqual(library.parseQuery(`a" b"c`).toString(), `((fuzzy "a\\\"") AND (fuzzy "b\\\"c"))`)
	//
	//strictEqual(library.parseQuery('not:A').toString(), '(not (fuzzy "a"))');
	//strictEqual(library.parseQuery('not:"A"').toString(), '(not (exact "A"))');
	//strictEqual(library.parseQuery('not:(a b)').toString(), '(not ((fuzzy "a") AND (fuzzy "b")))');
	//strictEqual(library.parseQuery('not:(a)').toString(), '(not (fuzzy "a"))');
	//strictEqual(library.parseQuery('not:not:a').toString(), '(not (not (fuzzy "a")))');
	//strictEqual(library.parseQuery('not:').toString(), '(fuzzy "not:")');
	//strictEqual(library.parseQuery('not: a').toString(), '((fuzzy "not:") AND (fuzzy "a"))');
	//strictEqual(library.parseQuery('not:)a').toString(), '((fuzzy "not:)") AND (fuzzy "a"))');
	//
	//strictEqual(library.parseQuery('or:()').toString(), '()');
	//strictEqual(library.parseQuery('or:(' ).toString(), '()');
	//strictEqual(library.parseQuery('or:').toString(), '(fuzzy "or:")');
	//strictEqual(library.parseQuery('or:(a)').toString(), '(fuzzy "a")');
	//strictEqual(library.parseQuery('or:(a' ).toString(), '(fuzzy "a")');
	//strictEqual(library.parseQuery('or:(a b)').toString(), '((fuzzy "a") OR (fuzzy "b"))');
	//strictEqual(library.parseQuery('or:(a b' ).toString(), '((fuzzy "a") OR (fuzzy "b"))');
	//strictEqual(library.parseQuery('or:(a (b c))').toString(), '((fuzzy "a") OR ((fuzzy "b") AND (fuzzy "c")))');
	//strictEqual(library.parseQuery('or:((a b) c)').toString(), '(((fuzzy "a") AND (fuzzy "b")) OR (fuzzy "c"))');
}
