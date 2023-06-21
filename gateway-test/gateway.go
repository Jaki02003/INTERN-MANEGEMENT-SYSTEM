package main

import (
	"fmt"
	"gomod/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/labstack/echo"
)

type Response struct {
	Message string `json:"message"`
}

type Response2 struct {
	Vivasoft string `json:"vivasoft"`
	TigerIT  string `json:"tigerit"`
	Cefalo   string `json:"cefalo"`
	Enosis   string `json:"enosis"`
}

var (
	serverArray   [4]string
	reqCapacity   map[string]int
	reqPassed     map[string]int
	whichURL      map[string]string
	capacityMutex sync.Mutex
	acceptedOrNot string
)

func sendRequestToServer(internStudent model.InternStudent, company string) (string, bool) {
	capacityMutex.Lock()
	if reqCapacity[company] < reqPassed[company] {
		capacityMutex.Unlock()
		return "Request Limit exceeded", false
	}

	reqPassed[company]++
	capacityMutex.Unlock()

	serverURL := whichURL[company]
	fmt.Println(internStudent.TotalSolved, company)

	params := url.Values{}
	params.Set("username", internStudent.UserName)
	params.Set("totalsolved", strconv.Itoa(internStudent.TotalSolved))
	params.Set("cgpa", strconv.FormatFloat(internStudent.CGPA, 'f', -1, 64))

	fullURL := fmt.Sprintf("%s?%s", serverURL, params.Encode())
	fmt.Println(fullURL)

	resp, err := http.PostForm(fullURL, params)
	if err != nil {
		fmt.Printf("Error connecting to the server: %v\n", err)
		return "Error connecting to the server", false
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Server is working.")

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return "Error reading response body", false
		}

		fmt.Println("Response Body:", string(body))

		if string(body) == "Not Accepted" {
			return string(body), false
		}
		return string(body), true
	} else {
		fmt.Printf("Server returned a non-OK status: %s\n", resp.Status)
		return "Server returned a non-OK status", false
	}
}

func handleRequest(c echo.Context) error {
	userName := c.QueryParam("username")
	totalSolved, _ := strconv.Atoi(c.QueryParam("totalsolved"))
	cgpa, _ := strconv.ParseFloat(c.QueryParam("cgpa"), 64)
	fmt.Println(userName, cgpa, totalSolved)

	intern := model.InternStudent{
		UserName:    userName,
		TotalSolved: totalSolved,
		CGPA:        cgpa,
	}

	acceptedOrNot = "Not Accepted"
	var ret string
	var retBool bool

	for i := 0; i < len(serverArray); i++ {
		ret, retBool = sendRequestToServer(intern, serverArray[i])
		if retBool {
			acceptedOrNot = ret
			break
		}
	}

	return c.JSON(http.StatusOK, Response{
		Message: acceptedOrNot,
	})
}

func handleWatch(c echo.Context) error {
	vivasoft := strconv.Itoa(reqPassed["vivasoft"])
	tigerit := strconv.Itoa(reqPassed["tigerit"])
	cefalo := strconv.Itoa(reqPassed["cefalo"])
	enosis := strconv.Itoa(reqPassed["enosis"])

	return c.JSON(http.StatusOK, Response2{
		Vivasoft: vivasoft,
		TigerIT:  tigerit,
		Cefalo:   cefalo,
		Enosis:   enosis,
	})
}

func main() {
	serverArray = [4]string{"vivasoft", "tigerit", "cefalo", "enosis"}
	reqCapacity = make(map[string]int)
	reqPassed = make(map[string]int)
	whichURL = make(map[string]string)

	for i := 0; i < len(serverArray); i++ {
		reqCapacity[serverArray[i]] = (i + 3) * 10

		whichURL[serverArray[i]] = "http://localhost:" + strconv.Itoa(8000+i) + "/apply"
	}

	e := echo.New()
	e.GET("/request", handleRequest)
	e.GET("/watch", handleWatch)
	e.Start(":8080")
}
