package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var flags struct {
	pkg string
	out string
	in  string
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [option] <src>\n", os.Args[0])
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&flags.out, "out", "", "Output file. If blank, it will be the input filename - ext + .go in current dir.")
	flag.StringVar(&flags.pkg, "pkg", "main", "Generated Go file will be in package `name`.")

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	flags.in = flag.Arg(0)

	if flags.out == "" {
		flags.out = filepath.Base(TrimExt(flags.in)) + ".go"
	}
}

func main() {
	name := filepath.Base(TrimExt(flags.out))

	fmt.Printf("Reading from %q...\n", flags.in)
	infile, err := os.Open(flags.in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open %q: %v", flags.in, err)
		os.Exit(1)
	}
	defer infile.Close()

	fmt.Printf("Writing to %q...\n", name+".go")
	outfile, err := os.Create(name + ".go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create %q: %v", name+".go", err)
		os.Exit(1)
	}
	defer outfile.Close()

	fmt.Println("Writing data...")
	err = tmpl.ExecuteTemplate(outfile, "data", &Data{
		Pkg:  flags.pkg,
		Name: name,
		In:   infile,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write data: %v", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
}
