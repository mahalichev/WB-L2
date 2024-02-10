package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
Посетитель - поведенческий паттерн проектирования, который отделяет алгоритм от структуры объекта.
Посетитель позволяет добавлять новые функции, не изменяя классы, над которыми эти операции выполняются.

Применимость:
Посетитель применяется в случаях, когда:
- Необходимо добавить операции к классам, при этом изменение классов нежелательно;
- Необходимо выполнить операцию над элементами сложной структуры объектов;
- Операция необходима только для некоторых классов.

Плюсы:
- Упрощение добавления операций для сложных элементов;
- Объединение родственных операций в одном классе;
- Возможно накапливание состояния при обходе структуры элементов.

Минусы:
- Возможно нарушение инкапсуляции;
- Посетитель не имеет смысла, если иерархия элементов часто меняется.

Пример использования:
Экспорт файлов в различных форматах.
*/

// Интерфейс, описывающий структуру, принимающую посетителя
type Acceptor interface {
	Accept(Visitor)
}

type Visitor interface {
	VisitForApple(*Apple)
	VisitForOrange(*Orange)
	VisitForWatermelon(*Watermelon)
}

// Конкретный класс яблоко
type Apple struct{}

func (apple *Apple) Accept(visitor Visitor) {
	visitor.VisitForApple(apple)
}

// Конкретный класс апельсин
type Orange struct{}

func (orange *Orange) Accept(visitor Visitor) {
	visitor.VisitForOrange(orange)
}

// Конкретный класс арбуз
type Watermelon struct{}

func (watermelon *Watermelon) Accept(visitor Visitor) {
	visitor.VisitForWatermelon(watermelon)
}

// Посетитель, добавляющий операцию мытья продуктов
type Washer struct{}

func (washer *Washer) VisitForApple(apple *Apple) {
	fmt.Println("Washing apple")
}

func (washer *Washer) VisitForOrange(orange *Orange) {
	fmt.Println("Washing orange")
}

func (washer *Washer) VisitForWatermelon(watermelon *Watermelon) {
	fmt.Println("Washing watermelon")
}

// Посетитель, добавляющий операцию поедания продуктов
type Eater struct{}

func (eater *Eater) VisitForApple(apple *Apple) {
	fmt.Println("Eating apple")
}

func (eater *Eater) VisitForOrange(orange *Orange) {
	fmt.Println("Eating orange")
}

func (eater *Eater) VisitForWatermelon(watermelon *Watermelon) {
	fmt.Println("Eating watermelon")
}

// func main() {
// 	food := []Acceptor{&Apple{}, &Orange{}, &Watermelon{}}

// 	washer := &Washer{}
// 	eater := &Eater{}
// 	for _, acceptor := range food {
// 		acceptor.Accept(washer)
// 		acceptor.Accept(eater)
// 	}
// }
