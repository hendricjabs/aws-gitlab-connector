package main

import (
	"archive/zip"
	"bytes"
)

// FileWithContent: the name of the file and the content
type FileWithContent struct {
	Name  string
	Bytes []byte
}

type ZipDataInput struct {
	Files []FileWithContent
}

type ZipDataOutput struct {
	Bytes []byte
	Size  int64
}

// ZipData: Create a Zip file
func ZipData(input *ZipDataInput) (output *ZipDataOutput, err error) {

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)

	// Add some files to the archive.
	for _, file := range input.Files {
		zipFile, err := zipWriter.Create(file.Name)
		if err != nil {
			return nil, err
		}
		_, err = zipFile.Write(file.Bytes)
		if err != nil {
			return nil, err
		}
	}
	// Make sure to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	//write the zipped file to the disk
	return &ZipDataOutput{Bytes: buf.Bytes(), Size: int64(buf.Len())}, err

}
