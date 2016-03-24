package extract

import (
	"strings"
	"testing"
)

func TestLoadRule(t *testing.T) {
	rule := `[{
  "name": "Yahoo New India",
  "updated_at": "2010-05-02T06:36:04+09:00",
  "resource_url": "http:\/\/wedata.net\/items\/32275",
  "created_by": "sameerb",
  "database_resource_url": "http:\/\/wedata.net\/databases\/LDRFullFeed",
  "data": {
    "url": "^http:\/\/in\\.news\\.yahoo\\.com",
    "type": "IND",
    "xpath": "\/\/div[@id='storybody']"
  },
  "created_at": "2010-05-02T06:19:40+09:00"
}]`
	ex := &Extractor{}
	if err := ex.Load(strings.NewReader(rule)); err != nil {
		t.Fatalf("load rule failed %s", err)
	}
	rules := ex.Dump()
	if len(rules) != 1 {
		t.Fatal("rule is not loaded.")
	}
	rawurl := "http://in.news.yahoo.com/hoge"
	r := ex.Match("http://in.news.yahoo.com/hoge")
	if r == nil {
		t.Fatalf("%s should match rule but not.", rawurl)
	}
}
