package angellist

import (
	"encoding/json"
	gq "github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func getJobsPage(searchWordsList []string) (jobsIDList []int, ifJobsIDFound bool) {

	angelistCookie, csrfToken, isLoginPageFetched := fetchLoginPage()

	if isLoginPageFetched {

		//time.Sleep(time.Duration(20) * time.Second)

		loggedInCookieString, loginSuccessful := sendLoginRequest(angelistCookie, csrfToken)

		//time.Sleep(time.Duration(5) * time.Second)

		if loginSuccessful {
			jobsIDList, ifJobsIDFound = setSearchKeywords(searchWordsList, loggedInCookieString)
		}
	}

	return
}

func sendLoginRequest(angelistCookie string, csrfToken string) (loggedInCookieString string, isLoginSuccessful bool) {

	loginURL := "https://angel.co/users/login"

	urlParams := url.Values{}

	urlParams.Add("utf8", "✓")
	urlParams.Add("authenticity_token", csrfToken)
	urlParams.Add("login_only", "true")
	urlParams.Add("user[email]", "dineshchander28@gmail.com")
	urlParams.Add("user[password]", "22dispareil22")

	loginDetails := urlParams.Encode()

	headers := make(map[string]string)

	headers["Content-Length"] = strconv.Itoa(len(loginDetails))
	headers["Referer"] = "https://angel.co/login"
	headers["Cookie"] = angelistCookie
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	loginResponse, loginErr := makeRequestToAngelListServer("POST", loginURL, loginDetails, headers, false)

	if loginErr != nil && loginResponse == nil {
		loggerInstance.Println(loginErr.Error())
		return
	}

	defer loginResponse.Body.Close()

	if loginResponse.StatusCode == 302 {

		cookies := loginResponse.Cookies()

		if len(cookies) > 0 {

			angelistLoggedInCookie := cookies[0]

			if angelistLoggedInCookie.Name == "_angellist" && angelistLoggedInCookie.Value != "" && angelistLoggedInCookie.Value != angelistCookie {

				loggedInCookieString = angelistLoggedInCookie.Name + "=" + angelistLoggedInCookie.Value
				isLoginSuccessful = true
			}
		}

	}

	return
}

func fetchLoginPage() (angelistCookieString string, csrfToken string, isLoginPageFetched bool) {

	loginPageResponse, loginPageFetchErr := makeRequestToAngelListServer("GET", "https://angel.co/login", "", nil, false)

	if loginPageFetchErr != nil {
		loggerInstance.Println(loginPageFetchErr.Error())
		return
	}

	angelistCookieString, csrfToken, isLoginPageFetched = getCSRFandCookie(loginPageResponse)

	loginPageResponse.Body.Close()

	return
}

func getCSRFandCookie(pageResponse *http.Response) (angelistCookieString string, csrfToken string, ifFetched bool) {

	cookies := pageResponse.Cookies()

	if len(cookies) > 0 {

		angelistCookie := cookies[0]

		doc, pageErr := gq.NewDocumentFromResponse(pageResponse)

		if pageErr != nil {
			loggerInstance.Println(pageErr.Error())
			return
		}

		var nodeFound bool

		csrfToken, nodeFound = doc.Find(`[name="csrf-token"]`).Attr("content")

		if !nodeFound {

			loggerInstance.Println("No CSRF token found")

		} else if csrfToken != "" && angelistCookie.Name == "_angellist" && angelistCookie.Value != "" {

			angelistCookieString = angelistCookie.Name + "=" + angelistCookie.Value
			ifFetched = true
		}
	}

	return
}

func setSearchKeywords(searchWordsList []string, loggedInCookieString string) (jobsIDList []int, ifJobsIDFound bool) {

	searchKeywordsURL := "https://angel.co/job_listings/startup_ids"

	params := url.Values{}

	params.Add("tab", "find")

	for _, searchWord := range searchWordsList {
		params.Add("filter_data[keywords][]", searchWord)
	}

	params.Add("filter_data[roles][]", "Software Engineer")

	requestQuery := params.Encode()

	headers := make(map[string]string)

	headers["Content-Length"] = strconv.Itoa(len(requestQuery))
	headers["Referer"] = "https://angel.co/jobs"
	headers["Cookie"] = loggedInCookieString
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["X-Requested-With"] = "XMLHttpRequest"

	queryResponse, queryErr := makeRequestToAngelListServer("POST", searchKeywordsURL, requestQuery, headers, false)

	if queryErr != nil {
		loggerInstance.Println(queryErr.Error())
		return
	}

	defer queryResponse.Body.Close()

	bufferedJSON, readErr := ioutil.ReadAll(queryResponse.Body)

	if readErr != nil {
		loggerInstance.Println(readErr.Error())
	} else {

		ids := gjson.GetBytes(bufferedJSON, "ids")

		unmarshalError := json.Unmarshal([]byte(ids.String()), &jobsIDList)

		if unmarshalError != nil {
			loggerInstance.Println(unmarshalError.Error())
		} else {
			ifJobsIDFound = true
		}
	}

	return
}
