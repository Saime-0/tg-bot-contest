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
` + CfgChannelUsername + ` \- @exampleChannelUsername
` + CfgChatUsername + ` \- @exampleChatUsername
` + CfgTopic + ` \- 1
` + "```" + `
*Описание параметров:*
__` + CfgMultiplicity + `__ \- обязательный числовой параметр, отвечает за необходимое количество приглашенных участников для получения номерков\.
__` + CfgKeyword + `__ \- обязательный текстовый параметр, слово которое необходимо написать в чат или топик для подсчета номерков\.
__` + CfgTopic + `__ \- опциональный параметр, идентификатор топика, в котором бот будет выдавать номерки, если не указан или равен 0, то писать ключевое слово надо в главный чат \(параметр __` + CfgChatUsername + `__\)\.
__` + CfgChatUsername + `__ \- обязательный, если не указан параметр __` + CfgChatID + `__\. Чат в котором будет проводиться конкурс\. Не всегда удается по этому параметру определить чат, в этом случае лучше использовать __` + CfgChatID + `__\.
__` + CfgChannelUsername + `__ \(либо __` + CfgChannelID + `__\) \- опциональный параметр\. Канал в котором будет проводиться конкурс\(вместо чата\)\. Участники должны будут приглашать друзей в него, а не в чат, но ключевое слово все также надо писать в чат\.
`
	ContestCreatePreviousNotOverYet = "в этом чате уже проходит конкурс, невозможно запустить еще один"
	RequestedDataNotFound           = "запрашиваемые данные не найдены"
	DintGetRightNumberOfInvitations = "Не набралось нужное количество приглашений"
	ParameterNotProvided            = "параметр не передан"
	ContestConfigBotCannotSendMsg   = "бот не может писать в этот чат (топик)"
)
