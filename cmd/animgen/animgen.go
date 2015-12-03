package main

import (
	"flag"
	"image"
	"image/color/palette"
	"log"
	"os"
	"path/filepath"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func loadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err = image.Decode(file)
	return
}

func writeMeta(dir string, data *Data) error {
	data.Name = "animgen"

	path := filepath.Join(dir, data.Name) + ".h"
	log.Printf("Writing metadata to %q...", path)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.ExecuteTemplate(file, "metaHeader", data)
	if err != nil {
		return err
	}

	path = filepath.Join(dir, data.Name) + ".c"
	log.Printf("Writing metadata to %q...", path)

	file, err = os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.ExecuteTemplate(file, "metaData", data)
	if err != nil {
		return err
	}

	return nil
}

func writeImage(dir string, data *Data) error {
	path := filepath.Join(dir, data.Name) + ".h"
	log.Printf("Writing header to %q...", path)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.ExecuteTemplate(file, "imageHeader", data)
	if err != nil {
		return err
	}

	path = filepath.Join(dir, data.Name) + ".c"
	log.Printf("Writing data to %q...", path)

	file, err = os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.ExecuteTemplate(file, "imageData", data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	data := &Data{
		Palette: palette.Plan9,
	}

	if flags.meta {
		err := writeMeta(flags.outdir, data)
		if err != nil {
			log.Fatalf("Failed to write metadata: %v\n", err)
		}

		return
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	img, err := loadImage(flag.Arg(0))
	if err != nil {
		log.Fatalf("Failed to load image from %q: %v", flag.Arg(0), err)
	}
	data.Image = img
	data.Name = TrimExt(filepath.Base(flag.Arg(0)))

	err = writeImage(flags.outdir, data)
	if err != nil {
		log.Fatalf("Failed to write header: %v", err)
	}
}
