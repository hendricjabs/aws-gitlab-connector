package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"

	"github.com/src-d/go-git/storage/memory"
)

// GitCloneAndZipInput: repositoryUrl and apiKey for cloning the repository
type GitCloneAndZipInput struct {
	repositoryUrl string
	apiKey        string
	username      string
	password      string
}

// GitCloneAndZipOutput: size and content of the zip file
type GitCloneAndZipOutput struct {
	Size int64
	File []byte
}

// GitCloneAndZip: Clone the given repository and put it into a zip file
func GitCloneAndZip(input *GitCloneAndZipInput) (GitCloneAndZipOutput, error) {
	log.WithField("input", input).Info("Validate Input")
	if isValid(*input) == false {
		return GitCloneAndZipOutput{}, errors.New("Input is invalid or values are missing!")
	}

	log.Info("Initialize In-Memory Filesystem")
	// Filesystem abstraction based on memory
	fs := memfs.New()

	// Git objects storer based on memory
	storer := memory.NewStorage()

	log.Infof("git clone %s", input.repositoryUrl)

    var auth transport.AuthMethod
    if input.apiKey != "" {
	   auth = &http.TokenAuth{Token: input.apiKey}
    } else {
       auth = &http.BasicAuth{Username: input.username, Password: input.password}
	}
	// Clones the given repository in memory
	_, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:  input.repositoryUrl,
		Auth: auth,
	})
	if CheckIfError(err) {
		return GitCloneAndZipOutput{}, err
	}

	files, err := readDir("/", fs)
	if CheckIfError(err) {
		return GitCloneAndZipOutput{}, err
	}

	log.Info("Create ZIP File")

	// Zip Files
	zipFile, err := ZipData(&ZipDataInput{Files: files})
	if CheckIfError(err) {
		return GitCloneAndZipOutput{}, err
	}

	return GitCloneAndZipOutput{File: zipFile.Bytes, Size: zipFile.Size}, err
}

// isValid: Check if input parameters are not empty
func isValid(event GitCloneAndZipInput) bool {
	if event.repositoryUrl == "" {
		log.Error("repositoryUrl may not be empty!")
		return false
	}
	if event.apiKey == "" {
		log.Warn("apiKey is empty!")
	}
	return true
}

// readDir: recursively read in-memory directory and return the files
func readDir(path string, fs billy.Filesystem) ([]FileWithContent, error) {
	fileInfos, err := fs.ReadDir(path)
	CheckIfError(err)

	var files []FileWithContent

	for _, fileInfo := range fileInfos {
		// Ignore path if it is equal to the path seperator char
		if path == string(os.PathSeparator) {
			path = ""
		}
		newPath := fmt.Sprintf("%s%c%s", path, os.PathSeparator, fileInfo.Name())

		if fileInfo.IsDir() {
			log.WithField("Path", newPath).Info("Read directory")

			fwc, err := readDir(newPath, fs)
			if CheckIfError(err) {
				return []FileWithContent{}, err
			}
			files = append(files, fwc...)
		} else {
			log.WithField("Path", newPath).Info("Read file")
			file, err := readFile(newPath, fs)
			if CheckIfError(err) {
				return []FileWithContent{}, err
			}
			files = append(files, *file)
		}
	}
	return files, err
}

// readFile: read a file and return the content
func readFile(path string, fs billy.Filesystem) (*FileWithContent, error) {
	reader, err := fs.Open(path)
	if CheckIfError(err) {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(reader)
	if CheckIfError(err) {
		return nil, err
	}

	return &FileWithContent{Name: path, Bytes: bytes}, err
}
