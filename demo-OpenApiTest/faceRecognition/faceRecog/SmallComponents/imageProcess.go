package SmallComponents

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

type imageInfo struct {
	ImageByte   []byte
	Format      string
	Width       float32
	Height      float32
	ImageObject image.Image
}

func ReadImage(path string) (img imageInfo) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file failed,err is", err)
		return
	}
	defer f.Close()

	img.ImageByte, _ = ioutil.ReadAll(f)

	img.ImageObject, img.Format, err = image.Decode(bytes.NewReader(img.ImageByte))
	//反复使用文件数据，需要将第一次读取的文件数据重新实例化（NewReader）。因为 go 中读取游标不会重置，导致读取数据为空或者其他异常
	if err != nil {
		log.Println("image decode failed,err is", err)
		return
	}
	img.Width = float32(img.ImageObject.Bounds().Dx())
	img.Height = float32(img.ImageObject.Bounds().Dy())

	return
}

func ImageCut(face faceDetect, img imageInfo) (imgRes []byte) {

	x0 := int(face.Result.List[0].Loc.Left)
	y0 := int(img.Height - face.Result.List[0].Loc.Top - face.Result.List[0].Loc.Height/2)

	x1 := int(face.Result.List[0].Loc.Left + face.Result.List[0].Loc.Width)
	y1 := int(img.Height - face.Result.List[0].Loc.Height)

	log.Printf("切割坐标：x0=%d,x1=%d,y0=%d,y1=%d", x0, x1, y0, y1)

	var buf []byte
	out := bytes.NewBuffer(buf)

	//out,_ := os.Create(outPath)
	//defer out.Close()

	switch img.Format {
	case "jpg":
		fallthrough
	case "jpeg":
		img1 := img.ImageObject.(*image.YCbCr)
		subImg := img1.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		jpeg.Encode(out, subImg, &jpeg.Options{100})
	case "png":
		switch img.ImageObject.(type) {
		case *image.NRGBA:
			img := img.ImageObject.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			png.Encode(out, subImg)
		case *image.RGBA:
			img := img.ImageObject.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			png.Encode(out, subImg)
		}
	case "gif":
		img := img.ImageObject.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		gif.Encode(out, subImg, &gif.Options{})
	default:
		errors.New("ERROR FORMAT")
	}

	imgRes = out.Bytes()

	return
}
