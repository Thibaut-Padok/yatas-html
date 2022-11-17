package main

import (
	"embed"
	"os"
	"strconv"

	"github.com/stangirard/yatas/plugins/commons"
)

//go:embed templates/*
var f embed.FS

func startTable(file *os.File, table_name string, attributes []string) {
	// To start an new html table
	// with the same number of column as the length of attributes
	file.WriteString("\n\t\t<table border=1 id=\"" + table_name + "\" align=center cellpadding=10>")
	file.WriteString("\n\t\t\t<thead>")
	file.WriteString("\n\t\t\t\t<tr>")
	for _, att := range attributes {
		file.WriteString("\n\t\t\t\t\t<th>" + att + "</th>")
	}
	file.WriteString("\n\t\t\t\t</tr>")
	file.WriteString("\n\t\t\t</thead>")
	file.WriteString("\n\t\t\t<tbody>")
}

func newLine(file *os.File, id, name, status string) {
	// To add a new table line (without message)
	file.WriteString("\n\t\t\t\t<tr>")
	file.WriteString("\n\t\t\t\t\t<td>" + id + "</td>")
	file.WriteString("\n\t\t\t\t\t<td>" + name + "</td>")
	file.WriteString("\n\t\t\t\t\t<td>" + status + "</td>")
	file.WriteString("\n\t\t\t\t</tr>")
}

func newLineWithMessage(file *os.File, item int, check commons.Check) {
	// To add a new table line with a result message
	status := "✅"
	if check.Status == "FAIL" {
		status = "❌"
	}
	key := strconv.Itoa(item)
	// Element
	file.WriteString("\n\t\t\t\t<tr onclick=\"showHideRow('hidden_row" + key + "');\">")
	file.WriteString("\n\t\t\t\t\t<td>" + check.Id + "</td>")
	file.WriteString("\n\t\t\t\t\t<td>" + check.Name + "</td>")
	file.WriteString("\n\t\t\t\t\t<td>" + status + "</td>")
	file.WriteString("\n\t\t\t\t</tr>")
	// Element hidden description
	file.WriteString("\n\t\t\t\t<tr id=\"hidden_row" + key + "\" class=\"hidden_row\">")
	file.WriteString("\n\t\t\t\t\t<td  colspan=3>")
	file.WriteString("\n\t\t\t\t\t\t<ul>")
	for _, res := range check.Results {
		file.WriteString("\n\t\t\t\t\t\t\t<li>" + res.Message + "</li>")
	}
	file.WriteString("\n\t\t\t\t\t\t</ul>")
	file.WriteString("\n\t\t\t\t\t</td>")
	file.WriteString("\n\t\t\t\t</tr>")
}

func closeTable(file *os.File) {
	// To close html table
	file.WriteString("\n\t\t\t</tbody>")
	file.WriteString("\n\t\t</table>")
}

func copyHeader(file *os.File) {
	// To start html file from header template
	data, _ := f.ReadFile("templates/template-header.html")
	content := string(data)
	file.WriteString(content)
}
func copyFooter(file *os.File) {
	// To end html file from header template
	data, _ := f.ReadFile("templates/template-footer.html")
	content := string(data)
	file.WriteString(content)
}

func WriteHtml(tests []commons.Tests) {
	// Write markdown report
	// Open the file for writing
	file, err := os.Create("report.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	copyHeader(file)
	// Write the header
	//file.WriteString("<html lang=\"fr\">\n<head>\n\t<meta charset=\"utf-8\">\n\t<title>YAMAS report</title>\n\t<link rel=\"shortcut icon\" href=\"https://www.padok.fr/hubfs/padok-favicon.png\">\n\t\n\t<script src=\"https://cdnjs.cloudflare.com/ajax/libs/jquery/3.3.1/jquery.min.js\"></script>\n\t<script src=\"https://stackpath.bootstrapcdn.com/bootstrap/4.1.2/js/bootstrap.min.js\"></script>\n\t<link rel=\"stylesheet\" href=\"https://stackpath.bootstrapcdn.com/bootstrap/4.1.2/css/bootstrap.min.css\">\n\t<link rel=\"stylesheet\"type=\"text/css\" href=\"https://use.fontawesome.com/releases/v5.6.3/css/all.css\">\n\t<script type=\"text/javascrip\t\">\n\t\tfunction showHideRow(row) {\n\t\t\t$(\"#\" + row).toggle();\n\t\t}\n\t</script><style>\n\ttable {\n\t\tborder-collapse: collapse;\n\t\tborder: 1px black solid;\n\t\tfont: 12px sans-serif;\n\t}\n\ttd {\n\t\tborder: 1px black solid;\n\t\tpadding: 5px;\n\t}\n\t</style>\n</head>\n<body>")

	// Write the checks in a table
	k := 1
	for _, test := range tests {
		// Title
		file.WriteString("\n\t\t<h1>" + test.Account + "</h1>")
		//Table
		categories := []string{"Id", "Name", "Statut"}
		startTable(file, "table_detail", categories)
		for _, check := range test.Checks {
			k++
			newLineWithMessage(file, k, check)
		}
		closeTable(file)
		// Add categories table
		WriteCategoriesSuccess(test, file)
		// Write the footer
		copyFooter(file)
	}
}

func WriteCategoriesSuccess(test commons.Tests, file *os.File) {
	// Find the categories
	categories := []string{}
	categoriesSuccess := map[string]int{}
	categoriesFailure := map[string]int{}
	for _, check := range test.Checks {
		for _, category := range check.Categories {
			if !contains(categories, category) {
				categories = append(categories, category)
				categoriesSuccess[category] = 0
				categoriesFailure[category] = 0
			}
			if check.Status == "OK" {
				categoriesSuccess[category]++
			} else {
				categoriesFailure[category]++
			}
		}
	}
	// Title
	file.WriteString("\n\t\t<h1>Categories</h1>")
	// Category Table
	categories2 := []string{"Category", "Completion"}
	startTable(file, "categories_table", categories2)
	for _, category := range categories {
		file.WriteString("\n\t\t\t\t<tr>")
		file.WriteString("\n\t\t\t\t\t<td>" + category + "</td>")
		file.WriteString("\n\t\t\t\t\t<td>" + CalculatePercent(categoriesSuccess[category], categoriesFailure[category]) + "</td>")
		file.WriteString("\n\t\t\t\t</tr>")
	}
	closeTable(file)
}

func CalculatePercent(success int, failure int) string {
	total := success + failure
	if total == 0 {
		return "0"
	}
	return strconv.Itoa((success * 100) / total)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
