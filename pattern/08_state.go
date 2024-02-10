package pattern

import "fmt"

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/*
Состояние - поведенческий паттерн проектирования, позволяющий объектам менять поведение в зависимости от внутреннего состояния.

Применимость:
Состояние применяется в случаях, когда:
- Существует объект, поведение которого полностью меняется от его внутреннего состояния, которые должны меняться динамически;
- Множество больших похожих условных операторов, выбирающих поведение от внутреннего состояния объекта.

Плюсы:
- Избавление от множества условных операторов;
- Упрощение кода контекста;
- Локализация состояния;
- Облегчение добавления нового состояния.

Минусы:
- Усложнение кода программы дополнительными классами.

Пример использования:
Компьютерная игра, в которой объекты могут находиться в различных состояниях (например, персонаж может бегать, ходить, плавать, говорить)
*/

// Интерфейс состояния
type State interface {
	MoveForward(*StateContext)
	MoveBack(*StateContext)
}

// Контекст состояния
type StateContext struct {
	State
}

func (context *StateContext) SetState(state State) {
	context.State = state
}

func (context *StateContext) MoveForward() {
	context.State.MoveForward(context)
}

func (context *StateContext) MoveBack() {
	context.State.MoveBack(context)
}

// Конкретное состояние - нахождение на старте
type StartState struct{}

func (state *StartState) MoveForward(context *StateContext) {
	fmt.Println("Moving forward from start")
	context.SetState(&OnTheWayState{})
}

func (state *StartState) MoveBack(context *StateContext) {
	fmt.Println("You are already at the start")
}

// Конкретное состояние - нахождение на середине пути
type OnTheWayState struct{}

func (state *OnTheWayState) MoveForward(context *StateContext) {
	fmt.Println("Moving towards the finish line")
	context.SetState(&FinishState{})
}

func (state *OnTheWayState) MoveBack(context *StateContext) {
	fmt.Println("Heading back to the start")
	context.SetState(&StartState{})
}

// Конкретное состояние - нахождение на финише
type FinishState struct{}

func (state *FinishState) MoveForward(context *StateContext) {
	fmt.Println("You are already at the finish")
}

func (state *FinishState) MoveBack(context *StateContext) {
	fmt.Println("Going back from finish")
	context.SetState(&OnTheWayState{})
}

// func main() {
// 	context := &StateContext{}
// 	context.SetState(&StartState{})
// 	context.MoveForward()
// 	context.MoveBack()
// 	context.MoveForward()
// 	context.MoveForward()
// 	context.MoveForward()
// 	context.MoveBack()
// 	context.MoveBack()
// 	context.MoveForward()
// }
