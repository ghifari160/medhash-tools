package main

import (
	"fmt"
	"strings"
)

// Summary generates a summary based on an array of rows and columns.
func Summary(tbl [][]string) string {
	return fmt.Sprintf("# Summary\n\n%s\n", GenTable(tbl))
}

// GenTable generates a Markdown table from an array of rows and columns.
func GenTable(tbl [][]string) string {
	length := make([]int, len(tbl[0]))

	for _, row := range tbl {
		for i, col := range row {
			if len(col) > length[i] {
				length[i] = len(col)
			}
		}
	}

	mdTbl := make([]string, 0)
	for i, row := range tbl {
		line := make([]string, len(row))

		for j, col := range row {
			line[j] = col

			for s := 0; s < length[j]-len(col)+1; s++ {
				line[j] += " "
			}
		}

		mdTbl = append(mdTbl, "| "+strings.Join(line, "| ")+"|")

		if i == 0 {
			line := make([]string, len(row))

			for j := range row {
				for s := 0; s < length[j]+1; s++ {
					line[j] += "-"
				}
			}

			mdTbl = append(mdTbl, "|-"+strings.Join(line, "|-")+"|")
		}
	}

	return strings.Join(mdTbl, "\n")
}
