package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/jung-kurt/gofpdf/internal/example"
)

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetHeaderFunc(func() {
		// pdf.Image(example.ImageFile("logo.png"), 10, 6, 30, 0, false, "", 0, "")
		pdf.Image("logo.png", 10, 6, 30, 0, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "B", 15)
		pdf.Cell(80, 0, "")
		pdf.CellFormat(30, 10, "Title", "1", 0, "C", false, 0, "")
		pdf.Ln(20)
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetFont("Times", "", 12)
	for j := 1; j <= 40; j++ {
		pdf.CellFormat(0, 10, fmt.Sprintf("Printing line number %d", j),
			"", 1, "", false, 0, "")
	}
	fileStr := "fpdf_addpage.pdf" // example.Filename("Fpdf_AddPage")
	err := pdf.OutputFileAndClose(fileStr)
	fmt.Println(err)
	// example.Summary(err, fileStr)
}
