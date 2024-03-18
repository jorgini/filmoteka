package filmoteka

import (
	"fmt"
	"strings"
	"time"
)

type Date time.Time

const layout = "02-01-2006"

func (d *Date) String() string {
	return fmt.Sprintf(time.Time.Format(time.Time(*d), layout))
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)

	tmp, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	*d = Date(tmp)
	return nil
}
