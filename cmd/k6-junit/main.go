package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/l-ross/k6-junit/pkg/summary"
)

func main() {
	in := flag.String("in", "", "location of the k6 json summary")
	out := flag.String("out", "", "location to write the JUnit summary to, if not specified prints to the console")

	flag.Parse()

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
