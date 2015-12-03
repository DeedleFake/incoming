package main

import (
	"image"
	"image/color"
	"strings"
	"text/template"
)

var (
	tmpl *template.Template
)

func init() {
	const (
		boiler = `// This file was auto-generated.
// You may not want to edit it manually.`

		metaHeader = `{{template "boiler"}}

#ifndef {{toupper .Name}}_H
#define {{toupper .Name}}_H

#define PALETTE_LENGTH {{len .Palette}}
extern const unsigned short palette_data[PALETTE_LENGTH];

#endif`

		metaData = `{{template "boiler"}}

#include "{{.Name}}.h"

const unsigned short palette_data[PALETTE_LENGTH] = {
{{range .Palette}}{{togba . | printf "\t0x%04X,"}}
{{end}}};`

		imageHeader = `{{template "boiler"}}

#ifndef {{toupper .Name}}_H
#define {{toupper .Name}}_H

#include "animgen.h"

#define {{toupper .Name}}_WIDTH {{.Width}}
#define {{toupper .Name}}_HEIGHT {{.Height}}
#define {{toupper .Name}}_FRAMES {{.NumFrames}}

extern const unsigned char {{.Name}}_data[{{.NumFrames}}][{{mult .Width .Height}}];

#endif`

		imageData = `{{template "boiler"}}

#include "{{.Name}}.h"

const unsigned char {{.Name}}_data[{{.NumFrames}}][{{mult .Width .Height}}] = {
	{{range .Frames}}{
		{{range .}}{{.}},
		{{end}}
	},{{end}}
};`
	)

	tmpl = new(template.Template)
	tmpl.Funcs(template.FuncMap{
		"toupper": strings.ToUpper,
		"togba":   GBAModel.Convert,

		"mult": func(i1, i2 int) int {
			return i1 * i2
		},
	})

	template.Must(tmpl.New("boiler").Parse(boiler))

	template.Must(tmpl.New("metaHeader").Parse(metaHeader))
	template.Must(tmpl.New("metaData").Parse(metaData))

	template.Must(tmpl.New("imageHeader").Parse(imageHeader))
	template.Must(tmpl.New("imageData").Parse(imageData))
}

type Data struct {
	Name    string
	Palette color.Palette
	Image   image.Image
}

func (d *Data) Width() int {
	if flags.width < 0 {
		return d.Image.Bounds().Dx()
	}

	return flags.width
}

func (d *Data) Height() int {
	if flags.height < 0 {
		return d.Image.Bounds().Dy()
	}

	return flags.height
}

func (d *Data) NumFrames() int {
	return flags.frames
}

func (d *Data) Frames() <-chan (<-chan int) {
	frames := make(chan (<-chan int), d.NumFrames())

	go func() {
		defer close(frames)

		for i := 0; i < d.NumFrames(); i++ {
			c := make(chan int)
			frames <- c

			go func(frame int) {
				defer close(c)

				for y := flags.start; y < flags.start+d.Height(); y++ {
					for x := d.Width() * frame; x < (d.Width()*frame)+d.Width(); x++ {
						col := d.Image.At(x, y)
						c <- d.Palette.Index(col)
					}
				}
			}(i)
		}
	}()

	return frames
}
