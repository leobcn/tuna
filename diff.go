package tuna

import "fmt"

// Diff runs a Extractor on the (x[i+1] - x[i]) version of a stream of
// values. This can be used in conjunction with a GroupBy to compute rolling
// statistics.
type Diff struct {
	Parse     func(Row) (float64, error)
	Extractor Extractor
	FieldName string
	seen      bool
	xi        float64
}

// Update Diff given a Row.
func (d *Diff) Update(row Row) error {
	var x, err = d.Parse(row)
	if err != nil {
		return err
	}
	if !d.seen {
		d.xi = x
		d.seen = true
		return nil
	}
	row[d.FieldName] = float2Str(x - d.xi)
	d.xi = x
	return d.Extractor.Update(row)
}

// Collect returns the current value.
func (d Diff) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		for r := range d.Extractor.Collect() {
			c <- r
		}
		close(c)
	}()
	return c
}

// Size is the size of the Extractor.
func (d Diff) Size() uint { return d.Extractor.Size() }

// NewDiff returns a Diff that applies a Extractor to the difference of
// a given field.
func NewDiff(field string, newExtractor func(s string) Extractor) *Diff {
	fn := fmt.Sprintf("%s_diff", field)
	return &Diff{
		Parse:     func(row Row) (float64, error) { return str2Float(row[field]) },
		Extractor: newExtractor(fn),
		FieldName: fn,
	}
}
