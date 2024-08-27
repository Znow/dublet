package main

import (
	"hash/crc32"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var selectedFile *widget.Label = widget.NewLabel("NoFileYet")
var selectedFolder *widget.Label = widget.NewLabel("NoFolderYet")
var fileURI fyne.URI
var folderPath string = "NoFolderYet"

func main() {
	a := app.NewWithID(uuid.New().String())
	w := a.NewWindow("Dublet")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Path", Widget: selectedFolder},
			{Text: "Find Folder", Widget: widget.NewButton("Folder", func() { showFilePicker(w) })},
			{Widget: widget.NewButton("Scan!", func() { FindDuplicates(selectedFolder.Text) })},
		}}

	w.SetContent(form)
	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()

	// // Example usage of CompareFiles
	// same, err := CompareFiles("path/to/file1", "path/to/file2")
	// if err != nil {
	// 	fmt.Printf("Error comparing files: %v\n", err)
	// 	return
	// }

	// if same {
	// 	fmt.Println("The files are the same.")
	// } else {
	// 	fmt.Println("The files are different.")
	// }
}

// Show file picker and return selected file
func showFilePicker(w fyne.Window) {
	dialog.ShowFolderOpen(func(f fyne.ListableURI, err error) {
		// saveFile := "NoFileYet"
		folderPath := "NoFolderYet"
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if f == nil {
			return
		}
		// saveFile = f.URI().Path()
		// fileURI = f.URI()
		folderPath = f.Path()
		// selectedFile.SetText(saveFile)
		selectedFolder.SetText(folderPath)
	}, w)
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

// FileInfo holds the file path and its CRC32 hash.
type FileInfo struct {
	Path string
	Hash uint32
}

// FindDuplicates scans a directory for duplicate files.
func FindDuplicates(dir string) (map[uint32][]string, error) {
	files := make(map[uint32][]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		hash, err := computeCRC32(path)
		if err != nil {
			return err
		}

		files[hash] = append(files[hash], path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	duplicates := make(map[uint32][]string)
	for hash, paths := range files {
		if len(paths) > 1 {
			duplicates[hash] = paths
		}
	}

	return duplicates, nil
}

// computeCRC32 computes the CRC32 hash of a file.
func computeCRC32(path string) (uint32, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	hash := crc32.NewIEEE()
	if _, err := io.Copy(hash, f); err != nil {
		return 0, err
	}

	return hash.Sum32(), nil
}
