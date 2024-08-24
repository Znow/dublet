package main

import (
	"hash/crc32"
	"io"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Dublet")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}

// CompareFiles checks if two files are the same by comparing their sizes and CRC32 hashes.
func CompareFiles(file1, file2 string) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	// Check file sizes
	fi1, err := f1.Stat()
	if err != nil {
		return false, err
	}
	fi2, err := f2.Stat()
	if err != nil {
		return false, err
	}

	if fi1.Size() != fi2.Size() {
		return false, nil
	}

	// Compute CRC32 hashes
	hash1 := crc32.NewIEEE()
	if _, err := io.Copy(hash1, f1); err != nil {
		return false, err
	}

	hash2 := crc32.NewIEEE()
	if _, err := io.Copy(hash2, f2); err != nil {
		return false, err
	}

	return hash1.Sum32() == hash2.Sum32(), nil
}

// You can use this method by calling CompareFiles("path/to/file1", "path/to/file2"). It will return true if the files are the same, and false otherwise.
