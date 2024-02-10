package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
Стратегия - поведенческий паттерн проектирования, который определяет набор алгоритмов, инкапсулирует их и обеспечивает взаимозаменяемость.
Сама замена алгоритмов происходит независимо от объекта, который использует данные алгоритмы.

Применимость:
Стратегия применяется в случаях, когда:
- Необходимо использовать разные вариации алгоритма внутри объекта;
- Есть множество схожих классов, отличающихся некоторым поведением;
- Необходимо менять алгоритмы во время выполнения программы;
- Необходимо инкапсулировать реализацию алгоритмов от других классов.

Плюсы:
- Переиспользование кода;
- Инкапсулирование кода и данных алгоритма от классов;
- Уменьшение условных операторов и ветвлений кода;
- Упрощение замены алгоритмов.

Минусы:
- Усложнение кода программы дополнительными классами;
- Клиент должен знать о различных стратегиях.

Пример использования:
Веб-разработка - выбор способа рендеринга сайта в зависимости от разрешения и устройства.
*/

// Интерфейс стратегии
type ConversationStrategy interface {
	Talk()
}

// Контекст для хранения стратегии
type ConversationContext struct {
	ConversationMethod ConversationStrategy
}

func (context *ConversationContext) SetStrategy(strategy ConversationStrategy) {
	context.ConversationMethod = strategy
}

func (context *ConversationContext) Execute() {
	context.ConversationMethod.Talk()
}

// Конкретная стратегия - разговор с использованием компьютера
type ConversationUsingComputer struct{}

func (conversation *ConversationUsingComputer) Talk() {
	fmt.Println("Conversation via computer")
}

// Конкретная стратегия - разговор с использованием телефона
type ConversationUsingPhone struct{}

func (conversation *ConversationUsingPhone) Talk() {
	fmt.Println("Conversation via phone")
}

// Конкретная стратегия - разговор с использованием рации
type ConversationUsingWalkieTalkie struct{}

func (conversation *ConversationUsingWalkieTalkie) Talk() {
	fmt.Println("Conversation via walkie-talkie")
}

// func main() {
// 	conversationUsingComputer := &ConversationUsingComputer{}
// 	conversationUsingPhone := &ConversationUsingPhone{}
// 	ConversationUsingWalkieTalkie := &ConversationUsingWalkieTalkie{}

// 	conversationContext := &ConversationContext{conversationUsingComputer}
// 	conversationContext.Execute()
// 	conversationContext.SetStrategy(conversationUsingPhone)
// 	conversationContext.Execute()
// 	conversationContext.SetStrategy(ConversationUsingWalkieTalkie)
// 	conversationContext.Execute()
// }
