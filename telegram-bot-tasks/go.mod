module taskbot

go 1.16

require (
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/mattn/go-sqlite3 v1.14.28

)

// это надо для переопределения адреса сервера
// в оригинале урл телеграма в константе, у меня там строка, которую я переопределяю в тесте
// replace gopkg.in/telegram-bot-api.v4 => ./local/telegram-bot-api.v4
replace github.com/go-telegram-bot-api/telegram-bot-api => ./local/telegram-bot-api
