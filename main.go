package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	//"sync"
	"time"

	"github.com/eiannone/keyboard"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear your terminal screen")
	}
}

func main() {
	// キーボードの設定を初期化
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	//var wg sync.WaitGroup
	//wg.Add(1)

	// File名一覧を取得
	fnames, err := getFileNameList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fnames)

	// チャネルを作成
	ch := make(chan string)

	// 入力文字列を非同期にキャプチャするゴルーチンを起動
	go captureInput(ch)

	for {
		select {
		// goroutineから受けとった文字列（キーボード入力）を表示
		case s := <- ch:
			CallClear()
			fmt.Printf("> ")
			fmt.Println(s)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func captureInput(ch chan string) {
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
			ch <- result
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

func getFileNameList() ([]string, error) {
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
		return []string{}, err
	}

	return filenames, nil
}
