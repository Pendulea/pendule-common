package pcommon

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type file struct{}

var File = file{}

func (f file) EnsureDir(path string) error {
	err := os.MkdirAll(path, 0755) // Creates the directory with rwx permissions for owner and rx for group and others
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return nil
}

// SortFolderFiles sorts ZIP files in a folder
func (f file) SortFolderFilesDesc(folderPath string) ([]string, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	zipFiles := []string{}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".zip" {
			zipFiles = append(zipFiles, file.Name())
		}
	}

	sort.Slice(zipFiles, func(i, j int) bool {
		return zipFiles[i] > zipFiles[j]
	})

	return zipFiles, nil
}

// getFileSize returns the size of the file
func (f file) GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// UnzipFile extracts a zip archive specified by zipPath into the directory outputPath.
func (f file) UnzipFile(zipPath, outputPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err // If opening the zip file fails, return the error
	}
	defer r.Close()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(outputPath, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(fpath, filepath.Clean(outputPath)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make directory if it is not present
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
		} else {
			// Make file's directory if it is not present (important if zip contains nested directories)
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				outFile.Close() // Close file handle on error opening zip content
				return err
			}

			_, err = io.Copy(outFile, rc) // Copy contents to the file
			// Close file handles regardless of io.Copy results
			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}

	time.Sleep(50 * time.Millisecond)
	return nil
}

// remove filepath
func (f file) RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to remove file: %v", err)
	}
	return nil
}

type FileCallback func(fileName string)

// Initializes the watcher and provides an initial list of zip files.
func (f file) InitFolderWatcher(folderPath string, callback FileCallback, watcher *fsnotify.Watcher) error {
	// First, list all existing zip files in the folder.

	callbackDownloadedFile := func(filePath string) bool {
		retry := 10
		if hasFileBeenModified(filePath, time.Minute) {
			return true
		}

		for i := 0; i < retry; i++ {
			if hasFileBeenModified(filePath, time.Minute) {
				return true
			}
			time.Sleep(1 * time.Minute)
		}
		return false
	}

	// Start a goroutine to handle the events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					if filepath.Ext(event.Name) == ".zip" {
						go func(fullPath string) {
							if callbackDownloadedFile(fullPath) {
								callback(fullPath)
							}
						}(event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add the directory to the watcher.
	err := watcher.Add(folderPath)
	if err != nil {
		return err
	}

	// Block forever (or until the watcher is stopped another way)
	select {}
}

func hasFileBeenModified(filePath string, duration time.Duration) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) >= duration
}

func (f file) ListenPairJSONFileChange(pairsPath string, callback func(path string)) {
	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start a goroutine to handle events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					callback(pairsPath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add the directory to the watcher
	err = watcher.Add(pairsPath)
	if err != nil {
		log.Fatal(err)
	}

	// Block forever
	select {}
}

func (f file) GetFolderSize(folderPath string) (int64, error) {
	var totalSize int64

	// Walk through the folder and its subdirectories
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Add the size of regular files to the total
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return totalSize, nil
}

type FileInfo struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
	Size int64  `json:"size"`
}

func (f file) GetSortedFilenamesByDate(directoryPath string) ([]FileInfo, error) {
	// Read all files in the directory
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold file info
	var fileInfos []FileInfo

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(directoryPath, file.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return nil, err
			}
			fileInfos = append(fileInfos, FileInfo{Name: file.Name(), Time: fileInfo.ModTime().Unix(), Size: fileInfo.Size()})
		}
	}

	// Sort files by modification time in descending order
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Time > fileInfos[j].Time
	})

	return fileInfos, nil
}

// addFileToZip adds a file or directory to the zip archive.
func addFileToZip(zipWriter *zip.Writer, filePath string, basePath string) error {
	// Open the file to be added
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a zip header from the file information
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	// Set the header name to be the relative path from the base path
	header.Name, err = filepath.Rel(basePath, filePath)
	if err != nil {
		return err
	}

	// Ensure the header name uses forward slashes for compatibility
	header.Name = strings.ReplaceAll(header.Name, string(os.PathSeparator), "/")

	// If the file is a directory, ensure the zip entry reflects that
	if fileInfo.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}

	// Create the zip file entry
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// If the file is a directory, return now
	if fileInfo.IsDir() {
		return nil
	}

	// Copy the file data to the zip entry
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}

// ZipDirectory zips the contents of the source directory to the target zip file.
func ZipDirectory(source string, target string) error {
	// Create a file to write the zip archive to
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a new zip archive writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk the directory tree and add each file to the zip archive
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return addFileToZip(zipWriter, path, source)
	})

	return err
}

func (f file) CopyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file for writing
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	// Copy the contents from the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy contents: %w", err)
	}

	// Ensure the contents are flushed to the destination file
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}
