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
	time.Sleep(1 * time.Second)

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

	files, err := f.SortFolderFilesDesc(folderPath)
	if err != nil {
		return err
	}

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

	for _, file := range files {
		if filepath.Ext(file) == ".zip" {
			fullPath := filepath.Join(folderPath, file)
			if hasFileBeenModified(fullPath, time.Minute) {
				callback(fullPath)
			} else {
				go func(fullPath string) {
					if callbackDownloadedFile(fullPath) {
						callback(fullPath)
					}
				}(fullPath)
			}
		}
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
	err = watcher.Add(folderPath)
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
