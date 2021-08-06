package main

import (
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	//bar := progressbar.Default(100)
	//for i:=0; i<100; i++ {
	//	err := bar.Add(1)
	//	if err != nil {
	//		return
	//	}
	//	time.Sleep(40 * time.Millisecond)
	//}

	//自定义样式
	//bar2 := progressbar.NewOptions(100,
	//	progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
	//	progressbar.OptionEnableColorCodes(true),
	//	progressbar.OptionShowBytes(true),
	//	progressbar.OptionSetWidth(15),
	//	progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
	//	progressbar.OptionSetTheme(progressbar.Theme{
	//		Saucer: "[green]=[reset]",
	//		SaucerHead: "[green]>[reset]",
	//		SaucerPadding: " ",
	//		BarStart: "[",
	//		BarEnd: "]",
	//	},
	//	),
	//)
	//for i:=0; i<100; i++ {
	//	bar2.Add(1)
	//	time.Sleep(50 * time.Millisecond)
	//}

	//根据文件大小显示bar
	req,_ := http.NewRequest("GET","https://dl.google.com/go/go1.14.2.src.tar.gz",nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	f,_ := os.OpenFile("go1.14.2.src.tar.gz", os.O_CREATE|os.O_WRONLY, 0664)
	defer f.Close()

	//bar3 := progressbar.DefaultBytes(
	//	resp.ContentLength,
	//	"downloading...",
	//	)
	bar3 := setBar(resp.ContentLength)	//自定义bar样式

	io.Copy(io.MultiWriter(f,bar3),resp.Body)

}

func setBar(contentLength int64) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		contentLength,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionSetWidth(15),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetDescription("downloading..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer: "[green]=[reset]",
			SaucerHead: "[green]>[reset]",
			SaucerPadding: " ",
			BarStart: "[",
			BarEnd: "]",
		}),
		)
	return bar
}
