package main

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.NewWithID("billingappv1.0.0")
	myWindow := myApp.NewWindow("Monthly Expenses Manager")

	billingData := &BillingData{}

	// Create fields for entering bill information
	labelEntry := widget.NewEntry()
	labelEntry.SetPlaceHolder("Expense Name")
	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Price")
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("Date (YYYY-MM-DD)")

	// Person input fields
	personNameEntry := widget.NewEntry()
	personNameEntry.SetPlaceHolder("Person's name")

	// Add a salary field for the person
	salaryEntry := widget.NewEntry()
	salaryEntry.SetPlaceHolder("Person Salary")

	// List of months
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	// Create a Select widget with the months
	monthSelect := widget.NewSelect(months, func(selected string) {})

	monthSelect.PlaceHolder = "Select a month"

	// Person data
	person := Person{}

	// Bills list
	var bills []Bill

	// Total label to display calculated amount
	totalLabel := widget.NewLabel("Total: $0.00")

	// Add Bill Button
	addBillButton := widget.NewButton("Add Bill", func() {
		// Convert price to float64
		price, err := strconv.ParseFloat(priceEntry.Text, 64)
		if err != nil {
			fmt.Println("Invalid price format")
			return
		}

		// Create a new bill
		bill := Bill{
			Label: labelEntry.Text,
			Price: price,
			Date:  dateEntry.Text,
		}
		bills = append(bills, bill)

		billingData.Bills = bills

		// Clear input fields
		labelEntry.SetText("")
		priceEntry.SetText("")
		dateEntry.SetText("")

		// Calculate total expenses and display
		total := calculateTotal(bills)

		// Refresh total value
		billingData.Total = total
		totalLabel.SetText(fmt.Sprintf("Total: R$%.2f", total))
	})

	// A label for displaying success messages
	successMessageLabel := canvas.NewText("", color.RGBA{0, 100, 0, 255})
	successMessageLabel.TextStyle = fyne.TextStyle{Bold: true}
	successMessageLabel.Hidden = true

	// Load data from JSON Button
	loadFromJSONButton := widget.NewButton("Load data from JSON", func() {
		// Create a file open dialog
		fileDialog := dialog.NewFileOpen(
			func(reader fyne.URIReadCloser, err error) {
				if err != nil || reader == nil {
					return
				}
				// Get the file name from the URI
				fileName := reader.URI().Name()
				var data, _ = LoadJSON(fileName)

				billingData = &data

				// Display success message
				successMessageLabel.Hidden = false
				successMessageLabel.Text = fmt.Sprintf("%s was loaded", fileName)
			}, myWindow)

		// Show only specific file types (optional)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fileDialog.Show()
	})

	// Save JSON Button
	saveJSONButton := widget.NewButton("Save to JSON", func() {
		// Set person data
		person.FullName = personNameEntry.Text
		salary, _ := strconv.ParseFloat(salaryEntry.Text, 64)
		person.Salary = salary

		// Create billing data
		billingData := BillingData{
			Person:    person,
			MonthName: monthSelect.Selected,
			Bills:     bills,
			Total:     calculateTotal(bills),
		}
		saveToJSON(billingData)
	})

	// Save CSV Button
	saveCSVButton := widget.NewButton("Save to CSV", func() {
		// Set person salary
		salary, _ := strconv.ParseFloat(salaryEntry.Text, 64)
		person.Salary = salary

		// Create billing data
		billingData := BillingData{
			Person:    person,
			MonthName: monthSelect.Selected,
			Bills:     bills,
			Total:     calculateTotal(bills),
		}
		saveToCSV(billingData)
	})

	// List all expenses
	listExpensesButton := widget.NewButton("Show detailed expenses", func() {
		// Set person data from entries
		if personNameEntry.Text != "" {
			billingData.Person.FullName = personNameEntry.Text
		}
		if salaryEntry.Text != "" {
			salary, _ := strconv.ParseFloat(salaryEntry.Text, 64)
			billingData.Person.Salary = salary
		}
		if monthSelect.Selected != "" {

			billingData.MonthName = monthSelect.Selected
		}

		listExpenses(myApp, billingData)
	})

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Enter Expense Information:"),
		labelEntry, priceEntry, dateEntry,
		addBillButton,
		widget.NewLabel("Enter Person Information:"),
		personNameEntry,
		salaryEntry,
		monthSelect,
		totalLabel,
		successMessageLabel,
		listExpensesButton,
		loadFromJSONButton,
		saveJSONButton,
		saveCSVButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 400))
	myWindow.ShowAndRun()
}
