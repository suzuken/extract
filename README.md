# Extract

Extract is HTML Extractor. This extractor is based on [wedata](http://wedata.net/).

## Acknowledgement

[items.json]() is originally from `http://wedata.net/databases/LDRFullFeed/items.json`.

## How to use

From [_example](_example),

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/suzuken/extract"
)

func main() {
	var (
		rawurl = flag.String("url", "http://example.com", "url for extract")
	)
	flag.Parse()
	ex := extract.New()
	if rule := ex.Match(*rawurl); rule == nil {
		log.Printf("%s doesn't match in rule", *rawurl)
		os.Exit(0)
	}
	c, err := ex.ExtractURL(*rawurl)
	if err != nil {
		log.Fatalf("extract failed: %s", err)
	}
	fmt.Printf("content: %v", c)
}
```

## LICENSE

MIT

All data in wedata are in the public domain. see also: [http://wedata.net/help/about](http://wedata.net/help/about) .

## Special Thanks

* Wedata project and members.
