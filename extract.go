package extract

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	xmlpath "gopkg.in/xmlpath.v2"
)

// IExtractor is interface of extractor.
type IExtractor interface {
	Extract(r io.Reader, rawurl string)
}

// Rule is
type Rule struct {
	ResourceURL         string    `json:"resource_url"`
	Name                string    `json:"name"`
	CreatedBy           string    `json:"created_by"`
	DatabaseResourceURL string    `json:"database_resource_url"`
	UpdatedAt           time.Time `json:"updated_at"`
	CreatedAt           time.Time `json:"created_at"`
	Data                struct {
		URL          string `json:"url"` // URL is the regex pattern of target pages
		Type         string `json:"type"`
		Enc          string `json:"enc"` // Enc is encoding of the page contents
		XPath        string `json:"xpath"`
		Base         string `json:"base"`
		MicroFormats string `json:"microformats"`
	}
	regexp *regexp.Regexp // regex is compiled regular expressions from Data field.
}

// Content is extracted contents.
type Content struct {
	Title       string
	Description string
	Text        string // Text is extracted body.
}

// XPath get XPath from rule.
func (r *Rule) XPath() string {
	return r.Data.XPath
}

// URL get raw URL from rule.
// URL is raw regular expression string.
func (r *Rule) URL() string {
	return r.Data.URL
}

// Extractor is actual extractor
type Extractor struct {
	rules []Rule
}

// Init initialize extractor.
func (e *Extractor) Init() {
	if err := e.LoadFile("./items.json"); err != nil {
		// it cannot recover. panic it.
		log.Fatalf("fail to initialize extractor: error %s", err)
	}
}

// New creates extrator
func New() *Extractor {
	e := &Extractor{}
	e.Init()
	return e
}

// Dump returns rules
func (e *Extractor) Dump() []Rule {
	return e.rules
}

// ExtractURL fetch contents from URL, and parse it, and extract.
func (e *Extractor) ExtractURL(rawurl string) (*Content, error) {
	resp, err := http.Get(rawurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return e.Extract(r, rawurl)
}

// Match returns if provided URL is match wedata rules.
//
// TODO: use prefix match for faster matching.
func (e *Extractor) Match(rawurl string) *Rule {
	for _, rl := range e.rules {
		if ok := rl.regexp.MatchString(rawurl); ok {
			return &rl
		}
	}
	return nil
}

// Extract load body from io.Reader and extract contents by rule.
// If URL is not match wedata rule, skip it and return nil.
// io.Reader should be UTF-8 byte stream.
func (e *Extractor) Extract(r io.Reader, rawurl string) (*Content, error) {
	rule := e.Match(rawurl)
	// if not found, skip it.
	if rule == nil {
		return nil, nil
	}

	// match by xpath
	// TODO should compile at init phase
	path, err := xmlpath.Compile(rule.XPath())
	if err != nil {
		return nil, err
	}
	node, err := xmlpath.ParseHTML(r)
	if err != nil {
		return nil, err
	}
	iter := path.Iter(node)

	body := ""
	for iter.Next() {
		body = body + strings.TrimSpace(iter.Node().String())
	}
	c := &Content{
		Text: body,
	}
	return c, nil
}

// Load loads rule from reader and extract rules.
func (e *Extractor) Load(r io.Reader) error {
	var rawRules []Rule
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return err
	}
	if err := json.Unmarshal(buf.Bytes(), &rawRules); err != nil {
		return err
	}
	rules := make([]Rule, 0, len(rawRules))
	// parse regexp. if failed, ignore rule.
	for _, r := range rawRules {
		regexp, err := regexp.Compile(r.URL())
		if err != nil {
			log.Printf("ignore %s because of illegal regexp.", r.Name)
			continue
		}
		r.regexp = regexp
		rules = append(rules, r)
	}
	e.rules = rules
	return nil
}

// LoadFile load LDRFullFeed JSON from filepath
func (e *Extractor) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return e.Load(f)
}
