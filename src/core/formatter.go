package main

import "github.com/fatih/color"

// Formatter class for colored output
type Formatter struct {
	Blue   func(...interface{}) string
	Yellow func(...interface{}) string
	Red    func(...interface{}) string
	Green  func(...interface{}) string
}

//NewFormatter constructor
func NewFormatter() Formatter {
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	formatter := Formatter{blue, yellow, red, green}
	return formatter
}
