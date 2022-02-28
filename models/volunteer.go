package models

import (
	"database/sql"
	"fmt"
	"go-thai-dialect/helper"
	"strings"
	"time"
)

type Province struct {
	ProvinceID   string `json:"province_id"`
	ProvinceName string `json:"province_name"`
}

type Volunteer struct {
	VolunteerID     string `json:"volunteerid"`
	NickName        string `json:"nick_name" form:"nick_name"`
	Gender          string `json:"gender"`
	Age             int    `json:"age"`
	OfficialAbility bool   `json:"official_ability" form:"official_ability"`
	UserId          string `json:"user_id" form:"user_id"`
	Zipcode         string `json:"zipcode"`
	DialectID       string `json:"dialect_id" form:"dialect_id"`
	AgreeTerm       bool   `json:"agree_term" form:"agree_term"`
	ProvinceID      string `json:"province_id" form:"province_id"`
	DistrictID      string `json:"district_id" form:"district_id"`
	IsFacebookGame  bool   `json:"is_facebook_game" form:"is_facebook_game"`
}

type VolunteerDialect struct {
	VolunteerID string `json:"volunteer_id"`
	DialectID   string `json:"dialect_id"`
	Name        string `json:"name"`
	RegionID    string `json:"region_id"`
}

type District struct {
	DistrictID   string `json:"district_id"`
	DistrictName string `json:"district_name"`
}

func GetAllProvince() (province_list []Province) {

	province := Province{}
	rows, err := helper.Conn.Query("SELECT province_id, province_name FROM province")

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&province.ProvinceID, &province.ProvinceName)
		if err == nil {
			province_list = append(province_list, province)
		}
	}

	return province_list
}

func GetDistrictByProvinceID(provinceid string) (district_list []District) {
	district := District{}
	rows, err := helper.Conn.Query("SELECT district_id, district_name FROM district WHERE province_id = $1", provinceid)

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&district.DistrictID, &district.DistrictName)
		if err == nil {
			district_list = append(district_list, district)
		}
	}

	return district_list
}

func CheckExistVolunteer(volunteer Volunteer) bool {
	volunteerid := ""
	err := helper.Conn.QueryRow(`
		SELECT volunteerid 
		FROM public.volunteer 
		WHERE nickname = $1 
		AND gender = $2 
		AND age = $3 
		AND userid = $4`,
		volunteer.NickName,
		volunteer.Gender,
		volunteer.Age,
		volunteer.UserId,
	).Scan(&volunteerid)

	fmt.Println(err)

	if err != nil {
		if err == sql.ErrNoRows {
			return true
		} else {
			fmt.Println("CheckExistVolunteer err : ", err)
			fmt.Println(volunteer.NickName, volunteer.Gender, volunteer.Age, volunteer.UserId)
			return true
		}
	}
	return false
}

func CreateVolunteer(volunteer Volunteer) (bool, Volunteer) {
	fmt.Println("CreateVolunteer")
	volunteer.VolunteerID = helper.GenerateRandomString(32)
	err := helper.Conn.QueryRow(`
		INSERT INTO public.volunteer (
			volunteerid, 
			nickname, 
			gender, 
			age, 
			registerdate, 
			officialability, 
			userid, 
			zipcode, 
			agree_term, 
			province_id, 
			district_id, 
			is_facebook_game
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		volunteer.VolunteerID,
		volunteer.NickName,
		volunteer.Gender,
		volunteer.Age,
		time.Now(),
		volunteer.OfficialAbility,
		volunteer.UserId,
		volunteer.Zipcode,
		volunteer.AgreeTerm,
		volunteer.ProvinceID,
		volunteer.DistrictID,
		volunteer.IsFacebookGame,
	).Scan()

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("volunteer.DialectID : ", volunteer.DialectID)
			for _, dialect := range strings.Split(volunteer.DialectID, ",") {
				err := helper.Conn.QueryRow(`
					INSERT INTO public.volunteer_dialect (
						volunteer_id, 
						dialect_id
					) VALUES ($1, $2)`,
					volunteer.VolunteerID,
					dialect,
				).Scan()
				if err != nil {
					fmt.Println("Insert volunteer dialect err : ", err)
				}
			}

			return true, volunteer
		} else {
			fmt.Println("CreateVolunteer err : ", err)
		}
	}
	return false, volunteer
}

func GetVolunteerDialect(volunteerid string) (volunteer_dialect []VolunteerDialect) {
	// fmt.Println("GetVolunteerDialect")
	rows, err := helper.Conn.Query(`
		SELECT vd.dialect_id, d.name, d.region_id
		FROM public.volunteer_dialect vd
		INNER JOIN public.dialect d ON vd.dialect_id = d.dialect_id
		WHERE vd.volunteer_id = $1`,
		volunteerid,
	)

	if err != nil {
		fmt.Println("GetVolunteerDialect err : ", err)
	}

	for rows.Next() {
		var dialect VolunteerDialect
		err = rows.Scan(&dialect.DialectID, &dialect.Name, &dialect.RegionID)
		if err == nil {
			volunteer_dialect = append(volunteer_dialect, dialect)
		}
	}

	return volunteer_dialect
}
