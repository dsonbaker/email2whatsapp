package automationWhatsapp

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(targetURL string, nameFile string, folderName string) {

	response, err := http.Get(targetURL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	os.MkdirAll(folderName, os.ModePerm)
	arquivo, err := os.Create(filepath.Join(folderName, nameFile))
	if err != nil {
		panic(err)
	}
	defer arquivo.Close()

	_, err = io.Copy(arquivo, response.Body)
	if err != nil {
		panic(err)
	}
}
