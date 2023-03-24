package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gin-gonic/gin"
)

func main() {
	var msg struct {
		NdId   string
		Name   string
		Status bool
	}

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

		getfile, err := os.Open("../assets/" + file.Filename)
		if err != nil {
			panic(err)
		}
		defer getfile.Close()

		// Decode the image
		img, err := jpeg.Decode(getfile)
		if err != nil {
			panic(err)
		}

		// Define the size of the pieces to cut
		width := img.Bounds().Max.X
		height := 500

		// Loop over the image and cut it into pieces
		for y := 0; y < 500; y += height {
			for x := 0; x < img.Bounds().Max.X; x += width {
				// Define the rectangle to cut
				rect := image.Rect(x, y, x+width, y+height)
				// Cut the image
				piece := img.(interface {
					SubImage(r image.Rectangle) image.Image
				}).SubImage(rect)

				// Save the piece to a file
				out, err := os.Create("../assets/" + file.Filename)
				if err != nil {
					panic(err)
				}
				defer out.Close()

				// Encode the piece and save it to the file
				jpeg.Encode(out, piece, nil)
			}
		}

		getFile, err := os.Open("../assets/" + file.Filename)
		if err != nil {
			fmt.Printf("Failed to read file: %v", err)
		}

		defer getFile.Close()
		image, err := vision.NewImageFromReader(getFile)
		if err != nil {
			fmt.Printf("Failed to create image: %v", err)
		}
		labels, err := client.DetectTexts(ctx, image, nil, 10)
		if err != nil {
			fmt.Printf("Failed to detect labels: %v", err)
		}

		lines := strings.Split(labels[0].Description, "\n")
		var joined = ""
		for index, label := range lines {
			fmt.Printf("'%d'. '%s'", index, label)
			fmt.Printf("\n")
			if index == 1 {
				substrings := strings.Split(label, " ")
				joined = strings.Join(substrings, "")
				if !CheckID(joined) {
					msg.Status = false
					break
				}
				msg.Status = true
				msg.NdId = joined
			}

			if index == 4 {
				if CheckContains(label) {
					st := strings.Split(label, "ชื่อตัวและชื่อสกุล")
					msg.Name = st[1]
				} else {
					msg.Status = false
					break
				}
			}
		}

		fmt.Println("msg", msg)
		if msg.Status {
			ctx.JSON(http.StatusOK, msg)
		} else {
			msg.Name = ""
			msg.NdId = ""
			os.Remove("../assets/" + file.Filename)
			ctx.JSON(http.StatusNotFound, msg)
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

func CheckContains(input string) bool {
	return strings.Contains(input, "ชื่อสกุล")
}
