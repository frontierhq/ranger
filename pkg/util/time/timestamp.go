package time

import (
	"strconv"
	"time"
)

func GetUnixTimestamp() string {
	n := time.Now()
	u := n.Unix()
	return strconv.FormatInt(u, 10)
}
