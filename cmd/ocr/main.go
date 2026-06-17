package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

type Result struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence,omitempty"`
}

func main() {
	jsonOutput := flag.Bool("json", false, "Output in JSON format")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		if *jsonOutput {
			json.NewEncoder(os.Stdout).Encode(Result{Text: "", Confidence: 0})
		} else {
			fmt.Println("Usage: ocr [-json] <image_path>")
		}
		os.Exit(1)
	}

	imagePath := args[0]

	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage("eng")
	client.SetPageSegMode(gosseract.PSM_SINGLE_LINE)
	client.SetWhitelist("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	client.SetBlacklist("!@#$%^&*()_+-=[]{}|;':\",./<>?`~")

	client.SetImage(imagePath)

	text, err := client.Text()
	if err != nil {
		if *jsonOutput {
			json.NewEncoder(os.Stdout).Encode(map[string]string{"error": err.Error()})
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}

	// 清理结果：去除首尾空白和换行，只保留字母数字
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, " ", "")
	text = regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(text, "")

	if *jsonOutput {
		result := Result{Text: text}
		json.NewEncoder(os.Stdout).Encode(result)
		fmt.Println() // 换行
	} else {
		fmt.Println(text) // 加换行
	}
}
