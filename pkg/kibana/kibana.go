package kibana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TimeInterval struct {
	Gte   time.Time
	Lte   time.Time
	NoMod bool
}

type TimeBucket struct {
	From  time.Time
	To    time.Time
	Count int
}

type KibanaFetchRequest struct {
	Url            string
	AcceptLanguage string
	Authorization  string
	ContentType    string
	KbnVersion     string
	Body           string
}

type KibanaResponseHeader struct {
	ID     int `json:"id"`
	Result struct {
		ID          string `json:"id"`
		RawResponse struct {
			Took     int  `json:"took"`
			TimedOut bool `json:"timed_out"`
			Shards   struct {
				Total      int `json:"total"`
				Successful int `json:"successful"`
				Skipped    int `json:"skipped"`
				Failed     int `json:"failed"`
			} `json:"_shards"`
			Hits struct {
				Total int `json:"total"`
			} `json:"hits"`
			Aggregations struct {
				Num2 struct {
					Buckets []struct {
						KeyAsString string `json:"key_as_string"`
						Key         int64  `json:"key"`
						DocCount    int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"2"`
			} `json:"aggregations"`
		} `json:"rawResponse"`
		IsPartial  bool   `json:"isPartial"`
		IsRunning  bool   `json:"isRunning"`
		Warning    string `json:"warning"`
		Total      int    `json:"total"`
		Loaded     int    `json:"loaded"`
		IsRestored bool   `json:"isRestored"`
	} `json:"result"`
}

type KibanaResult struct {
	ID     int `json:"id"`
	Result struct {
		ID          string `json:"id"`
		RawResponse struct {
			Took     int  `json:"took"`
			TimedOut bool `json:"timed_out"`
			Shards   struct {
				Total      int `json:"total"`
				Successful int `json:"successful"`
				Skipped    int `json:"skipped"`
				Failed     int `json:"failed"`
			} `json:"_shards"`
			Hits struct {
				Hits []struct {
					Index   string         `json:"_index"`
					Type    string         `json:"_type"`
					ID      string         `json:"_id"`
					Version int            `json:"_version"`
					Fields  map[string]any `json:"fields"`
					Sort    []int64        `json:"sort"`
				} `json:"hits"`
			} `json:"hits"`
		} `json:"rawResponse"`
		IsPartial  bool   `json:"isPartial"`
		IsRunning  bool   `json:"isRunning"`
		Warning    string `json:"warning"`
		Total      int    `json:"total"`
		Loaded     int    `json:"loaded"`
		IsRestored bool   `json:"isRestored"`
	} `json:"result"`
}

type KibanaResponse struct {
	KibanaResponseHeader KibanaResponseHeader
	KibanaResult         KibanaResult
}

const kibanaMaxResult = 500

func POST(reguest KibanaFetchRequest) (error, KibanaResponse) {
	err, response := internalPOST(reguest)

	if 0 == len(response.KibanaResult.Result.RawResponse.Hits.Hits) {
		println("sleep and retry")
		time.Sleep(15 * time.Second) // dont stress the server
		err, response := internalPOST(reguest)
		return err, response
	}
	return err, response
}

func internalPOST(request KibanaFetchRequest) (error, KibanaResponse) {
	client := &http.Client{}

	size := RawParse(`\"size\":\d+`, request.Body)
	newSize := `"size":` + strconv.Itoa(kibanaMaxResult)
	request.Body = strings.Replace(request.Body, size, newSize, 1)

	req, err := http.NewRequest("POST", request.Url, bytes.NewBuffer([]byte(request.Body)))
	if nil != err {
		return err, KibanaResponse{}
	}

	req.Header.Add("accept-language", request.AcceptLanguage)
	req.Header.Add("authorization", request.Authorization)
	req.Header.Add("content-type", request.ContentType)
	req.Header.Add("kbn-version", request.KbnVersion)

	resp, err := client.Do(req)
	if nil != err {
		return err, KibanaResponse{}
	}

	body, err := io.ReadAll(resp.Body)
	if nil != err {
		return err, KibanaResponse{}
	}

	response := strings.Split(string(body[:]), "\n")
	header := response[0]
	result := response[1]
	var kibanaResponseHeaderLocal KibanaResponseHeader
	err = json.Unmarshal([]byte(header), &kibanaResponseHeaderLocal)
	if err != nil {
		return err, KibanaResponse{}
	}

	var kibanaResult KibanaResult
	err = json.Unmarshal([]byte(result), &kibanaResult)
	if err != nil {
		return err, KibanaResponse{}
	}

	defer resp.Body.Close()

	kibanaResponse := KibanaResponse{
		KibanaResponseHeader: kibanaResponseHeaderLocal,
		KibanaResult:         kibanaResult}

	return nil, kibanaResponse
}

func LoadFromFetchRequest(f string) (KibanaFetchRequest, error) {
	fetchFile, err := file.Open(f)
	if nil != err {
		return KibanaFetchRequest{}, err
	}

	fetch := string(fetchFile)

	bodyPrefix := `"body": "`
	bodySuffix := `}",`
	bodyLine := RawParse(bodyPrefix+`{\\".*?`+bodySuffix, fetch)
	bodyWithSuffix := strings.Replace(bodyLine, bodyPrefix, "", 1)
	rawBody := strings.Replace(bodyWithSuffix, bodySuffix, "", 1) + "}"
	bodyWithoutEscaping := strings.Replace(rawBody, `\"`, `"`, -1)
	bodyWithoutEscapingErr := strings.Replace(bodyWithoutEscaping, `\\`, `\`, -1)

	kibanaRequest := KibanaFetchRequest{
		Url:            strings.Replace(ParseForValueInQuote(`fetch\(".*?\"`, 0, fetch), "compress=true", "compress=false", 1),
		Authorization:  ParseForValueInQuote(`\"authorization\": \".*?\"`, 1, fetch),
		AcceptLanguage: ParseForValueInQuote(`\"accept-language\": \".*?\"`, 1, fetch),
		ContentType:    ParseForValueInQuote(`\"content-type\": \".*?\"`, 1, fetch),
		KbnVersion:     ParseForValueInQuote(`\"kbn-version\": \".*?\"`, 1, fetch),
		Body:           bodyWithoutEscapingErr,
	}
	return kibanaRequest, nil
}

func ParseForValueInQuote(subStringRegExp string, postition int, fetch string) string {
	rawValue := RawParse(subStringRegExp, fetch)

	quotedPartRegExp := regexp.MustCompile(`"[^"]+"`)
	quotedPart := quotedPartRegExp.FindAllString(rawValue, postition+1)
	value := quotedPart[postition]
	return strings.Replace(value, "\"", "", 2)
}

func RawParse(subStringRegExp string, text string) string {
	regExp := regexp.MustCompile(subStringRegExp)
	rawValue := regExp.Find([]byte(text))
	return string(rawValue)
}

// FilterAndConvertToJson filter map[remappedFieldName,resultFieldName]
func FilterAndConvertToJson(
	kibanaResult KibanaResult,
	filter map[string]string,
	resultExits map[string]bool) (error, []string, map[string]bool, TimeInterval) {

	var result []string
	for _, hit := range kibanaResult.Result.RawResponse.Hits.Hits {

		filteredHit := make(map[string]interface{})
		allFieldsFound := true
		for key, field := range filter {
			value := hit.Fields[field]
			if nil == value {
				println("field missing: " + field + ", skip result")
				allFieldsFound = false
				break
			}
			filteredHit[key] = value.([]interface{})[0]
		}
		if allFieldsFound {
			hitAsJsonBytes, err := json.Marshal(filteredHit)
			if err != nil {
				return err, result, resultExits, TimeInterval{}
			}

			hitAsJson := string(hitAsJsonBytes)
			if !resultExits[hitAsJson] {
				result = append(result, hitAsJson)
				resultExits[hitAsJson] = true
			}
		}
	}

	totalHits := len(kibanaResult.Result.RawResponse.Hits.Hits)
	if 0 < totalHits {
		firstHit := kibanaResult.Result.RawResponse.Hits.Hits[0]
		firstHitTimestamp := fmt.Sprintf("%v", firstHit.Fields["@timestamp"].([]interface{})[0])

		lastHit := kibanaResult.Result.RawResponse.Hits.Hits[totalHits-1]
		lastHitTimestamp := fmt.Sprintf("%v", lastHit.Fields["@timestamp"].([]interface{})[0])

		gte, _ := time.Parse(time.RFC3339Nano, lastHitTimestamp)
		lte, _ := time.Parse(time.RFC3339Nano, firstHitTimestamp)
		return nil, result, resultExits, TimeInterval{Gte: gte, Lte: lte}
	}

	return nil, result, resultExits, TimeInterval{}
}

func ExecuteKibanaQuery(
	kibanaRequest KibanaFetchRequest,
	requestTimeInterval TimeInterval,
	filter map[string]string,
	resultExits map[string]bool,
	outputPadding string) (error, KibanaResponse, []string, map[string]bool) {

	request := CreateRequestForInterval(requestTimeInterval, kibanaRequest)

	err, kibanaResponse := POST(request)
	if err != nil {
		return err, KibanaResponse{}, nil, resultExits
	}

	err, result, updatedResultExits, responseTimeInterval := FilterAndConvertToJson(kibanaResponse.KibanaResult, filter, resultExits)
	if err != nil {
		return err, KibanaResponse{}, nil, resultExits
	}

	fmt.Println(outputPadding, "-------------------------------------")
	fmt.Println(outputPadding, "requestTimeInterval : ", requestTimeInterval.Gte.Format(time.RFC3339Nano), requestTimeInterval.Lte.Format(time.RFC3339Nano))
	fmt.Println(outputPadding, "responseTimeInterval: ", responseTimeInterval.Gte.Format(time.RFC3339Nano), requestTimeInterval.Lte.Format(time.RFC3339Nano))
	fmt.Println(outputPadding, "Url:", kibanaRequest.Url)
	fmt.Println(outputPadding, "Result.took:", kibanaResponse.KibanaResult.Result.RawResponse.Took)
	fmt.Println(outputPadding, "Hits:", kibanaResponse.KibanaResponseHeader.Result.RawResponse.Hits.Total)
	fmt.Println(outputPadding, "Hits delivered:", len(kibanaResponse.KibanaResult.Result.RawResponse.Hits.Hits))
	fmt.Println(outputPadding, "Unique hits:", len(result))

	if kibanaMaxResult == len(kibanaResponse.KibanaResult.Result.RawResponse.Hits.Hits) {
		newRequestTimeInterval := TimeInterval{
			Gte: requestTimeInterval.Gte,
			Lte: responseTimeInterval.Gte,
		}

		requestTimeInterval.Lte = responseTimeInterval.Gte

		err, kibanaResponse, subQueryResult, _ := ExecuteKibanaQuery(kibanaRequest, newRequestTimeInterval, filter, updatedResultExits, outputPadding+"  ")
		if err != nil {
			return err, KibanaResponse{}, nil, resultExits
		}

		aggregatedResult := append(result, subQueryResult...)
		return nil, kibanaResponse, aggregatedResult, resultExits
	}

	return nil, kibanaResponse, result, resultExits
}

func ExtractTimeIntervalFrom(kibanaRequest KibanaFetchRequest) (TimeInterval, error) {
	gte := ParseForValueInQuote(`\"gte\":\".*?\"`, 1, kibanaRequest.Body)
	lte := ParseForValueInQuote(`\"lte\":\".*?\"`, 1, kibanaRequest.Body)

	gteParsed, err := time.Parse(time.RFC3339Nano, gte)
	if err != nil {
		return TimeInterval{}, err
	}
	lteParsed, err := time.Parse(time.RFC3339Nano, lte)
	if err != nil {
		return TimeInterval{}, err
	}

	timeInterval := TimeInterval{
		Gte: gteParsed,
		Lte: lteParsed,
	}
	return timeInterval, nil
}

func CreateRequestForInterval(timeInterval TimeInterval, kibanaRequest KibanaFetchRequest) KibanaFetchRequest {
	gte := ParseForValueInQuote(`\"gte\":\".*?\"`, 1, kibanaRequest.Body)
	lte := ParseForValueInQuote(`\"lte\":\".*?\"`, 1, kibanaRequest.Body)

	newGte := timeInterval.Gte.Format(time.RFC3339)
	newLte := timeInterval.Lte.Format(time.RFC3339)

	kibanaRequest.Body = strings.Replace(kibanaRequest.Body, gte, newGte, 1)
	kibanaRequest.Body = strings.Replace(kibanaRequest.Body, lte, newLte, 1)
	return kibanaRequest
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func CreateFilter(extractFieldsInput string, fieldRemapInput string) map[string]string {
	extractFields := strings.Split(extractFieldsInput, ",")
	fieldRemap := strings.Split(fieldRemapInput, ",")

	filter := make(map[string]string)
	for i, field := range extractFields {
		if len(extractFields) == len(fieldRemap) {
			filter[fieldRemap[i]] = field
		} else {
			filter[field] = field
		}
	}
	return filter
}
