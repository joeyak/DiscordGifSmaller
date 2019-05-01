package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/c2h5oh/datasize"
)

const maxSize = 256 * datasize.KB

type imgInfo struct {
	x, y int
	size datasize.ByteSize
}

func main() {
	for _, arg := range os.Args[1:] {
		_, err := os.Stat(arg)
		if err == nil {
			fmt.Printf("Processing: %s\n", arg)
			ext := filepath.Ext(arg)
			err = processFile(arg, arg[:len(arg)-len(ext)]+"Small"+ext, arg, 0, 0, 0)
			if err != nil {
			}
		}

		if err != nil {
			fmt.Printf("Could not process \"%s\": %v", arg, err)
		}
	}
}

func processFile(old, new, sizeFile string, x, y, pass int) error {
	size, err := getDataSize(sizeFile)
	if err != nil {
		return fmt.Errorf("could not get file sizes: %v", err)
	}

	if (x == 0 || y == 0) && pass > 0 {
		return errors.New("a width or height of 0 was given, either the file is very long in one direction, or the developer messed up")
	}

	if size > maxSize {
		if x == 0 || y == 0 {
			x, y, err = getFileSizes(old)
			if err != nil {
				return fmt.Errorf("could not get file sizes: %v", err)
			}
		}

		x--
		y--
		dim := fmt.Sprintf("%dx%d", x, y)

		fmt.Printf("[%d] Current size %s: Creating %s with a size of %s\n", pass, datasize.ByteSize(size).HR(), new, dim)
		cmd := exec.Command("gifsicle.exe", old, "--resize", dim, "--colors", "256", "-o", new)

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("could not run process: %v", err)
		}
		return processFile(old, new, new, x, y, pass+1)
	}
	return nil
}

func getDataSize(path string) (datasize.ByteSize, error) {
	r, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("could not open file for size: %v", err)
	}
	defer r.Close()

	info, err := r.Stat()
	if err != nil {
		return 0, err
	}
	return datasize.ByteSize(info.Size()), nil
}

func getFileSizes(path string) (int, int, error) {
	r, err := os.Open(path)
	if err != nil {
		return 0, 0, fmt.Errorf("could not open file for size: %v", err)
	}
	defer r.Close()

	img, _, err := image.DecodeConfig(r)
	if err != nil {
		return 0, 0, err
	}
	return img.Width, img.Height, nil
}
