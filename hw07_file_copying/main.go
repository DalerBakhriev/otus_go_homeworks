package main

import (
	"flag"
	"log"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	file, err := os.Open(from)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, limit)
	_, err = file.ReadAt(buf, offset)
	if err != nil {
		log.Fatal(err)
	}

	dst, err := os.Create(to)
	if err != nil {
		log.Fatal(err)
	}
	_, err = dst.Write(buf)
	if err != nil {
		log.Fatal(err)
	}
}
