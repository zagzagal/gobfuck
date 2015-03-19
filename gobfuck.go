package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
)

var verbose = flag.Bool("v", false, "Verbose mode")
var debug = flag.Bool("d", false, "Debug mode")
var file = flag.String("file", "", "File to read")

func main() {
	flag.Parse()
	in := bufio.NewReader(os.Stdin)
	var data []byte
	if *file != "" {
		data, _ = ioutil.ReadFile(*file)
	} else {
		data, _, _ = in.ReadLine()
	}
	bf := NewBrainFuck()
	bf.instructions(data)
	bf.exec(in, os.Stdout, *verbose, *debug)
}
