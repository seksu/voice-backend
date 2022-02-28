package models

import "go-thai-dialect/helper"

type User struct {
	UserID       string `json:"userid"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Email        string `json:"email"`
	OauthToken   string `json:"oauthtoken"`
	PhoneNumber  string `json:"phonenumber"`
	RegisterDate string `json:"registerdate"`
	AgentModel   string `json:"agentmodel"`
	RegionID     string `json:"regionid"`
	Terms        bool   `json:"terms"`
}

func InsertUser(user User) (user_id string, err error) {

	user_id = helper.GenerateRandomString(32)
	user.UserID = user_id

	stmt, err := Conn.Prepare("INSERT INTO public.user (userid, firstname, lastname, email, oauthtoken, phonenumber, registerdate, agentmodel, regionid, terms ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
	if err != nil {
		return user_id, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.UserID, user.FirstName, user.LastName, user.Email, user.OauthToken, user.PhoneNumber, user.RegisterDate, user.AgentModel, user.RegionID, user.Terms)
	if err != nil {
		return user_id, err
	}

	return user_id, nil
}
