package gitbsl

import (
	"catcher/app/internal/testutil/logging"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExternalModule(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"ВнешняяОбработкаSomething", true},
		{"ВнешнийОтчетSomething", true},
		{"Обработка", false},
		{"Отчет", false},
		{"", false},
	}

	for _, tt := range tests {
		got := IsExternalModule(tt.input)
		assert.Equal(t, tt.want, got, "IsExternalModule(%q)", tt.input)
	}
}

func TestAbsPath(t *testing.T) {
	logger := &logging.TestLogger{}

	tests := []struct {
		input          string
		sourceCodeRoot string
		wantPrefix     string
		wantErr        bool
	}{
		{
			input:          "ОбщийМодуль.ОбменСSentry.Модуль",
			sourceCodeRoot: "config",
			wantPrefix:     "config/CommonModules/ОбменСSentry/Ext/Module.bsl",
			wantErr:        false,
		},
		{
			input:          "Документ._Шаблон.Форма.ФормаДокумента.Форма",
			sourceCodeRoot: "config",
			wantPrefix:     "config/Documents/_Шаблон/Forms/ФормаДокумента/Ext/Form/Module.bsl",
			wantErr:        false,
		},
		{
			input:          "Обработка.ОбменЭлектроннымиДокументамиСБанком.Команда.СписокЭД.МодульКоманды",
			sourceCodeRoot: "config",
			wantPrefix:     "config/DataProcessors/ОбменЭлектроннымиДокументамиСБанком/Commands/СписокЭД/Ext/CommandModule.bsl",
			wantErr:        false,
		},
		{
			input:          "ОбщаяФорма.АдресУчастникаОбменаЭД.Форма",
			sourceCodeRoot: "config",
			wantPrefix:     "config/CommonForms/АдресУчастникаОбменаЭД/Ext/Form/Module.bsl",
			wantErr:        false,
		},
		{
			input:          "Обработка.АРМПриемосдатчика.Форма.ПереченьСД.Форма",
			sourceCodeRoot: "config",
			wantPrefix:     "config/DataProcessors/АРМПриемосдатчика/Forms/ПереченьСД/Ext/Form/Module.bsl",
			wantErr:        false,
		},
		{
			input:          "МодульУправляемогоПриложения",
			sourceCodeRoot: "config",
			wantPrefix:     "config/Ext/ManagedApplicationModule.bsl",
			wantErr:        false,
		},
		{
			input:          "ОбщаяКоманда.ВыгрузитьДанныеВФайл.МодульКоманды",
			sourceCodeRoot: "config",
			wantPrefix:     "config/CommonCommands/ВыгрузитьДанныеВФайл/Ext/CommandModule.bsl",
			wantErr:        false,
		},
		{
			input:          "ВнешнийИсточникДанных.Шлюз.Таблица.PDT_PartialOperations.МодульМенеджера",
			sourceCodeRoot: "config",
			wantPrefix:     "config/ExternalDataSources/Шлюз/Tables/PDT_PartialOperations/Ext/ManagerModule.bsl",
			wantErr:        false,
		},
		{
			input:          "НеизвестныйТип.Объект",
			sourceCodeRoot: "config",
			wantPrefix:     "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		p := NewPath(tt.input, tt.sourceCodeRoot, logger)
		got, err := p.AbsPath()
		if tt.wantErr {
			assert.Error(t, err, "AbsPath(%q) expected error", tt.input)
		} else {
			assert.NoError(t, err, "AbsPath(%q) unexpected error", tt.input)
			assert.True(t, strings.HasPrefix(got, tt.wantPrefix), "AbsPath(%q) = %q, want prefix %q", tt.input, got, tt.wantPrefix)
		}
	}
}

func TestTranslateBaseType(t *testing.T) {
	logger := &logging.TestLogger{}
	p := NewPath("", "", logger)

	tests := []struct {
		input string
		want  string
	}{
		{"ОбщийМодуль", "CommonModules"},
		{"Документ", "Documents"},
		{"Неизвестно", ""},
	}

	for _, tt := range tests {
		got := p.translateBaseType(tt.input, false)
		assert.Equal(t, tt.want, got, "translateBaseType(%q)", tt.input)
	}
}
