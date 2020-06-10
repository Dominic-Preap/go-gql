package scalar

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// Write custom scalar for javascript datetime
// https://morioh.com/p/040fc7ab1854
// https://gist.github.com/gigo1980/e77549d19aa807b5e3bfc84762049cf9

// DateTime .
type DateTime time.Time

// ISOString .
const ISOString = "2006-01-02T15:04:05.999Z07:00"

// MarshalDateTime .
func MarshalDateTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, e := io.WriteString(w, fmt.Sprintf("%s%s%s", "\"", t.Format(ISOString), "\""))
		if e != nil {
			panic(e)
		}
	})
}

// UnmarshalDateTime .
func UnmarshalDateTime(v interface{}) (time.Time, error) {
	str, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("datetime must be strings")
	}
	withoutQuotes := strings.ReplaceAll(str, "\"", "")
	i, err := time.Parse(ISOString, withoutQuotes)
	if err != nil {
		i, err = time.Parse("2006-01-02T15:04:05", withoutQuotes)
	}
	return i, err
}
