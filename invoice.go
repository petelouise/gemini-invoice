package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/signintech/gopdf"
)

type Invoice struct {
	Id         string    `json:"id" yaml:"id"`
	Title      string    `json:"title" yaml:"title"`
	Logo       string    `json:"logo" yaml:"logo"`
	From       string    `json:"from" yaml:"from"`
	To         string    `json:"to" yaml:"to"`
	Date       string    `json:"date" yaml:"date"`
	Due        string    `json:"due" yaml:"due"`
	Items      []string  `json:"items" yaml:"items"`
	Quantities []int     `json:"quantities" yaml:"quantities"`
	Rates      []float64 `json:"rates" yaml:"rates"`
	Tax        float64   `json:"tax" yaml:"tax"`
	Discount   float64   `json:"discount" yaml:"discount"`
	Currency   string    `json:"currency" yaml:"currency"`
	Note       string    `json:"note" yaml:"note"`
}

func DefaultInvoice() Invoice {
	return Invoice{
		Id:         time.Now().Format("20060102"),
		Title:      "INVOICE",
		Rates:      []float64{25},
		Quantities: []int{2},
		Items:      []string{"Paper Cranes"},
		From:       "Project Folded, Inc.",
		To:         "Untitled Corporation, Inc.",
		Date:       time.Now().Format("Jan 02, 2006"),
		Due:        time.Now().AddDate(0, 0, 14).Format("Jan 02, 2006"),
		Tax:        0,
		Discount:   0,
		Currency:   "USD",
	}
}

var file = Invoice{}

func GenerateInvoice(invoice Invoice, output string) error {
	pdf := gopdf.GoPdf{}

	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	pdf.SetMargins(40, 40, 40, 40)
	pdf.AddPage()
	err := pdf.AddTTFFontData("Inter", interFont)
	if err != nil {
		return err
	}

	err = pdf.AddTTFFontData("Inter-Bold", interBoldFont)
	if err != nil {
		return err
	}

	writeLogo(&pdf, file.Logo, file.From)
	writeTitle(&pdf, file.Title, file.Id, file.Date)
	writeBillTo(&pdf, file.To)
	writeHeaderRow(&pdf)
	subtotal := 0.0
	for i := range file.Items {
		q := 1
		if len(file.Quantities) > i {
			q = file.Quantities[i]
		}

		r := 0.0
		if len(file.Rates) > i {
			r = file.Rates[i]
		}

		writeRow(&pdf, file.Items[i], q, r)
		subtotal += float64(q) * r
	}
	if file.Note != "" {
		writeNotes(&pdf, file.Note)
	}
	writeTotals(&pdf, subtotal, subtotal*file.Tax, subtotal*file.Discount)
	if file.Due != "" {
		writeDueDate(&pdf, file.Due)
	}
	writeFooter(&pdf, file.Id)
	output = strings.TrimSuffix(output, ".pdf") + ".pdf"
	err = pdf.WritePdf(output)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %s\n", output)

	return nil
}
