package main

import (
	"bufio"
	"io"
	"os"
	"text/template"
)

const (
	dataTmpl = `// This file was auto-generated.

package {{.Pkg}}

var {{.Name}}Data = [{{.Size}}]byte{
{{range .Data}}	{{printf "0x%02X" .}},
{{end}}}`
)

var tmpl = new(template.Template)

func init() {
	template.Must(tmpl.New("data").Parse(dataTmpl))
}

type Data struct {
	Pkg  string
	Name string
	In   *os.File

	stat os.FileInfo
}

func (data *Data) Size() int {
	if data.stat == nil {
		data.stat, _ = data.In.Stat()
	}

	return int(data.stat.Size())
}

func (data *Data) Data() <-chan byte {
	out := make(chan byte)
	go func() {
		defer close(out)

		r := bufio.NewReader(data.In)
		for {
			c, err := r.ReadByte()
			if err != nil {
				if err == io.EOF {
					break
				}

				// Hmmm... I can't think of a good way to handle an error
				// here, so how about crashing?
				panic(err)
			}

			out <- c
		}
	}()

	return out
}
