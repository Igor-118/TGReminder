package service

import ()

// Reply — что показать пользователю.
type Reply struct {
	Text    string
	Options []Option // варианты на выбор; пусто — просто текст
}

// Option — один вариант ответа.
type Option struct {
	Label string // что видит пользователь: "Москва"
	Value string // что вернётся в сервис при выборе: "Europe/Moscow"
}
