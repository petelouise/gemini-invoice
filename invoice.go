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

	writeLogo(&pdf, invoice.Logo, invoice.From)
	writeTitle(&pdf, invoice.Title, invoice.Id, invoice.Date)
	writeBillTo(&pdf, invoice.To)
	writeHeaderRow(&pdf)
	subtotal := 0.0
	for i := range invoice.Items {
		q := 1
		if len(invoice.Quantities) > i {
			q = invoice.Quantities[i]
		}

		r := 0.0
		if len(invoice.Rates) > i {
			r = invoice.Rates[i]
		}

		writeRow(&pdf, invoice, invoice.Items[i], q, r)
		subtotal += float64(q) * r
	}
	if invoice.Note != "" {
		writeNotes(&pdf, invoice.Note)
	}
	writeTotals(&pdf, invoice, subtotal, subtotal*invoice.Tax, subtotal*invoice.Discount)
	if invoice.Due != "" {
		writeDueDate(&pdf, invoice.Due)
	}
	writeFooter(&pdf, invoice.Id)
	output = strings.TrimSuffix(output, ".pdf") + ".pdf"
	err = pdf.WritePdf(output)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %s\n", output)

	return nil
}
