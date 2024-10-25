// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package translations

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"be_BY": &dictionary{index: be_BYIndex, data: be_BYData},
		"en_GB": &dictionary{index: en_GBIndex, data: en_GBData},
		"pl_PL": &dictionary{index: pl_PLIndex, data: pl_PLData},
		"ru_RU": &dictionary{index: ru_RUIndex, data: ru_RUData},
	}
	fallback := language.MustParse("en-GB")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"Create new wishlist": 1,
	"Wishlist not found":  0,
}

var be_BYIndex = []uint32{ // 3 elements
	0x00000000, 0x00000032, 0x00000066,
} // Size: 36 bytes

const be_BYData string = "" + // Size: 102 bytes
	"\x02Спіс жаданняў не знойдзены\x02Стварыць новы спіс жаданняў"

var en_GBIndex = []uint32{ // 3 elements
	0x00000000, 0x00000013, 0x00000027,
} // Size: 36 bytes

const en_GBData string = "\x02Wishlist not found\x02Create new wishlist"

var pl_PLIndex = []uint32{ // 3 elements
	0x00000000, 0x0000001e, 0x0000003c,
} // Size: 36 bytes

const pl_PLData string = "" + // Size: 60 bytes
	"\x02Lista życzeń nie znaleziona\x02Utwórz nową listę życzeń"

var ru_RUIndex = []uint32{ // 3 elements
	0x00000000, 0x0000002e, 0x00000064,
} // Size: 36 bytes

const ru_RUData string = "" + // Size: 100 bytes
	"\x02Список желаний не найден\x02Создать новый список желаний"

	// Total table size 445 bytes (0KiB); checksum: CACD2550