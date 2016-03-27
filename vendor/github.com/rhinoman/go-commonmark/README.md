go-commonmark
=======


[![Build Status](https://travis-ci.org/rhinoman/go-commonmark.svg?branch=master)](https://travis-ci.org/rhinoman/go-commonmark)

Description
-----------

go-commonmark is a [Go](http://golang.org) (golang) wrapper for the [CommonMark](http://commonmark.org/) C library


Installation
------------

```
go get github.com/rhinoman/go-commonmark
```

**Note:** The [cmark](https://github.com/jgm/cmark) C reference implementation has been folded into this repository, no need to install it separately.  It will be built automagically by cgo.

Documentation
-------------

See the Godoc: http://godoc.org/github.com/rhinoman/go-commonmark


Example Usage
-------------
If all you need is to convert CommonMark text to Html, just do this:

```go

import "github.com/rhinoman/go-commonmark"

...

	htmlText := commonmark.Md2Html(mdText)  

```
