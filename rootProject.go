package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// function สำหรับเริ่มต้นการ build template
func InitProject(data PackageInfo) {
	// กำหนด path ที่เก็บ template เอาไว้
	pathTemplate := "./template"
	fullTemplatePath, errAbs := filepath.Abs(pathTemplate)
	if errAbs != nil {
		log.Fatalln(errAbs)
	}
	// กำหนด root path ของ project
	rootTemplateDir := filepath.Dir(fullTemplatePath)
	fileTemplates := []Template{}
	// list file ทั้งหมดใน template path
	err := filepath.Walk(pathTemplate, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		// ถ้า path นั้นไม่ใช่ directory
		if !info.IsDir() {
			isTemplete := false
			// ถ้าไฟล์นามสกุล .tmpl กำหนดว่าเป็น Templete ที่ต้องใช้ variable/logic
			if strings.Contains(path, ".tmpl") {
				isTemplete = true
			}
			// แทนที่นามสกุลไฟล์ที่ไม่ต้องการออก แล้วเก็บใส่ array
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
	// loop ไฟล์ทั้งหมดใน template
	for _, v := range fileTemplates {
		fmt.Printf("process : %+v\n", v.TemplatePath)
		var fileData []byte
		// ถ้าใช่ไฟล์ template ให้ทำการเรียก function BundleTempleFile
		if v.IsTemplateFile {
			b, errR := BundleTempleFile(v, data)
			if errR != nil {
				log.Fatalf("Error BundleTempleFile: %+v", errR)
			}
			fileData = b
		} else {
			// ถ้าเป็นไฟล์ธรรมดาให้อ่านไฟล์แล้วส่งกลับเป็น byte
			b, errR := ioutil.ReadFile(v.TemplatePath)
			if errR != nil {
				log.Fatalf("Error ReadFile: %+v", errR)
			}
			fileData = b
		}
		// เขียนลงไฟล์
		WriteToFile(fileData, v.OutPath)
	}
}

// function สำหรับรับ input ข้อมูลต่างๆ และสร้าง project จาก template
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
