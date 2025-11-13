package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kost0/L4/internal/client"
	"github.com/Kost0/L4/internal/cut"
	"github.com/spf13/pflag"
)

func main() {
	start := time.Now()

	fields := pflag.StringP("fields", "f", "", "specify fields to output")
	delimiter := pflag.StringP("delimiter", "d", "\t", "symbol for separation")
	separated := pflag.BoolP("separated", "s", false, "only separated")

	pflag.Parse()

	inputFiles := make([]string, 0)

	if pflag.NArg() > 0 {
		inputFiles = pflag.Args()
	} else {
		_, err := fmt.Fprintln(os.Stderr, "Not enough arguments")
		if err != nil {
			os.Exit(1)
		}
		os.Exit(1)
	}

	numFields, err := cut.ParseFields(*fields)
	if err != nil {
		log.Fatal(err)
	}

	opts := cut.CutOptions{
		Fields:    numFields,
		Delimiter: *delimiter,
		Separated: *separated,
	}

	res := client.StartClient(inputFiles, &opts)

	for _, line := range res {
		for i := range len(line) - 1 {
			fmt.Print(line[i] + opts.Delimiter)
		}
		fmt.Print(line[len(line)-1] + "\n")
	}

	elapsed := time.Since(start)
	log.Printf("Completed in %s", elapsed)
}
