package main

import (
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("CLI Interface")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter command...")

	content := container.NewVBox(
		widget.NewLabel("Enter CLI command:"),
		input,
		widget.NewButton("Execute", func() {
			go executeCLICommand(input.Text)
		}),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.ShowAndRun()
}

func executeCLICommand(command string) {
	cmd := exec.Command("your-cli-program", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(string(output))
}
