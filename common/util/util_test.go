/**  Copyright (c) 2014-present James Adam.  All rights reserved.
*
*		 This file is part of Wikifeat.
*
*    Wikifeat is free software: you can redistribute it and/or modify
*    it under the terms of the GNU General Public License as published by
*    the Free Software Foundation, either version 2 of the License, or
*    (at your option) any later version.
*
*    This program is distributed in the hope that it will be useful,
*    but WITHOUT ANY WARRANTY; without even the implied warranty of
*    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*    GNU General Public License for more details.
*
*    You should have received a copy of the GNU General Public License
*    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package util

import (
	"testing"
)

func TestParseArrayParam(t *testing.T) {
	strQuery := "[bill,ted,steven]"
	queryArray, err := ParseArrayParam(strQuery)
	if err != nil {
		t.Error(err)
	}
	for i, str := range queryArray {
		t.Logf("%d: %s", i, str)
	}
	if len(queryArray) != 3 {
		t.Errorf("Should be 3 records, not %v!", len(queryArray))
	}
}

func TestBadArrayParam(t *testing.T) {
	strQuery := "scotch, wine, gin]"
	_, err := ParseArrayParam(strQuery)
	if err == nil {
		t.Fail()
	}
	t.Logf("Error: %v", err)
}

func TestGenToken(t *testing.T) {
	tok := GenToken()
	if tok == "" {
		t.Fail()
	}
	t.Logf("Token: %v", tok)
}
