package main

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Downloader struct {
	concurrentcyN int
	resume bool
	bar *progressbar.ProgressBar	//进度条
}

func NewDownloader(concurrencyN int,resume bool) *Downloader {
	return &Downloader{concurrentcyN: concurrencyN,resume: resume}
}

func (d *Downloader) Download(strURL,fileName string) error {
	if fileName == "" {
		fileName = path.Base(strURL)
	}

	resp, err := http.Head(strURL)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK && resp.Header.Get("Accept-Ranges") == "bytes" {
		return d.multiDownload(strURL,fileName,int(resp.ContentLength))
	}

	return d.singleDownload(strURL,fileName)
}

/**
并发下载
 */
func (d *Downloader) multiDownload(strURL,fileName string, contentLen int) error {
	d.setBar(contentLen)

	partSize := contentLen / d.concurrentcyN
	partDir := d.getPartDir(fileName)
	os.Mkdir(partDir,0777)
	defer os.RemoveAll(partDir)

	var wg sync.WaitGroup
	wg.Add(d.concurrentcyN)

	rangeStart := 0

	for i:=0; i<d.concurrentcyN; i++ {
		go func(i, rangeStart int) {
			defer wg.Done()
			rangeEnd := rangeStart + partSize
			if i == d.concurrentcyN-1 {
				rangeEnd = contentLen
			}

			downloaded := 0
			if d.resume {
				partFileName := d.getPartFileName(fileName,i)
				content, err := ioutil.ReadFile(partFileName)
				if err == nil {
					downloaded = len(content)
				}
				d.bar.Add(downloaded)
			}

			d.downloadPartial(strURL, fileName, rangeStart+downloaded, rangeEnd,i)
		}(i, rangeStart)

		rangeStart += partSize + 1
	}

	wg.Wait()

	d.merge(fileName)

	return nil
}

func (d *Downloader) merge(fileName string) error {
	destFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for i := 0; i < d.concurrentcyN; i++ {
		partFileName := d.getPartFileName(fileName,i)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}

		io.Copy(destFile,partFile)
		partFile.Close()

		os.Remove(partFileName)
	}

	return nil
}

/**
下载整个文件
 */
func (d *Downloader) singleDownload(strURL, fileName string) error {
	resp, err := http.Get(strURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(f,resp.Body,buf)
	fmt.Println("singleDownload end ...")
	return nil
}

func (d *Downloader) downloadPartial(strURL, fileName string, rangeStart, rangeEnd, i int) {
	if rangeStart > rangeEnd {
		return
	}

	req, err := http.NewRequest("GET",strURL,nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("range",fmt.Sprintf("bytes=%d-%d",rangeStart,rangeEnd))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	partFile, err := os.OpenFile(d.getPartFileName(fileName,i),flags,0666)
	if err != nil {
		log.Fatal(err)
	}
	defer partFile.Close()

	buf := make([]byte,32*1024)
	//_, err = io.CopyBuffer(partFile,resp.Body,buf)
	_, err = io.CopyBuffer(io.MultiWriter(partFile,d.bar),resp.Body,buf)	//加进度条
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
}

func (d *Downloader) getPartDir(fileName string) string {
	return strings.SplitN(fileName,".",2)[0]
}
func (d *Downloader) getPartFileName(fileName string, partNum int) string {
	dirName := d.getPartDir(fileName)
	return fmt.Sprintf("%s/%s-%d",dirName,fileName,partNum)
}

func (d *Downloader) setBar(contentLength int) {
	d.bar = progressbar.NewOptions(
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
}