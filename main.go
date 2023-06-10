package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

func main() {
	// キーボードの設定を初期化
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	getFileNameList()

	// 入力文字列を非同期にキャプチャするゴルーチンを起動
	go captureInput()

	// メインゴルーチンは終了しないようにする
	select {}
}

func captureInput() {
	str := []int32{}
	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			break
		}

		var builder strings.Builder

		if char != 0 {
			str = append(str, char)
			for _, c := range str {
				builder.WriteRune(c)
			}
			result := builder.String()
			fmt.Println(result)
		}

		time.Sleep(50 * time.Millisecond) // 50ミリ秒ごとにキャプチャする間隔を設定
	}
}

func getFileContentsList() {
	files := []string{}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, string(content))
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, file := range files {
		fmt.Println(file)
	}
}

func getFileNameList() {
	filenames := []string{}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			filenames = append(filenames, info.Name())
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// print slice
	for _, filename := range filenames {
		fmt.Println(filename)
	}
}
