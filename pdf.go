package main

import (
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

const (
	quantityColumnOffset = 360
	rateColumnOffset     = 405
	amountColumnOffset   = 480
)

const (
	subtotalLabel = "Subtotal"
	discountLabel = "Discount"
	taxLabel      = "Tax"
	totalLabel    = "Total"
)

func WriteLogo(pdf *gopdf.GoPdf, logo string, from string) {
	startY := pdf.GetY()
	
	// Write 'from' information first
	pdf.SetTextColor(55, 55, 55)
	formattedFrom := strings.ReplaceAll(from, `\n`, "\n")
	fromLines := strings.Split(formattedFrom, "\n")

	for i := 0; i < len(fromLines); i++ {
		if i == 0 {
			_ = pdf.SetFont("Inter", "", 12)
			_ = pdf.Cell(nil, fromLines[i])
			pdf.Br(18)
		} else {
			_ = pdf.SetFont("Inter", "", 10)
			_ = pdf.Cell(nil, fromLines[i])
			pdf.Br(15)
		}
	}

	// Add logo to the right side
	if logo != "" {
		width, height, err := getImageDimension(logo)
		if err != nil {
			fmt.Printf("Warning: Unable to get image dimensions for %s: %v\n", logo, err)
		} else {
			scaledWidth := 100.0
			scaledHeight := float64(height) * scaledWidth / float64(width)
			pageWidth := gopdf.PageSizeA4.W
			logoX := pageWidth - scaledWidth - 40 // 40 is right margin
			err = pdf.Image(logo, logoX, startY, &gopdf.Rect{W: scaledWidth, H: scaledHeight})
			if err != nil {
				fmt.Printf("Warning: Unable to add logo to PDF: %v\n", err)
			}
		}
	}

	pdf.Br(21)
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX(), pdf.GetY(), 260, pdf.GetY())
	pdf.Br(36)
}

func WriteTitle(pdf *gopdf.GoPdf, title, id, date string) {
	_ = pdf.SetFont("Inter-Bold", "", 24)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.Cell(nil, title)
	pdf.Br(36)
	_ = pdf.SetFont("Inter", "", 12)
	pdf.SetTextColor(100, 100, 100)
	_ = pdf.Cell(nil, "#")
	_ = pdf.Cell(nil, id)
	pdf.SetTextColor(150, 150, 150)
	_ = pdf.Cell(nil, "  ·  ")
	pdf.SetTextColor(100, 100, 100)
	_ = pdf.Cell(nil, date)
	pdf.Br(48)
}

func WriteDueDate(pdf *gopdf.GoPdf, due string) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, "Due Date")
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.SetFontSize(11)
	pdf.SetX(amountColumnOffset - 15)
	_ = pdf.Cell(nil, due)
	pdf.Br(12)
}

func WriteBillTo(pdf *gopdf.GoPdf, to string, toAddress string) {
	pdf.SetTextColor(75, 75, 75)
	_ = pdf.SetFont("Inter", "", 9)
	_ = pdf.Cell(nil, "BILL TO")
	pdf.Br(18)
	pdf.SetTextColor(75, 75, 75)

	_ = pdf.SetFont("Inter", "", 15)
	_ = pdf.Cell(nil, to)
	pdf.Br(20)

	_ = pdf.SetFont("Inter", "", 10)
	formattedAddress := strings.ReplaceAll(toAddress, `\n`, "\n")
	addressLines := strings.Split(formattedAddress, "\n")

	for _, line := range addressLines {
		_ = pdf.Cell(nil, line)
		pdf.Br(15)
	}
	pdf.Br(64)
}

func WriteHeaderRow(pdf *gopdf.GoPdf) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, "ITEM")
	pdf.SetX(quantityColumnOffset)
	_ = pdf.Cell(nil, "QTY")
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, "RATE")
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, "AMOUNT")
	pdf.Br(24)
}

func WriteNotes(pdf *gopdf.GoPdf, notes string) {
	pdf.SetY(pdf.GetY() + 20) // Add some space before notes

	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, "NOTES")
	pdf.Br(18)
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(0, 0, 0)

	formattedNotes := strings.ReplaceAll(notes, `\n`, "\n")
	notesLines := strings.Split(formattedNotes, "\n")

	for i := 0; i < len(notesLines); i++ {
		_ = pdf.Cell(nil, notesLines[i])
		pdf.Br(15)
	}

	pdf.Br(48)
}
func WriteFooter(pdf *gopdf.GoPdf, id string) {
	pageHeight := 841.89 // A4 height in points
	pdf.SetY(pageHeight - 40) // Position footer 40 points from bottom

	_ = pdf.SetFont("Inter", "", 10)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, id)
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX()+10, pdf.GetY()+6, 550, pdf.GetY()+6)
}

func WriteRow(pdf *gopdf.GoPdf, item string, quantity int, rate float64, currency string) {
	_ = pdf.SetFont("Inter", "", 11)
	pdf.SetTextColor(0, 0, 0)

	total := float64(quantity) * rate
	amount := strconv.FormatFloat(total, 'f', 2, 64)

	_ = pdf.Cell(nil, item)
	pdf.SetX(quantityColumnOffset)
	_ = pdf.Cell(nil, strconv.Itoa(quantity))
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, currencySymbols[currency]+strconv.FormatFloat(rate, 'f', 2, 64))
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, currencySymbols[currency]+amount)
	pdf.Br(24)
}

func WriteTotals(pdf *gopdf.GoPdf, subtotal float64, tax float64, discount float64, currency string) {
	pdf.SetY(pdf.GetY() + 20) // Add some space before totals

	WriteTotal(pdf, subtotalLabel, subtotal, currency)
	if tax > 0 {
		WriteTotal(pdf, taxLabel, tax, currency)
	}
	if discount > 0 {
		WriteTotal(pdf, discountLabel, discount, currency)
	}
	WriteTotal(pdf, totalLabel, subtotal+tax-discount, currency)
}

func WriteTotal(pdf *gopdf.GoPdf, label string, total float64, currency string) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, label)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.SetFontSize(12)
	pdf.SetX(amountColumnOffset - 15)
	if label == totalLabel {
		_ = pdf.SetFont("Inter-Bold", "", 11.5)
	}
	_ = pdf.Cell(nil, currencySymbols[currency]+strconv.FormatFloat(total, 'f', 2, 64))
	pdf.Br(24)
}

func getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return image.Width, image.Height, nil
}
func WritePaymentInstructions(pdf *gopdf.GoPdf, instructions, accountNumber, routingNumber string) {
	pdf.SetY(pdf.GetY() + 20) // Add some space before payment instructions

	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, "PAYMENT INSTRUCTIONS")
	pdf.Br(18)
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(0, 0, 0)

	_ = pdf.Cell(nil, instructions)
	pdf.Br(15)
	_ = pdf.Cell(nil, "Account Number: "+accountNumber)
	pdf.Br(15)
	_ = pdf.Cell(nil, "Routing Number: "+routingNumber)

	pdf.Br(48)
}
