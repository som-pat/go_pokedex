package imagegen

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"net/http"
	"strings"

	"github.com/nfnt/resize"
)

// func main(){
// 	imurl := "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/208.png"
// 	aimg, err := fetchandConvert(imurl)
// 	if err != nil{fmt.Println("Error:",err)
// 	}else{fmt.Println(aimg)}
// }

func AsciiGen(imageURL string,reqwidth int) (string, error) {
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
	case strings.Contains(fileType, "gif"):
		img, err = gif.Decode(bytes.NewReader(body))
	default:
		img, _, err = image.Decode(bytes.NewReader(body))
	}
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}
	var newid = reqwidth //96 best 
	img = resize.Resize(uint(newid),uint(newid), img, resize.Lanczos3)
	cropimg := image_cropping(img)
	ascii := convertToAscii(cropimg)
	ascii = trimAndPadAscii(ascii)
	return ascii,nil
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

func convertToAscii(img image.Image) string{
	asciiChars := " .:-=+*#%@"
	var asciiArt strings.Builder
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y+=2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
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
