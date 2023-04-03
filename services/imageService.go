package services

import (
	"api/utils"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func ScanAndConvert() {
	utils.ReportMessage("- Scanning and converting images -")
	files, err := listFilesInFolder("./data")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
mainloop:
	for _, file := range files {
		if !strings.Contains(file, ".webp") {
			fileOld := strings.Split(file, ".")
			fileNew := strings.ReplaceAll(file, fmt.Sprintf(".%s", fileOld[1]), ".webp")
			utils.ReportMessage(fmt.Sprintf("Converting %s to %s", file, fileNew))
			if fileExists(fmt.Sprintf("./data_final/%s", fileNew)) {
				continue mainloop
			}
			o, err := exec.Command("bash", "-c", fmt.Sprintf("cwebp -q 95 -resize 500 500 ./data/%s -o ./data_final/%s", file, fileNew)).Output()
			if err != nil {
				utils.WrapErrorLog(err.Error())
				continue mainloop
			}
			utils.WrapErrorLog(string(o))
		} else {
			fileOld := strings.Split(file, ".")
			fileNew := strings.ReplaceAll(file, fmt.Sprintf(".%s", fileOld[1]), ".webp")
			utils.ReportMessage(fmt.Sprintf("Converting %s to %s", file, fileNew))
			if fileExists(fmt.Sprintf("./data_final/%s", fileNew)) {
				continue mainloop
			}
			o, err := exec.Command("bash", "-c", fmt.Sprintf("cwebp -q 95 -resize 500 500 ./data/%s -o ./data_final/%s", file, fileNew)).Output()
			if err != nil {
				utils.WrapErrorLog(err.Error())
				continue mainloop
			}
			utils.WrapErrorLog(string(o))
		}
	}
	//convert files to webp
}

func listFilesInFolder(folderPath string) ([]string, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}
	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}
	return filenames, nil
}

func copyFile(source string, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
