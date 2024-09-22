package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"your-module-name/invoice"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Invoice Generator")

	inv := invoice.DefaultInvoice()

	idEntry := widget.NewEntry()
	idEntry.SetText(inv.Id)
	titleEntry := widget.NewEntry()
	titleEntry.SetText(inv.Title)
	fromEntry := widget.NewEntry()
	fromEntry.SetText(inv.From)
	toEntry := widget.NewEntry()
	toEntry.SetText(inv.To)
	itemEntry := widget.NewEntry()
	itemEntry.SetText(strings.Join(inv.Items, ", "))
	outputEntry := widget.NewEntry()
	outputEntry.SetText("invoice.pdf")

	content := container.NewVBox(
		widget.NewLabel("Invoice Generator"),
		widget.NewLabel("ID:"),
		idEntry,
		widget.NewLabel("Title:"),
		titleEntry,
		widget.NewLabel("From:"),
		fromEntry,
		widget.NewLabel("To:"),
		toEntry,
		widget.NewLabel("Items (comma-separated):"),
		itemEntry,
		widget.NewLabel("Output filename:"),
		outputEntry,
		widget.NewButton("Generate Invoice", func() {
			inv.Id = idEntry.Text
			inv.Title = titleEntry.Text
			inv.From = fromEntry.Text
			inv.To = toEntry.Text
			inv.Items = strings.Split(itemEntry.Text, ",")
			for i := range inv.Items {
				inv.Items[i] = strings.TrimSpace(inv.Items[i])
			}
			output := outputEntry.Text
			if !strings.HasSuffix(output, ".pdf") {
				output += ".pdf"
			}
			err := invoice.GenerateInvoice(inv, output)
			if err != nil {
				fmt.Println("Error generating invoice:", err)
			} else {
				fmt.Println("Invoice generated successfully!")
			}
		}),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 500))
	myWindow.ShowAndRun()
}
