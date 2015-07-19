package slugification_test

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/go-slugification"
	"testing"
)

func TestSlugification(t *testing.T) {
	//Simple string
	testString1 := "Page Title"
	slugged := slugification.Slugify(testString1)
	if slugged != "page-title" {
		t.Errorf("slugged string is wrong: %v", slugged)
	}
	//Not so simple string
	testString2 := "golang-convert-iso8859-1-to-utf8"
	slugged = slugification.Slugify(testString2)
	if slugged != "golang-convert-iso8859-1-to-utf8" {
		t.Errorf("slugged string is wrong: %v", slugged)
	}
	//Url query string
	testString3 := "view?id=22145749&trk=nav_responsive_tab_profile"
	slugged = slugification.Slugify(testString3)
	if slugged != "viewid22145749trknav_responsive_tab_profile" {
		t.Errorf("slugged string is wrong: %v", slugged)
	}
	//pathological string
	testString4 := "1 < 3 > 5 #2234?x=5&4%20{poo}|p\\or2^6~42[11]`d;/f:f@me$3.50"
	slugged = slugification.Slugify(testString4)
	if slugged != "1--3--5-2234x5420poopor264211dffme350" {
		t.Errorf("slugged string is wrong: %v", slugged)
	}
}
