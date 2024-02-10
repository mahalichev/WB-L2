package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
Фабричный метод - порождающий паттерн проектирования, определяющий интерфейс создания объектов, при том что конкретный тип создаваемого объекта
определяется в подклассах.

Применимость:
Фабричный метод применяется в случаях, когда:
- Заранее неизвестно, объекты каких типов необходимо создавать;
- Система не должна быть зависима от создания новых объектов и должна быть расширяемой новыми типами данных;
- Необходимо делегировать создание объекта подклассу.

Плюсы:
- Устранение привязки к конкретным классам продуктов;
- Инкапсулирование код создания продукта;
- Упрощение расширения новыми типами данных.

Минусы:
- Усложнение кода новыми классами-фабриками;

Пример использования:
Приложения транспортной логистики. При добавлении новой логистики будет достаточно создать класс, удовлетворяющий интерфейс, и фабрику для данного типа.
*/

// Интерфейс продукта
type Messenger interface {
	SendMessage()
	RecieveMessage()
}

// Конкретный продукт - VK
type VK struct{}

func (vk *VK) SendMessage() {
	fmt.Println("VK: Send message")
}

func (vk *VK) RecieveMessage() {
	fmt.Println("VK: Recieve message")
}

// Конкретный продукт - Telegram
type Telegram struct{}

func (telegram *Telegram) SendMessage() {
	fmt.Println("Telegram: Send message")
}

func (telegram *Telegram) RecieveMessage() {
	fmt.Println("Telegram: Recieve message")
}

// Интерфейс создателя объектов
type MessengerCreator interface {
	CreateMessanger() Messenger
}

// Конкретный создатель VK
type VKCreator struct{}

func (vkCreator VKCreator) CreateMessanger() Messenger {
	return &VK{}
}

// Конкретный создатель Telegram
type TelegramCreator struct{}

func (telegramCreator TelegramCreator) CreateMessanger() Messenger {
	return &Telegram{}
}

// func main() {
// 	var creator MessengerCreator = VKCreator{}
// 	vk := creator.CreateMessanger()
// 	vk.SendMessage()
// 	vk.RecieveMessage()

// 	creator = TelegramCreator{}
// 	telegram := creator.CreateMessanger()
// 	telegram.SendMessage()
// 	telegram.RecieveMessage()
// }
