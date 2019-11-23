package cmd

import (
	"bytes"
	"fmt"
	"github.com/gofuncchan/ginger-cli/util"
	"github.com/urfave/cli"
	"io"
	"text/template"
)

// 创建redis cache方法
var CacheCommand = cli.Command{
	Name:        "cache",
	Usage:       "generate cache function code",
	UsageText:   "ginger-cli cache [option]",
	Description: "generate cache function code ",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "module, m", Usage: "cache module name",},
		cli.StringSliceFlag{Name: "func, f", Usage: "cache function name,one or more"},
	},
	Action: cacheCommandAction,
}

type cacheTmplData struct {
	FuncName string
}

func cacheCommandAction(c *cli.Context) error {

	module := c.String("module")
	fs := c.StringSlice("func")

	var buffs bytes.Buffer
	for _, f := range fs {
		var buff bytes.Buffer
		// handler函数模板
		err := template.Must(template.ParseFiles("./tmpl/cache.tmpl")).Execute(&buff, cacheTmplData{
			FuncName:       f,
		})
		if err != nil {
			return err
		}
		io.Copy(&buffs, &buff)
	}

	// 设置输出
	out, err := util.OutputFile(util.CacheOutput, module)
	if err != nil {
		util.OutputWarn(err.Error())
	}

	_, err = io.Copy(out, &buffs)
	if err != nil {
		return err
	}

	// stdout 输出router代码设置
	util.OutputInfo("Generate Successful", outputCacheTips(fs))

	return nil
}

func outputCacheTips(funcNames []string) string {
	header := "Please reset input or output params of cache function.\n"
	footer := "You should handling errors in cache function,and return data or result to caller.\n"
	examples := "For example:\n"

	for _, f := range funcNames {
		examples += fmt.Sprintf(`
		func %s(arg1 int, arg2 string) int64 {
			// TODO model logic code
			// ...
			}
	`,  f)
	}

	return header + examples + footer

}
