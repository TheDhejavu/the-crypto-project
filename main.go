package main

import (
	"io"
	"os"
)

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := Copy("/the-crypto-project/", "/init_1")
	if err != nil {
		panic(err)
	}
}
