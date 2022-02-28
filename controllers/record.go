package controllers

import (
	"encoding/json"
	"fmt"
	"go-thai-dialect/helper"
	"go-thai-dialect/models"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/labstack/echo"
)

func postRecord(c echo.Context) error {

	// fmt.Println("test")

	var snr models.Snr

	dialect_id := c.FormValue("dialect_id")
	dialect_code := c.FormValue("dialect_code")
	volunteer_id := c.FormValue("volunteer_id")
	sentence := c.FormValue("sentence")

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("upload/" + file.Filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println(err)
		return err
	}

	cmd := exec.Command("/usr/bin/python3", "/root/go-thai-dialect/snr-main/snr.py", "/root/go-thai-dialect/upload/"+file.Filename)
	// cmd := exec.Command("python3 print('test')")
	out, err := cmd.Output()

	if err != nil {
		println("err : ", err.Error())
	}

	fmt.Println(string(out))

	out_str := strings.ReplaceAll(string(out), `'`, `"`)

	if err := json.Unmarshal([]byte(out_str), &snr); err != nil {
		fmt.Println("snr err : ", err.Error())
	}

	err = models.InsertRecord(models.RecordDataset{
		RecordID:    helper.GenerateRandomString(32),
		DialectID:   helper.StringToInt(dialect_id[4:]),
		DialectType: dialect_id[:4],
		DialectCode: dialect_code,
		VolunteerID: volunteer_id,
		Sentence:    sentence,
		Url:         file.Filename,
		RecordTime:  helper.CurrentTime(),
		Vad:         snr.Vad.Value,
		Snr:         snr.Snr.Value,
		Energy:      snr.Energy.Value,
		Clipping:    snr.Clipping.Value,
		Pass:        (snr.Snr.Status == "OK" && snr.Vad.Status == "OK" && snr.Energy.Status == "OK"),
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"pass":                  (snr.Snr.Status == "OK" && snr.Vad.Status == "OK" && snr.Energy.Status == "OK"),
		"transcription":         "",
		"similar":               1,
		"transcription_is_pass": true,
		"vad":                   snr.Vad.Status == "OK",
		"snr":                   snr.Snr.Status == "OK",
		"energy":                snr.Energy.Status == "OK",
		"clipping":              true,
	})

}

func RecordDock(_echo *echo.Group) {
	_echo.POST("/record", postRecord)

	// records, err := models.GetRecordList("2021-12-01 00:00:00", "2021-12-22 23:59:59")

	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	for _, record := range records {
	// 		var snr models.Snr
	// 		fmt.Println(record.Url)
	// 		cmd := exec.Command("/usr/bin/python3", "/root/go-thai-dialect/snr-main/snr.py", "/root/go-thai-dialect/upload/"+record.Url)
	// 		// cmd := exec.Command("python3 print('test')")
	// 		out, err := cmd.Output()

	// 		if err != nil {
	// 			println("err : ", err.Error())
	// 		}

	// 		fmt.Println(string(out))

	// 		out_str := strings.ReplaceAll(string(out), `'`, `"`)

	// 		if err := json.Unmarshal([]byte(out_str), &snr); err != nil {
	// 			fmt.Println("snr err : ", err.Error())
	// 		} else {
	// 			err := models.UpdateSnr(record.Url, snr)
	// 			if err != nil {
	// 				fmt.Println(err)
	// 			}
	// 		}
	// 	}
	// }

	// config := helper.GetConfig()

	// for _, dialect := range config.DialectCode { // count and update to db
	// 	ecom_dialect, _ := models.GetRecordDialect(dialect, "ECOM")
	// 	println(len(ecom_dialect))
	// 	models.ClerDataset("composed_ecommerce_" + dialect)
	// 	for _, record := range ecom_dialect {
	// 		println("ecom : ", record.DialectID, record.Dataset)
	// 		models.UpdateDataset("composed_ecommerce_"+dialect, helper.IntToString(record.DialectID), record.Dataset)
	// 	}

	// 	surv_dialect, _ := models.GetRecordDialect(dialect, "SERV")
	// 	println(len(surv_dialect))

	// 	models.ClerDataset("composed_survival_" + dialect)
	// 	for _, record := range surv_dialect {
	// 		println("surv : ", record.DialectID, record.Dataset)
	// 		models.UpdateDataset("composed_survival_"+dialect, helper.IntToString(record.DialectID), record.Dataset)
	// 	}
	// }

}
