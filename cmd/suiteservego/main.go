package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/suiteserve/suiteserve-go/internal/client"
	"github.com/suiteserve/suiteserve-go/internal/sstesting"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	projectFlag = flag.String("project", "",
		"The project of the suite, defaults to the current directory")
	reprintFlag = flag.Bool("reprint", false,
		"Reprint test event output")
	tagsFlag = flag.String("tags", "",
		"A comma-separated list of tags for the suite")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: suiteservego [options] url\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		flag.Usage()
		return
	}

	err := func() (err error) {
		project := *projectFlag
		if project == "" {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			project = filepath.Base(wd)
		}
		c := client.Open(url, project, strings.Split(*tagsFlag, ","))
		defer safeClose(c, &err)

		dec := json.NewDecoder(os.Stdin)
		for {
			var e sstesting.TestEvent
			if err := dec.Decode(&e); err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("read test event: %v", err)
			}
			if *reprintFlag {
				fmt.Print(e.Output)
			}

			if err := c.OnEvent(&e); err != nil {
				return err
			}
		}
	}()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func safeClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}
