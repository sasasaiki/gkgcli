package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gokigen-cli"
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

	projectPathWithGoRoot := getGoSrcRoot() + "/" + projectPath

	err := makeDir(projectPathWithGoRoot)
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
}

func getGoSrcRoot() string {
	// lookup go path
	gopath := build.Default.GOPATH
	if gopath == "" {
		fmt.Println("GOPATH が設定されていません.bash_profileなど$GOPATHを設定してください")
	}
	//  取得した$GOPATHが:つなぎなどで複数設定されていたら一番先頭を使う。
	srcRoot := filepath.Join(filepath.SplitList(gopath)[0], "src")
	return srcRoot
}

func makeDir(path string) error {
	fmt.Println("make " + path)
	if existFile(path) {
		return fmt.Errorf("failed mkdir. dir has been exist : %v", path)
	}
	err := exec.Command("mkdir", path).Run()
	if err != nil {
		printlnf("failed mkdir path : %v", path)
	}
	return err
}

func cloneGokigen(path string) error {
	printlnf("cloning gokigen to : %v ...", path)
	const repo = "https://github.com/sasasaiki/gokigen.git"
	err := exec.Command("git", "clone", repo, path).Run()
	if err != nil {
		printlnf("failed clone : %v", path)
		return err
	}
	return err
}

func existFile(path string) bool {
	_, e := os.Stat(path)
	if e != nil {
		return false
	}
	return true
}

func printlnf(s string, param interface{}) {
	fmt.Println(fmt.Sprintf(s, param))
}

func replacePathInFiles(path, projectPath string) {
	fileNameList := dirwalk(path, []string{})
	for _, file := range fileNameList {
		fmt.Println(file)
		replacePathInFile(file, projectPath)
	}
}

// 対象のファイルを取得する
func getFileNames(dir string) (matches []string) {
	files, _ := filepath.Glob(dir + "/*.go")
	files2, _ := filepath.Glob(dir + "/*.yaml")

	return append(files, files2...)
}

func dirwalk(dir string, paths []string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			paths = dirwalk(filepath.Join(dir, file.Name()), paths)
		}
	}

	f := getFileNames(filepath.Join(dir))
	return append(paths, f...)
}

// 書き込み処理を行う
func replacePathInFile(filename, projectPath string) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	const origin = "github.com/sasasaiki/gokigen"
	fmt.Println(projectPath)
	for i, line := range lines {
		fmt.Println(line)
		lines[i] = strings.Replace(line, origin, projectPath, -1)
		fmt.Println(lines[i])
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
