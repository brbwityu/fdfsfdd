package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cookie = ""
var mainfakecookie = ""
var secondfakecookie = "_|WARNING:-DO-NOT-SHARE-THIS.--Sharing-this-will-allow-someone-to-log-in-as-you-and-to-steal-your-ROBUX-and-items.|_73A1D0B0279905F605B521EEDAD73564BD0D57EF68457BDD41FE7AA21959E76061C61ED927C4004CB26CFFC2F79C73604CEEC0893AEBB503F7ACF4BD1DA140595B7DB16A048F6CF5AC1BF6900D3A76D4925E29EA40476FE5F58B993036785DDAFCC252D226602313235509D717EDC8F0060943A38D88EC53C38A876D1AAD510DD0610912D780113E48755F0992DECEFCAACB7CE132B092F5F7B77C93845EB5B835CBE41F8EF6E61F6C73715DC27D02137AEDE4129EA5C700567513535368FC31980B7525F99E5BA7FB3F71FA23306F42B2F7611464A4C7F87F94C7D460CE042A8B9F05FFD69B7D138F05C6E0B468775B642A011478983A04E5F30549ACB33EB45A45D3D98C3B5A90D16CFC5C50E2B8562DAED35928C6CF292F958A927ED2172443E4B035B59443B453C963183AF09F8EDFCB30C5CCCA45CDCA0076C2CF48F8D97F65D769616D2200D7287DA8BDEC9CB9863A6A2E903937CB4B634DB39352592AE9E7C5AFC784AC3C4BBFDD4A539188914313D13720B41ECC5C06EB1DBAF8044B871DA7DFD89CCC0183C17327D30E16301EFAF8DD20D4CE717CA6DDDD579291F4122AF9D98CC3B0441D24BAC80514318D780373321D7B961AAFFD127698C19594B9288CB8841D91AE71BFBC639A0190FCB5CD92BF75402CED7DF5024AFEC21E5F055DC3F211B8F63D128A46BDF4C1988E63D1A233D3DB3E7B2FAD9F4A29EE68AF48347190004B5945CBA335DC20FCC640F6527757675FDDF5DC0B4C37724BDC0F192607436638F7BD723F50B2E58736FCFB4A4DD08B497147AB1286ED36B21028A6525266F522D768AF7983DE2537067C2CE1C0708545D040A2F8748AEF6D263F8D0BF33860A6E4DE872C80CEDE36DAB7B56EAD65F8B8D155B6F3430804AB4F542D5D55CA468BB55A8F4B0913B2FF4A9FB59780DA5F879F3F53BCD40838E994C8E1966914CF67399E99763E2551F5DC7A43B2ADCF5D241C00D7C3DCD5C24CA0041188CC6C6E7BABC1C019432083B2AEF40E5EAEF2"

var vouchhook = "https://discord.com/api/webhooks/1238639691297853521/zDukH6Vu5GeF53ZDgrOx78EEVGlBbvfb9tdKL31uUsMTUk2RqviwdG4xFo_lAXwM_eLQ"

var snipehook string = ""
var discordid string = ""

var licenseKey string = ""

var workercount int = 5

var profitmargin float64 = 0

var filterproxies = false
var currentproxy = ""
var currentchecks = 0
var filters = 0
var currentlychanging = false

var proxyless = "false"

var debug = "false"
var debughook = "https://discord.com/api/webhooks/1238675797607321766/06-BYDWyB3Nm_8a86Urauvhlzh-vuOFqn0cHA5vmbSfPRaIUyRJfAh0O0TNIJn11WY21"

var unusedproxies []string

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getRandomProxy() (string, error) {
	if len(unusedproxies) == 0 {
		return "", errors.New("No proxies left in the list")
	}

	returnLine := unusedproxies[rand.Intn(len(unusedproxies))]

	return returnLine, nil
}

func siftProxies() {
	fmt.Println("Sifting proxies...")

	file, err := os.Open("proxies.txt")
	if err != nil {
		log.Println("Couldn't open proxies.txt")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)

	// Read lines from the file and append into the lines slice
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Error handling for scanner error
	if scanErr := scanner.Err(); scanErr != nil {
		log.Printf("Error scanning file: %v", scanErr)
	}

	// If there are no lines, return an error
	if len(lines) == 0 {
		fmt.Println("NO PROXIES FOUND!")
	}

	// Generate a random number and fetch the proxy
	rand.Seed(time.Now().Unix())

	unusedproxies = []string{}

	for _, line := range lines {
		proxyParts := strings.Split(line, ":")
		if len(proxyParts) == 4 {
			ip := proxyParts[0]
			port := proxyParts[1]
			username := proxyParts[2]
			password := proxyParts[3]
			rl := fmt.Sprintf("%s:%s:%s:%s", username, password, ip, port)

			unusedproxies = append(unusedproxies, rl)
		} else {
			unusedproxies = append(unusedproxies, line)
		}
	}

	fmt.Println(unusedproxies)
}

func filterRequests() {
	filterproxies = true
	currentlychanging = true
	filters = filters + 1

	if len(unusedproxies) <= 1 {
		siftProxies()
	}

	if filters >= 6 {
		filterproxies = false
		filters = 0
	}
	fmt.Println("Rotating proxy...")
	currentproxy2, err := getRandomProxy()
	if err != nil {
		siftProxies()
		currentproxy2, err = getRandomProxy()
		currentproxy = currentproxy2
	} else {
		currentproxy = currentproxy2
	}

	currentlychanging = false
}

var rap = make(map[int64]interface{})

func getCSRF(icookie string) string {
	// Create an HTTP client
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	csrfURL := "https://auth.roblox.com/v2/login"

	// Create a new POST request to the login endpoint
	req, err := http.NewRequest("POST", csrfURL, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	cookie := &http.Cookie{Name: ".ROBLOSECURITY", Value: icookie}

	// Add the cookie to the request headers
	req.AddCookie(cookie)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	// Find the X-CSRF-TOKEN in the response headers
	csrfToken := resp.Header.Get("X-CSRF-TOKEN")
	if csrfToken == "" {
		fmt.Println("X-CSRF-TOKEN not found in the header")
		return ""
	}

	return csrfToken
}

func loadConfig() {
	file, err := os.Open("config.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)

		// Check if parts has at least 2 elements
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "cookie":
			cookie = value
		case "mainfakecookie":
			mainfakecookie = value
		case "proxyless":
			proxyless = value
		case "snipehook":
			snipehook = value
		case "Key":
			licenseKey = value
		case "discordid":
			discordid = value
		case "debug":
			debug = value
		case "profitmargin":
			id, _ := strconv.ParseFloat(value, 64)
			profitmargin = id
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func getRandomCookie() string {
	file, err := os.Open("falsecookies.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	randomLine := lines[r.Intn(len(lines))]

	return randomLine
}

func sendWebhook(itemid string, price string, statuscode int, reason string) {
	webhookURL := snipehook
	content := `{
		"content": "<@` + discordid + `> Successfully sniped item https://www.roblox.com/catalog/` + itemid + ` at the price of: ` + price + ` robux!"
	}`

	if statuscode != 200 {
		content = `{
		"content": "Failed to snipe item https://www.roblox.com/catalog/` + itemid + ` at the price of: ` + price + ` robux! Reason: ` + reason + `"
	}`
	} else {
		req, _ := http.NewRequest("POST", vouchhook, strings.NewReader(content))
		req.Header.Set("Content-Type", "application/json")

		http.DefaultClient.Do(req)
	}

	req, _ := http.NewRequest("POST", webhookURL, strings.NewReader(content))
	req.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(req)
}

func sendDebugHook(itemid string, price string) {
	webhookURL := debughook
	content := `{
		"content": "<` + discordid + `> has opportunity for: https://www.roblox.com/catalog/` + itemid + ` at the price of: ` + price + ` robux!"
	}`

	req, _ := http.NewRequest("POST", webhookURL, strings.NewReader(content))
	req.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(req)
}

func getRecentAveragePrices() {
	rap = make(map[int64]interface{})
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Println(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "www.rolimons.com",
	}, []*http.Cookie{})

	req, err := http.NewRequest("GET", "https://www.rolimons.com/itemapi/itemdetails", nil)
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	type Response struct {
		Items map[string][]interface{} `json:"items"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var res Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println(err)
		return
	}

	for id, item := range res.Items {
		idStr := id // assert that the key is a string

		if len(item) < 8 {
			continue
		}

		_, ok := item[7].(float64)
		idInt64, err := strconv.ParseInt(idStr, 10, 64) // convert string id to int64

		if ok && err == nil && item[7].(float64) == -1 {
			var f = item[2]
			var newrap = int64(f.(float64))
			rap[idInt64] = newrap
		} else {
			fmt.Println("PROJECTED!")
		}
	}
}

func checkId(id int64) (string, error) {
	var rv string = ""
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}

	var idstring = strconv.FormatInt(id, 10)

	url, _ := url.Parse("https://www.roblox.com/catalog/" + idstring)

	cookies := []*http.Cookie{
		&http.Cookie{Name: ".ROBLOSECURITY", Value: mainfakecookie},
	}
	jar.SetCookies(url, cookies)

	client := &http.Client{
		Jar: jar,
	}

	firstRequest, err := http.NewRequest("GET", url.String(), strings.NewReader("1"))

	if err != nil {
		return "", err
	}

	response, err := client.Do(firstRequest)

	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	attributes := make(map[string]interface{})

	var productId int64 = 0

	doc.Find("div.content>div.page-content").Each(func(i int, s *goquery.Selection) {
		price, exists := s.Attr("data-expected-price")

		if exists {
			if reflect.TypeOf(price).Kind() == reflect.String {
				attributes["expectedPrice"] = price
				attributes["expectedCurrency"] = "1"
			}
		}

		seller, exists := s.Attr("data-expected-seller-id")

		if exists {
			if reflect.TypeOf(seller).Kind() == reflect.String {
				attributes["expectedSellerId"] = seller
			}
		}

		productid, exists := s.Attr("data-product-id")

		if exists {
			newid, err := strconv.ParseInt(productid, 10, 64)
			checkErr(err)

			productId = newid
		}

		lowestassetid, exists := s.Attr("data-lowest-private-sale-userasset-id")

		if exists {
			if reflect.TypeOf(lowestassetid).Kind() == reflect.String {
				attributes["userAssetId"] = lowestassetid
			}
		}

		// add more attributes as needed
	})

	//convert the map to a JSON
	jsonPayload, err := json.Marshal(attributes)
	if err != nil {
		log.Fatal(err)
	}

	//Test for snipe level
	if attributes["expectedPrice"] != nil {
		priceAttr, ok := attributes["expectedPrice"].(string)

		if !ok {
			log.Fatal("expectedPrice attribute is not a string!")
		}

		newval, _ := strconv.ParseInt(priceAttr, 10, 64)

		valueInterface, exist := rap[id]
		value, ok := valueInterface.(int64)
		value = value

		if !exist {
			fmt.Println("ID does not exist in map!")
		}

		fmt.Println(newval)

		if float64(newval) != float64(0) {
			if float64(newval) <= (profitmargin * float64(rap[id].(int64))) {
				url, _ := url.Parse("https://economy.roblox.com/v1/purchases/products/" + strconv.FormatInt(productId, 10))

				cookies := []*http.Cookie{
					&http.Cookie{Name: ".ROBLOSECURITY", Value: cookie},
					&http.Cookie{Name: "x-csrf-token", Value: getCSRF(cookie)},
					&http.Cookie{Name: "_gcl_au", Value: "1.1.435265357.1710645925"},
					&http.Cookie{Name: "GuestData", Value: "UserID=-999927455"},
					&http.Cookie{Name: "authority", Value: "auth.roblox.com"},
					&http.Cookie{Name: "content-type", Value: "application/json"},
					&http.Cookie{Name: "accept", Value: "application/json, text/plain, */*"},
				}
				jar.SetCookies(url, cookies)

				client := http.Client{
					Jar: jar,
				}

				secondRequest, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonPayload))
				checkErr(err)

				secondRequest.Header.Set("authority", "auth.roblox.com")
				secondRequest.Header.Set("accept", "application/json, text/plain, */*")
				secondRequest.Header.Set("content-type", "application/json")
				secondRequest.Header.Set("_gcl_au", "1.1.435265357.1710645925")
				secondRequest.Header.Set("GuestData", "UserID=-999927455")
				secondRequest.Header.Set("x-csrf-token", getCSRF(cookie))

				req, err := client.Do(secondRequest)
				checkErr(err)

				bodyBytes, err := ioutil.ReadAll(req.Body)
				if err != nil {
					log.Fatal(err)
				}

				// Always close the response body
				defer req.Body.Close()

				bodyString := string(bodyBytes)
				fmt.Println(bodyString)

				var dat map[string]interface{}

				// Decoding/Unmarshalling the JSON string
				if err := json.Unmarshal([]byte(bodyString), &dat); err != nil {
					log.Fatal(err)
				}
				var statuscode int = 400
				// Getting the value of "purchased" from the map
				purchased, ok := dat["purchased"]
				if !ok {
					log.Println("Key 'purchased' not found in JSON")
				} else {
					if purchased.(bool) == true {
						statuscode = 200
					}
				}

				var sendreason string = ""

				reason, ok := dat["reason"].(string)
				if !ok {
					log.Println("reason attribute is not a string!")
				} else {
					sendreason = reason
				}

				sendWebhook(idstring, priceAttr, statuscode, sendreason)
			}
		}
	}
	return rv, nil
	// now jsonPayload contains your desired JSON payload
}

func checkIds(proxy string, newcookie string, useragent string, link string, cursor string) (string, error) {
	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		log.Println(err)
		return "", errors.New("Error creating proxy URL")
	}

	var returncursor string = ""

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	url, _ := url.Parse(link + cursor)

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Println(err)
		return "", err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 4 * time.Second,
	}

	if filterproxies == true && proxyless == "false" {
		client.Transport = transport
	}

	cookies := []*http.Cookie{}
	jar.SetCookies(url, cookies)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Printf("HTTP POST request failed: %s\n", err)
		return "", err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("User-Agent", useragent)

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to catalog: %v\n", err)
		return "", errors.New("Error sending request to catalog")
	}

	if resp != nil {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %s\n", err)
			return "", err
		}

		if debug == "true" {
			fmt.Println(time.Since(start).Milliseconds(), "ms")
		}

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Error unmarshaling JSON response: %s\n", err)
			return "", err
		}

		var wg sync.WaitGroup

		results, ok := result["data"].([]interface{})
		if !ok {
			log.Println("Unable to cast data to array")
			if proxyless == "false" {
				filterRequests()
			}
			return "", errors.New("Error collecting cast data from catalog. RATE LIMITED??")
		}

		workerChannel := make(chan int, workercount)

		if result["nextPageCursor"] != nil {
			returncursor = result["nextPageCursor"].(string)
		}

		for i := range results {
			wg.Add(1)
			workerChannel <- 1
			go func(i int) {
				defer wg.Done()

				item, ok := results[i].(map[string]interface{})
				if !ok {
					log.Printf("Unable to cast item to map: %v\n", results[i])
					return
				}

				if idVal, ok := item["id"].(float64); ok {
					newid := int64(idVal)

					if ResalePrice, ok := item["lowestResalePrice"].(float64); ok {
						newrp := int64(ResalePrice)

						if debug == "false" {
							log.Println("Checked: " + strconv.FormatInt(newid, 10))
						}

						if float64(newrp) != float64(0) {
							newrap, ok := rap[newid].(int64)

							if ok == true {
								if float64(newrp) <= float64(newrap)*profitmargin {
									checkId(newid)
									go checkId(newid)
									checkId(newid)
									nID := strconv.FormatInt(newid, 10)
									nrp := strconv.FormatInt(newrp, 10)
									sendDebugHook(nID, nrp)
								}
							}
						}
					}
				}

				<-workerChannel
			}(i)
		}

		wg.Wait()
	} else {
		log.Println("Response is nil")
		return "", errors.New("Response is nil")
	}

	return returncursor, nil
}

type License struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type ValidationParams struct {
	Key string `json:"key"`
}

type ValidationResult struct {
	Valid bool   `json:"valid"`
	Code  string `json:"code"`
}

type ValidationRequest struct {
	Meta ValidationParams `json:"meta"`
}

type ValidationResponse struct {
	Result  ValidationResult `json:"meta"`
	License *License         `json:"data,omitempty"`
}

func validateLicenseKey(key string) (*ValidationResponse, error) {
	req, err := json.Marshal(ValidationRequest{ValidationParams{key}})
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(req)
	res, err := http.Post(
		fmt.Sprintf("https://api.keygen.sh/v1/accounts/%s/licenses/actions/validate-key", "13026efc-dca0-4993-b5cf-cf24370bb168"),
		"application/vnd.api+json",
		body,
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)

		return nil, fmt.Errorf("An API error occurred: %s", body)
	}

	var v *ValidationResponse
	json.NewDecoder(res.Body).Decode(&v)

	return v, nil
}

func promptForLicenseKey() string {
	return fmt.Sprintf("%s", licenseKey)
}

func main() {
	loadConfig()
	// Validate Keygen license
	licenseKey := promptForLicenseKey()
	validation, err := validateLicenseKey(licenseKey)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	if validation.Result.Valid {
		fmt.Println("License Key is valid!")
	} else {
		fmt.Printf("License key is invalid: code=%s\n", validation.Result.Code)
		return
	}

	fmt.Println("Getting config and initializing Devious Lick")
	time.Sleep(time.Second * 1)

	fmt.Println("Getting proxies...")

	var useragent1 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	time.Sleep(time.Second * 1)
	fmt.Println("Proxies retrieved.")
	time.Sleep(time.Second * 2)
	fmt.Println("Retrieving RAP Info.")
	time.Sleep(time.Second * 1)

	getRecentAveragePrices()

	time.Sleep(time.Second * 1)

	fmt.Println("RAP Info retrieved. Devious Lick initialized. Starting Devious Lick.")
	time.Sleep(time.Second * 3)

	var link1 string = "https://catalog.roblox.com/v2/search/items/details?Category=1&MinPrice=1&MaxPrice=5000&limit=120&CreatorName=ROBLOX&SortType=4&salesTypeFilter=2&cursor="
	var cursor1 string = ""
	var link2 string = "https://catalog.roblox.com/v2/search/items/details?Category=1&MinPrice=5001&MaxPrice=50000&limit=120&CreatorName=ROBLOX&SortType=4&salesTypeFilter=2&cursor="
	var cursor2 string = ""

	siftProxies()

	var clearraplevel = 0

	if proxyless == "false" {
		// Create a channel to send jobs to workers.
		jobChan := make(chan string, 100)

		// Spawn workers.
		for i := 0; i < 3; i++ {
			if clearraplevel >= 5 {
				go getRecentAveragePrices()
				clearraplevel = 0
			}

			go func() {
				for proxy := range jobChan {
					// Make 3 requests using workers
					cursor, _ := checkIds(proxy, cookie, useragent1, link1, cursor1)
					cursor1 = cursor
				}
			}()

			time.Sleep(time.Millisecond * 25)

			go func() {
				for proxy := range jobChan {
					// Make 3 requests using workers
					cursor, _ := checkIds(proxy, cookie, useragent1, link2, cursor2)
					cursor2 = cursor
				}
			}()

			time.Sleep(time.Millisecond * 25)

			go func() {
				for proxy := range jobChan {
					// Make 3 requests using workers
					checkIds(proxy, cookie, useragent1, link1, "")
				}
			}()

			time.Sleep(time.Millisecond * 25)
			clearraplevel++
		}

		// Send jobs to workers.
		for true {
			proxy, err := getRandomProxy()
			if err != nil {
				log.Println(err)
				siftProxies()
				continue
			}
			jobChan <- proxy
		}
	}
}
