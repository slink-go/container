package container

import (
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	str := d.Format("2006-01-02")
	return []byte(fmt.Sprintf("\"%s\"", str)), nil
}
