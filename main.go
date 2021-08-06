package main
/**
https://github.com/polaris1119/downloader

运行:
go build
 ./downloader --url https://apache.claz.org/zookeeper/zookeeper-3.7.0/apache-zookeeper-3.7.0-bin.tar.gz

 */
import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
)

func main() {
	concurrencyN := runtime.NumCPU()

	app := &cli.App{
		Name: "downloader",
		Usage: "File concurrency downloader",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "url",
				Aliases: []string{"u"},
				Usage: "`url` to download",
				Required: true,
			},
			&cli.StringFlag{
				Name: "output",
				Aliases: []string{"o"},
				Usage: "output `filename`",
			},
			&cli.IntFlag{
				Name: "concurrencyN",
				Aliases: []string{"n"},
				Usage: "Concurrency `number`",
				Value: concurrencyN,
			},
		},
		Action: func(c *cli.Context) error {
			strURL := c.String("url")
			fileName := c.String("output")
			concurrencyN := c.Int("concurrencyN")

			return NewDownloader(concurrencyN,false).Download(strURL,fileName)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
