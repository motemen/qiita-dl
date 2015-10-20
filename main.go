// qiita-dl is a simple tool that donwloads snippets published on Qiita <http://qiita.com>.
//
//   qiita-dl [-x] [-f] [-o <name>] [-d <directory>] <url>
//
// Example
//
// To download a snippet under ~/bin and make it executable:
//
//   $ qiita-dl -x -d ~/bin http://qiita.com/uasi/items/57da2e4268d348b371fb
//   Title: "git commit --fixup で fixup する対象を peco/fzf で選べるスクリプト書いた"
//   Saved to ~/bin/git-fixup
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	log.SetFlags(0)

	var (
		flagExecutable = flag.Bool("x", false, "mark downloaded snippet as executable")
		flagFilename   = flag.String("o", "", "output filename")
		flagDirname    = flag.String("d", "", "output directory")
		flagIndex      = flag.Uint("n", 0, "specify snippet index")
		flagForce      = flag.Bool("f", false, "overwrite existing file")
	)

	flag.Usage = func() {
		log.Printf("usage: %s [-x] [-o <name>] [-d <directory>] <url>", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		flag.Usage()
		os.Exit(1)
	}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	log.Println("Title:", doc.Find("title").Text())

	snippets := doc.Find("section[itemprop=articleBody] .code-frame")
	if snippets.Size() == 0 {
		log.Fatal("error: No snippets found")
	}

	var snippet *goquery.Selection
	if *flagIndex >= 1 {
		snippet = snippets.Eq(int(*flagIndex - 1))
	} else {
		snippet = snippets.Filter(":has(.code-lang)")
	}
	if snippet == nil || snippet.Size() == 0 {
		log.Fatal("error: No snippets found")
	} else if snippet.Size() > 1 {
		log.Print("error: Too many snippets are there:")
		snippets.Each(func(n int, s *goquery.Selection) {
			body := s.Find("pre").Text()
			if len(body) > 60 {
				body = body[0:60] + "..."
			}
			body = strconv.Quote(body)
			body = body[1 : len(body)-1]

			log.Printf("[%d] %q\t%s", n+1, s.Find(".code-lang").Text(), body)
		})
		log.Fatal("Specify one with -n")
	}

	filename := *flagFilename
	if filename == "" {
		filename = strings.TrimSpace(snippet.Find(".code-lang").Text())
	}
	if filename == "" {
		log.Fatal("error: Could not detect filename; specify with -o")
	}
	if *flagDirname != "" {
		filename = filepath.Join(*flagDirname, filename)
	}

	content := snippet.Find("pre").Text()
	if content == "" {
		log.Fatal("error: Could not find content")
	}

	perm := os.FileMode(0666)
	if *flagExecutable {
		perm = 0777
	}

	if _, err := os.Stat(filename); err == nil {
		if *flagForce {
			err := os.Remove(filename)
			if err != nil {
				log.Fatalf("error: rm %s: %s", filename, err)
			}
		} else {
			log.Fatalf("error: File %s exists; overwrite with -f", filename)
		}
	}

	err = ioutil.WriteFile(filename, []byte(content), perm)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Saved to %s", filename)
}
