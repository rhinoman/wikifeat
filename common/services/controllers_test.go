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

package services_test

import (
	"github.com/rhinoman/wikifeat/common/services"
	"testing"
)

func TestGenIndexLinks(t *testing.T) {
	userRoles := []string{"books:admin", "videos:write", "music:read"}
	booksUri := "/api/books"
	bookLinks := services.GenIndexLinks(userRoles, "books", booksUri)
	if bookLinks.Self.Href != booksUri || bookLinks.Self.Method != "GET" {
		t.Errorf("Self link not set!")
	}
	if bookLinks.Create.Href != booksUri || bookLinks.Create.Method != "POST" {
		t.Errorf("Create link not set!")
	}
}
