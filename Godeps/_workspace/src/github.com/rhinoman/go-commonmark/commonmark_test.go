package commonmark_test

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/go-commonmark"
	"testing"
	"time"
)

func TestMd2Html(t *testing.T) {
	htmlText := commonmark.Md2Html("Boo\n===", 0)
	if htmlText != "<h1>Boo</h1>\n" {
		t.Errorf("Html text is not as expected :(")
	}
	t.Logf("Html Text: %v", htmlText)
}

func TestCMarkVersion(t *testing.T) {
	version := commonmark.CMarkVersion()
	t.Logf("\nVersion: %v", version)
}

func TestCMarkParser(t *testing.T) {
	parser := commonmark.NewCmarkParser(commonmark.CMARK_OPT_DEFAULT)
	if parser == nil {
		t.Error("Parser is nil!")
	}
	parser.Feed("Boo\n")
	parser.Feed("===\n")
	document := parser.Finish()
	if document == nil {
		t.Error("Document is nil!")
	}
	//Call it twice to make sure it doesn't crash :)
	parser.Free()
	parser.Free()
	htmlText := document.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	if htmlText != "<h1>Boo</h1>\n" {
		t.Error("Html text is not as expected :(")
	}
	t.Logf("Html Text: %v", htmlText)
	document.RenderXML(commonmark.CMARK_OPT_DEFAULT)
	document.Free()

	document2 := commonmark.ParseDocument("Foobar\n------", 0)
	htmlText = document2.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	document2.RenderXML(commonmark.CMARK_OPT_DEFAULT)
	if htmlText != "<h2>Foobar</h2>\n" {
		t.Error("Html text 2 is not as expected :(")
	}
	t.Logf("Html Text2: %v", htmlText)
	document2.Free()
	document2.Free()
}

func TestParseFile(t *testing.T) {
	node, err := commonmark.ParseFile("test_data/test_file.md", 0)
	if err != nil {
		t.Error(err)
	}
	if node == nil {
		t.Error(err)
	}
	htmlText := node.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	if htmlText != "<h1>Test File</h1>\n<h2>Description</h2>\n<p>This is just a test file.</p>\n" {
		t.Error("Html text is not as expected :(")
	}
	t.Logf("Html Text: %v", htmlText)
	node.Free()
	//try to parse a non-existent file
	eNode, err := commonmark.ParseFile("notafile.md", 0)
	if err == nil {
		t.Errorf("Should have been an error!")
	}
	t.Logf("error string: %v", err.Error())
	if eNode != nil {
		t.Errorf("Node should be nil!")
	}
}

func TestCMarkNodeOps(t *testing.T) {
	root := commonmark.NewCMarkNode(commonmark.CMARK_NODE_DOCUMENT)
	if root == nil {
		t.Error("Root is nil!")
	}
	if root.GetNodeType() != commonmark.CMARK_NODE_DOCUMENT {
		t.Error("Root is wrong type!")
	}
	if root.GetNodeTypeString() != "document" {
		t.Error("Root is wrong type string!")
	}
	header1 := commonmark.NewCMarkNode(commonmark.CMARK_NODE_HEADER)
	if header1.GetNodeType() != commonmark.CMARK_NODE_HEADER {
		t.Error("header1 is wrong type!")
	}
	header1.SetHeaderLevel(1)
	if header1.SetLiteral("boo") != false {
		t.Error("SetLiteral should return false for header node")
	}
	header1str := commonmark.NewCMarkNode(commonmark.CMARK_NODE_TEXT)
	header1str.SetLiteral("I'm the main header!")
	if header1str.GetLiteral() != "I'm the main header!" {
		t.Error("header1str content is wrong!")
	}
	header1.AppendChild(header1str)
	header2 := commonmark.NewCMarkNode(commonmark.CMARK_NODE_HEADER)
	header2str := commonmark.NewCMarkNode(commonmark.CMARK_NODE_TEXT)
	if header2str.SetLiteral("Another header!") == false {
		t.Error("SetLiteral returned false for valid input")
	}
	header2.AppendChild(header2str)
	header2.SetHeaderLevel(2)
	if root.PrependChild(header1) == false {
		t.Error("Couldn't prepend header to root")
	}
	root.AppendChild(header2)
	t.Logf("\nXML: %v", root.RenderXML(commonmark.CMARK_OPT_DEFAULT))

	htmlStr := root.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	if htmlStr != "<h1>I'm the main header!</h1>\n<h2>Another header!</h2>\n" {
		t.Error("htmlStr is wrong!")
	}
	t.Logf("Html Text: %v", htmlStr)
	//Rearrange...
	header1.InsertBefore(header2)
	t.Logf("\nXML: %v", root.RenderXML(commonmark.CMARK_OPT_DEFAULT))
	htmlStr = root.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	if htmlStr != "<h2>Another header!</h2>\n<h1>I'm the main header!</h1>\n" {
		t.Error("htmlStr is wrong!")
	}
	t.Logf("Html Text: %v", htmlStr)
	//removing something
	header2.Unlink()
	t.Logf("\nXML: %v", root.RenderXML(commonmark.CMARK_OPT_DEFAULT))
	htmlStr = root.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	if htmlStr != "<h1>I'm the main header!</h1>\n" {
		t.Error("htmlStr is wrong!")
	}
	latexStr := root.RenderLatex(commonmark.CMARK_OPT_DEFAULT, 80)
	t.Logf("\nLatex: %v", latexStr)
	manStr := root.RenderMan(commonmark.CMARK_OPT_DEFAULT, 80)
	t.Logf("\nMAN: %v", manStr)
	cmStr := root.RenderCMark(commonmark.CMARK_OPT_DEFAULT, 0)
	t.Logf("\nCMARK: %v", cmStr)
	root.ConsolidateTextNodes()
	t.Logf("\nXML: %v", root.RenderXML(commonmark.CMARK_OPT_DEFAULT))
	root.SetNodeUserData("STRING!")
	x := root.GetNodeUserData()
	t.Logf("X: %v", x)
	//header2.Free()
	root.Free()
}

func TestCMarkLists(t *testing.T) {
	root := commonmark.NewCMarkNode(commonmark.CMARK_NODE_DOCUMENT)
	list := commonmark.NewCMarkNode(commonmark.CMARK_NODE_LIST)
	list.SetListType(commonmark.CMARK_ORDERED_LIST)
	listItem1 := commonmark.NewCMarkNode(commonmark.CMARK_NODE_ITEM)
	listItem2 := commonmark.NewCMarkNode(commonmark.CMARK_NODE_ITEM)
	li1para := commonmark.NewCMarkNode(commonmark.CMARK_NODE_PARAGRAPH)
	li1str := commonmark.NewCMarkNode(commonmark.CMARK_NODE_TEXT)
	li1str.SetLiteral("List Item 1")
	li1para.AppendChild(li1str)
	if listItem1.AppendChild(li1para) == false {
		t.Error("Couldn't append paragraph to list item")
	}
	list.AppendChild(listItem1)
	list.AppendChild(listItem2)
	list.SetListTight(true)
	root.AppendChild(list)
	t.Logf("\nXML: %v", root.RenderXML(commonmark.CMARK_OPT_DEFAULT))
	htmlString := root.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	if htmlString != "<ol>\n<li>List Item 1</li>\n<li></li>\n</ol>\n" {
		t.Error("htmlString is wrong!")
	}
	t.Logf("\nHtmlString: \n%v", htmlString)
	t.Logf("\nList start: %v", list.GetListStart())
	t.Logf("\nList tight: %v", list.GetListTight())
	root.Free()
}

func TestCMarkCodeBlocks(t *testing.T) {
	root := commonmark.NewCMarkNode(commonmark.CMARK_NODE_DOCUMENT)
	cb := commonmark.NewCMarkNode(commonmark.CMARK_NODE_CODE_BLOCK)
	cb.SetLiteral("int main(){\n return 0;\n }")
	cb.SetFenceInfo("c")
	if cb.GetFenceInfo() != "c" {
		t.Error("Fence info isn't c")
	}
	if cb.GetLiteral() != "int main(){\n return 0;\n }" {
		t.Error("Code has changed somehow")
	}
	if root.AppendChild(cb) == false {
		t.Error("Couldn't append code block to document")
	}
	t.Logf("\nXML: %v", root.RenderXML(commonmark.CMARK_OPT_DEFAULT))
	htmlString := root.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	t.Logf("\nHtml String: %v\n", htmlString)
	if htmlString != "<pre><code>int main(){\n return 0;\n }</code></pre>\n" {
		t.Error("htmlString isn't right!")
	}
	root.Free()
}

func TestCMarkUrls(t *testing.T) {
	root := commonmark.NewCMarkNode(commonmark.CMARK_NODE_DOCUMENT)
	para := commonmark.NewCMarkNode(commonmark.CMARK_NODE_PARAGRAPH)
	link := commonmark.NewCMarkNode(commonmark.CMARK_NODE_LINK)
	root.AppendChild(para)
	if para.AppendChild(link) == false {
		t.Error("Couldn't append link node to paragraph!")
	}
	if link.SetUrl("http://duckduckgo.com") == false {
		t.Error("Couldn't set URL!!!")
	}
	if link.GetUrl() != "http://duckduckgo.com" {
		t.Error("Url doesn't match")
	}
	t.Logf("\nXML: %v", root.RenderXML(commonmark.CMARK_OPT_DEFAULT))
	htmlString := root.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	t.Logf("\nHtml String: %v\n", htmlString)
	if htmlString != "<p><a href=\"http://duckduckgo.com\"></a></p>\n" {
		t.Error("htmlString isn't right!")
	}
	root.Free()
}

func TestCMarkIter(t *testing.T) {
	root := commonmark.NewCMarkNode(commonmark.CMARK_NODE_DOCUMENT)
	list := commonmark.NewCMarkNode(commonmark.CMARK_NODE_LIST)
	list.SetListType(commonmark.CMARK_ORDERED_LIST)
	listItem1 := commonmark.NewCMarkNode(commonmark.CMARK_NODE_ITEM)
	listItem2 := commonmark.NewCMarkNode(commonmark.CMARK_NODE_ITEM)
	li1para := commonmark.NewCMarkNode(commonmark.CMARK_NODE_PARAGRAPH)
	li1str := commonmark.NewCMarkNode(commonmark.CMARK_NODE_TEXT)
	li1str.SetLiteral("List Item 1")
	li1para.AppendChild(li1str)
	if listItem1.AppendChild(li1para) == false {
		t.Error("Couldn't append paragraph to list item")
	}
	list.AppendChild(listItem1)
	list.AppendChild(listItem2)
	list.SetListTight(true)
	root.AppendChild(list)
	t.Logf("\nXML: %v", root.RenderXML(commonmark.CMARK_OPT_DEFAULT))
	iter := commonmark.NewCMarkIter(root)
	for {
		ne := iter.Next()
		t.Logf("NodeEvent: %v", ne)
		iNode := iter.GetNode()
		if iNode == nil {
			t.Error("iter node was nil!")
		}
		if ne == commonmark.CMARK_EVENT_DONE {
			break
		}

	}
	iter.Reset(listItem2, commonmark.CMARK_EVENT_DONE)
	iter.Free()
	root.Free()
}

func createTree() *commonmark.CMarkNode {
	root := commonmark.NewCMarkNode(commonmark.CMARK_NODE_DOCUMENT)
	header1 := commonmark.NewCMarkNode(commonmark.CMARK_NODE_HEADER)
	header2 := commonmark.NewCMarkNode(commonmark.CMARK_NODE_HEADER)
	header1str := commonmark.NewCMarkNode(commonmark.CMARK_NODE_TEXT)
	header2str := commonmark.NewCMarkNode(commonmark.CMARK_NODE_TEXT)
	header1str.SetLiteral("Header 1!")
	header2str.SetLiteral("Header 2!")
	root.AppendChild(header1)
	root.AppendChild(header2)
	header1.AppendChild(header1str)
	header2.AppendChild(header2str)
	return root

}

//Checking mem management functions
func TestMem(t *testing.T) {
	tree := createTree()
	time.Sleep(3 * time.Second)
	t.Logf("\nXML: %v", tree.RenderXML(commonmark.CMARK_OPT_DEFAULT))
	iter := commonmark.NewCMarkIter(tree)
	i := 1
	for {
		ne := iter.Next()
		t.Logf("NodeEvent: %v", ne)
		if ne == commonmark.CMARK_EVENT_DONE {
			break
		}
		i += 1
	}
	if i < 9 {
		t.Errorf("Lost some nodes somewhere: %v", i)
	}
	tree.Free()
}
