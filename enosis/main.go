package main

import (
	"fmt"
	"gomod/model"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	capacity int = 40

	capacityMutex sync.Mutex
)

func main() {
	e := echo.New()

	e.POST("/apply", func(c echo.Context) error {
		intern := new(model.InternStudent)

		intern.UserName = c.QueryParam("username")
		tempSolved, _ := strconv.Atoi(c.QueryParam("totalsolved"))

		intern.TotalSolved = tempSolved
		tempCGPA, _ := strconv.ParseFloat(c.QueryParam("cgpa"), 64)
		intern.CGPA = tempCGPA

		if intern.CGPA >= 3 && intern.TotalSolved >= 300 {
			fmt.Println(tempCGPA)
			capacityMutex.Lock()
			if capacity > 0 {
				capacity--
				capacityMutex.Unlock()

				return c.String(http.StatusOK, "Enosis Accepted")
			}
			capacityMutex.Unlock()
		}

		return c.String(http.StatusOK, "Not Accepted")
	})

	e.Start(":8003")
}
