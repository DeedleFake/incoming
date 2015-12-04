package main

import (
	"flag"
	"fmt"
	"os"
)

var flags struct {
	in string
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [option] <src>\n", os.Args[0])
		// NOTE: Enable if options are added.
		//fmt.Fprintln(os.Stderr)
		//fmt.Printf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	flags.in = flag.Arg(0)
}

func main() {
	name := TrimExt(flags.in)
	fmt.Printf("Name: %q\n", name)

	infile, err := os.Open(flags.in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open %q: %v", flags.in, err)
		os.Exit(1)
	}
	defer infile.Close()

	outfile, err := os.Create(name + ".go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create %q: %v", name+".go", err)
		os.Exit(1)
	}
	defer outfile.Close()

	fmt.Println("Writing data...")
	err = tmpl.ExecuteTemplate(outfile, "data", &Data{
		Name: name,
		In:   infile,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write data: %v", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
}
