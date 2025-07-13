package normalize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestId(t *testing.T) {
	input := "123-456-789"
	expected := "123456789"
	got := Id(input)
	assert.Equal(t, expected, got, "Id(%q) should return %q", input, expected)
}

func TestRelease(t *testing.T) {
	id := "module"
	version := "1.0"
	expected := "module@1.0"
	got := Release(id, version)
	assert.Equal(t, expected, got, "Release(%q, %q) should return %q", id, version, expected)
}

func TestMessage(t *testing.T) {
	input := " {test} message {test} "
	expected := "message"
	got := Message(input)
	assert.Equal(t, expected, got, "Message(%q) should return %q", input, expected)
}

func TestIp(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"192.168.0.1, 10.0.0.1", "192.168.0.1"},
		{"  127.0.0.1  , localhost", "127.0.0.1"},
		{"255.255.255.255", "255.255.255.255"},
		{" 10.10.10.10 ,", "10.10.10.10"},
		{" , 1.1.1.1", ""},
		{"", ""},
	}

	for _, tt := range tests {
		result := Ip(tt.input)
		assert.Equal(t, tt.expected, result, "Ip(%q) должен вернуть %q", tt.input, tt.expected)
	}
}

func TestSplitString(t *testing.T) {
	input := "Infostart Toolkit PROF (2023.4.08)"
	sStart := "("
	sEnd := ")"
	expected := [2]string{"Infostart Toolkit PROF", "2023.4.08"}
	got := SplitString(input, sStart, sEnd)
	assert.Equal(t, expected, got, "SplitString(%q, %q, %q) should return %v", input, sStart, sEnd, expected)
}

func TestRemoveModuleNameSuffix(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"example.Модуль", "example"},
		{"example.Форма", "example"},
		{"example", "example"},
	}
	for _, c := range cases {
		got := RemoveModuleNameSuffix(c.input)
		assert.Equal(t, c.expected, got, "RemoveModuleNameSuffix(%q) should return %q", c.input, c.expected)
	}
}

func TestRemoveBraces(t *testing.T) {
	input := "{test}: some text }"
	expected := " some text "
	got := removeBraces(input)
	assert.Equal(t, expected, got, "removeBraces(%q) should return %q", input, expected)
}

func TestRemoveFromSecondBrace(t *testing.T) {
	input := `Ошибка при вызове метода контекста (ПредопределенноеЗначение)
{Обработка.ОтчетПоРейсу.МодульОбъекта(115)}:ОценкаПроизводительностиВызовСервераПолныеПрава.ЗакончитьЗамерВремениНаСервере(ПредопределенноеЗначение("Справочник.КлючевыеОперации.Рейс_АД_Отчет_Услуги_Сформировать")
{Обработка.ОтчетПоРейсу.МодульОбъекта(32)}:СформироватьОтчетПоУслугамРейса(ВидОтчета, Результат, ПараметрыСКД, ДополнительныеПараметры);
{Справочник.Рейсы.МодульМенеджера(1672)}:ОбработкаОтчета.СформироватьОтчет(ВидОтчета, Результат, ПараметрыСКД, ДополнительныеПараметры);
{Справочник.Рейсы.Форма.ФормаЭлементаАДКД.Форма(3074)}:Справочники.Рейсы.СформироватьОтчет("УслугиРейса", ТабличныйДокументУслуги, ПараметрыСКД, ДополнительныеПараметры);
{Справочник.Рейсы.Форма.ФормаЭлементаАДКД.Форма(3026)}:СформироватьУслуги();

[ОшибкаВоВремяВыполненияВстроенногоЯзыка]
по причине:
Предопределенный элемент не существует`

	expected := `Ошибка при вызове метода контекста (ПредопределенноеЗначение)
{Обработка.ОтчетПоРейсу.МодульОбъекта(115)}:ОценкаПроизводительностиВызовСервераПолныеПрава.ЗакончитьЗамерВремениНаСервере(ПредопределенноеЗначение("Справочник.КлючевыеОперации.Рейс_АД_Отчет_Услуги_Сформировать")}

[ОшибкаВоВремяВыполненияВстроенногоЯзыка]
по причине:
Предопределенный элемент не существует`

	result := RemoveFromSecondBrace(input, true)
	assert.Equal(t, expected, result, "Строка должна быть обрезана начиная со второй '{' до пустой строки, запятая удалена, добавлена '}'")

	// Тест: если меньше двух '{', строка не меняется
	input2 := "Текст без фигурных скобок\n"
	expected2 := input2
	result2 := RemoveFromSecondBrace(input2, true)
	assert.Equal(t, expected2, result2, "Если меньше двух '{', строка не меняется")

	// Тест: если нет пустой строки после второй '{', и есть запятая в конце
	input3 := "abc {first} def {second}, xyz"
	expected3 := "abc {first} def"
	result3 := RemoveFromSecondBrace(input3, false)
	assert.Equal(t, expected3, result3, "Если пустой строки нет после второй '{', удаляем запятую и добавляем '}'")

	// Тест: если нет пустой строки после второй '{', и запятая отсутствует
	input4 := "abc {first} def {second} xyz"
	expected4 := "abc {first} def"
	result4 := RemoveFromSecondBrace(input4, false)
	assert.Equal(t, expected4, result4, "Если пустой строки нет после второй '{', и запятая отсутствует, просто добавляем '}'")
}
