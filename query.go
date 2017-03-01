package MusicLibraryIndex

import (
	"strings"
	"errors"
)

const (
	WHITESPACE        = iota
	OPEN_PARENTHESIS
	CLOSE_PARENTHESIS
	NOT
	OR
	LABEL
	QUOTED_THING
	QUOTED_THING_OPEN
	QUOTED_THING_CLOSE
	NORMAL_WORD
)

type token struct {
	typeVal int
	text string
}

type tokens struct {
	index int
	tokens []*token
}

func (ts *tokens) add(t *token)  {
	ts.tokens = append(ts.tokens, t)
}

func (ts *tokens) len() int {
	return len(ts.tokens)
}

func (ts *tokens) get(i int) *token {
	if i >=0 && i < ts.len() {
		return ts.tokens[i]
	}
	return nil
}

func newTokens(size int) *tokens {
	t := &tokens{
		index: 0,
		tokens: make([]*token, 0, size),
	}
	return t
}


///// MakeAndMatcher /////
type MakeAndMatcher interface {
	toString() string
	matcher(track *Track) bool
}

type MakeAndMatchers []MakeAndMatcher

func (ms MakeAndMatchers) toString() string {
	switch len(ms) {
	case 0:
		return "()"
	case 1:
		return ms[0].toString()
	}

	arrStr := make([]string, len(ms))
	for i, m :=range ms {
		arrStr[i] = m.toString()
	}
	return "(" + strings.Join(arrStr," AND ") + ")"
}

func (ms MakeAndMatchers) matcher(track *Track) bool {
	rez := true
	for _, val := range ms {
		rez = val.matcher(track)
		if !rez {
			break
		}
	}
	return rez
}

type FuzzyTextMatche struct {
	fuzzyTerm string
}

type ExactTextMatcher struct {
	exactTerm string
}

type NotMatcher struct {
	subMatcher MakeAndMatcher
}

type OrMatcher struct {
	subMatchers MakeAndMatchers
}

type LabelMatcher struct {
	id string
}

func makeFuzzyTextMatcher(str string) *FuzzyTextMatche {
	return &FuzzyTextMatche{
		fuzzyTerm: str,
	}
}

func makeExactTextMatcher(str string) *ExactTextMatcher {
	return &ExactTextMatcher{
		exactTerm: str,
	}
}

func (m *FuzzyTextMatche) toString() string {
	return "(fuzzy \"" + m.fuzzyTerm + "\")"
}

func (m *ExactTextMatcher) toString() string {
	return "(exact \"" + m.exactTerm + "\")"
}

func (m *FuzzyTextMatche) matcher(track *Track) bool {
	if m.fuzzyTerm == "" {
		return true
	}
	srtFind := formatSearchable(m.fuzzyTerm)
	index := strings.Index(track.fuzzySearchTags, srtFind)
	return index >= 0
}

func (m *ExactTextMatcher) matcher(track *Track) bool {
	if m.exactTerm == "" {
		return true
	}
	srtFind := m.exactTerm
	index := strings.Index(track.exactSearchTags, srtFind)
	return index >= 0
}

///// Library /////
func (this *Library) parseQuery (query string) MakeAndMatchers {
	tokens := tokenizeQuery(query)
	return tokens.parseList(-1)
}

func (t *tokens) parseList(waitForTokenType int) MakeAndMatchers {
	matchers := make(MakeAndMatchers, 0)
	justSawWhitespace := true
	for t.index < len(t.tokens) {
		token := t.get(t.index)
		t.index++
		switch token.typeVal {
		case OPEN_PARENTHESIS:
			subMatcher := t.parseList(CLOSE_PARENTHESIS)
			matchers = append(matchers, subMatcher...)
		case CLOSE_PARENTHESIS:
			if (waitForTokenType == CLOSE_PARENTHESIS) {
				return matchers
			}
			var previousMatcher= matchers[len(matchers)-1].(*FuzzyTextMatche)
			if (!justSawWhitespace && previousMatcher.fuzzyTerm != "") {
				previousMatcher.fuzzyTerm += token.text;
			} else {
				matchers = append(matchers, makeFuzzyTextMatcher(token.text))
			}
		case NOT:
			matchers = append(matchers, t.parseNot())
		case OR:
			subMatchers := t.parseList(-1)
			matchers = append(matchers, makeOrMatcher(subMatchers))
		case LABEL:
			matchers = append(matchers, t.parseLabel())
		case QUOTED_THING:
			if len(token.text) != 0 {
				matchers = append(matchers, makeExactTextMatcher(token.text))
			}
		case NORMAL_WORD:
			matchers = append(matchers, makeFuzzyTextMatcher(token.text))
		}
		justSawWhitespace = (token.typeVal == WHITESPACE)
	}
	return matchers
}

func (t *tokens) parseLabel() MakeAndMatcher {
	if t.index >= t.len() {
		return makeFuzzyTextMatcher(t.get(t.index - 1).text)
	}
	token := t.get(t.index)
	t.index++

	switch (token.typeVal) {
	case WHITESPACE:
	case CLOSE_PARENTHESIS:
		t.index--
		return makeFuzzyTextMatcher(t.get(t.index - 1).text)
	case OPEN_PARENTHESIS,
		NOT,
		OR,
		LABEL,
		QUOTED_THING,
		NORMAL_WORD:
		return makeLabelMatcher(token.text)
	}
	panic(errors.New("unreachable"))
}

func (t *tokens) parseNot() MakeAndMatcher {
	if t.index >= t.len() {
		return makeFuzzyTextMatcher(t.get(t.index-1).text)
	}
	token := t.get(t.index)
	t.index++

	switch (token.typeVal) {
	case WHITESPACE:
	case CLOSE_PARENTHESIS:
		t.index--
		return makeFuzzyTextMatcher(t.get(t.index-1).text);
	case OPEN_PARENTHESIS:
		return makeNotMatcher(t.parseList(CLOSE_PARENTHESIS))
	case NOT:
		return makeNotMatcher(t.parseNot())
	case OR:
		return makeNotMatcher(t.parseList(CLOSE_PARENTHESIS))
	case LABEL:
		return makeNotMatcher(t.parseLabel())
   	case QUOTED_THING:
		return makeNotMatcher(makeExactTextMatcher(token.text))
	case NORMAL_WORD:
		return makeNotMatcher(makeFuzzyTextMatcher(token.text))
	}
	panic(errors.New("unreachable"))
}

func (m *NotMatcher) toString() string {
	return "(not " + m.subMatcher.toString() + ")"
}

func makeNotMatcher(subMatcher MakeAndMatcher) *NotMatcher {
	return &NotMatcher{subMatcher: subMatcher}
}

func (m *NotMatcher) matcher(track *Track) bool {
	return !m.subMatcher.matcher(track)
}

func makeOrMatcher(subMatchers MakeAndMatchers) *OrMatcher {
	return &OrMatcher{subMatchers: subMatchers}
}

func (m *OrMatcher) toString() string {
	arr := make([]string, 0, len(m.subMatchers))
	for _, subMatcher := range m.subMatchers {
		arr = append(arr, subMatcher.toString())
	}
	return strings.Join(arr, " OR ")
}

func (m *OrMatcher) matcher(track *Track) bool {
	val := false
	for _, subMatcher := range m.subMatchers {
		val = subMatcher.matcher(track)
		if val {
			break
		}
	}
	return val
}

func (m *LabelMatcher) toString() string {
	return"(label " + m.id + ")"
}

func makeLabelMatcher(text string) *LabelMatcher {
	return &LabelMatcher{id: text}
}

func (m *LabelMatcher) matcher(track *Track) bool {
	_, ok := track.Labels[m.id]
	return ok
}

func tokenizeQuery(query string) *tokens {
	arr := strings.Split(query, " ")
	tokens := newTokens(len(arr))

	for i, str := range arr{
		if len(str)>1 {
			startStr :=0
			endStr := 0
			typeVal := NORMAL_WORD

			bFind := true
			beginArray := make([]*token, 0)
			endArray := make([]*token, 0)
			for bFind {
				if strings.HasPrefix(str,`"`) {
					typeVal = QUOTED_THING_OPEN
					startStr++
					break
				}
				if strings.HasSuffix(str,`"`) && !strings.HasSuffix(str,`\"`) {
					typeVal = QUOTED_THING_CLOSE
					endStr = len(str)-1
					break
				}

				bFind = false
				if strings.HasPrefix(str, `(`) {
					t := &token{
						typeVal: OPEN_PARENTHESIS,
						text:    `(`,
					}
					beginArray = append(beginArray, t)
					str = str[1:]
					bFind = true
					continue
				}else if strings.HasPrefix(str, `or:`) {
					t := &token{
						typeVal: OR,
						text:    `or:`,
					}
					beginArray = append(beginArray, t)
					str = str[3:]
					bFind = true
					continue
				}else if strings.HasPrefix(str, `not:`) {
					t := &token{
						typeVal: NOT,
						text:    `not:`,
					}
					beginArray = append(beginArray, t)
					str = str[4:]
					bFind = true
					continue
				}else if strings.HasPrefix(str, `label:`) {
					t := &token{
						typeVal: LABEL,
						text:    `label:`,
					}
					beginArray = append(beginArray, t)
					str = str[6:]
					bFind = true
					continue
				}else if strings.HasSuffix(str, `)`) {
					t := &token{
						typeVal: CLOSE_PARENTHESIS,
						text:    `)`,
					}
					endArray = append(endArray, t)
					str = str[:len(str)-1]
					bFind = true
					continue
				}

			}

			if typeVal != QUOTED_THING_OPEN && strings.HasPrefix(str,`"`) {
				typeVal = QUOTED_THING_OPEN
				startStr++
			}

			if typeVal != QUOTED_THING_CLOSE && strings.HasSuffix(str,`"`) &&
												!strings.HasSuffix(str,`\"`){
				endStr = len(str)-1
				if typeVal == QUOTED_THING_OPEN {
					typeVal = QUOTED_THING
				}else{
					typeVal = QUOTED_THING_CLOSE
				}
			}

			if startStr > 0 && endStr > 0 {
				str = str[startStr:endStr]
			}else if startStr > 0  {
				str = str[startStr:]
			}else if endStr > 0 {
				str = str[:endStr]
			}

			for _, t := range beginArray {
				tokens.add(t)
			}
			t := &token{
				typeVal: typeVal,
				text:    str,
			}
			tokens.add(t)
			for j := len(endArray); j > 0; j-- {
				tokens.add(endArray[j-1])
			}
		}else if len(str)==1 {
			t := &token{
				typeVal: NORMAL_WORD,
				text:    str,
			}
			tokens.add(t)
		}

		if i < len(arr)-1 {
			t := &token{
				typeVal: WHITESPACE,
				text:    " ",
			}
			tokens.add(t)
		}

	}

	tokensRez := newTokens(len(arr))
	for i := 0; i<tokens.len(); i++ {
		token_i := tokens.get(i)
		if token_i.typeVal == QUOTED_THING_OPEN {
			str := make([]string, 0)
			j := i
			for ; j<tokens.len(); j++ {
				token_j := tokens.get(j)
				str = append(str, token_j.text)
				if token_j.typeVal == QUOTED_THING_CLOSE {
					break
				}
			}
			token_i.typeVal = QUOTED_THING
			token_i.text = strings.Join(str, "")
			tokensRez.add(token_i)
			i = j
		}else if token_i.typeVal == WHITESPACE {
			continue
		}else {
			tokensRez.add(token_i)
		}
	}
	return tokensRez
}
