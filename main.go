package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type PackageInfo struct {
	PACKAGE_NAME  string
	GolangVersion string
	PathOutput    string
}
type GenModelInfo struct {
	RootProjectPath string
	ModelDir        string
	DSN             string
}
type Template struct {
	TemplatePath   string
	OutPath        string
	IsTemplateFile bool
}

func BundleTempleFile(TemplateInfo Template, data PackageInfo) ([]byte, error) {
	tmpl, err := template.ParseFiles(TemplateInfo.TemplatePath)
	if err != nil {
		log.Fatalf("Can not ParseFiles:%+v\n", err)
		return nil, err
	}
	var process bytes.Buffer
	err = tmpl.Execute(&process, data)
	if err != nil {
		log.Fatalf("Can not Execute main\n")
		return nil, err
	}
	if strings.Contains(TemplateInfo.OutPath, ".go") {
		formatt, err := format.Source(process.Bytes())
		if err != nil {
			log.Fatalf("Can not Execute format:%v\n", err)
			return nil, err
		}
		return formatt, nil
	} else {
		return process.Bytes(), nil
	}

}
func InitProject(data PackageInfo) {
	pathTemplate := "./template"
	fullTemplatePath, errAbs := filepath.Abs(pathTemplate)
	if errAbs != nil {
		log.Fatalln(errAbs)
	}
	rootTemplateDir := filepath.Dir(fullTemplatePath)
	fileTemplates := []Template{}
	err := filepath.Walk(pathTemplate, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			isTemplete := false
			if strings.Contains(path, ".tmpl") {
				isTemplete = true
			}
			fileOutName := strings.ReplaceAll(path, ".tmpl", "")
			fileOutName = strings.ReplaceAll(fileOutName, ".txt", "")
			fileOutName = strings.ReplaceAll(fileOutName, "template", "")
			fileTemplates = append(fileTemplates, Template{
				TemplatePath:   filepath.Join(rootTemplateDir, path),
				OutPath:        filepath.Join(data.PathOutput, fileOutName),
				IsTemplateFile: isTemplete,
			})
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range fileTemplates {
		fmt.Printf("process : %+v\n", v.TemplatePath)
		var fileData []byte
		if v.IsTemplateFile {
			b, errR := BundleTempleFile(v, data)
			if errR != nil {
				log.Fatalf("Error BundleTempleFile: %+v", errR)
			}
			fileData = b
		} else {
			b, errR := ioutil.ReadFile(v.TemplatePath)
			if errR != nil {
				log.Fatalf("Error ReadFile: %+v", errR)
			}
			fileData = b
		}
		WriteToFile(fileData, v.OutPath)
	}
}
func WriteToFile(fileData []byte, outpath string) {
	var permission fs.FileMode = 0700
	recusiveDir := filepath.Dir(outpath)
	if _, err := os.Stat(recusiveDir); os.IsNotExist(err) {
		err := os.MkdirAll(recusiveDir, permission)
		if err != nil {
			log.Fatalf("Error Mkdir :%+v\n", err)
		}
	}
	errW := ioutil.WriteFile(outpath, fileData, permission)
	if errW != nil {
		log.Fatalf("Error WriteFile :%+v\n", errW)
	}
	fmt.Printf("generate: %s\n", outpath)
}
func RootProject() {
	var outputPath string = "../out"
	var pInfo PackageInfo = PackageInfo{
		GolangVersion: "1.18",
	}
	fmt.Println("===================== This program is generate golang templete project =====================")
	fmt.Print("Golang package name: ")
	fmt.Scanln(&pInfo.PACKAGE_NAME)
	fmt.Print("Golang Version: ")
	fmt.Scanln(&pInfo.GolangVersion)
	fmt.Print("Output Path: ")
	fmt.Scanln(&outputPath)
	fullOutputPath, errAbsOut := filepath.Abs(outputPath)
	if errAbsOut != nil {
		log.Fatalln(errAbsOut)
	}
	pInfo.PathOutput = fullOutputPath
	fmt.Printf("PackageInfo : %+v", pInfo)
	InitProject(pInfo)
}
func GenerateModel() {
	var outputPath string = "../out"
	var modelDirPath string = "/internal/v1/models"
	var dsn string
	var pInfo GenModelInfo = GenModelInfo{
		RootProjectPath: outputPath,
		ModelDir:        modelDirPath,
	}

	fmt.Println("===================== This program is generate golang model =====================")
	fmt.Print("Root Project Path: ")
	fmt.Scanln(&outputPath)
	fmt.Print("Model Dir Path (ex. /internal/v1/models):")
	fmt.Scanln(&modelDirPath)
	if modelDirPath == "" {
		pInfo.ModelDir = "/internal/v1/models"
	} else {
		pInfo.ModelDir = modelDirPath
	}
	fmt.Print("DSN: ")
	fmt.Scanln(&dsn)
	if dsn == "" {
		pInfo.DSN = os.Getenv("DSN")
	} else {
		pInfo.DSN = dsn
	}
	fullOutputPath, errAbsOut := filepath.Abs(outputPath)
	if errAbsOut != nil {
		log.Fatalln(errAbsOut)
	}
	pInfo.RootProjectPath = fullOutputPath
	fmt.Printf("PackageInfo : %+v", pInfo)
	db := InitDB(pInfo.DSN)
	selectTable := SelectTable(db)
	for _, v := range selectTable {
		GeneratorModel(pInfo, db, v)
	}
}
func main() {
	cmd := ""
	fmt.Println("=========== This command is generate golang project ================")
	fmt.Print("Select Command (root,model) : ")
	fmt.Scanln(&cmd)
	switch cmd {
	case "root":
		RootProject()
	case "model":
		GenerateModel()
	}

}
