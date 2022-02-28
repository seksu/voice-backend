package controllers

import (
	"context"
	"fmt"
	"go-thai-dialect/helper"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"go-thai-dialect/models"

	"github.com/labstack/echo"
)

var (
	composedEcommerce [][]models.ComposedSentence
	composedSurvival  [][]models.ComposedSentence

	dialectComposedEcommerce []models.DialectComposedSentence
	dialectComposedSurvival  []models.DialectComposedSentence

	count = 0
	ctx   = context.Background()

	composedEcommerceP *[][]models.ComposedSentence
	composedSurvivalP  *[][]models.ComposedSentence

	index     = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	indexEcom = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	mu sync.Mutex

	cacheSurv = []string{}
)

func shuffle(s []models.ComposedSentence) []models.ComposedSentence {
	rand.Seed(time.Now().UnixNano())
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func getDialectOption(c echo.Context) error {

	return c.JSON(http.StatusOK, models.DialectOption{
		{
			Name: "อีสาน",
			DialectList: []models.Dialect{
				{
					ID:      "1",
					Code:    "laos",
					Dialect: "ลาว",
				},
				{
					ID:      "2",
					Code:    "korat",
					Dialect: "โคราช",
				},
				{
					ID:      "3",
					Code:    "khamen",
					Dialect: "เขมร",
				},
			},
		},
		{
			Name: "ใต้",
			DialectList: []models.Dialect{
				{
					ID:      "4",
					Code:    "krabi",
					Dialect: "กระบี่",
				},
				{
					ID:      "5",
					Code:    "pattani",
					Dialect: "ปัตตานี",
				},
				{
					ID:      "6",
					Code:    "phangnga",
					Dialect: "พังงา",
				},
				{
					ID:      "10",
					Code:    "songkhla",
					Dialect: "สงขลา",
				},
			},
		},
		{
			Name: "เหนือ",
			DialectList: []models.Dialect{
				{
					ID:      "7",
					Code:    "khummuang",
					Dialect: "คำเมือง",
				},
				{
					ID:      "8",
					Code:    "nan",
					Dialect: "น่าน",
				},
				{
					ID:      "9",
					Code:    "yno",
					Dialect: "ยอง",
				},
			},
		},
		{
			Name: "กลาง",
			DialectList: []models.Dialect{
				{
					ID:      "0",
					Code:    "official",
					Dialect: "กลาง",
				},
			},
		},
	})
}

func getDialect(c echo.Context) error {

	is_odd := (time.Now().Unix() % 2) == 0

	config := helper.GetConfig()

	dialect_list := models.GetVolunteerDialect(c.QueryParam(`volunteer_id`))

	responseData := echo.Map{}

	for _, dialect := range dialect_list {

		responseList := []models.ComposedSentence{}

		id, err := strconv.Atoi(dialect.DialectID)

		if err == nil {
			list_template_ecom := models.DialectTemplateList(config.ComposedEcommerce[id])

			// size_ecom_composed := models.GetSize(config.ComposedEcommerce[id])

			fmt.Println(config.DialectSentence)
			ecomList := []models.ComposedSentence{}
			for i := 0; i < config.DialectSentence; i++ {
				fmt.Println("1 : ", i)
				temp_index := helper.RandomCon(is_odd, len(list_template_ecom))
				dialect_template_ecom_id := list_template_ecom[temp_index].DialectTemplateID
				dialect_template_ecom_count := list_template_ecom[temp_index].Count

				list_template_ecom = append(list_template_ecom[:temp_index], list_template_ecom[temp_index+1:]...)

				rand := helper.RandomCon(is_odd, dialect_template_ecom_count)
				// rand := helper.RandomCon(is_odd, size_ecom_composed)
				dialect_template_ecom_list := models.ComposedSentenceList(config.ComposedEcommerce[id], "ECOM", dialect_template_ecom_id, rand)

				// if len(dialect_template_ecom_list) == 0 {
				// 	fmt.Println(i, config.ComposedEcommerce[id], dialect_template_ecom_id, rand)
				// }

				ecomList = append(ecomList, dialect_template_ecom_list...)
			}
			responseList = append(responseList, ecomList...)
			fmt.Println("ecomList : ", len(ecomList))

			// list_template_serv := models.DialectTemplateList(config.ComposedSurvival[id])
			fmt.Println(config.SurvivalSentence)
			servList := []models.ComposedSentence{}
			for i := 0; i < config.SurvivalSentence; i++ {
				fmt.Println("2 : ", i)
				// temp_index := helper.RandomCon(is_odd, len(list_template_serv))
				// dialect_template_serv_id := list_template_serv[temp_index].DialectTemplateID
				// dialect_template_serv_count := list_template_serv[temp_index].Count
				dialect_template_serv_id := "0"
				// list_template_serv = append(list_template_serv[:temp_index], list_template_serv[temp_index+1:]...)

				size_surv_composed := models.GetSize(config.ComposedSurvival[id])

				dialect_template_serv_list := models.ComposedSentenceList(config.ComposedSurvival[id], "SERV", dialect_template_serv_id, helper.RandomCon(is_odd, size_surv_composed))
				servList = append(servList, dialect_template_serv_list...)
			}
			responseList = append(responseList, servList...)
			fmt.Println("servList : ", len(servList))

			responseList = shuffle(responseList)
			responseData[dialect.Name] = responseList

		}

	}

	return c.JSON(http.StatusOK, responseData)
}

func getDialectNew(c echo.Context) error {

	// fmt.Println("getDialectNew")

	// is_odd := (time.Now().Unix() % 2) == 0

	config := helper.GetConfig()

	dialect_list := models.GetVolunteerDialect(c.QueryParam(`volunteer_id`))

	responseData := echo.Map{}

	for _, dialect := range dialect_list {

		responseList := []models.ComposedSentence{}

		id, err := strconv.Atoi(dialect.DialectID)

		// fmt.Println("len(composedEcommerce[id]) : ", len(composedEcommerce[id]))
		// fmt.Println("len(composedSurvival[id] : ", len(composedSurvival[id]))

		if err == nil {
			ecomList := []models.ComposedSentence{}
			for i := 0; i < config.DialectSentence; i++ {
				mu.Lock()
				fmt.Println("1 : ", i)
				// temp_index := helper.RandomCon(is_odd, len(composedEcommerce[id]))

				// dialect_template_ecom_list := []models.ComposedSentence{composedEcommerce[id][temp_index]}

				ecomList = append(ecomList, composedEcommerce[id][indexEcom[id]])

				indexEcom[id] = indexEcom[id] + 1

				if indexEcom[id]%(len(composedEcommerce[id])-1) == 0 {
					indexEcom[id] = 0
				}

				mu.Unlock()
			}
			responseList = append(responseList, ecomList...)

			servList := []models.ComposedSentence{}

			composedSurvival[id] = shuffle(composedSurvival[id])

			for i := 0; i < config.SurvivalSentence; i++ {
				mu.Lock()
				fmt.Println("2 : ", i)
				// temp_index := helper.RandomCon(is_odd, len(composedSurvival[id]))

				// sort.SliceStable(composedSurvival[id], func(i, j int) bool {
				// 	return composedSurvival[id][i].Dataset < composedSurvival[id][j].Dataset
				// })

				// // fmt.Println("composedSurvival : ", composedSurvival[id][0].Dataset)

				// // dialect_template_serv_list := []models.ComposedSentence{composedSurvival[id][0]}
				// fmt.Println("Data : ", composedSurvival[id][0])
				// composedSurvival[id][0].Dataset = composedSurvival[id][0].Dataset + 1

				servList = append(servList, composedSurvival[id][index[id]])

				// if index[id]%len(composedSurvival[id]) == 0 {
				// 	index[id] = 0
				// } else {
				// 	index[id] = index[id] + 1
				// }

				index[id] = index[id] + 1

				if index[id]%(len(composedSurvival[id])-1) == 0 {
					index[id] = 0
				}

				mu.Unlock()
			}

			responseList = append(responseList, servList...)
			// fmt.Println("servList : ", len(servList))

			responseList = shuffle(responseList)
			responseData[dialect.Name] = responseList

		}

	}

	return c.JSON(http.StatusOK, responseData)
}

func getDialectNew2(c echo.Context) error {

	fmt.Println("getDialectNew2")

	config := helper.GetConfig()

	dialect_list := models.GetVolunteerDialect(c.QueryParam(`volunteer_id`))

	responseData := echo.Map{}

	for _, dialect := range dialect_list {

		responseList := []models.ComposedSentence{}

		id, err := strconv.Atoi(dialect.DialectID)

		if err == nil {
			ecomList := []models.ComposedSentence{}
			servList := []models.ComposedSentence{}
			// mu.Lock()

			for i := 0; i < config.DialectSentence; i++ {

				dialectComposedEcommerce[id].Mu.Lock()

				// fmt.Println("dialectComposedEcommerce : ", dialectComposedEcommerce[id].ComposedSentence)

				ecomList = append(ecomList, dialectComposedEcommerce[id].ComposedSentence[indexEcom[id]])

				dialectComposedEcommerce[id].ComposedSentence[indexEcom[id]].Dataset = dialectComposedEcommerce[id].ComposedSentence[indexEcom[id]].Dataset + 1
				dialectComposedEcommerce[id].ComposedSentence[indexEcom[id]].Parity = true

				indexEcom[id] = indexEcom[id] + 1

				if indexEcom[id]%(len(dialectComposedEcommerce[id].ComposedSentence)-1) == 0 {
					indexEcom[id] = 0
				}

				dialectComposedEcommerce[id].Mu.Unlock()

			}
			// mu.Unlock()
			responseList = append(responseList, ecomList...)

			// mu.Lock()
			for i := 0; i < config.SurvivalSentence; i++ {

				dialectComposedSurvival[id].Mu.Lock()

				fmt.Println("len : ", len(dialectComposedSurvival), len(dialectComposedSurvival[id].ComposedSentence), id)

				servList = append(servList, dialectComposedSurvival[id].ComposedSentence[index[id]])

				dialectComposedSurvival[id].ComposedSentence[index[id]].Dataset = dialectComposedSurvival[id].ComposedSentence[index[id]].Dataset + 1
				dialectComposedSurvival[id].ComposedSentence[index[id]].Parity = true

				index[id] = index[id] + 1

				if index[id]%(len(dialectComposedSurvival[id].ComposedSentence)-1) == 0 {
					index[id] = 0
				}

				dialectComposedSurvival[id].Mu.Unlock()

			}

			// mu.Unlock()

			responseList = append(responseList, servList...)
			// fmt.Println("servList : ", len(servList))

			responseList = shuffle(responseList)
			responseData[dialect.Name] = responseList

		}

	}

	return c.JSON(http.StatusOK, responseData)
}

func test(c echo.Context) error {

	fmt.Println(composedEcommerce[0])
	composedEcommerce[0][0].Dataset = composedEcommerce[0][0].Dataset + 1

	sort.SliceStable(composedEcommerce[0], func(i, j int) bool {
		return composedEcommerce[0][i].Dataset < composedEcommerce[0][j].Dataset
	})
	return c.JSON(http.StatusOK, "test")

}

func countAllDataset() int {
	var countDataset int
	for dialect_id, dialect := range dialectComposedEcommerce {
		for sentence_id, _ := range dialect.ComposedSentence {
			if dialectComposedEcommerce[dialect_id].ComposedSentence[sentence_id].Parity {
				countDataset++
			}
		}
	}

	for dialect_id, dialect := range dialectComposedSurvival {
		for sentence_id, _ := range dialect.ComposedSentence {
			if dialectComposedSurvival[dialect_id].ComposedSentence[sentence_id].Parity {
				countDataset++
			}
		}
	}

	return countDataset
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func DialectDock(_echo *echo.Group) {

	config := helper.GetConfig()

	go func() {
		for i, dialect := range config.ComposedEcommerce {
			fmt.Println(dialect, i)

			// composedEcommerce = append(composedEcommerce, shuffle(models.ComposedSentenceListAll(dialect, "ECOM")))

			dialectComposedEcommerce = append(dialectComposedEcommerce, models.DialectComposedSentence{ComposedSentence: shuffle(models.ComposedSentenceListAll(dialect, "ECOM"))})
		}
		fmt.Println("end Ecom")
	}()

	go func() {
		for i, dialect := range config.ComposedSurvival {
			fmt.Println(dialect, i)
			// composedSurvival = append(composedSurvival, shuffle(models.ComposedSentenceListAll(dialect, "SERV")))
			dialectComposedSurvival = append(dialectComposedSurvival, models.DialectComposedSentence{ComposedSentence: shuffle(models.ComposedSentenceListAll(dialect, "SERV"))})
		}
		fmt.Println("end Serv")
	}()

	go func() {
		for range time.Tick(time.Second * 10) {
			fmt.Print("Syncing ", countAllDataset())
			for dialect_id, _ := range config.ComposedEcommerce {
				if len(dialectComposedEcommerce) != 11 {
					break
				}
				for id, _ := range dialectComposedEcommerce[dialect_id].ComposedSentence {
					if dialectComposedEcommerce[dialect_id].ComposedSentence[id].Parity {
						models.UpdateDataset(config.ComposedEcommerce[dialect_id], dialectComposedEcommerce[dialect_id].ComposedSentence[id].DialectID[4:], dialectComposedEcommerce[dialect_id].ComposedSentence[id].Dataset)
						dialectComposedEcommerce[dialect_id].ComposedSentence[id].Parity = false
					}
				}
			}
			for dialect_id, _ := range config.ComposedSurvival {
				if len(dialectComposedSurvival) != 11 {
					break
				}
				for id, _ := range dialectComposedSurvival[dialect_id].ComposedSentence {
					if dialectComposedSurvival[dialect_id].ComposedSentence[id].Parity {
						models.UpdateDataset(config.ComposedSurvival[dialect_id], dialectComposedSurvival[dialect_id].ComposedSentence[id].DialectID[4:], dialectComposedSurvival[dialect_id].ComposedSentence[id].Dataset)
						dialectComposedSurvival[dialect_id].ComposedSentence[id].Parity = false
					}
				}
			}
			fmt.Println("Complete")
		}
	}()

	_echo.GET("/dialect/option", getDialectOption)
	_echo.GET("/dialect", getDialectNew2)

}
