package files

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
)

type SharedStrings struct {
	SI []struct {
		T string `xml:"t"`
	} `xml:"si"`
}

type Cell struct {
	R string `xml:"r,attr"`
	T string `xml:"t,attr,omitempty"`
	V string `xml:"v"`
}

type Row struct {
	C []Cell `xml:"c"`
}

type SheetData struct {
	Row []Row `xml:"row"`
}

type Worksheet struct {
	SheetData SheetData `xml:"sheetData"`
}

func ReadExcelData(filename string) {
	// if len(os.Args) < 2 {
	//     fmt.Println("Usage: go run main.go <xlsx file>")
	//     return
	// }

	// filename := os.Args[1]

	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer zipReader.Close()

	var sharedStrings SharedStrings
	var worksheet Worksheet

	for _, file := range zipReader.File {
		if file.Name == "xl/sharedStrings.xml" {
			f, err := file.Open()
			if err != nil {
				fmt.Println("Error opening sharedStrings.xml:", err)
				return
			}
			defer f.Close()
			data, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Println("Error reading sharedStrings.xml:", err)
				return
			}
			xml.Unmarshal(data, &sharedStrings)
		} else if file.Name == "xl/worksheets/sheet1.xml" {
			f, err := file.Open()
			if err != nil {
				fmt.Println("Error opening sheet1.xml:", err)
				return
			}
			defer f.Close()
			data, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Println("Error reading sheet1.xml:", err)
				return
			}
			xml.Unmarshal(data, &worksheet)
		}
	}

	dataMap := make(map[string]string)
	for _, row := range worksheet.SheetData.Row {
		for _, cell := range row.C {
			var value string
			if cell.T == "s" { // shared string
				index, _ := strconv.Atoi(cell.V)
				value = sharedStrings.SI[index].T
			} else {
				value = cell.V
			}
			dataMap[cell.R] = value
		}
	}

	fmt.Println("Data Map:")
	for k, v := range dataMap {
		fmt.Printf("%s: %s\n", k, v)
	}
}
