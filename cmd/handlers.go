package cmd

import (
	"bytes"
	"fmt"
	"github.com/gofuncchan/ginger-cli/util"
	"github.com/urfave/cli"
	"io"
	"text/template"
)

// 创建handlers方法
var HandlerCommand = cli.Command{
	Name:        "handler",
	Usage:       "generate handler function code",
	UsageText:   "ginger-cli handler [option]",
	Description: "generate handler function code and request params validator struct",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "module, m", Usage: "module name",},
		cli.StringSliceFlag{Name: "func, f", Usage: "handler function name,one or more"},
	},
	Action: handlerCommandAction,
}

type handlerTmplData struct {
	PackageName string
	ModuleName    string
	FuncName    string
	StructName  string
}

func handlerCommandAction(c *cli.Context) error {
	module := c.String("module")
	fs := c.StringSlice("func")

	var buffs bytes.Buffer
	for _, f := range fs {
		var buff bytes.Buffer
		// handler函数模板
		err := template.Must(template.ParseFiles("./tmpl/handler.tmpl")).Execute(&buff, handlerTmplData{
			PackageName: "handler",
			ModuleName:    module,
			StructName:  f + "Params",
			FuncName:    f,
		})
		if err != nil {
			return err
		}
		io.Copy(&buffs, &buff)
	}

	// 设置输出
	out, err := util.OutputFile(util.HandlerOutput, module)
	if err != nil {
		util.OutputWarn(err.Error())
	}

	_, err = io.Copy(out, &buffs)
	if err != nil {
		return err
	}

	// stdout 输出router代码设置
	util.OutputInfo("Generate Successful", outputHandlerTips(module, fs))

	return nil
}

func outputHandlerTips(moduleName string, funcNames []string) string {
	header := "binding this handler function to app router. \n"
	footer := "According to the http method what you need,copy the code to router/router.go.\n"
	examples := "For example:\n"
	for _, f := range funcNames {
		snake := util.SnakeString(f)
		examples += fmt.Sprintf(`
		r.POST("%s/%s", handler.%s)
	`, moduleName, snake, f)
	}

	return header + examples + footer

}