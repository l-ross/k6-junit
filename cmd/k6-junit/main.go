package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/l-ross/k6-junit/pkg/summary"
)

var version = "UNVERSIONED"

func main() {
	in := flag.String("in", "", "location of the k6 json summary")
	out := flag.String("out", "", "location to write the JUnit summary to, if not specified prints to the console")
	ver := flag.Bool("version", false, "print ver and exit")

	flag.Parse()

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	f, err := os.Open(*in)
	if err != nil {
		log.Fatalf("failed to open input file %q: %s", *in, err)
	}
	defer f.Close()

	s, err := summary.NewSummaryFromReader(f)
	if err != nil {
		log.Fatalf("failed to parse summary: %s", err)
	}

	j, err := s.JUnit()
	if err != nil {
		log.Fatalf("failed to marshal summary to junit: %s", err)
	}

	if *out == "" {
		fmt.Println(string(j))
	} else {
		err = ioutil.WriteFile(*out, j, 0600)
		if err != nil {
			log.Fatalf("failed to write to out file %q: %s", *out, err)
		}
	}
}
