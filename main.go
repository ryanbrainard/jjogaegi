package main

import (
	"github.com/ryanbrainard/jjogaegi/jjogaegi"
	"os"
)

func main() {
	jjogaegi.Parse(os.Stdin, os.Stdout)
}
