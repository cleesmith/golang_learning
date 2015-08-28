package main

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

func main() {
	excelFileName := "MyXLSXFile.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf(err.Error())
	}
	// for _, sheet := range xlFile.Sheets {
	// 	for _, row := range sheet.Rows {
	// 		for _, cell := range row.Cells {
	// 			fmt.Printf("%s\n", cell.String())
	// 		}
	// 	}
	// }

	sheet1 := xlFile.Sheets[0]
	fmt.Printf("sheet1=%v\n", sheet1)
	row1 := sheet1.Rows[0]
	fmt.Printf("row1=%v\n", row1)
	cell1 := row1.Cells[0]
	fmt.Printf("cell1=%v\n", cell1)

	// var sheet *xlsx.Sheet
	// var row *xlsx.Row
	// var cell *xlsx.Cell
	// var err error

	// sheet = file.AddSheet("Sheet1")
	row := sheet1.AddRow()
	cell := row.AddCell()
	cell.Value = "row2"
	err = xlFile.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}

}
