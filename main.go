package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	octal   bool
	warnCsc bool
)

func getByte(s string) (byte, error) {
	i, err := strconv.ParseUint(s, 16, 8)
	return byte(i), err
}

func checksum(bs []byte) (bool, int) {
	var ret int
	for _, b := range bs {
		ret += int(b)
	}
	return (ret & 0xFF) == 0, ret
}

func getData(bytesAndChecksum string) []byte {
	var bs []byte
	for i := 0; i < len(bytesAndChecksum); i += 2 {
		b, err := getByte(bytesAndChecksum[i : i+2])
		if err != nil {
			// Bad behaviour, should probably warn or fail, but cba.
			// Will probably fail checksum at this point anyways.
			bs = append(bs, 0)
		} else {
			bs = append(bs, b)
		}
	}
	return bs
}

func printHex(r io.Reader) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		ll := len(line)
		if ll < 11 { // Haha, hope you have a decent font
			continue
		}
		if line[0] != ':' {
			continue
		}

		nbytes, err := strconv.ParseUint(line[1:3], 16, 8)
		nberr := err != nil
		addr, err := strconv.ParseUint(line[3:7], 16, 8)
		addrerr := err != nil
		csc := line[ll-2:]
		bs := getData(line[1:])
		goodChecksum, _ := checksum(bs)
		databytes := bs[4 : len(bs)-1]

		// Checksum warning
		if !goodChecksum && warnCsc {
			fmt.Fprintf(os.Stderr, "Bad checksum %s (%s)\n", csc, line)
		}

		// Number of bytes
		if nberr {
			color.Set(color.FgRed)
			fmt.Printf("?%s", line[1:3])
			color.Set(color.Reset)
		} else {
			color.Set(color.ReverseVideo)
			color.Set(color.FgMagenta)
			fmt.Printf("%db", nbytes)
		}
		color.Set(color.Reset)
		fmt.Print("\t")

		// Address
		if addrerr {
			color.Set(color.FgRed)
			fmt.Printf("?%s", line[1:3])
		} else {
			color.Set(color.ReverseVideo)
			color.Set(color.FgBlue)
			if octal {
				fmt.Printf("%06o", addr)
			} else {
				fmt.Printf("%04x", addr)
			}
		}
		color.Set(color.Reset)
		fmt.Print("\t")

		// Record type
		color.Set(color.ReverseVideo)
		color.Set(color.FgRed)
		fmt.Printf("%02x", bs[3])
		color.Set(color.Reset)
		fmt.Print("\t")

		// Data
		for _, by := range databytes {
			color.Set(color.ReverseVideo)
			color.Set(color.FgCyan)
			if octal {
				fmt.Printf("%03o", by)
			} else {
				fmt.Printf("%02x", by)
			}
			color.Set(color.Reset)
			fmt.Print(" ")
		}
		fmt.Print("\t")

		// Checksum
		if goodChecksum {
			color.Set(color.ReverseVideo)
			color.Set(color.FgGreen)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Print(csc)
		color.Set(color.Reset)
		fmt.Print("\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

}

func main() {
	flag.BoolVar(&octal, "o", false, "print an octal representation (default hex)")
	flag.BoolVar(&warnCsc, "w", false, "warn on bad checksum")
	flagColor := flag.Bool("c", false, "color output")
	flagHeader := flag.Bool("H", false, "print a header")
	flag.Parse()
	if *flagHeader {
		fmt.Printf("NBytes\tAddr\tRecType\tData\tChecksum\n")
	}
	color.NoColor = !*flagColor
	args := flag.Args()

	if len(args) < 1 {
		printHex(os.Stdin)
		return
	}

	for _, arg := range args {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} else {
			defer f.Close()
			printHex(f)
		}
	}
}
