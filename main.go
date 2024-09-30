package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/signintech/gopdf"
	"gopkg.in/yaml.v2"
)

type Item struct {
	Name     string
	Quantity int
	Price    float64
}

var currencySymbols = map[string]string{
	"USD": "$",
	// Add other currency symbols as needed
}

//go:embed "fonts/Inter.ttf"
var interFont []byte

//go:embed "fonts/Inter-Bold.ttf"
var interBoldFont []byte

type Invoice struct {
	Id                  string    `json:"id" yaml:"id"`
	Title               string    `json:"title" yaml:"title"`
	Logo                string    `json:"logo" yaml:"logo"`
	From                string    `json:"from" yaml:"from"`
	To                  string    `json:"to" yaml:"to"`
	ToAddress           string    `json:"to_address" yaml:"to_address"`
	Date                string    `json:"date" yaml:"date"`
	Due                 string    `json:"due" yaml:"due"`
	Items               []string  `json:"items" yaml:"items"`
	Quantities          []int     `json:"quantities" yaml:"quantities"`
	Rates               []float64 `json:"rates" yaml:"rates"`
	Tax                 float64   `json:"tax" yaml:"tax"`
	Discount            float64   `json:"discount" yaml:"discount"`
	Currency            string    `json:"currency" yaml:"currency"`
	Note                string    `json:"note" yaml:"note"`
	AccountNumber       string    `json:"account_number" yaml:"account_number"`
	RoutingNumber       string    `json:"routing_number" yaml:"routing_number"`
	PaymentInstructions string    `json:"payment_instructions" yaml:"payment_instructions"`
}

type Config struct {
	Title               string `yaml:"title"`
	From                string `yaml:"from"`
	Logo                string `yaml:"logo"`
	AccountNumber       string `yaml:"account_number"`
	RoutingNumber       string `yaml:"routing_number"`
	PaymentInstructions string `yaml:"payment_instructions"`
}

func LoadConfig(filename string) (*Config, error) {
	// Try to load from current directory first
	if data, err := ioutil.ReadFile(filename); err == nil {
		var config Config
		if err := yaml.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	}

	// If not found in current directory, try other locations
	locations := []string{
		filename,
		filepath.Join(".", filename),
		filepath.Join("..", filename),
		filepath.Join("..", "Resources", filename),
	}

	execPath, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(execPath)
		locations = append(locations,
			filepath.Join(dir, filename),
			filepath.Join(dir, "..", filename),
			filepath.Join(dir, "..", "Resources", filename),
		)
	}

	for _, path := range locations {
		data, err := ioutil.ReadFile(path)
		if err == nil {
			var config Config
			if err := yaml.Unmarshal(data, &config); err == nil {
				fmt.Printf("Loaded config from: %s\n", path)
				return &config, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find or parse %s in any of the searched locations", filename)
}

func DefaultInvoice(config *Config) Invoice {
	return Invoice{
		Id:                  time.Now().Format("20060102"),
		Title:               config.Title,
		Logo:                config.Logo,
		Rates:               []float64{0},
		Quantities:          []int{1},
		Items:               []string{""},
		From:                config.From,
		To:                  "",
		Date:                time.Now().Format("Jan 02, 2006"),
		Due:                 time.Now().AddDate(0, 0, 14).Format("Jan 02, 2006"),
		Tax:                 0,
		Discount:            0,
		Currency:            "USD",
		AccountNumber:       config.AccountNumber,
		RoutingNumber:       config.RoutingNumber,
		PaymentInstructions: config.PaymentInstructions,
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

	WriteLogo(&pdf, invoice.Logo, invoice.From)
	WriteTitle(&pdf, invoice.Title, invoice.Id, invoice.Date)
	WriteBillTo(&pdf, invoice.To, invoice.ToAddress)
	WriteHeaderRow(&pdf)
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

		WriteRow(&pdf, invoice.Items[i], q, r, invoice.Currency)
		subtotal += float64(q) * r
	}

	if invoice.Note != "" {
		WriteNotes(&pdf, invoice.Note)
	}
	WriteTotals(&pdf, subtotal, subtotal*invoice.Tax, subtotal*invoice.Discount, invoice.Currency)
	if invoice.Due != "" {
		WriteDueDate(&pdf, invoice.Due)
	}
	WritePaymentInstructions(&pdf, invoice.PaymentInstructions, invoice.AccountNumber, invoice.RoutingNumber)
	WriteFooter(&pdf, invoice.Id)

	err = pdf.WritePdf(output)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %s\n", output)

	return nil
}

func createSampleInvoice() Invoice {
	return Invoice{
		Id:                  "SAMPLE-001",
		Title:               "Sample Invoice",
		Logo:                "",
		From:                "Your Company\n123 Your Street\nYour City, State 12345",
		To:                  "Sample Customer",
		ToAddress:           "456 Customer Street\nCustomer City, State 67890",
		Date:                time.Now().Format("Jan 02, 2006"),
		Due:                 time.Now().AddDate(0, 0, 30).Format("Jan 02, 2006"),
		Items:               []string{"Item 1", "Item 2", "Item 3"},
		Quantities:          []int{2, 1, 3},
		Rates:               []float64{100.00, 200.00, 50.00},
		Tax:                 0.08,
		Discount:            0.05,
		Currency:            "USD",
		Note:                "Thank you for your business!",
		AccountNumber:       "1234567890",
		RoutingNumber:       "987654321",
		PaymentInstructions: "Please make payment within 30 days.",
	}
}

func main() {
	testFlag := flag.Bool("test", false, "Run a quick test of PDF generation")
	flag.Parse()

	if *testFlag {
		sampleInvoice := createSampleInvoice()
		err := GenerateInvoice(sampleInvoice, "sample_invoice.pdf")
		if err != nil {
			fmt.Printf("Error generating sample invoice: %v\n", err)
		} else {
			fmt.Println("Sample invoice generated: sample_invoice.pdf")
		}
		return
	}

	config, err := LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Println("Current working directory:", getCurrentDirectory())
		return
	}

	myApp := app.NewWithID("com.gemini.invoice")
	myApp.Settings().SetTheme(NewPinkTheme())
	myWindow := myApp.NewWindow("Gemini Invoice")

	inv := DefaultInvoice(config)

	idEntry := widget.NewEntry()
	idEntry.SetText(inv.Id)
	toEntry := widget.NewEntry()
	toEntry.SetPlaceHolder("Customer Name")
	toAddressEntry := widget.NewMultiLineEntry()
	toAddressEntry.SetPlaceHolder("Customer Address")

	items := []Item{{Name: "", Quantity: 1, Price: 0.0}}
	
	itemsContainer := container.NewVBox()

	var updateItemsContainer func()
	updateItemsContainer = func() {
		itemsContainer.Objects = nil
		for i := range items {
			index := i // Capture the index in a local variable
			nameEntry := widget.NewEntry()
			nameEntry.SetText(items[i].Name)
			nameEntry.OnChanged = func(value string) {
				items[index].Name = value
			}

			quantityEntry := widget.NewEntry()
			quantityEntry.SetText(strconv.Itoa(items[i].Quantity))
			quantityEntry.OnChanged = func(value string) {
				quantity, _ := strconv.Atoi(value)
				items[index].Quantity = quantity
			}

			priceEntry := widget.NewEntry()
			priceEntry.SetText(fmt.Sprintf("$%.2f", items[i].Price))
			priceEntry.OnChanged = func(value string) {
				// Remove the dollar sign and any commas
				value = strings.TrimPrefix(value, "$")
				value = strings.ReplaceAll(value, ",", "")
				price, _ := strconv.ParseFloat(value, 64)
				items[index].Price = price
				// Update the display with proper formatting
				priceEntry.SetText(fmt.Sprintf("$%.2f", price))
			}

			removeButton := widget.NewButton("Remove", func() {
				items = append(items[:index], items[index+1:]...)
				updateItemsContainer()
			})

			itemContainer := container.NewVBox(
				nameEntry,
				container.NewGridWithColumns(2,
					quantityEntry,
					priceEntry,
				),
				removeButton,
			)

			itemsContainer.Add(itemContainer)
		}
		itemsContainer.Refresh()
	}

	updateItemsContainer()

	addItemButton := widget.NewButton("Add Item", func() {
		items = append(items, Item{Name: "", Quantity: 1, Price: 0.0})
		updateItemsContainer()
	})

	var outputDir string
	outputDirButton := widget.NewButton("Select Output Directory", nil)

	title := canvas.NewText("Gemini Invoice", color.NRGBA{R: 219, G: 112, B: 147, A: 255})
	title.TextSize = 24
	title.Alignment = fyne.TextAlignCenter

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: idEntry},
			{Text: "To", Widget: toEntry},
			{Text: "Address", Widget: toAddressEntry},
			{Text: "Items", Widget: container.NewVBox(itemsContainer, addItemButton)},
			{Text: "Output Directory", Widget: outputDirButton},
		},
	}

	generateButton := widget.NewButton("Generate Invoice", nil)
	generateButton.Importance = widget.HighImportance

	outputDirButton.OnTapped = func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println("Error selecting directory:", err)
				return
			}
			if uri == nil {
				return
			}
			outputDir = uri.Path()
			outputDirButton.SetText(outputDir)
			myApp.Preferences().SetString("lastOutputDir", outputDir)
		}, myWindow)
	}

	generateButton.OnTapped = func() {
		if outputDir == "" {
			dialog.ShowInformation("Error", "Please select an output directory", myWindow)
			return
		}

		inv.Id = idEntry.Text
		inv.To = toEntry.Text
		inv.ToAddress = toAddressEntry.Text
		inv.Items = make([]string, 0, len(items))
		inv.Quantities = make([]int, 0, len(items))
		inv.Rates = make([]float64, 0, len(items))

		for _, item := range items {
			if item.Name != "" {
				inv.Items = append(inv.Items, item.Name)
				inv.Quantities = append(inv.Quantities, item.Quantity)
				inv.Rates = append(inv.Rates, item.Price)
			}
		}

		// Ensure the currency is set
		inv.Currency = "USD"

		baseFilename := "invoice"
		extension := ".pdf"
		counter := 1
		var output string
		for {
			if counter == 1 {
				output = filepath.Join(outputDir, baseFilename+extension)
			} else {
				output = filepath.Join(outputDir, fmt.Sprintf("%s_%d%s", baseFilename, counter, extension))
			}
			if _, err := os.Stat(output); os.IsNotExist(err) {
				break
			}
			counter++
		}

		err := GenerateInvoice(inv, output)
		if err != nil {
			dialog.ShowError(fmt.Errorf("error generating invoice: %v", err), myWindow)
		} else {
			dialog.ShowInformation("Success", fmt.Sprintf("Invoice generated successfully: %s", output), myWindow)
		}
	}

	content := container.NewVBox(
		title,
		form,
		generateButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 700))
	myWindow.ShowAndRun()
}

type PinkTheme struct{}

var _ fyne.Theme = (*PinkTheme)(nil)

func (m PinkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 255, G: 240, B: 245, A: 255} // Light pink background
	case theme.ColorNameButton:
		return color.NRGBA{R: 219, G: 112, B: 147, A: 255} // Medium pink for buttons
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 139, G: 0, B: 139, A: 255} // Dark pink for text
	case theme.ColorNameHover:
		return color.NRGBA{R: 255, G: 182, B: 193, A: 255} // Light pink for hover
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 199, G: 21, B: 133, A: 255} // Medium violet red for placeholders
	case theme.ColorNamePressed:
		return color.NRGBA{R: 199, G: 21, B: 133, A: 255} // Medium violet red for pressed state
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 219, G: 112, B: 147, A: 255} // Medium pink as primary color
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 255, G: 182, B: 193, A: 255} // Light pink for scrollbar
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m PinkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m PinkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m PinkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameText:
		return 14
	case theme.SizeNameInputBorder:
		return 2
	default:
		return theme.DefaultTheme().Size(name)
	}
}

func NewPinkTheme() fyne.Theme {
	return &PinkTheme{}
}

func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return "Unable to get current directory"
	}
	return dir
}
