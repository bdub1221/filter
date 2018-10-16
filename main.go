package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gomarkdown/markdown/parser"
)

var (
	flagVersion = flag.Bool("v", false, "show filter version")
	flagList    = flag.Bool("l", false, "list all available plugins")
	flagPlugins = flag.String("p", "noop, emph", "comma separated list of plugins to load")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "SYNOPSIS: %s [OPTIONS] %s\n", os.Args[0], "[FILE...]")
		fmt.Println("\nOPTIONS:")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"os.Stdin"}
	}
	if *flagVersion {
		fmt.Println(Version)
		return
	}
	if *flagList {
		for name, _ := range Plugins {
			fmt.Println(name)
		}
		return
	}

	if *flagPlugins == "" {
		return
	}

	r := &Renderer{}
	requested := strings.Split(*flagPlugins, ",")
	for _, plugin := range requested {
		impl, ok := Plugins[plugin]
		if !ok {
			log.Fatalf("Plugin %q not found", plugin)
		}

		r.RegisterPlugin(impl)
	}

	for _, fileName := range args {
		var (
			d   []byte
			err error
		)
		if fileName == "os.Stdin" {
			d, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Printf("Couldn't read %q: %q", fileName, err)
				continue
			}
		} else {
			d, err = ioutil.ReadFile(fileName)
			if err != nil {
				log.Printf("Couldn't open %q: %q", fileName, err)
				continue
			}
		}

		x := r.Render(d)
		fmt.Print(string(x))
		continue
	}
}

// Extensions is exported to we can use it in tests. (copied from mmark)
var Extensions = parser.Tables | parser.FencedCode | parser.Autolink | parser.Strikethrough |
	parser.SpaceHeadings | parser.HeadingIDs | parser.BackslashLineBreak | parser.SuperSubscript |
	parser.DefinitionLists | parser.MathJax | parser.AutoHeadingIDs | parser.Footnotes |
	parser.Strikethrough | parser.OrderedListStart | parser.Attributes | parser.Mmark
