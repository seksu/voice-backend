package models

import (
	"database/sql"
	"fmt"
	"go-thai-dialect/db"
	"sync"
)

var (
	Conn = db.ConnectDB()
)

type DialectOption []struct {
	Name        string    `json:"name"`
	DialectList []Dialect `json:"dialectList"`
}

type Dialect struct {
	ID      string `json:"id"`
	Dialect string `json:"dialect"`
	Code    string `json:"code"`
}

type DialectTemplate struct {
	DialectTemplateID string `json:"dialect_template_id"`
	Count             int    `json:"count"`
}

type ComposedSentence struct {
	DialectID        string `json:"dialectid"`
	LocalSentence    string `json:"localsentence"`
	OfficialSentence string `json:"official_sentence"`
	Dataset          int    `json:"dataset"`
	Parity           bool   `json:"parity"`
}

type DialectComposedSentence struct {
	ComposedSentence []ComposedSentence `json:"composed_sentence"`
	Mu               sync.Mutex
}

func DialectTemplateList(tb_name string) []DialectTemplate {

	var dialectTemplateList []DialectTemplate

	rows, err := Conn.Query(`
		SELECT 
			dialect_template_id, 
			COUNT(*) AS count 
		FROM ` + tb_name + ` 
		WHERE active = true 
		GROUP BY dialect_template_id`)

	if err != nil {
		fmt.Println("GetVolunteerDialect err : ", err)
	}

	for rows.Next() {
		var dialectTemplate DialectTemplate
		err = rows.Scan(
			&dialectTemplate.DialectTemplateID,
			&dialectTemplate.Count,
		)
		if err == nil {
			dialectTemplateList = append(dialectTemplateList, dialectTemplate)
		} else {
			fmt.Println("err : ", err)
		}
	}

	return dialectTemplateList

}

func GetSize(tb_name string) (size int) {

	err := Conn.QueryRow(`
		SELECT COUNT(*) AS count 
		FROM ` + tb_name + ` 
		WHERE active = true`).Scan(&size)

	if err != nil {
		fmt.Println("GetVolunteerDialect err : ", err)
	}

	return size

}

func ComposedSentenceList(tb_name string, sentence_type string, dialect_template_id string, offset int) []ComposedSentence {

	var composedSentenceList []ComposedSentence

	rows, err := Conn.Query(`
		SELECT 
			concat('`+sentence_type+`',composed_sentence_id) as dialectid, 
			composed_sentence as localsentence,
			composed_sentence_official AS official_sentence 
		FROM	`+tb_name+` 
		WHERE active = true
		AND dialect_template_id = $1
		ORDER BY dataset
		OFFSET $2
		LIMIT 1`,
		dialect_template_id,
		offset,
	)

	if err != nil {
		fmt.Println("GetVolunteerDialect err : ", err)
	} else {
		for rows.Next() {
			var composedSentence ComposedSentence
			err = rows.Scan(
				&composedSentence.DialectID,
				&composedSentence.LocalSentence,
				&composedSentence.OfficialSentence,
			)
			if err == nil {
				UpdateDatasetInc(tb_name, composedSentence.DialectID[4:])
				composedSentenceList = append(composedSentenceList, composedSentence)
			} else {
				fmt.Println("err : ", err)
			}
		}
	}
	return composedSentenceList
}

func UpdateDatasetInc(tb_name string, dialect_id string) {
	err := Conn.QueryRow(`
		UPDATE `+tb_name+`
		SET dataset = dataset + 1
		WHERE composed_sentence_id = $1`,
		dialect_id,
	).Scan()
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("updateDataset err : ", err)
		}
	}
}

func UpdateDataset(tb_name string, dialect_id string, dataset int) {
	err := Conn.QueryRow(`
		UPDATE `+tb_name+`
		SET dataset = $1
		WHERE composed_sentence_id = $2`,
		dataset,
		dialect_id,
	).Scan()
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("updateDataset err : ", err)
		}
	}
}

func ClerDataset(tb_name string) {
	_, err := Conn.Exec(`
		UPDATE ` + tb_name + `
		SET dataset = 0
		WHERE active = true`)
	if err != nil {
		fmt.Println("updateDataset err : ", err)
	}
}

func ComposedSentenceListO(tb_name string, sentence_type string, dialect_template_id string, offset int) []ComposedSentence {

	var composedSentenceList []ComposedSentence

	rows, err := Conn.Query(`
		SELECT 
			concat('`+sentence_type+`',composed_sentence_id) as dialectid, 
			composed_sentence as localsentence,
			composed_sentence_official AS official_sentence 
		FROM	`+tb_name+` 
		WHERE active = true
		ORDER BY dataset
		OFFSET $1
		LIMIT 1`,
		offset,
	)

	if err != nil {
		fmt.Println("GetVolunteerDialect err : ", err)
	} else {
		for rows.Next() {
			var composedSentence ComposedSentence
			err = rows.Scan(
				&composedSentence.DialectID,
				&composedSentence.LocalSentence,
				&composedSentence.OfficialSentence,
			)
			if err == nil {
				composedSentenceList = append(composedSentenceList, composedSentence)
			} else {
				fmt.Println("err : ", err)
			}
		}
	}
	return composedSentenceList
}

func ComposedSentenceListAll(tb_name string, sentence_type string) []ComposedSentence {

	var composedSentenceList []ComposedSentence

	rows, err := Conn.Query(`
		SELECT 
			concat('` + sentence_type + `',composed_sentence_id) as dialectid, 
			composed_sentence as localsentence,
			composed_sentence_official AS official_sentence,
			dataset
		FROM	` + tb_name + ` 
		WHERE active = true
		ORDER BY dataset
		LIMIT 1000000`)

	if err != nil {
		fmt.Println("GetVolunteerDialect err : ", err)
	} else {
		for rows.Next() {
			var composedSentence ComposedSentence
			err = rows.Scan(
				&composedSentence.DialectID,
				&composedSentence.LocalSentence,
				&composedSentence.OfficialSentence,
				&composedSentence.Dataset,
			)
			if err == nil {
				composedSentenceList = append(composedSentenceList, composedSentence)
			} else {
				fmt.Println("err : ", err)
			}
		}
	}
	// fmt.Println("finish fetching official sentence")
	return composedSentenceList
}

func GetSizeZero(tb_name string) (size float64) {
	err := Conn.QueryRow(`
		SELECT count(*) 
		FROM ` + tb_name + ` 
		WHERE dataset = 0`).Scan(&size)
	if err != nil {
		fmt.Println("GetSizeZero err : ", err)
	}
	return size
}

func GetSizeTotal(tb_name string) (size float64) {
	err := Conn.QueryRow(`
		SELECT count(*) 
		FROM ` + tb_name).Scan(&size)
	if err != nil {
		fmt.Println("GetSizeTotal err : ", err)
	}
	return size
}
