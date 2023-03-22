package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/check", func(ctx *gin.Context) {
		client, err := vision.NewImageAnnotatorClient(ctx)
		if err != nil {
			fmt.Println("err", err)
		}

		file, _ := ctx.FormFile("file")
		if err := ctx.SaveUploadedFile(file, "../assets/"+file.Filename); err != nil {
			fmt.Println("err", err)
		}

		getFile, err := os.Open("../assets/" + file.Filename)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		defer getFile.Close()
		image, err := vision.NewImageFromReader(getFile)
		if err != nil {
			log.Fatalf("Failed to create image: %v", err)
		}
		labels, err := client.DetectTexts(ctx, image, nil, 10)
		if err != nil {
			log.Fatalf("Failed to detect labels: %v", err)
		}

		fmt.Println("Labels: ", labels)
		fmt.Println("============================================================")
		lines := strings.Split(labels[0].Description, "\n")
		// var joined = ""
		// var name = ""
		for index, label := range lines {
			fmt.Printf("'%d'. '%s'", index, label)
			fmt.Printf("\n")
			// if index == 2 {
			// 	substrings := strings.Split(label, " ")
			// 	joined = strings.Join(substrings, "")
			// 	fmt.Println("index == ", index, " ", joined)
			// }
			// if index == 4 {
			// 	st := strings.Split(label, "ชื่อตัวและชื่อสกุล")
			// 	fmt.Println("index == ", index, " ", st[1])
			// 	name = st[1]
			// }
		}
		var msg struct {
			NdId   string
			Name   string
			Status bool
		}
		// msg.Name = name
		// msg.NdId = joined
		msg.Status = true
		fmt.Println("msg", msg)
		ctx.JSON(http.StatusOK, msg)
	})
	r.GET("/watch", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pass")
	})
	r.Run(":4000")
}
