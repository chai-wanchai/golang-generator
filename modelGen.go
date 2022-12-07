package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"gorm.io/driver/mysql"
	"gorm.io/gen"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func render(tmpl string, wr io.Writer, data interface{}) error {
	t, err := template.New(tmpl).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}
func initGormGen(db *gorm.DB) *gen.Generator {
	g := gen.NewGenerator(gen.Config{
		// if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable: true,
		// if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		FieldCoverable: true,
		// if you want generate field with unsigned integer type, set FieldSignable true
		FieldSignable: true,
		// if you want to generate index tags from database, set FieldWithIndexTag true
		FieldWithIndexTag: true,
		// if you want to generate type tags from database, set FieldWithTypeTag true
		FieldWithTypeTag: true,
	})
	g.UseDB(db)
	return g
}
func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	return db
}
func tableShow(input []string, cols int) {
	maxWidth := 0
	for _, s := range input {
		if len(s) > maxWidth {
			maxWidth = len(s)
		}
	}
	format := fmt.Sprintf("%%d) %%-%ds%%s", maxWidth)
	rows := (len(input) + cols - 1) / cols
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			i := col*rows + row
			if i >= len(input) {
				break // This means the last column is not "full"
			}
			padding := ""
			if i < 9 {
				padding = " "
			}
			fmt.Printf(format, i, input[i], padding)
		}
		fmt.Println()
	}
}
func SelectTable(db *gorm.DB) []string {
	tableList, err := db.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("get all tables fail: %w", err))
	}
	fmt.Printf("\n===================================\n")
	tableShow(tableList, 5)
	fmt.Println("===================================")
	var selectStr string
	fmt.Print("Select Table to gen model (ex. 0,1,7,8,9) :")
	fmt.Scanln(&selectStr)
	var selectTable []string
	for _, v := range strings.Split(selectStr, ",") {
		index, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		selectTable = append(selectTable, tableList[index])
	}
	return selectTable
}

type ModelData struct {
	ModelInfo    interface{}
	DatabaseName string
}

func GeneratorModel(packageInfo GenModelInfo, db *gorm.DB, modelName string) {
	g := initGormGen(db)
	model := g.GenerateModel(modelName)
	model.StructInfo.Package = "models"
	tmpl, err := template.ParseFiles("./database/model.go.tmpl")
	if err != nil {
		log.Fatalf("Can not ParseFiles:%+v\n", err)
	}
	var process bytes.Buffer
	dbName := ""
	db.Raw("SELECT DATABASE()").Scan(&dbName)
	data := ModelData{
		ModelInfo:    model,
		DatabaseName: dbName,
	}
	err = tmpl.Execute(&process, data)
	if err != nil {
		log.Fatalf("Can not Execute main:%+v\n", err)
	}
	formatt, err := format.Source(process.Bytes())
	if err != nil {
		log.Fatalf("Can not Execute format:%v\n", err)
	}
	modelFile := fmt.Sprintf("%s%s/%s.go", packageInfo.RootProjectPath, packageInfo.ModelDir, model.FileName)
	WriteToFile(formatt, modelFile)
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
	db := InitDB(pInfo.DSN)
	selectTable := SelectTable(db)
	for _, v := range selectTable {
		GeneratorModel(pInfo, db, v)
	}
}
