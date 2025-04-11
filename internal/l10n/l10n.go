package l10n

const (
	CmdDescStart                       = "открыть меню"
	KBBtnRunConfig                     = "Запустить конкурс из конфигурации"
	CfgDelimiter                       = " - "
	CfgMultiplicity                    = "кратность"
	CfgKeyword                         = "ключевое слово"
	CfgChatUsername                    = "чат"
	CfgChatID                          = "ид чата"
	CfgChannelUsername                 = "канал"
	CfgChannelID                       = "ид канала"
	CfgTopic                           = "топик"
	DefaultKeyword                     = "готов"
	YourTicketNumbers                  = "Ваши номера участника - "
	YourTicketNumbersDelimiter         = ", "
	ReactErrorPrefix                   = "Ошибка! "
	ReactErrorSuffix                   = "🫠" // emoji
	ContestStopNotFound                = "в этом чате нечего останавливать"
	ChatTakeNotFound                   = "чат не найден"
	ContestStopUsage                   = "Пример:\n/contestStop @exampleChatUsername"
	CreateContestNoAdminRights         = "требуются права администратора"
	CreateContestCantVerifyAdminRights = "невозможно проверить права администратора"
	ContestConfigRunSuccess            = "Конкурс запущен🎉"    // emoji
	ContestStopSuccess                 = "Конкурс остановлен👏" // emoji
	ContestConfigRunUsage              = `

*Для запуска конкурса, боту надо отправить сообщение строго в таком формате:*
/contestConfigRun
\[название параметра\] \- \[значение параметра\]

*Пример сообщения:*
` + "```" + `
/contestConfigRun
` + CfgMultiplicity + ` \- 10
` + CfgKeyword + ` \- Готово
` + CfgChatUsername + ` \- @exampleChatUsername
` + CfgTopic + ` \- 1
` + "```" + `
*Описание параметров:*
__` + CfgMultiplicity + `__ \- обязательный числовой параметр, отвечает за необходимое количество приглашенных участников для получения билета\.
__` + CfgKeyword + `__ \- обязательный текстовый параметр, слово которое необходимо написать в чат или топик для подсчета билетов\.
__` + CfgChatID + `__ \- обязательный, если не указан параметр "` + CfgTopic + `"\. Числовой параметр, позволяет определить чат, в котором будет производиться конкурс\.
__` + CfgChatUsername + `__ \- обязательный, если не указан параметр "` + CfgChatID + `"\. Параметр вида @exampleChatUsername, чат в котором будет проводиться конкурс\. Не всегда удается по этому параметру определить чат, в этом случае лучше использовать "` + CfgChatID + `"\.
__` + CfgTopic + `__ \- опциональный параметр, идентификатор топика, в котором бот будет выдавать номерки, если не указан или равен 0, то писать ключевое слово надо в главный чат`
	ContestCreatePreviousNotOverYet = "в этом чате уже проходит конкурс, невозможно запустить еще один"
	RequestedDataNotFound           = "запрашиваемые данные не найдены"
	DintGetRightNumberOfInvitations = "Не набралось нужное количество приглашений"
	ParameterNotProvided            = "параметр не передан"
	ContestConfigBotCannotSendMsg   = "бот не может писать в этот чат (топик)"
)
