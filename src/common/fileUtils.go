package common

import (
	"os"
	"log"
	"strings"
)

func GiveFileDescpt(fileLocation string) *os.File{
	file, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func DirectoryPath(filePath string) string{
	sl := strings.Split(filePath, "/")
	return strings.Join(sl[:len(sl)-1],"/")
}