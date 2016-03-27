go-slugification
================

[![Build Status](https://travis-ci.org/rhinoman/go-slugification.svg?branch=master)](https://travis-ci.org/rhinoman/go-slugification)

Description
-----------
Creates Slugified versions of strings suitable for use in URLs

Installation
------------

```
go get github.com/rhinoman/go-slugification
```

Usage
-----

```go
import "github.com/rhinoman/go-slugification"
...

slugification.Slugify("Page Title") //Returns "page-title"

```


