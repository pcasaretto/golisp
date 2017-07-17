// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golisp

import (
	"fmt"
	"testing"
)

// Make the types prettyprint.
var itemName = map[itemType]string{
	itemError:        "error",
	itemBool:         "bool",
	itemChar:         "char",
	itemCharConstant: "charconst",
	itemComplex:      "complex",
	itemEOF:          "EOF",
	itemIdentifier:   "identifier",
	itemLeftParen:    "(",
	itemNumber:       "number",
	itemRawString:    "raw string",
	itemRightParen:   ")",
	itemSpace:        "space",
	itemString:       "string",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

type lexTest struct {
	name  string
	input string
	items []item
}

func mkItem(typ itemType, text string) item {
	return item{
		typ: typ,
		val: text,
	}
}

var (
	tEOF        = mkItem(itemEOF, "")
	tLpar       = mkItem(itemLeftParen, "(")
	tQuote      = mkItem(itemString, `"abc \n\t\" "`)
	tRpar       = mkItem(itemRightParen, ")")
	tSpace      = mkItem(itemSpace, " ")
	raw         = "`" + `abc\n\t\" ` + "`"
	rawNL       = "`now is{{\n}}the time`" // Contains newline inside raw quote.
	tRawQuote   = mkItem(itemRawString, raw)
	tRawQuoteNL = mkItem(itemRawString, rawNL)
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"left paren", "(", []item{tLpar, tEOF}},
	{"left right paren", "()", []item{tLpar, tRpar, tEOF}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest, left, right string) (items []item) {
	l := lex(t.name, t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test, "", "")
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
