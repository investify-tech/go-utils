package file

import (
	"github.com/investify-tech/go-utils/must"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Exists checks if a file or dir is existing
func Exists(fileOrDirPath string) bool {
	if _, err := os.Stat(fileOrDirPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// WriteFile writes a file and encapsulates error check and file rights setting
func WriteFile(dstFilePath string, fileContent string, executable bool) {
	var fileMode os.FileMode
	if executable {
		fileMode = 0755
	} else {
		fileMode = 0644
	}
	contentBytes := []byte(fileContent)
	must.Void(os.WriteFile(dstFilePath, contentBytes, fileMode))
}

// ReadFile reads and returns the content of a file and encapsulates error handling
func ReadFile(filePath string) string {
	fileBytes := must.AnySlice(os.ReadFile(filePath))
	return string(fileBytes)
}

// ReadFileLines ReadFile reads and returns the content of a file as an arrays of lines and encapsulates error handling
func ReadFileLines(filePath string) []string {
	fileBytes := must.AnySlice(os.ReadFile(filePath))
	return strings.Split(string(fileBytes), "\n")
}

// CopyFile copies a single file from src to dst
func CopyFile(srcPath, dstPath string) error {
	return CopyFileAndReplaceContent(srcPath, dstPath, nil)
}

// CopyFileAndReplaceContent copies a single file from src to dst and does replacements in the file's content
// as provided in the given map
func CopyFileAndReplaceContent(srcFilePath, dstFilePath string, replacements map[string]string) error {
	var err error
	var srcFileRef *os.File
	var dstFileRef *os.File
	var srcFileInfo os.FileInfo

	if srcFileRef, err = os.Open(srcFilePath); err != nil {
		return err
	}
	defer srcFileRef.Close()

	if dstFileRef, err = os.Create(dstFilePath); err != nil {
		return err
	}
	defer dstFileRef.Close()

	if _, err = io.Copy(dstFileRef, srcFileRef); err != nil {
		return err
	}
	if srcFileInfo, err = os.Stat(srcFilePath); err != nil {
		return err
	}

	for old, new := range replacements {
		ReplaceFileContent(dstFilePath, old, new)
	}

	return os.Chmod(dstFilePath, srcFileInfo.Mode())
}

// CopyDir copies a whole directory recursively
func CopyDir(srcDirPath, dstDirPath string) error {
	return CopyDirAndReplaceContentInFiles(srcDirPath, dstDirPath, nil)
}

// CopyDirAndReplaceContentInFiles copies a whole directory recursively and does replacements in all files' contents
// as provided in the given map
func CopyDirAndReplaceContentInFiles(srcDirPath, dstDirPath string, replacements map[string]string) error {
	var err error
	var srcDirContent []os.FileInfo
	var srcDirInfo os.FileInfo

	if srcDirInfo, err = os.Stat(srcDirPath); err != nil {
		return err
	}
	if err = os.MkdirAll(dstDirPath, srcDirInfo.Mode()); err != nil {
		return err
	}
	if srcDirContent, err = ioutil.ReadDir(srcDirPath); err != nil {
		return err
	}

	for _, fileInfo := range srcDirContent {
		srcPath := path.Join(srcDirPath, fileInfo.Name())
		dstPath := path.Join(dstDirPath, fileInfo.Name())

		if fileInfo.IsDir() {
			must.Void(CopyDirAndReplaceContentInFiles(srcPath, dstPath, replacements))
		} else {
			must.Void(CopyFileAndReplaceContent(srcPath, dstPath, replacements))
		}
	}
	return nil
}

func CopyDirAndRenameFiles(srcDirPath, dstDirPath string, replacements map[string]string) error {
	err := CopyDir(srcDirPath, dstDirPath)
	if err != nil {
		return err
	}

	err = RenameFilesInDir(dstDirPath, replacements)
	return err
}

func RenameFilesInDir(dirPath string, replacements map[string]string) error {
	err := filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		fileName := f.Name()
		for key, val := range replacements {
			if strings.Contains(fileName, key) {
				dir := filepath.Dir(path)
				newFileName := strings.Replace(fileName, key, val, 1)
				newFilePath := filepath.Join(dir, newFileName)
				os.Rename(path, newFilePath)
			}
		}
		return nil
	})
	return err
}

// ReplaceFileContent replaces value in the file's content by replacement
func ReplaceFileContent(filePath string, value string, replacement string) {
	input := must.AnySlice(os.ReadFile(filePath))
	fileInfo, err := os.Stat(filePath)
	must.Void(err)

	output := strings.ReplaceAll(string(input), value, replacement)

	must.Void(os.WriteFile(filePath, []byte(output), fileInfo.Mode().Perm()))
}

// CleanDir "cleans" the given directory by deleting and recreating it with "allow everything" rights
func CleanDir(dirPath string) {
	must.Void(os.RemoveAll(dirPath))
	must.Void(os.MkdirAll(dirPath, 0777))
}

// CleanDirFromFiles "cleans" the given directory from files, i.e. deletes all files (but no directories) from it
func CleanDirFromFiles(dirPath string) {
	for _, dirEntry := range must.AnySlice(os.ReadDir(dirPath)) {
		if dirEntry.IsDir() {
			continue
		}
		must.Void(os.Remove(path.Join(dirPath, dirEntry.Name())))
	}
}

// CreateDir creates the given directory with "allow everything" rights
func CreateDir(dirPath string) {
	must.Void(os.MkdirAll(dirPath, 0777))
}

// RemoveDir removes the given directory - just there for convenience
func RemoveDir(dirPath string) {
	must.Void(os.RemoveAll(dirPath))
}
