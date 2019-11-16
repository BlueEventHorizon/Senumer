package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"./analyzer"
)

func main() {

	// flagの使い方
	// https://qiita.com/Yaruki00/items/7edc04720a24e71abfa2

	var (
		topDir     string
		enumString string
		enumImage  string
		enumColor  string
		useDefault bool
		//	enableImage = flag.Bool("image", true, "enable scan for image assets")
		//	enableColor = flag.Bool("color", true, "enable scan for color assets")
	)

	flag.StringVar(&topDir, "dir", "./", "dir to scan")
	flag.StringVar(&enumString, "string", "", "enum name for Localizable.strings. If blank, disable output")
	flag.StringVar(&enumImage, "image", "", "enum name for Image Assets. If blank, disable output")
	flag.StringVar(&enumColor, "color", "", "enum name for Color Assets. If blank, disable output")
	flag.BoolVar(&useDefault, "default", false, "if true, enable all output with default file name")
	flag.Parse()

	if useDefault {
		enumString = "LocalizableStrings"
		enumImage = "AppResource.ImageResource"
		enumColor = "AppResource.ColorResource"
	}

	// ---- LocalizableStrings ----
	if enumString != "" {
		stringOutput := new(Output)
		stringOutput.Open(fmt.Sprintf("%s/%s.swift", topDir, enumString))
		stringOutput.Print("import Foundation\n\n")
		stringOutput.Print(fmt.Sprintf("enum %s: String {\n", enumString))

		texts := make([]string, 100, 500)
		ScanFile(topDir, analyzer.LocalisableStringsAnalyzer, &texts)
		for _, text := range texts {
			if text == "" {
				continue
			}
			// 空白はアンダースコアに置換
			keyword := strings.Replace(text, " ", "_", -1)
			// ピリオドはアンダースコアに置換
			keyword = strings.Replace(keyword, " ", ".", -1)
			// ハイフンはアンダースコアに置換
			keyword = strings.Replace(keyword, " ", "-", -1)

			keyword = convertToCamelCase(keyword)
			stringOutput.Print(fmt.Sprintf("    case %s = \"%s\",\n", keyword, text))
		}
		stringOutput.Print("}\n")
		stringOutput.Close()
		fmt.Printf("Completed to generate %s\n", enumString)
	} else {
		fmt.Println("Skipped to scan Localizable.strings")
	}

	// ---- imageAssets ----
	if enumImage != "" {
		imageOutput := new(Output)
		imageOutput.Open(fmt.Sprintf("%s/%s", topDir, enumImage))
		imageAssets := make([]string, 0, 500)
		ScanDir(topDir, analyzer.ImageAssetAnalyzer, &imageAssets)
		for _, asset := range imageAssets {
			if asset == "" {
				continue
			}
			imageOutput.Print(fmt.Sprintf("imageAssets = \"%s\",\n", asset))
		}
		imageOutput.Close()
		fmt.Printf("Completed to generate %s\n", enumImage)
	} else {
		fmt.Println("Skipped to scan Image Assets")
	}

	// ---- colorAssets ----
	if enumColor != "" {
		colorOutput := new(Output)
		colorOutput.Open(fmt.Sprintf("%s/%s", topDir, enumColor))
		colorAssets := make([]string, 0, 500)
		ScanDir(topDir, analyzer.ColorAssetAnalyzer, &colorAssets)
		for _, asset := range colorAssets {
			if asset == "" {
				continue
			}
			colorOutput.Print(fmt.Sprintf("colorAssets = \"%s\",\n", asset))
		}
		colorOutput.Close()
		fmt.Printf("Completed to generate %s\n", enumColor)
	} else {
		fmt.Println("Skipped to scan Color Assets")
	}
}

func convertToCamelCase(text string) string {
	var keyword string
	var foundUnderScore = false
	for i := 0; i < len(text); i++ {
		letter := text[i : i+1]
		if letter == "_" {
			foundUnderScore = true
			continue
		}
		if foundUnderScore {
			foundUnderScore = false
			keyword = keyword + strings.ToUpper(letter)
		} else {
			keyword = keyword + letter
		}
	}

	return keyword
}

type Output struct{}

var fd *os.File
var err error

func (t Output) Open(path string) {

	if path == "" {
		return
	}

	fd, err = os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	fd.Seek(0, 0)

}

func (t Output) Print(text string) {
	if fd != nil {
		fd.WriteString(text)
	} else {
		fmt.Print(text)
	}
}

func (t Output) Close() {
	if fd == nil {
		return
	}
	fd.Close()
}
