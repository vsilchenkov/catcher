package gitbsl

import (
	"catcher/app/internal/lib/logging"
	"errors"
	"strings"
)

const (
	ExtProcessing = "ВнешняяОбработка"
	ExtReport     = "ВнешнийОтчет"
	Expansion     = "bsl"
	PathSeparator = "/"
)

type Config struct {
	SourceCodeRoot string
}

type Path struct {
	Config
	Value  string
	Logger logging.Logger
}

func IsExternalModule(m string) bool {
	return strings.HasPrefix(m, ExtProcessing) || strings.HasPrefix(m, ExtReport)
}

func IsExpansion(m string) bool {
	return strings.Contains(m, " ") // имя расширения и путь к объекту разделяются пробелом
}

func NewConfig(sourceCodeRoot string) Config {
	return Config{
		SourceCodeRoot: sourceCodeRoot,
	}
}

func NewPath(value string, sourceCodeRoot string, logger logging.Logger) Path {
	return Path{
		Config: NewConfig(sourceCodeRoot),
		Value:  strings.TrimSpace(value),
		Logger: logger}
}

func (p Path) AbsPath() (string, error) {

	value := p.Value
	if value == "" {
		return "", nil
	}

	path := strings.TrimSpace(value)

	paths := make([]string, 0)
	modules := make([]string, 0)

	parts := strings.Split(path, ".")
	lenParts := len(parts)
	position := 0

	for _, segment := range parts {

		position++

		switch position {
		case 1: //Базовый тип

			base := p.translateBaseType(segment, true)
			if base == "" {
				return "", errors.New("базовый путь не определен")
			}
			paths = append(paths, base)
			modules = p.addModulesPath(modules, segment, false)

		case 2: //Имя объекта метаданных

			paths = append(paths, segment)
			modules = p.addModulesPath(modules, segment, false)

		case 3: //Форма, команда, модуль объекта, модуль менеджера,

			if lenParts > 3 {
				// Добавляем только, если есть еще часть
				base := p.translateObjectName(segment, false)
				if base != "" {
					paths = append(paths, base)
				}
			}

			modules = p.addModulesPath(modules, segment, false)

		case 4: //Имя формы, имя команды

			paths = append(paths, segment)
			modules = p.addModulesPath(modules, segment, false)

		}
	}

	paths = append(paths, modules...)

	if len(paths) == 0 {
		return "", nil
	}

	result := p.SourceCodeRoot +
		PathSeparator +
		strings.Join(paths, PathSeparator) +
		"." +
		Expansion

	p.Logger.Debug("AbsPath",
		p.Logger.Str("input", p.Value),
		p.Logger.Str("output", result))
	return result, nil

}

func (p Path) translateBaseType(base string, warn bool) string {

	const op = "gitbsl.translateBaseType"

	result, ok := mapBasesTypes()[base]
	if !ok {
		if warn {
			p.Logger.Warn("Неизвестное имя типа",
				p.Logger.Op(op),
				p.Logger.Str("name", base),
				p.Logger.Str("value", p.Value),
				p.Logger.Err(errors.New("тип не найден в mapBasesTypes")))
		}
		return ""
	}
	return result
}

func (p Path) translateObjectName(base string, warn bool) string {

	const op = "gitbsl.translateObjectName"

	result, ok := mapObjectNames()[base]
	if !ok {
		if warn {
			p.Logger.Warn("Неизвестное имя объекта",
				p.Logger.Op(op),
				p.Logger.Str("name", base),
				p.Logger.Str("value", p.Value),
				p.Logger.Err(errors.New("тип не найден в mapObjectNames")))
		}
		return ""
	}
	return result
}
func (p Path) addModulesPath(pathModules []string, name string, warn bool) []string {

	const op = "gitbsl.addPartPathModulesGit"

	if name == "" {
		return pathModules
	}

	modules, ok := mapModulesPath()[name]
	if !ok {
		if warn {
			p.Logger.Warn("Неизвестное имя модуля",
				p.Logger.Op(op),
				p.Logger.Str("name", name),
				p.Logger.Str("value", p.Value),
				p.Logger.Err(errors.New("имя модуля не найдено в mapModulesPath")))
		}
		return pathModules
	}

	pathModules = append(pathModules, modules...)

	return pathModules

}

func mapBasesTypes() map[string]string {

	return map[string]string{

		"Подсистема":     "Subsystems",
		"ОбщийМодуль":    "CommonModules",
		"ПараметрСеанса": "SessionParameters",

		"Роль":                "Roles",
		"ОбщийРеквизит":       "CommonAttributes",
		"ПланОбмена":          "ExchangePlans",
		"КритерийОтбора":      "FilterCriteria",
		"ПодпискаНаСобытия":   "EventSubscriptions",
		"РегламентноеЗадание": "ScheduledJobs",

		"ФункциональнаяОпция":          "FunctionalOptions",
		"ПараметрыФункциональнойОпции": "FunctionalOptionsParameters",

		"ОпределяемыйТип":   "DefinedTypes",
		"ХранилищеНастроек": "SettingsStorages",

		"ОбщаяФорма":    "CommonForms",
		"ОбщаяКоманда":  "CommonCommands",
		"ГруппаКоманды": "CommandGroups",
		"ОбщийМакет":    "CommonTemplates",
		"ОбщаяКартинка": "CommonPictures",

		"ПакетXDTO":  "XDTOPackages",
		"WebСервис":  "WebServices",
		"HTTPСервис": "HTTPServices",
		"WSСсылка":   "WSReferences",

		"ЭлементСтиля": "StyleItems",
		"Язык":         "Languages",

		"Константа": "Constants",

		"Справочник": "Catalogs",

		"НумераторДокументов": "DocumentNumerators",
		"Документ":            "Documents",

		"ЖурналДокументов": "DocumentJournals",

		"Перечисление": "Enums",

		"Отчет":     "Reports",
		"Обработка": "DataProcessors",

		"ПланВидовХарактеристик ": "ChartsOfCharacteristicTypes",

		"РегистрСведений":   "InformationRegisters",
		"РегистрНакопления": "AccumulationRegisters",

		"БизнесПроцесс": "BusinessProcesses",
		"Задача":        "Tasks",
	}

}

func mapObjectNames() map[string]string {

	return map[string]string{

		"Команда": "Commands",
		"Форма":   "Forms",
		"Макет":   "Templates",
	}

}

func mapModulesPath() map[string][]string {

	return map[string][]string{
		"МодульУправляемогоПриложения": {"Ext", "ManagedApplicationModule"},
		"МодульСеанса":                 {"Ext", "SessionModule"},
		"ОбщийМодуль":                  {"Ext", "Module"},
		"HTTPСервис":                   {"Ext", "Module"},
		"WebСервис":                    {"Ext", "Module"},
		"СервисИнтеграции":             {"Ext", "Module"},
		"МодульМенеджера":              {"Ext", "ManagerModule"},
		"МодульОбъекта":                {"Ext", "ObjectModule"},
		"Форма":                        {"Ext", "Form", "Module"},
		"Команда":                      {"Ext", "CommandModule"},
	}

}
