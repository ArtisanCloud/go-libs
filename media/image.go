package media

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"os"
)

func SaveImage(imgByte []byte, imgPath string, opts jpeg.Options) {

	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln("image decode error", err)
	}

	out, _ := os.Create(imgPath)
	defer out.Close()

	err = jpeg.Encode(out, img, &opts)
	//jpeg.Encode(out, img, nil)
	if err != nil {
		log.Println("encode error:", err)
	}

}
