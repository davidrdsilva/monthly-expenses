package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Monthly Expenses Manager")

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

		// Clear input fields
		labelEntry.SetText("")
		priceEntry.SetText("")
		dateEntry.SetText("")

		// Calculate total expenses and display
		total := calculateTotal(bills)
		totalLabel.SetText(fmt.Sprintf("Total: R$%.2f", total))
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

	// Test it
	listExpensesButton := widget.NewButton("Show detailed expenses", func() {
		title := fmt.Sprintf("%s's expenses | %s", personNameEntry.Text, monthSelect.Selected)
		listExpenses(myApp, title, bills, salaryEntry)
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
		listExpensesButton,
		saveJSONButton,
		saveCSVButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 400))
	myWindow.ShowAndRun()
}
