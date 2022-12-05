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

// structure กำหนดชุดตัวแปรที่จะใช้ใน template
type PackageInfo struct {
	PACKAGE_NAME  string
	GolangVersion string
	PathOutput    string
}

// structure กำหนดชุดตัวแปรที่จะใช้ในการสร้าง model จาก database
type GenModelInfo struct {
	RootProjectPath string
	ModelDir        string
	DSN             string
}

// structure กำหนดชุดตัวแปรที่ใช้เก็บข้อมูลของ template
type Template struct {
	TemplatePath   string
	OutPath        string
	IsTemplateFile bool
}

// function สำหรับการอ่านและใส่ variable/logic ไปใน template
func BundleTempleFile(TemplateInfo Template, data PackageInfo) ([]byte, error) {
	// อ่านไฟล์
	tmpl, err := template.ParseFiles(TemplateInfo.TemplatePath)
	if err != nil {
		log.Fatalf("Can not ParseFiles:%+v\n", err)
		return nil, err
	}
	var process bytes.Buffer
	// ทำการรวม variable/logic ไปใน template
	err = tmpl.Execute(&process, data)
	if err != nil {
		log.Fatalf("Can not Execute main\n")
		return nil, err
	}
	// เช็คว่าชื่อ template มีนามสกุล .go หรือไม่
	if strings.Contains(TemplateInfo.OutPath, ".go") {
		// ถ้าเป็นนามสกุล .go ให้ format ไฟล์เป็น golang แล้วแปลงข้อมูลทั้งหมดเป็น Bytes
		formatt, err := format.Source(process.Bytes())
		if err != nil {
			log.Fatalf("Can not Execute format:%v\n", err)
			return nil, err
		}
		return formatt, nil
	} else {
		// แปลงข้อมูลทั้งหมดเป็น Bytes
		return process.Bytes(), nil
	}
}

// function สำหรับเขียนข้อมูลลงไฟล์
func WriteToFile(fileData []byte, outpath string) {
	var permission fs.FileMode = 0700
	// หา parent directory ของไฟล์
	recusiveDir := filepath.Dir(outpath)
	// เช็คว่า parent directory มีอยู่หรือไม่
	if _, err := os.Stat(recusiveDir); os.IsNotExist(err) {
		// ถ้าไม่มี parent directory ให้สร้างแบบทั้งหมด
		err := os.MkdirAll(recusiveDir, permission)
		if err != nil {
			log.Fatalf("Error Mkdir :%+v\n", err)
		}
	}
	// เขียนจ้อมูลลงไฟล์
	errW := ioutil.WriteFile(outpath, fileData, permission)
	if errW != nil {
		log.Fatalf("Error WriteFile :%+v\n", errW)
	}
	fmt.Printf("generate: %s\n", outpath)
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
