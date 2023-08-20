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
			if strings.Contains(file, ".txt") || strings.Contains(file, ".html") {
				err := copyFile(fmt.Sprintf("./data/%s", file), fmt.Sprintf("./data_final/%s", file))
				if err != nil {
					utils.WrapErrorLog(err.Error())
				}
				continue mainloop
			}
			fileOld := strings.Split(file, ".")
			fileNew := strings.ReplaceAll(file, fmt.Sprintf(".%s", fileOld[1]), ".webp")
			if utils.FileExists(fmt.Sprintf("./data_final/%s", fileNew)) {
				continue mainloop
			}
			utils.ReportMessage(fmt.Sprintf("Converting %s to %s", file, fileNew))
			err := exec.Command("bash", "-c", fmt.Sprintf("cwebp -q 50 -resize 500 0 ./data/%s -o ./data_final/%s", file, fileNew)).Run()
			if err != nil {
				utils.WrapErrorLog(err.Error())
				continue mainloop
			}
		} else {
			fileOld := strings.Split(file, ".")
			fileNew := strings.ReplaceAll(file, fmt.Sprintf(".%s", fileOld[1]), ".webp")
			if utils.FileExists(fmt.Sprintf("./data_final/%s", fileNew)) {
				continue mainloop
			}
			utils.ReportMessage(fmt.Sprintf("Converting %s to %s", file, fileNew))
			err := exec.Command("bash", "-c", fmt.Sprintf("cwebp -q 50 -resize 500 0 ./data/%s -o ./data_final/%s", file, fileNew)).Run()
			if err != nil {
				utils.WrapErrorLog(err.Error())
				continue mainloop
			}
		}
	}
	_ = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s/api/data/*", utils.GetHomeDir()))).Run()
	utils.ReportMessage("- Done -")
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
