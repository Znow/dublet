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

func main() {
	dir := "path/to/directory"
	duplicates, err := FindDuplicates(dir)
	if err != nil {
		fmt.Printf("Error finding duplicates: %v\n", err)
		return
	}

	for hash, paths := range duplicates {
		fmt.Printf("Duplicate files with hash %x:\n", hash)
		for _, path := range paths {
			fmt.Println(path)
		}
	}
}