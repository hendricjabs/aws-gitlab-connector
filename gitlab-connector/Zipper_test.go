package main

import (
	"crypto/rand"
	"testing"
)

func TestZipper(t *testing.T) {
	t.Run("Repository not available", func(t *testing.T) {
		files := make([]FileWithContent, 1)
		fileContent := make([]byte, 64)
		_, err := rand.Read(fileContent)

		file := &FileWithContent{Name: "testfile", Bytes: fileContent}
		files = append(files, *file)

		testData := &ZipDataInput{Files: files}

		out, err := ZipData(testData)
		if err != nil {
			t.Fatal("Error: Should Zip Data without error")
		}

		if out.Bytes == nil || out.Size <= 0 {
			t.Fatal("Error: Zip should not be empty")
		}

	})

}
