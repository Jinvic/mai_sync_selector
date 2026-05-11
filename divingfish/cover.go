package divingfish

import (
	"fmt"
	"strconv"
)

var defaultCoverURL = "https://www.diving-fish.com/covers"

func GetCoverUrl(id string) string {
	return fmt.Sprintf("%s/%s.png", defaultCoverURL, formatIDForCover(id))
}

func formatIDForCover(idstr string) string {
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return idstr
	}
	if id > 10000 && id < 11000 {
		id -= 10000
	}
	return fmt.Sprintf("%05d", id)
}
