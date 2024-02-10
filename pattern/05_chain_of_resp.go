package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Цепочка обязанностей - поведенческий паттерн проектирования, который позволяет передавать запросы последовательно по цепочке обработчиков.
Каждый обработчик содержит логику, определяющую какие запросы он может обработать, остальные передаются дальше по цепочке.

Применимость:
Цепочка обязанностей применяется в случаях, когда:
- Программа должна обрабатывать различные запросы, но неизвестно, какие именно будут приходить запросы и какие обработчики понадобятся;
- Важен порядок выполнения обработчиков;
- Набор обработчиков должен задаваться динамически.

Плюсы:
- Уменьшение зависимости между клиентом и обработчиком;
- Упрощённое добавление новых обработчиков;
- Гибкое назначение обязанностей.

Минусы:
- Запрос может остаться необработанным;
- Увеличение времени обработки запроса.

Пример использования:
Обработка http-запросов, где необходимы такие обработчики, как аутентификация, авторизация, валидация, кеширование и т.д.
*/

// Интерфейс обработчика
type Handler interface {
	Handle(string)
	SetNext(Handler)
}

// Конкретный обработчик - senior
type SeniorHandler struct {
	Next Handler
}

func (handler *SeniorHandler) Handle(question string) {
	fmt.Println("Senior successfully answers the question")
}

func (handler *SeniorHandler) SetNext(nextHandler Handler) {
	handler.Next = nextHandler
}

// Конкретный обработчик - middle
type MiddleHandler struct {
	Next Handler
}

func (handler *MiddleHandler) Handle(question string) {
	if question == "how to create a high load service?" {
		fmt.Println("Middle successfully answers the question")
	} else {
		handler.Next.Handle(question)
	}
}

func (handler *MiddleHandler) SetNext(nextHandler Handler) {
	handler.Next = nextHandler
}

// Конкретный обработчик - junior
type JuniorHandler struct {
	Next Handler
}

func (handler *JuniorHandler) Handle(question string) {
	if question == "how to set up grpc?" {
		fmt.Println("Junior successfully answers the question")
	} else {
		handler.Next.Handle(question)
	}
}

func (handler *JuniorHandler) SetNext(nextHandler Handler) {
	handler.Next = nextHandler
}

// func main() {
// 	juniorHandler := &JuniorHandler{}
// 	middleHandler := &MiddleHandler{}
// 	seniorHandler := &SeniorHandler{}
// 	juniorHandler.SetNext(middleHandler)
// 	middleHandler.SetNext(seniorHandler)

// 	juniorHandler.Handle("how to set up grpc?")
// 	juniorHandler.Handle("how to create a high load service?")
// 	juniorHandler.Handle("very hard question")
// }
