package model

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

var DifficultyList = []string{"basic", "advanced", "expert", "master", "re_master"}
var BanquetDifficultyList = []string{"banquet", "banquet2"}

type DS map[string]float64
type Level map[string]string
type Cids map[string]int

// 将数组按难度映射为 map
func mapByDifficulty[T float64 | string | int](row []T, genre string) map[string]T {
	result := make(map[string]T)

	var difficultyList []string
	if genre == "宴会場" {
		difficultyList = BanquetDifficultyList
	} else {
		difficultyList = DifficultyList
	}
	for i := range row {
		difficulty := difficultyList[i]
		result[difficulty] = row[i]
	}
	return result
}

func (d DS) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d)
}

func (d *DS) Scan(value interface{}) error {
	if value == nil {
		*d = make(DS)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, d)
}

func (l Level) Value() (driver.Value, error) {
	if l == nil {
		return nil, nil
	}
	return json.Marshal(l)
}

func (l *Level) Scan(value interface{}) error {
	if value == nil {
		*l = make(Level)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, l)
}

func (c Cids) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *Cids) Scan(value interface{}) error {
	if value == nil {
		*c = make(Cids)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, c)
}

// ------------------------------------------------------------

var LevelList = []string{}
var LevelSortMap map[string]int

func init() {
	LevelList = []string{}
	LevelSortMap = make(map[string]int)
	LevelSortMap["0"] = 0

	for i := 0; i <= 15; i++ {
		level := strconv.Itoa(i)
		levelPlus := level + "+"
		LevelList = append(LevelList, level)
		LevelList = append(LevelList, levelPlus)
		LevelSortMap[level] = i*2 + 1     // 1,3,5...
		LevelSortMap[levelPlus] = i*2 + 2 // 2,4,6...
	}
}

func LevelSortFunc(level1 string, level2 string) int {
	return LevelSortMap[level1] - LevelSortMap[level2]
}

func GetLevelRange(minLevel string, maxLevel string) []string {
	if minLevel == "" {
		minLevel = "0"
	}
	if maxLevel == "" {
		maxLevel = "15+"
	}

	levelList := make([]string, 0)
	for level := range LevelSortMap {
		if LevelSortFunc(level, minLevel) >= 0 && LevelSortFunc(level, maxLevel) <= 0 {
			levelList = append(levelList, level)
		}
	}
	return levelList
}
