package imagegen

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"time"

	// "image/gif"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

// func hmain(){
// 	imurl := "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png"
// 	slash1 :="imagegen/sprites/2.png"
// 	aimg, _ := fetchandConvert(imurl)
// 	slsmg,err := normfetchandConvert(slash1)
// 	if err != nil{fmt.Println("Error:",err)
// 	}else{fmt.Println(aimg)
// 		  fmt.Println(slsmg)}
// }

// fetchandConvert(imurl string)
func AsciiGen(imageURL string,reqwidth int) (string, error) {
// func fetchandConvert(imageURL  string) (string, error) {
	res, err := http.Get(imageURL)
	if err != nil || res.StatusCode != http.StatusOK {
		return "[Image Unavailable]", fmt.Errorf("failed to fetch image: %v", err)
	}
	defer res.Body.Close()
	fileType := res.Header.Get("Content-Type")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %v", err)
	}
	var img image.Image
	switch {
	case strings.Contains(fileType, "png"):
		img, err = png.Decode(bytes.NewReader(body))
	default:
		img, _, err = image.Decode(bytes.NewReader(body))
	}
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}
	var newid = reqwidth //96 best
	asciiChars := " .:-=+*#%@" 
	img = resize.Resize(uint(newid),uint(newid), img, resize.Lanczos3)
	cropimg := image_cropping(img)
	ascii := rgbconvertToAscii(cropimg,asciiChars,2)
	ascii = trimAndPadAscii(ascii)
	return ascii,nil
}

func AttackGen(path string) (string){
	var img image.Image
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	// Decode the PNG image
	img, err = png.Decode(file)
	if err != nil {
		return ""
	}
	var newid = 96 //96 best 
	asciiChars := " .:-=+*#"
	img = resize.Resize(uint(newid),uint(newid), img, resize.Lanczos3)
	cropimg := image_cropping(img)
	ascii := rgbconvertToAscii(cropimg,asciiChars,2)
	ascii = trimAndPadAscii(ascii)
	return ascii

}

func BgMaker(path string) (string){
	var img image.Image
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	// Decode the PNG image
	img, err = png.Decode(file)
	if err != nil {
		return ""
	}
	var newid = 40 //96 best 
	asciiChars := " .:-=+*#"
	img = resize.Resize(uint(newid),uint(newid), img, resize.Lanczos3)
	// cropimg := image_cropping(img)
	ascii := convertToAscii(img,asciiChars,5)
	ascii = trimAndPadAscii(ascii)
	// ascii = removelines(ascii)
	return ascii
	//trim last 4-5 lines
}

func GifGen(imageURL string,reqwidth int) ([]string,time.Duration,error){
	res, err := http.Get(imageURL)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil,0,fmt.Errorf("failed to fetch gif: %v", err)
	}
	defer res.Body.Close()
	fileType := res.Header.Get("Content-Type")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil,0,fmt.Errorf("failed to read gif data: %v", err)
	}
	var frames []string
	switch {
	case strings.Contains(fileType, "gif"):
		gifImg, err := gif.DecodeAll(bytes.NewReader(body))
		if err != nil {
			return nil,0,fmt.Errorf("failed to decode GIF: %v", err)
		}
		for _, frame := range gifImg.Image {
			asciiFrame := processImageFrame(frame, reqwidth)
			frames = append(frames, asciiFrame)
		}
	default:
		return nil,0, fmt.Errorf("unsupported image format")
	}
	delay := 150 * time.Millisecond
	// displayAsciiAnimation(frames, delay)
	return frames,delay ,nil
} 

func processImageFrame(img image.Image, reqwidth int) string {
	asciiChars := " .:-=+*#"
	img = resize.Resize(uint(reqwidth), 0, img, resize.Lanczos3)
	cropimg := image_cropping(img)
	ascii := rgbconvertToAscii(cropimg, asciiChars,2)
	ascii = trimAndPadAscii(ascii)
	return ascii
}

func displayAsciiAnimation(frames []string, delay time.Duration) {
	for {
		if time.Since(time.Now()) > (500*time.Millisecond) {
			break
		}	
		for _, frame := range frames {
			//fmt.Print("\033[H\033[2J") // Clear screen
			fmt.Println(frame)
			time.Sleep(delay)
		}
	 }
}

func image_cropping(img image.Image)image.Image{
	bounds := img.Bounds()

	var minX, minY, maxX, maxY int
	minX, minY = bounds.Max.X, bounds.Max.Y

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if !(r == 0xffff && g == 0xffff && b == 0xffff){
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}}}
	
	rect := image.Rect(minX, minY, maxX+1, maxY+1)
	croppedImg := image.NewRGBA(rect)
	draw.Draw(croppedImg, rect.Bounds(), img, rect.Min, draw.Src)
	// outFile, err := os.Create("imagegen/crop_images/output.png")
	// if err != nil {
	// 	panic(err)
	// }
	// defer outFile.Close()
	// png.Encode(outFile, croppedImg)
	return croppedImg
}

func rgbToAnsi(r, g, b uint32) string {
	r, g, b = r>>8, g>>8, b>>8 // scale down to 8-bit values
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func convertToAscii(img image.Image,asciiChars string,inc int) string {
	var asciiArt strings.Builder
	bounds := img.Bounds()
	count := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y += inc {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			charIndex := int(brightness * float64(len(asciiChars)-1) / 65535)
			char := (asciiChars[charIndex])
			asciiArt.WriteByte(char)
		}		
		asciiArt.WriteString("\n")
		count ++
		
	}
	return asciiArt.String()
}

func rgbconvertToAscii(img image.Image,asciiChars string,yinc int ) string{
	var asciiArt strings.Builder
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y+=yinc {
		for x := bounds.Min.X; x < bounds.Max.X; x+=1 {
			r, g, b, _ := img.At(x, y).RGBA()
			// brightness := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
			brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			charIndex := int(brightness * float64(len(asciiChars)-1) / 65535)
			// charInd := int((brightness/255.0) * float64(len(asciiChars)-1))
			char := asciiChars[charIndex]
			colorCode := rgbToAnsi(r, g, b)
			if char != ' '{
				asciiArt.WriteString(colorCode + string(char) + "\x1b[0m")
			}else{
				asciiArt.WriteByte(char)
			}
		}
		asciiArt.WriteString("\n")
	}
	return asciiArt.String()
}

// func downscaleAscii(ascii string, originalWidth, targetWidth int) string {
// 	lines := strings.Split(ascii, "\n")
// 	var downscaled strings.Builder
// 	ratio := originalWidth / targetWidth

// 	for i := 0; i < len(lines); i += ratio {
// 		for j := 0; j < originalWidth; j += ratio {
// 			if j < len(lines[i]) {
// 				downscaled.WriteByte(lines[i][j])
// 			}
// 		}
// 		downscaled.WriteString("\n")
// 	}
// 	return downscaled.String()
// }


func trimAndPadAscii(ascii string) string {
    lines := strings.Split(ascii, "\n")

    for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
        lines = lines[1:]
    }

    for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
        lines = lines[:len(lines)-1]
    }
    minLeadingSpaces := len(lines[0])
    for _, line := range lines {
        trimmedLine := strings.TrimSpace(line)
        if len(trimmedLine) > 0 {
            leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
            if leadingSpaces < minLeadingSpaces {
                minLeadingSpaces = leadingSpaces
            }
        }
    }

    for i, line := range lines {
        trimmed := line[minLeadingSpaces:]
        lines[i] = "  " + trimmed + "  " 
    }

    lines = append([]string{""}, lines...)
    lines = append(lines, "")

    return strings.Join(lines, "\n")
}

// func removelines(ascii string)(string){
// 	lines := strings.Split("\n", ascii)
// 	var newlines strings.Builder
// 	for i, line := range lines{
// 		if i < 2{
// 			newlines.WriteString(line+"\n")
// 		}else if i == 2{
// 			newlines.WriteString(line)
// 		} 
// 	}
// 	return newlines.String()
// }

// package imagegen

