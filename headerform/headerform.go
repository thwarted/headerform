package headerform

import (
	"fmt"
	"io"
	"strings"
)

// TODO works on strings at the byte level, passing multibyte utf8
//      strings is undefined

type alignment uint8

const (
	alignNone alignment = iota
	alignLeft
	alignCenter
	alignRight
)

type HeaderFormatter struct {
	DividingLine    string // Characater to use for dividing line
	InputDelimiter  string // Field delimiter in input (defaults to \t)
	OutputDelimiter string // Column delimiter in output
	Buffer          bool   // Should we buffer

	to       io.Writer   // Output destination
	headers  []string    // The names of the headers
	widths   []int       // Column widths
	flexCols []int       // List of columns that are flexiblly sized
	aligns   []alignment // Column alignments
	_buffer  [][]string
}

// New returns a new formatter object.
func New(w io.Writer) *HeaderFormatter {
	return &HeaderFormatter{
		to:              w,
		OutputDelimiter: " ",
		InputDelimiter:  "\t",
		Buffer:          false,
		DividingLine:    "",
	}
}

// FormatAs sets the layout and headers for the output.
//
// `s` contains a layout specification of the names of headers delimited by
// pipes.  Place whitespace around the header names to indicate alignment
// for that column:
//  - whitespace on the right will left justify
//  - whitespace on the left will right justify
//  - whitespace on both sides will center
//  - the width of the header field, including whitespace, determines the
//    width of the outpu
//
//    "Name    | RightJustified|  Centered  |LeftJustified   |Remaining"
//
// A trailing pipe will set the width of the final column (and the alignment
// of the header).
//
// Whitespace around the last field (without a trailing pipe) will set the
// width and alignment for the header and for shorter values, but truncation
// will not be performed.
//
// Output fields that are longer than the header specification are truncated
// and indicated with a …
//
// If more fields are passed than there are headers they will be output space
// delimited and no alignment or truncation will be performed.
//
// If anything was buffered (by calling Print* before calling FormatAs), then
// FormatAs will call Headers() and then Flush() to output the buffer.
//
// If < and/or > is used in place of, respectively, leading or trailing
// whitespace above, then the column width will be flexible based on the
// of buffered data for that column.
//
//    "Left>|<Right|<Centered>|Left>|Remaining"
//
func (b *HeaderFormatter) FormatAs(s string) {
	h := strings.Split(s, "|")
	var a alignment
	for i, header := range h {
		p := strings.HasPrefix(header, " ") || strings.HasPrefix(header, "<")
		s := strings.HasSuffix(header, " ") || strings.HasSuffix(header, ">")
		flex := strings.HasPrefix(header, "<") || strings.HasSuffix(header, ">")
		if p {
			a = alignRight
		}
		if s {
			a = alignLeft
		}
		if p && s {
			a = alignCenter
		}
		b.widths = append(b.widths, len([]byte(header)))
		b.aligns = append(b.aligns, a)
		b.headers = append(b.headers, strings.Trim(header, " <>"))
		if flex {
			b.Buffer = true
			b.flexCols = append(b.flexCols, i)
		}
	}
	if !b.Buffer {
		b.Headers()
	}
}

// Headers outputs the headers and flushes any buffered prints.
//
// Output will be buffered until Flush() is called
func (b *HeaderFormatter) Headers() {
	if len(b.headers) == 0 {
		return
	}

	h := [][]string{b.headers}
	b.recalculateFlexibleColumns()
	if len(b.DividingLine) > 0 {
		h = append(h, b.outputDivider())
	}
	if b.Buffer {
		b._buffer = append(h, b._buffer...)
	} else {
		b.output(h[0])
		if len(h) == 2 {
			b.output(h[1])
		}
	}
}

func (b *HeaderFormatter) outputDivider() []string {
	divider := string(b.DividingLine[0])
	a := make([]string, len(b.widths))
	for i, w := range b.widths {
		a[i] = strings.Repeat(divider, w)
	}
	return a[:]
}

// Flush outputs anything that was buffered, printing the header
// if it wasn't already immediately previously printed.
func (b *HeaderFormatter) Flush() {
	if len(b._buffer) > 0 {
		for _, fields := range b._buffer {
			b.outputImmediately(fields)
		}
	}
	b._buffer = nil
}

// Printf accepts a fmt.Printf-style format specification.
// Embedded fields should be delimited with tabs to designate, within the
// resulting string, the content of each field.
func (b *HeaderFormatter) Printf(s string, i ...interface{}) {
	b.Print(fmt.Sprintf(s, i...))
}

// Print accepts a single string with embedded tabs to designate the content
// of each field.
func (b *HeaderFormatter) Print(s string) {
	id := b.InputDelimiter
	if len(id) == 0 {
		id = "\t"
	}
	b.output(strings.Split(s, id))
}

// PrintAny accepts an array of values of any type, converts those values to
// strings through fmt.Sprintf("%v") (which may use the value's fmt.Stringer
// interface), then outputs the final array of strings.
func (b *HeaderFormatter) PrintAny(a ...interface{}) {
	s := make([]string, len(a))
	for i, x := range a {
		s[i] = fmt.Sprintf("%v", x)
	}
	b.output(s)
}

// PrintStrings accepts an array of strings and outputs them per the format.
func (b *HeaderFormatter) PrintStrings(fields ...string) {
	b.output(fields)
}

// output accepts an array of strings and either buffers them if a format
// layout hasn't been specified yet or we're explicitly buffering, otherwise
// outputs them immediately.
func (b *HeaderFormatter) output(fields []string) {
	if len(b.headers) == 0 || b.Buffer {
		b._buffer = append(b._buffer, fields)
	} else {
		b.outputImmediately(fields)
	}
	return
}

// outputImmediately is the meat of the functionality, aligning and truncating
// the fields as necessary and sending them to the formatter's io.Writer.
func (b *HeaderFormatter) outputImmediately(fields []string) {
	if len(b.widths) != len(b.aligns) {
		panic("widths and alignment don't match")
	}
	last := len(b.headers) - 1
	for i, f := range fields {
		if i > 0 {
			// print the output delimiter between columns
			if i <= last {
				fmt.Fprint(b.to, b.OutputDelimiter)
			} else {
				// unless this field is beyond the last column in the layout
				// then we delimit with a single space
				fmt.Fprint(b.to, " ")
			}
		}
		if i >= len(b.widths) || b.widths[i] == 0 {
			// no formatting/layout information for this column, output
			// it directly
			fmt.Fprint(b.to, f)
			continue
		}
		switch b.aligns[i] {
		case alignRight:
			if i < last && len([]byte(f)) > b.widths[i] {
				f = f[len(f)-b.widths[i]+1 : len(f)]
				f = "…" + f
			}
			fmt.Fprintf(b.to, fmt.Sprintf("%%%ds", b.widths[i]), f)
		case alignCenter:
			if len(f) < b.widths[i] {
				f = strings.Repeat(" ", (b.widths[i]-len(f))/2) + f + strings.Repeat(" ", b.widths[i])
				f = f[0:b.widths[i]]
			}
			fallthrough // if it's too wide for centering in the given width, then align left
		case alignLeft, alignNone:
			if i < last && len([]byte(f)) > b.widths[i] {
				f = f[0 : b.widths[i]-1]
				f = f + "…"
			}
			fmt.Fprintf(b.to, fmt.Sprintf("%%-%ds", b.widths[i]), f)
		default:
			panic("unimplemented alignment")
		}
	}
	fmt.Fprintln(b.to)
}

// a simple type that keeps track of the maximum value
type mmax int

func (m *mmax) Max(n int) {
	if n > int(*m) {
		*m = mmax(n)
	}
}

func (b *HeaderFormatter) recalculateFlexibleColumns() {
	if len(b.flexCols) == 0 || len(b._buffer) == 0 {
		return
	}
	m := new(mmax)
	for _, c := range b.flexCols {
		*m = mmax(len(b.headers[c]))
		for i := 0; i < len(b._buffer); i++ {
			if len(b._buffer[i]) > c {
				m.Max(len(b._buffer[i][c]))
			}
		}
		b.widths[c] = int(*m)
	}
}
