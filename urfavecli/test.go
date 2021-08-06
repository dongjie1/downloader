package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
	"strconv"
)

var language string

//直接运行命令  go run test.go -pt 80 -l zh
//可以 go install 安装命令再运行 test -p 80 -l zh
func main() {
	app := &cli.App{
		Name: "boom",
		Usage: "make an explosive entrance",
		Action: func(c *cli.Context) error {
			fmt.Println("boom! I Say")
			fmt.Println(c.String("lang"),c.Int("port"))
			fmt.Println(language)
			if language == "zh" {
				fmt.Println("你好")
			}
			return nil
		},
	}

	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Aliases: []string{"p","pt"},//--port的简写,也就是说"--port"和"-p"是等价的
			Value: 8000,             //可以指定flag的默认值。
			Usage: "listening port", //flag的描述信息
		},
		&cli.StringFlag{
			Name: "lang",
			Aliases: []string{"l","la"},//--lang的简写-l
			Value: "english",
			Usage: "read from `File`",//如果你想给用户增加一些属性值类型的提示，可以通过占位符（placeholder）来实现，比如上面的"--lang FILE"。占位符通过``符号来标识。
			Destination: &language,//可以为该flag指定一个接收者，比如上面的language变量。解析完"--lang"这个flag后会自动存储到这个变量里，后面的代码就可以直接使用这个变量的值了。
		},
	}

	app.Commands = []*cli.Command{
		{
			Name: "add",
			Aliases: []string{"a"},
			Usage: "cal 1+1",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				a,_ := strconv.Atoi(c.Args().Get(0))
				b,_ := strconv.Atoi(c.Args().Get(1))
				fmt.Println("cal sum:",(a+b))
				return nil
			},
		},
		{
			Name: "sub",
			Aliases: []string{"s"},
			Usage: "cal 5-3",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				a,_ := strconv.Atoi(c.Args().Get(0))
				b,_ := strconv.Atoi(c.Args().Get(1))
				fmt.Println("cal sub=",a-b)
				return nil
			},
		},
		{
			Name: "db",
			Usage: "db operations",
			Category: "db",
			Subcommands: []*cli.Command{
				{
					Name: "insert",
					Usage: "insert db data",
					Action: func(c *cli.Context) error {
						fmt.Println("insert db data")
						return nil
					},
				},
				{
					Name: "delete",
					Usage: "delete db data",
					Action: func(c *cli.Context) error {
						fmt.Println("delete db data")
						return nil
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
