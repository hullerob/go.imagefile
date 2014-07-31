package main

import (
	"github.com/hullerob/go.imagefile"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		usage()
		os.Exit(1)
	}
	m, _, err := image.Decode(os.Stdin)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	err = imagefile.Encode(os.Stdout, m)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func usage() {
	os.Stderr.WriteString("usage: if2png < pic.png > pic.if\n")
}
