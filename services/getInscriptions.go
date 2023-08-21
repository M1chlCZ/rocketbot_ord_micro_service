package services

import (
	"api/cmd"
	"api/db"
	"api/models"
	"api/utils"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"strconv"
	"strings"
	"sync"
)

var s sync.Mutex

func GetInscriptions() {
	s.Lock()
	defer s.Unlock()
	//dm := GetDaemon()
	utils.ReportMessage(" = = Get Inscriptions = = ")
	callString, err := cmd.CallArrayJSON[models.Inscriptions]("bash", "-c", "/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300 --wallet ord wallet inscriptions")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	for _, ins := range callString {
		txid := strings.Split(ins.Inscription, "i")
		currentTX := strings.Split(ins.Location, ":")[0]
		voutArr := strings.Split(ins.Location, ":")
		vout, err := strconv.Atoi(strings.Split(ins.Location, ":")[len(voutArr)-1])
		if err != nil {
			vout = 0
		}

		utils.ReportMessage("--------------------")
		utils.ReportMessage(fmt.Sprintf("ins: %s location: %s", txid, currentTX))
		info, err := utils.GETRequest[models.RawTX](fmt.Sprintf("https://blockstream.info/api/tx/%s", currentTX))
		if err != nil {
			utils.WrapErrorLog(err.Error())
			continue
		}
		utils.ReportMessage(fmt.Sprintf("Address: %s", info.Vout[vout].ScriptpubkeyAddress))
		contentLink := fmt.Sprintf("https://ordinals.com/content/%s", ins.Inscription)

		r, err := getWitnessData(txid[0], vout)
		if err != nil {
			utils.WrapErrorLog(err.Error())
			continue
		}

		_, errNsfw := utils.SaveInscription(*r)

		if errNsfw != nil {
			utils.WrapErrorLog(errNsfw.Error())
			_, _ = db.InsertSQl(`INSERT INTO NSFW_ORD (tx_id, file_format, ord_id, bc_address, link, content_link) 
									VALUES (?, ?, ?, ?, ?, ?)`, currentTX, r.FileType, ins.Inscription, info.Vout[vout].ScriptpubkeyAddress, ins.Explorer, contentLink)
			continue
		} else {
			_, err = db.InsertSQl(`INSERT INTO TRANSACTIONS_ORD (tx_id, file_format, ord_id, bc_address, link, content_link) 
									VALUES (?,?, ?, ?, ?, ?)`, currentTX, r.FileType, ins.Inscription, info.Vout[vout].ScriptpubkeyAddress, ins.Explorer, contentLink)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE constraint failed") {
					continue
				}
				utils.WrapErrorLog(err.Error())
				continue
			}
		}

	}
	go ScanAndConvert()
	utils.ReportMessage("Inscriptions saved into db")
}

func getWitnessData(txID string, vout int) (*models.WitnessData, error) {
	res, err := utils.GETRequest[models.RawTX](fmt.Sprintf("https://blockstream.info/api/tx/%s", txID))
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return nil, err
	}

	witnessBytes, err := hex.DecodeString(res.Vin[vout].Witness[1])
	if err != nil {
		return nil, err
	}

	// Parse the script
	disasm, err := txscript.DisasmString(witnessBytes)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return nil, err
	}

	// Split into operations
	ops := strings.Split(disasm, " ")
	b64 := ""
	fileType := ""
	var err2 error
	// Print the operations
	for i := 0; i < len(ops); i++ {
		op := ops[i]
		if op == "OP_IF" {
			// The next four items should be 'ord', version, MIME type, and the data
			if i+4 < len(ops) {
				// Decode the MIME type
				mimeTypeBytes, err := hex.DecodeString(ops[i+3])
				if err != nil {
					fmt.Printf("Failed to decode MIME type for operation %s: %v\n", ops[i+3], err)
					continue
				}

				mimeType := string(mimeTypeBytes)
				fmt.Printf("MIME type: %s\n", mimeType)

				// Decode the data
				var dataBytes []byte
				for j := i + 5; j < len(ops); j++ {
					op2 := ops[j]
					if op2 == "OP_ENDIF" {
						break
					}

					if op2 == "OP_PUSHDATA1" || op2 == "OP_PUSHDATA2" || op2 == "OP_PUSHDATA4" {
						j++
					}

					dataByte, err := hex.DecodeString(op2)
					if err != nil {
						fmt.Printf("Failed to decode data byte for operation %s: %v\n", op2, err)
						continue
					}

					dataBytes = append(dataBytes, dataByte...)
				}
				utils.ReportMessage(fmt.Sprintf("MimeType: %s", mimeType))
				// Handle the data differently depending on the MIME type
				switch mimeType {
				case "text/plain;charset=utf-8":
					fileType = mimeType
					b64 = base64.StdEncoding.EncodeToString(dataBytes)
					break
				case "text/html;charset=utf-8":
					fileType = mimeType
					b64 = base64.StdEncoding.EncodeToString(dataBytes)
					break
				case "image/png":
					fileType = mimeType
					b64 = base64.StdEncoding.EncodeToString(dataBytes)
					break
				case "image/jpeg", "image/jpg":
					fileType = mimeType
					b64 = base64.StdEncoding.EncodeToString(dataBytes)
					break
				case "image/webp":
					fileType = mimeType
					b64 = base64.StdEncoding.EncodeToString(dataBytes)
					break
				case "image/gif":
					fileType = mimeType
					b64 = base64.StdEncoding.EncodeToString(dataBytes)
					break
				default:
					err2 = errors.New("unknown file type")
					fmt.Printf("Data: Unknown MIME type, length %d\n", len(dataBytes))
				}

				// Skip the next four items, as we've already processed them
				i += 4
			}
		}
	}
	if err2 != nil {
		return nil, err2
	}
	fl := strings.Split(fileType, "/")
	if strings.Contains("plain;charset=utf-8", fl[1]) {
		fileType = "plain/txt"
	}

	r := &models.WitnessData{
		Txid:     txID,
		B64:      b64,
		FileType: fileType,
	}
	return r, nil
}
