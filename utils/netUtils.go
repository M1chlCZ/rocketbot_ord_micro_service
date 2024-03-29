package utils

import (
	"api/grpcClient"
	"api/grpcModels"
	"api/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/image/webp"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GETRequest[T any](url string) (T, error) {
	var data T
	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	if resp.StatusCode != http.StatusOK {
		if body != nil {
			err = json.Unmarshal(body, &data)
			if err != nil {
				return data, err
			}
			return data, errors.New("GET request failed with status code: " + resp.Status)
		}
		return data, errors.New("GET request failed with status code: " + resp.Status)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func POSTRequest[T any](urlPost string, data *fiber.Map) (T, error) {
	var responseData T
	jsonData, err := json.Marshal(data)
	if err != nil {
		return responseData, err
	}

	req, err := http.NewRequest("POST", urlPost, bytes.NewBuffer(jsonData))
	if err != nil {
		return responseData, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return responseData, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseData, err
	}

	if resp.StatusCode != http.StatusOK {
		if body != nil {
			err = json.Unmarshal(body, &responseData)
			if err != nil {
				return responseData, err
			}
			return responseData, errors.New("GET request failed with status code: " + resp.Status)
		}
		return responseData, errors.New("GET request failed with status code: " + resp.Status)
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return responseData, err
	}

	return responseData, nil
}

func DownloadInscription(insciptID string) (string, error) {
	// Check if file exist in data_final folder
	if FileExists(fmt.Sprintf("%s/api/data_final/%s.webp", GetHomeDir(), insciptID[:8])) {
		return fmt.Sprintf("%s/api/data_final/%s.webp", GetHomeDir(), insciptID[:8]), nil
	}
	contentLink := fmt.Sprintf("https://ordinals.com/content/%s", insciptID)
	// Send HTTP GET request to the URL
	resp, err := http.Get(contentLink)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Extract the file format from the Content-Type header
	contentType := resp.Header.Get("Content-Type")
	dataType := strings.Split(contentType, "/")[0]
	fileFormat := strings.Split(contentType, "/")[1]
	if dataType != "image" {
		return "", errors.New("URL does not point to an image")
	}
	// Content Type: eg: image/png, decide if it is a picture or not
	format := strings.Split(contentType, "/")[1]

	// Create a new file with a unique name in the current directory
	filename := fmt.Sprintf("%s/api/data/%s.%s", GetHomeDir(), insciptID[:8], format)

	//check for file existing and return filename
	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Copy the image data from the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	fileToOpen := filename

	if fileFormat == "webp" {
		f0, err := os.Open(filename)
		if err != nil {
			WrapErrorLog(err.Error())
			return "", err
		}
		defer f0.Close()
		img0, err := webp.Decode(f0)
		if err != nil {
			WrapErrorLog(err.Error())
			return "", err
		}
		filepng := fmt.Sprintf("%s/api/data/%s.%s", GetHomeDir(), insciptID[:8], "png")
		pngFile, err := os.Create(filepng)
		if err != nil {
			fmt.Println(err)
		}
		err = png.Encode(pngFile, img0)
		if err != nil {
			fmt.Println(err)
		}
		defer func() {
			err := os.Remove(filepng)
			if err != nil {
				log.Println(err.Error())
			}
		}()
		fileToOpen = filepng
	}

	fileBytes, err := os.ReadFile(fileToOpen)
	if err != nil {
		return "", err
	}

	base64 := EncodePayload(fileBytes)

	tx := &grpcModels.NSFWRequest{
		Base64:   base64,
		Filename: fmt.Sprintf("pic.%s", fileFormat),
	}
	res, err := grpcClient.DetectNSFW(tx)
	if err != nil {
		WrapErrorLog(err.Error())
		return "", err
	}

	if res.NsfwPicture {
		//err = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", filename))).Run()
		//if err != nil {
		//    WrapErrorLog("Can't delete NSFW file in data")
		//}
		return "", ReturnError("NSFW image")
	}

	if res.NsfwText {
		//err = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", filename))).Run()
		//if err != nil {
		//    WrapErrorLog("Can't delete file in data")
		//}
		return "", ReturnError("NSFW Text in the image")
	}

	return filename, nil
}

func SendBackupToServer(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Prepare a buffer to hold the multipart request data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a new file field in the multipart request
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy the file content to the multipart request
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the multipart writer to finalize the request body
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Send the HTTP POST request with the multipart file upload
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:7100/backup", ServerUrl), &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Content-Type header to the correct value for a multipart request
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Length", strconv.Itoa(requestBody.Len()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status code %d", resp.StatusCode)
	}

	return nil
}

func SaveInscription(in models.WitnessData) (string, error) {
	errChain := make(chan error, 1)
	returnedFilename := make(chan string, 1)
	go func() {
		defer close(errChain)
		defer close(returnedFilename)

		contentType := in.FileType
		dataType := strings.Split(contentType, "/")[0]
		fileFormat := strings.Split(contentType, "/")[1]
		// Content Type: eg: image/png, decide if it is a picture or not
		format := strings.Split(contentType, "/")[1]

		// Create a new file with a unique name in the current directory
		filename := fmt.Sprintf("%s/api/data/%s.%s", GetHomeDir(), in.Txid[:8], format)
		// Create the directory if it doesn't exist
		dir := filepath.Dir(filename)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				errChain <- err
				return
			}
		}

		// Check if file exists
		if _, err := os.Stat(filename); err == nil {
			errChain <- err
			return
		}

		file, err := os.Create(filename)
		if err != nil {
			errChain <- err
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				//errChain <- err
				return
			}
		}(file)

		fl, err := DecodePayload([]byte(in.B64))
		if err != nil {
			errChain <- err
			return
		}

		_, err = io.Copy(file, bytes.NewReader(fl))
		if err != nil {
			errChain <- err
			return
		}

		if dataType != "image" {
			returnedFilename <- "txt"
			return
		}
		fileToOpen := filename

		if fileFormat == "webp" {
			f0, err := os.Open(filename)
			if err != nil {
				WrapErrorLog(err.Error())
				errChain <- err
				return
				//return "", err
			}
			defer func(f0 *os.File) {
				err := f0.Close()
				if err != nil {
					//errChain <- err
					return
				}
			}(f0)
			img0, err := webp.Decode(f0)
			if err != nil {
				WrapErrorLog(err.Error())
				errChain <- err
				return
				//return "", err
			}
			filepng := fmt.Sprintf("%s/api/data/%s.%s", GetHomeDir(), in.Txid[:8], "png")
			pngFile, err := os.Create(filepng)
			if err != nil {
				fmt.Println(err)
				errChain <- err
				return
			}
			err = png.Encode(pngFile, img0)
			if err != nil {
				fmt.Println(err)
				errChain <- err
				return
			}
			defer func() {
				err := os.Remove(filepng)
				if err != nil {
					log.Println(err.Error())
					//errChain <- err
					return
				}
			}()
			fileToOpen = filepng
		}

		fileBytes, err := os.ReadFile(fileToOpen)
		if err != nil {
			errChain <- err
			return
		}

		base64 := EncodePayload(fileBytes)

		tx := &grpcModels.NSFWRequest{
			Base64:   base64,
			Filename: fmt.Sprintf("pic.%s", fileFormat),
		}
		res, err := grpcClient.DetectNSFW(tx)
		if err != nil {
			WrapErrorLog(err.Error())
			errChain <- err
			return
			//return "", err
		}

		if res.NsfwPicture {
			//err = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", filename))).Run()
			//if err != nil {
			//    WrapErrorLog("Can't delete NSFW file in data")
			//}
			errChain <- ReturnError("NSFW image")
			return
			//return "", ReturnError("NSFW image")
		}

		if res.NsfwText {
			//err = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", filename))).Run()
			//if err != nil {
			//    WrapErrorLog("Can't delete file in data")
			//}
			errChain <- ReturnError("NSFW Text in the image")
			return
			//return "", ReturnError("NSFW Text in the image")
		}
	}()

	select {
	case err := <-errChain:
		return "", err
	case filename := <-returnedFilename:
		return filename, nil
	}
}
