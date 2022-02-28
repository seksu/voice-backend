package models

import (
	"database/sql"
	"go-thai-dialect/helper"
)

type RecordRes struct {
	Pass                bool   `json:"pass"`
	Transcription       string `json:"transcription"`
	Similar             int    `json:"similar"`
	TranscriptionIsPass bool   `json:"transcription_is_pass"`
	Vad                 bool   `json:"vad"`
	Snr                 bool   `json:"snr"`
	Energy              bool   `json:"energy"`
	Clipping            bool   `json:"clipping"`
}

type RecordDataset struct {
	RecordID     string  `json:"recordid"`
	Url          string  `json:"url"`
	RecordTime   string  `json:"recordtime"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	VolunteerID  string  `json:"volunteerid"`
	DialectID    int     `json:"dialectid"`
	NoiseRatio   string  `json:"noiseratio"`
	DialectCode  string  `json:"dialect_code"`
	DialectType  string  `json:"dialect_type"`
	Inactive     bool    `json:"inactive"`
	Sentence     string  `json:"sentence"`
	Vad          float64 `json:"vad"`
	Snr          float64 `json:"snr"`
	Energy       float64 `json:"energy"`
	Clipping     float64 `json:"clipping"`
	Pass         bool    `json:"pass"`
	Warning      string  `json:"warning"`
	Src          string  `json:"src"`
	NoSegment    bool    `json:"no_segment"`
	SoxFailed    bool    `json:"sox_failed"`
	NoSound      bool    `json:"no_sound"`
	Transcript   string  `json:"transcript"`
	SimilarScore float64 `json:"similar_score"`
	Dataset      int     `json:"dataset"`
}

type Snr struct {
	ID     string  `json:"id"`
	Length float64 `json:"length"`
	Vad    struct {
		Value  float64 `json:"value"`
		Status string  `json:"status"`
	} `json:"VAD"`
	Snr struct {
		Value  float64 `json:"value"`
		Status string  `json:"status"`
	} `json:"SNR"`
	Energy struct {
		Value  float64 `json:"value"`
		Status string  `json:"status"`
	} `json:"energy"`
	Clipping struct {
		Value  float64 `json:"value"`
		Status string  `json:"status"`
	} `json:"clipping"`
}

func InsertRecord(record RecordDataset) error {
	err := Conn.QueryRow(`
		INSERT INTO public.record_dataset (
			recordid,
			url,
			recordtime,
			latitude,
			longitude,
			volunteerid,
			dialectid,
			noiseratio,
			dialect_code,
			dialect_type,
			inactive,
			sentence,
			vad,
			snr,
			energy,
			clipping,
			pass,
			warning,
			src,
			no_segment,
			sox_failed,
			no_sound,
			transcript,
			similar_score
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)`,
		record.RecordID,
		record.Url,
		record.RecordTime,
		record.Latitude,
		record.Longitude,
		record.VolunteerID,
		record.DialectID,
		record.NoiseRatio,
		record.DialectCode,
		record.DialectType,
		record.Inactive,
		record.Sentence,
		record.Vad,
		record.Snr,
		record.Energy,
		record.Clipping,
		record.Pass,
		record.Warning,
		record.Src,
		record.NoSegment,
		record.SoxFailed,
		record.NoSound,
		record.Transcript,
		record.SimilarScore,
	).Scan()

	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	return nil
}

func UpdateSnr(url string, snr Snr) error {
	_, err := Conn.Exec(`
		UPDATE public.record_dataset
		SET
			recordid = $7,
			vad = $2,
			snr = $3,
			energy = $4,
			clipping = $5,
			pass = $6
		WHERE url = $1`,
		url,
		snr.Vad.Value,
		snr.Snr.Value,
		snr.Energy.Value,
		snr.Clipping.Value,
		(snr.Snr.Status == "OK" && snr.Vad.Status == "OK" && snr.Energy.Status == "OK"),
		helper.GenerateRandomString(32),
	)

	if err != nil {
		return err
	}

	return nil
}

func GetRecordList(start_date string, end_date string) (records []RecordDataset, err error) {
	rows, err := Conn.Query(`
		SELECT
			recordid,
			url
		FROM public.record_dataset
		WHERE recordtime between $1 and $2
		AND recordid = ''`,
		start_date,
		end_date,
	)

	if err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var record RecordDataset
			err = rows.Scan(
				&record.RecordID,
				&record.Url,
			)

			if err != nil {
				return nil, err
			}

			records = append(records, record)
		}
	}

	if err != nil {
		if err != sql.ErrNoRows {
			return records, err
		}
	}
	return records, nil
}

func GetRecordDialect(dialect_code string, dialect_type string) (records []RecordDataset, err error) {
	rows, err := Conn.Query(`
		SELECT
			dialectid,
			count(*)
		FROM public.record_dataset
		WHERE dialect_code = $1
		AND dialect_type = $2
		GROUP BY dialectid`,
		dialect_code,
		dialect_type,
	)

	if err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var record RecordDataset
			err = rows.Scan(
				&record.DialectID,
				&record.Dataset,
			)

			if err != nil {
				return nil, err
			}

			records = append(records, record)
		}
	}

	if err != nil {
		if err != sql.ErrNoRows {
			return records, err
		}
	}
	return records, nil
}
