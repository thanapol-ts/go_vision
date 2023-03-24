package main

import (
	"fmt"
	"github/go_vision/response"
	"net/http"
	"os"
	"strconv"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gin-gonic/gin"
)

func main() {

	type Infomation struct {
		IdCard  string `json:"id_card"`
		Name    string `json:"name"`
		Address string `json:"address"`
		Dob     string `json:"dob"`
	}

	r := gin.Default()
	r.POST("/check", func(ctx *gin.Context) {
		var result bool = false
		client, err := vision.NewImageAnnotatorClient(ctx)
		if err != nil {
			fmt.Println("err", err)
		}

		file, _ := ctx.FormFile("file")
		if err := ctx.SaveUploadedFile(file, "../assets/"+file.Filename); err != nil {
			fmt.Println("err", err)
		}
		getfile, err := os.Open("../assets/" + file.Filename)
		if err != nil {
			panic(err)
		}
		defer getfile.Close()
		image, err := vision.NewImageFromReader(getfile)
		if err != nil {
			fmt.Printf("Failed to create image: %v", err)
		}
		labels, err := client.DetectTexts(ctx, image, nil, 20)
		if err != nil {
			fmt.Printf("Failed to detect labels: %v", err)
		}
		lines := strings.Split(labels[0].Description, "\n")
		info := Infomation{}
		for index, label := range lines {
			fmt.Printf("'%d'. '%s'", index, label)
			fmt.Printf("\n")
			substrings := strings.Split(label, " ")
			id := strings.Join(substrings, "")
			if CheckID(id) {
				info.IdCard = id
			}

			if CheckContains(label, "ชื่อสกุล") {
				st := strings.Replace(label, "ชื่อตัวและชื่อสกุล ", "", 1)
				info.Name = st
			}

			if CheckContains(label, "ที่อยู่") {
				st := strings.Replace(label, "ที่อยู่ ", "", 1)
				info.Address += st + lines[index+1]
			}

			if CheckContains(label, "เกิดวันที่") {
				st := strings.Replace(label, "เกิดวันที่ ", "", 1)
				info.Dob = st
			}
		}

		if (info != Infomation{}) {
			result = true
		} else {
			result = false
		}

		// os.Remove("../assets/" + file.Filename)
		res := response.Response{
			Result:  result,
			Status:  http.StatusOK,
			Message: "success",
			Data:    info,
		}
		if result {
			ctx.JSON(http.StatusOK, res)
		} else {
			ctx.JSON(http.StatusNotFound, res)
		}
	})
	r.GET("/watch", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pass")
	})
	r.Run(":4000")
}

func CheckID(id string) bool {
	if len(id) != 13 {
		return false
	}
	sum := 0
	for i := 0; i < 12; i++ {
		digit, err := strconv.ParseFloat(string(id[i]), 64)
		if err != nil {
			return false
		}
		sum += int(digit) * (13 - i)
	}
	return (11-sum%11)%10 == int(id[12]-'0')
}

func CheckContains(input string, condition string) bool {
	return strings.Contains(input, condition)
}
