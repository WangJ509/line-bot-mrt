package mrt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type StationTimeTable struct {
	RouteID     string `json:"RouteID"`
	LineID      string `json:"LineID"`
	StationID   string `json:"StationID"`
	StationName struct {
		ZhTw string `json:"Zh_tw"`
		En   string `json:"En"`
	} `json:"StationName"`
	Direction              int    `json:"Direction"`
	DestinationStaionID    string `json:"DestinationStaionID"`
	DestinationStationName struct {
		ZhTw string `json:"Zh_tw"`
		En   string `json:"En"`
	} `json:"DestinationStationName"`
	Timetables []struct {
		Sequence      int    `json:"Sequence"`
		ArrivalTime   string `json:"ArrivalTime"`
		DepartureTime string `json:"DepartureTime"`
	} `json:"Timetables"`
	ServiceDay struct {
		ServiceTag       string `json:"ServiceTag"`
		Monday           bool   `json:"Monday"`
		Tuesday          bool   `json:"Tuesday"`
		Wednesday        bool   `json:"Wednesday"`
		Thursday         bool   `json:"Thursday"`
		Friday           bool   `json:"Friday"`
		Saturday         bool   `json:"Saturday"`
		Sunday           bool   `json:"Sunday"`
		NationalHolidays bool   `json:"NationalHolidays"`
	} `json:"ServiceDay"`
	SrcUpdateTime time.Time `json:"SrcUpdateTime"`
	UpdateTime    time.Time `json:"UpdateTime"`
	VersionID     int       `json:"VersionID"`
}

var (
	baseURL             string
	stationTimeTableURL string
)

func init() {
	baseURL := os.Getenv("MRT_BASE_URL")
	if len(baseURL) == 0 {
		baseURL = "https://ptx.transportdata.tw/MOTC/v2/Rail/Metro"
	}

	stationTimeTableURL = baseURL + "/StationTimeTable/TRTC"
}

type MRTService struct {
	client http.Client
}

func NewMRTService() *MRTService {
	return &MRTService{client: http.Client{}}
}

func (s *MRTService) GetUpcomingTimeTable(station string, destination string, number int) ([]string, error) {
	req, err := getStationTimeTableRequest(station, destination)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("StationName/Zh_tw eq '%s' and DestinationStationName/Zh_tw eq '%s'", station, destination)
	query := req.URL.Query()
	query.Add("$filter", filter)
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("failed to do http request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("failed to read response body:", err)
		return nil, nil
	}

	var timeTables []StationTimeTable
	err = json.Unmarshal(body, &timeTables)
	if err != nil {
		log.Println("failed to decode response body:", err, "body:", string(body))
		return nil, err
	}

	for _, timeTable := range timeTables {
		if isValidTimeTable(timeTable) {
			return getUpcomingArrivalTime(timeTable, number), nil
		}
	}

	return nil, nil
}

func getStationTimeTableRequest(station string, destination string) (*http.Request, error) {
	req, err := http.NewRequest("GET", stationTimeTableURL, nil)
	if err != nil {
		log.Println("failed to get request:", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")

	return req, nil
}

func getNowTime() time.Time {
	location, _ := time.LoadLocation("Asia/Taipei")
	return time.Now().In(location).Add(time.Hour * 8)
}

func isValidTimeTable(timeTable StationTimeTable) bool {
	t := getNowTime()
	switch t.Weekday() {
	case time.Monday:
		return timeTable.ServiceDay.Monday
	case time.Tuesday:
		return timeTable.ServiceDay.Tuesday
	case time.Wednesday:
		return timeTable.ServiceDay.Wednesday
	case time.Thursday:
		return timeTable.ServiceDay.Thursday
	case time.Friday:
		return timeTable.ServiceDay.Friday
	case time.Saturday:
		return timeTable.ServiceDay.Saturday
	case time.Sunday:
		return timeTable.ServiceDay.Sunday
	}

	return false
}

func getUpcomingArrivalTime(timeTable StationTimeTable, number int) []string {
	nowHour, nowMinute, _ := getNowTime().Clock()
	if nowHour == 0 {
		nowHour = 24
	}

	offset := len(timeTable.Timetables)
	for i, table := range timeTable.Timetables {
		var hour, minute int
		fmt.Sscanf(table.ArrivalTime, "%d:%d", &hour, &minute)
		if hour == 0 {
			hour = 24
		}
		if 60*hour+minute >= 60*nowHour+nowMinute {
			offset = i
			break
		}
	}

	var result []string
	for i := 0; i < number; i++ {
		if i+offset >= len(timeTable.Timetables) {
			break
		}
		result = append(result, timeTable.Timetables[i+offset].ArrivalTime)
	}

	return result
}
