package helper

import (
	"encoding/json"
	"fmt"
	"go-thai-dialect/db"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/appleboy/go-fcm"

	"github.com/go-redis/redis/v8"
)

var (
	loc, _ = time.LoadLocation("Asia/Bangkok")
	// timeLog        = time.Now().In(loc)
	timeInc        = time.Hour * 7
	Conn           = db.ConnectDB()
	dateLayout     = "2006-01-02 15:04:05"
	address        = "010555206244101"
	jwtSecret      = "imsek"
	machineTest    = []string{"E0E2E672B9A0", "DC:4F:22:47:F8:98"}
	fcm_api_key    = "AAAAd-Ak7Oc:APA91bFptjWN9FqrJsQmdYeC0X9UznodJEUC3lE7MVfpbd497q-aEqT3uYHa1x2DnrJDha68RDaZJo2UKVHwDSuhwc39ZAd_kH8Fitpm99l7yW-XEmpBsL8vRLuuA_OX79igHw3yHOTi"
	config_dialect = `{
		"dialectSentence" : 72,
		"survivalSentence" : 48,
		"dialect_code" : ["official","laos","korat","khamen","krabi","pattani","phangnga","khummuang","nan","yno","songkhla"],
		"composed_survival" : ["composed_survival_official", "composed_survival_laos","composed_survival_korat","composed_survival_khamen","composed_survival_krabi","composed_survival_pattani","composed_survival_phangnga","composed_survival_khummuang","composed_survival_nan","composed_survival_yno","composed_survival_songkhla"],
		"composed_ecommerce" : ["composed_ecommerce_official","composed_ecommerce_laos","composed_ecommerce_korat","composed_ecommerce_khamen","composed_ecommerce_krabi","composed_ecommerce_pattani","composed_ecommerce_phangnga","composed_ecommerce_khummuang","composed_ecommerce_nan","composed_ecommerce_yno","composed_ecommerce_songkhla"]
	}`

	composedEcommerce [][]ComposedSentence
	composedSurvival  [][]ComposedSentence
	count             = 0

	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
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
	DialectID        string `json:"dialect_id"`
	LocalSentence    string `json:"local_sentence"`
	OfficialSentence string `json:"official_sentence"`
	Dataset          int    `json:"dataset"`
}

type ConfigDialect struct {
	DialectSentence   int      `json:"dialectSentence"`
	SurvivalSentence  int      `json:"survivalSentence"`
	DialectCode       []string `json:"dialect_code"`
	ComposedSurvival  []string `json:"composed_survival"`
	ComposedEcommerce []string `json:"composed_ecommerce"`
}

func IntToString(input int) string {
	return strconv.Itoa(input)
}

func GetCount() int {
	return count
}

func SetCount() {
	count++
}

func StringToInt(input string) int {
	i, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("err : ", err)
		return 0
	}
	return i
}

func RandomCon(is_odd bool, max_value int) (random_number int) {

	timeNano := time.Now().UnixNano()

	// fmt.Println("timeNano : ", timeNano)

	rand.Seed(timeNano)
	random_number = rand.Intn(max_value)
	if (random_number%2 == 0) == is_odd {
		if random_number == max_value {
			random_number = random_number - 1
		}
	} else {
		if random_number == max_value {
			random_number = random_number - 1
		}
	}
	// fmt.Println("random_number : ", random_number)
	return random_number
}

func GetConfig() ConfigDialect {
	var config ConfigDialect
	err := json.Unmarshal([]byte(config_dialect), &config)
	fmt.Println(err)
	return config
}

func CheckNullArrString(input []string, index int) string {
	if index < len(input) {
		return input[index]
	} else {
		return ""
	}
}

func Notification(token string, title string, body string) {
	// Create the message to be sent.
	msg := &fcm.Message{
		To: token,
		Data: map[string]interface{}{
			"foo": "bar",
		},
		Notification: &fcm.Notification{
			Title: title,
			Body:  body,
		},
	}

	// Create a FCM client to send the message.
	client, err := fcm.NewClient(fcm_api_key)
	if err != nil {
		fmt.Println(err)
	}

	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("response : ", response)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func stringToInt(input string) int {
	i, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("err : ", err)
		return 0
	}
	return i
}

func stringToDatetime(input string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}
	return t
}

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// func negative(x int) int {
// 	if x < 0 {
// 		return x
// 	}
// 	return -x
// }

func crc(str string) (scrc string) {

	crc := uint16(0xffff)
	crc_table := [256]uint16{
		0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
		0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef,
		0x1231, 0x0210, 0x3273, 0x2252, 0x52b5, 0x4294, 0x72f7, 0x62d6,
		0x9339, 0x8318, 0xb37b, 0xa35a, 0xd3bd, 0xc39c, 0xf3ff, 0xe3de,
		0x2462, 0x3443, 0x0420, 0x1401, 0x64e6, 0x74c7, 0x44a4, 0x5485,
		0xa56a, 0xb54b, 0x8528, 0x9509, 0xe5ee, 0xf5cf, 0xc5ac, 0xd58d,
		0x3653, 0x2672, 0x1611, 0x0630, 0x76d7, 0x66f6, 0x5695, 0x46b4,
		0xb75b, 0xa77a, 0x9719, 0x8738, 0xf7df, 0xe7fe, 0xd79d, 0xc7bc,
		0x48c4, 0x58e5, 0x6886, 0x78a7, 0x0840, 0x1861, 0x2802, 0x3823,
		0xc9cc, 0xd9ed, 0xe98e, 0xf9af, 0x8948, 0x9969, 0xa90a, 0xb92b,
		0x5af5, 0x4ad4, 0x7ab7, 0x6a96, 0x1a71, 0x0a50, 0x3a33, 0x2a12,
		0xdbfd, 0xcbdc, 0xfbbf, 0xeb9e, 0x9b79, 0x8b58, 0xbb3b, 0xab1a,
		0x6ca6, 0x7c87, 0x4ce4, 0x5cc5, 0x2c22, 0x3c03, 0x0c60, 0x1c41,
		0xedae, 0xfd8f, 0xcdec, 0xddcd, 0xad2a, 0xbd0b, 0x8d68, 0x9d49,
		0x7e97, 0x6eb6, 0x5ed5, 0x4ef4, 0x3e13, 0x2e32, 0x1e51, 0x0e70,
		0xff9f, 0xefbe, 0xdfdd, 0xcffc, 0xbf1b, 0xaf3a, 0x9f59, 0x8f78,
		0x9188, 0x81a9, 0xb1ca, 0xa1eb, 0xd10c, 0xc12d, 0xf14e, 0xe16f,
		0x1080, 0x00a1, 0x30c2, 0x20e3, 0x5004, 0x4025, 0x7046, 0x6067,
		0x83b9, 0x9398, 0xa3fb, 0xb3da, 0xc33d, 0xd31c, 0xe37f, 0xf35e,
		0x02b1, 0x1290, 0x22f3, 0x32d2, 0x4235, 0x5214, 0x6277, 0x7256,
		0xb5ea, 0xa5cb, 0x95a8, 0x8589, 0xf56e, 0xe54f, 0xd52c, 0xc50d,
		0x34e2, 0x24c3, 0x14a0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
		0xa7db, 0xb7fa, 0x8799, 0x97b8, 0xe75f, 0xf77e, 0xc71d, 0xd73c,
		0x26d3, 0x36f2, 0x0691, 0x16b0, 0x6657, 0x7676, 0x4615, 0x5634,
		0xd94c, 0xc96d, 0xf90e, 0xe92f, 0x99c8, 0x89e9, 0xb98a, 0xa9ab,
		0x5844, 0x4865, 0x7806, 0x6827, 0x18c0, 0x08e1, 0x3882, 0x28a3,
		0xcb7d, 0xdb5c, 0xeb3f, 0xfb1e, 0x8bf9, 0x9bd8, 0xabbb, 0xbb9a,
		0x4a75, 0x5a54, 0x6a37, 0x7a16, 0x0af1, 0x1ad0, 0x2ab3, 0x3a92,
		0xfd2e, 0xed0f, 0xdd6c, 0xcd4d, 0xbdaa, 0xad8b, 0x9de8, 0x8dc9,
		0x7c26, 0x6c07, 0x5c64, 0x4c45, 0x3ca2, 0x2c83, 0x1ce0, 0x0cc1,
		0xef1f, 0xff3e, 0xcf5d, 0xdf7c, 0xaf9b, 0xbfba, 0x8fd9, 0x9ff8,
		0x6e17, 0x7e36, 0x4e55, 0x5e74, 0x2e93, 0x3eb2, 0x0ed1, 0x1ef0}

	var bs = []byte(str)
	for i := 0; i < len(str); i++ {
		crc = (crc_table[(uint16(bs[i])^(crc>>8))&0xff] ^ (crc << 8)) & 0xffff
	}

	scrc = strings.ToUpper(strconv.FormatInt(int64(crc), 16))
	return
}

func genIntZero(input int) string {
	if input <= 9 {
		return "0" + strconv.Itoa(input)
	}
	return strconv.Itoa(input)
}

func PadZero(input int, number int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(number)+"d", input)
}

func StringPadZero(input string, number int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(number)+"s", input)
}

func GenQrcode(amount string, ref2 string, ref3 string) string { //"00020101021129"
	merchantIdent := "0016A00000067701011201" + genIntZero(len(address)) + address //"0016A0000006770101110113" + address
	merchantIdent += "02" + genIntZero(len(ref2)) + ref2
	merchantIdent += "03" + genIntZero(len(ref3)) + ref3

	qrcodeString := "00020101021230"
	qrcodeString += genIntZero(len(merchantIdent)) + merchantIdent // Merchant Identity
	qrcodeString += "53" + "03" + "764"                            // THB
	qrcodeString += "54" + genIntZero(len(amount)) + amount        // Amount
	qrcodeString += "58" + "02" + "TH"                             // Country
	qrcodeString += "62" + "10" + "0706TDW001"
	qrcodeString += "63" + "04" // CRC Checking

	scrc := crc(qrcodeString)
	for i := len(scrc); i < 4; i++ {
		scrc = "0" + scrc
	}
	fmt.Println(scrc)
	qrcodeString += scrc
	return qrcodeString
}

func weekday(day int) string {
	switch day {
	case 0:
		return "sunday"
	case 1:
		return "monday"
	case 2:
		return "tuesday"
	case 3:
		return "wednesday"
	case 4:
		return "thursday"
	case 5:
		return "friday"
	case 6:
		return "saturday"
	}
	return "sunday"
}

func intToHour(n int) string {
	if n < 9 {
		return "0" + strconv.Itoa(n)
	} else {
		return strconv.Itoa(n)
	}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

//GenerateRandomString : ..
func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// const letters = "123456789"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return ""
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func paymentEnum(input string) string {
	switch input {
	case "qr_test":
		return "TW Pay"
	case "tw_pay":
		return "TW Pay"
	case "linepay":
		return "Line Pay"
	case "truemoney":
		return "Truemoney Wallet"
	case "omise":
		return "Credit Card"
	case "alipay":
		return "Alipay"
	case "scb":
		return "Thai QR-Code"
	case "admin":
		return "Admin"
	case "promotion":
		return "Promotion"
	default:
		return input
	}
}

func machineStatusEnum(input string) string {
	switch input {
	case "0":
		return "Online"
	case "1":
		return "Working"
	case "-1":
		return "Delete"
	default:
		return "Loss Connection"
	}
}

func isNan(input float64) float64 {
	if math.IsNaN(input) || input == math.Inf(-1) || input == math.Inf(1) {
		return 0.0
	}
	return float64(int(input*100)) / 100
}

func checkIntString(input string) string {
	if input == "" {
		return "0"
	}
	return input
}
