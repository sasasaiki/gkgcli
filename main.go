package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/sasasaiki/gkgfiler"
)

const gokigenRepo = "https://github.com/sasasaiki/gokigen.git"

func main() {
	app := cli.NewApp()
	app.Name = "gkgcli"
	app.Usage = "Webアプリケーションテンプレートgokigenのコマンドラインツールです"
	app.Version = "1.0.0"

	// グローバルオプション設定
	// app.Flags = []cli.Flag{
	// 	cli.BoolFlag{
	// 		Name:  "dryrun, d", // 省略指定 => d
	// 		Usage: "グローバルオプション dryrunです。",
	// 	},
	// }

	app.Commands = []cli.Command{
		// コマンド設定
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "gokigenのスケルトンを$GOPATH/指定されたディレクトリに展開します",
			Action:  newAction,
		},
	}

	app.Before = func(c *cli.Context) error {
		// 開始前の処理をここに書く
		fmt.Println("開始")
		return nil // error を返すと処理全体が終了
	}

	app.After = func(c *cli.Context) error {
		// 終了時の処理をここに書く
		fmt.Println("終了")
		return nil
	}

	app.Run(os.Args)
}

func newAction(c *cli.Context) {

	// グローバルオプション
	// var isDry = c.GlobalBool("dryrun")
	// if isDry {
	// 	fmt.Println("this is dry-run")
	// }

	// パラメータ
	var projectPath = ""
	if len(c.Args()) <= 0 {
		fmt.Println("作成するプロジェクトPathを渡してください")
		return
	}
	projectPath = c.Args().First() // c.Args()[0] と同じ意味

	gosrc, err := gkgfiler.GetGoSrcPath()
	if err != nil {
		fmt.Println(err)
		return
	}

	projectPathWithGoRoot := gosrc + "/" + projectPath

	err = makeDir(projectPathWithGoRoot)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cloneGokigen(projectPathWithGoRoot)
	if err != nil {
		fmt.Println(err)
		return
	}

	replacePathInFiles(projectPathWithGoRoot, projectPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := os.RemoveAll(projectPathWithGoRoot + "/.git"); err != nil {
		fmt.Println(err)
		return
	}

	if err := overWriteLience(projectPathWithGoRoot, projectPath); err != nil {
		fmt.Println(err)
		return
	}
}

func makeDir(path string) error {
	fmt.Println("make " + path)
	if gkgfiler.Exist(path) {
		return fmt.Errorf("failed mkdir. dir has been exist : %v", path)
	}
	err := os.MkdirAll(path, 0777)
	if err != nil {
		printlnf("failed mkdir path : %v", path)
	}
	return err
}

func cloneGokigen(path string) error {
	printlnf("cloning gokigen to : %v ...", path)

	err := exec.Command("git", "clone", gokigenRepo, path).Run()
	if err != nil {
		printlnf("failed clone : %v", path)
		return err
	}
	return err
}

func printlnf(s string, param interface{}) {
	fmt.Println(fmt.Sprintf(s, param))
}

func replacePathInFiles(path, projectPath string) error {
	fileNameList, e := gkgfiler.GetPathsRecurcive(path, false, "*")
	if e != nil {
		return e
	}

	for _, file := range fileNameList {
		e = gkgfiler.ReplaceText(file, "github.com/sasasaiki/gokigen", projectPath, 0777)
	}
	if e != nil {
		return e
	}

	return nil
}

func overWriteLience(path, project string) error {
	e := gkgfiler.WriteText(path+"/README.md", "# "+project, false, 0777)
	if e != nil {
		return e
	}

	return nil
}
