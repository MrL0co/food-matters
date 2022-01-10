package Markdown

import (
	"strings"
)

type Alignment int64

const (
	Left Alignment = iota
	Center
	Right
)

type table struct {
	columns int
	header  []string
	align   []Alignment
	rows    [][]string
}

func NewTable(header ...string) *table {
	t := table{}
	t.SetHeader(header...)
	return &t
}

func (table *table) SetHeader(h ...string) {
	table.header = h
	table.align = make([]Alignment, len(h))
	table.rows = [][]string{}

	table.columns = len(h)
}

func (table *table) SetAlign(align ...Alignment) {
	if len(align) != table.columns {
		panic("mismatch table data")
	}

	table.align = align
}

func (table *table) AddRow(v ...string) {
	if len(v) != table.columns {
		panic("mismatch table data")
	}

	table.rows = append(table.rows, v)
}

func (table *table) ToString() string {
	var sb strings.Builder

	table.renderRow(&sb, table.header)

	sb.WriteString("| ")
	for i, a := range table.align {
		if i > 0 {
			sb.WriteString(" | ")
		}

		switch a {
		case Left:
			sb.WriteString(":--")
		case Center:
			sb.WriteString(":--:")
		case Right:
			sb.WriteString("--:")
		}
	}
	sb.WriteString(" |\n")

	for _, row := range table.rows {
		table.renderRow(&sb, row)
	}

	return sb.String()
}

func (table *table) renderRow(sb *strings.Builder, row []string) {
	sb.WriteString("| ")
	for i, h := range row {
		if i > 0 {
			sb.WriteString(" | ")
		}
		sb.WriteString(h)
	}
	sb.WriteString(" |\n")
}
