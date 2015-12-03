package main

import (
	"flag"
	"fmt"
	"os"
)

var flags struct {
	outdir string

	start  int
	width  int
	height int
	frames int

	meta bool
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] <image>\n", os.Args[0])
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&flags.outdir, "out", ".", "Output to dir/name.h and `dir`/name.c.")
	flag.IntVar(&flags.start, "start", 0, "Look for the image at (0, `y`).")
	flag.IntVar(&flags.width, "width", -1, "The width of a frame.")
	flag.IntVar(&flags.height, "height", -1, "The height of a frame.")
	flag.IntVar(&flags.frames, "frames", 1, "The number of frames.")

	flag.BoolVar(&flags.meta, "meta", false, "Generate metadata in dir/animgen.h and dir/animgen.c instead of image data.")

	flag.Parse()
}
