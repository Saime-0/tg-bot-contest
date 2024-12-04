package l10n

const (
	CmdDescStart               = "открыть меню"
	KBBtnRunConfig             = "Запустить конкурс из конфигурации"
	CfgDelimiter               = " - "
	CfgMultiplicity            = "кратность"
	CfgKeyword                 = "ключевое слово"
	CfgChatUsername            = "чат"
	CfgTopic                   = "топик"
	DefaultKeyword             = "готов"
	YourTicketNumbers          = "Ваши номера участника - "
	YourTicketNumbersDelimiter = ", "
	ReactErrorPrefix           = "Ошибка! "
	ReactErrorSuffix           = "🫠" // emoji
	ContestStopNotFound        = "в этом чате нечего останавливать"
	ContestStopUsage           = "Пример:\n/contestStop @exampleChatUsername"
	CreateContestNoAdminRights = "требуются права администратора"
	ContestConfigRunSuccess    = "Конкурс запущен🎉"    // emoji
	ContestStopSuccess         = "Конкурс остановлен👏" // emoji
	ContestConfigRunUsage      = `

Для запуска конкурса, боту надо отправить сообщение строго в таком формате:
/contestConfigRun
[название параметра] - [значение параметра]

Пример сообщения:
/contestConfigRun
кратность - 10
ключевое слово - Готово
чат - @exampleChatUsername
топик - 1

Описание параметров:
кратность - обязательный числовой параметр, отвечает за необходимое количество приглашенных участников для получения билета.
ключевое слово - обязательный текстовый параметр, слово которое необходимо написать в чат или топик для подсчета билетов.
чат - обязательный параметр вида @exampleChatUsername, чат в котором будет проводиться конкурс
топик - опциональный параметр, идентификатор топика, в котором бот будет выдавать номерки, если не указан или равен 0, то писать ключевое слово надо в главный чат`
	ContestCreatePreviousNotOverYet = "в этом чате уже проходит конкурс, невозможно запустить еще один"
	RequestedDataNotFound           = "запрашиваемые данные не найдены"
	DintGetRightNumberOfInvitations = "не набралось нужное количество приглашений"
)
