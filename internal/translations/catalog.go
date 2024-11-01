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
	"Add a wish":                    8,
	"Belarusian":                    6,
	"Copy link to share":            7,
	"Create a new wishlist":         12,
	"Edit":                          10,
	"English":                       3,
	"Got it":                        15,
	"Max items: %d. You added %d.":  1,
	"No items provided":             2,
	"Polish":                        5,
	"Russian":                       4,
	"Save":                          9,
	"This item is already reserved": 13,
	"This item is already taken by another user":                               14,
	"Wishlist not found":                                                       11,
	"You are making too many requests!":                                        16,
	"You have reached the limit for creating wishlists. Please wait a moment.": 17,
	"You have reached the limit for the number or size of items in the wishlist. Please remove some items.": 19,
	"Your request is too large!":            18,
	"Your text is %d characters; max is %d": 0,
}

var be_BYIndex = []uint32{ // 21 elements
	0x00000000, 0x0000004b, 0x00000090, 0x000000b9,
	0x000000ce, 0x000000db, 0x000000ec, 0x00000101,
	0x00000141, 0x00000163, 0x00000174, 0x00000189,
	0x000001bf, 0x000001f7, 0x0000022d, 0x00000293,
	0x000002a6, 0x000002ea, 0x00000375, 0x000003a4,
	0x00000462,
} // Size: 108 bytes

const be_BYData string = "" + // Size: 1122 bytes
	"\x02Даўжыня тэксту %[1]d знакаў; максімум — %[2]d\x02Максімум элементаў:" +
	" %[1]d. Вы дадалі %[2]d.\x02Элементы не зададзены\x02Англійская\x02Руска" +
	"я\x02Польская\x02Беларуская\x02Скапіраваць спасылку для адпраўкі\x02Дад" +
	"айце пажаданне\x02Захаваць\x02Рэдагаваць\x02Спіс пажаданняў не знойдзен" +
	"ы\x02Стварыць новы спіс пажаданняў\x02Гэты пункт ужо зарэзерваваны\x02Г" +
	"эты падарунак ужо зарэзерваваны іншым карыстальнікам\x02Зразумела\x02Вы" +
	" адпраўляеце занадта шмат запытаў!\x02Вы дасягнулі мяжы стварэння спісаў" +
	" пажаданняў. Калі ласка, пачакайце крыху.\x02Ваш запыт занадта вялікі!" +
	"\x02Вы дасягнулі ліміту па колькасці або памеры элементаў у спісе. Калі " +
	"ласка, выдаліце некаторыя элементы."

var en_GBIndex = []uint32{ // 21 elements
	0x00000000, 0x0000002c, 0x0000004f, 0x00000061,
	0x00000069, 0x00000071, 0x00000078, 0x00000083,
	0x00000096, 0x000000a1, 0x000000a6, 0x000000ab,
	0x000000be, 0x000000d4, 0x000000f2, 0x0000011d,
	0x00000124, 0x00000146, 0x0000018f, 0x000001aa,
	0x00000210,
} // Size: 108 bytes

const en_GBData string = "" + // Size: 528 bytes
	"\x02Your text is %[1]d characters; max is %[2]d\x02Max items: %[1]d. You" +
	" added %[2]d.\x02No items provided\x02English\x02Russian\x02Polish\x02Be" +
	"larusian\x02Copy link to share\x02Add a wish\x02Save\x02Edit\x02Wishlist" +
	" not found\x02Create a new wishlist\x02This item is already reserved\x02" +
	"This item is already taken by another user\x02Got it\x02You are making t" +
	"oo many requests!\x02You have reached the limit for creating wishlists. " +
	"Please wait a moment.\x02Your request is too large!\x02You have reached " +
	"the limit for the number or size of items in the wishlist. Please remove" +
	" some items."

var pl_PLIndex = []uint32{ // 21 elements
	0x00000000, 0x00000030, 0x0000005d, 0x00000073,
	0x0000007d, 0x00000086, 0x0000008d, 0x00000099,
	0x000000b8, 0x000000c8, 0x000000cf, 0x000000d6,
	0x000000f4, 0x00000112, 0x00000136, 0x0000016e,
	0x00000177, 0x00000196, 0x000001dd, 0x000001fe,
	0x00000263,
} // Size: 108 bytes

const pl_PLData string = "" + // Size: 611 bytes
	"\x02Twój tekst ma %[1]d znaków; maksimum to %[2]d\x02Maksimum przedmiotó" +
	"w: %[1]d. Dodano: %[2]d.\x02Nie podano elementów\x02Angielski\x02Rosyjsk" +
	"i\x02Polski\x02Białoruski\x02Skopiuj link do udostępnienia\x02Dodaj życz" +
	"enie\x02Zapisz\x02Edytuj\x02Lista życzeń nie znaleziona\x02Utwórz nową l" +
	"istę życzeń\x02Ten prezent jest już zarezerwowany\x02Ten prezent jest ju" +
	"ż wzięty przez innego użytkownika\x02Rozumiem\x02Wysyłasz zbyt wiele żą" +
	"dań!\x02Osiągnięto limit tworzenia list życzeń. Proszę chwilę poczekać." +
	"\x02Twoje żądanie jest zbyt duże!\x02Osiągnięto limit liczby lub rozmiar" +
	"u przedmiotów na liście. Proszę usunąć kilka przedmiotów."

var ru_RUIndex = []uint32{ // 21 elements
	0x00000000, 0x0000004b, 0x00000092, 0x000000b5,
	0x000000ca, 0x000000d9, 0x000000ea, 0x00000101,
	0x0000013d, 0x00000161, 0x00000174, 0x0000018f,
	0x000001c1, 0x000001fb, 0x00000233, 0x00000297,
	0x000002a6, 0x000002ee, 0x0000037e, 0x000003b1,
	0x00000472,
} // Size: 108 bytes

const ru_RUData string = "" + // Size: 1138 bytes
	"\x02Длина текста %[1]d символов; максимум — %[2]d\x02Максимум элементов:" +
	" %[1]d. Добавлено: %[2]d.\x02Элементы не заданы\x02Английский\x02Русский" +
	"\x02Польский\x02Белорусский\x02Скопировать ссылку для отправки\x02Добави" +
	"ть пожелание\x02Сохранить\x02Редактировать\x02Список пожеланий не найде" +
	"н\x02Создать новый список пожеланий\x02Этот пункт уже зарезервирован" +
	"\x02Этот подарок уже зарезервирован другим пользователем\x02Понятно\x02В" +
	"ы отправляете слишком много запросов!\x02Вы достигли лимита создания сп" +
	"исков пожеланий. Пожалуйста, подождите немного.\x02Ваш запрос слишком б" +
	"ольшой!\x02Вы достигли лимита по количеству или размеру элементов в спи" +
	"ске. Пожалуйста, удалите некоторые элементы."

	// Total table size 3831 bytes (3KiB); checksum: A5BAE445
