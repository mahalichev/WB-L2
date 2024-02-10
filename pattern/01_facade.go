package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
Фасад - структурный паттерн проектирования, предоставляющий простой интерфейс к сложной системе.
Фасад применяется там, где необходимо скрыть сложную систему и свести все вызовы к одному объекту.

Применимость:
Фасад применяется в случаях, когда:
- Необходимо упростить работу с сложной системой, предоставив только то, что необходимо клиенту;
- Необходимо упростить взаимодействие подсистем и при этом минимизировать зависимости - разложение подсистемы на слои.

Плюсы:
- Изолирование клиента от компонентов сложной подсистемы;
- Минимизация зависимостей.

Минусы:
- Риск создания божественного объекта;

Пример использования:
Взаимодействие с библиотекой/API - фасад может облегчить взаимодействие с библиотекой/API благодаря предоставлению упрощенного интерфейса
*/

// Подсистема "Камера"
type Camera struct{}

func (camera *Camera) TakePhoto() string {
	fmt.Println("Camera took a photo")
	return "<<photo>>"
}

func (camera *Camera) FilmVideo() string {
	fmt.Println("Camera took a video")
	return "<<video>>"
}

// Подсистема "Телефон"
type Phone struct{}

func (phone *Phone) Call(optionalData ...string) {
	fmt.Println("Phone called")
	for _, data := range optionalData {
		fmt.Println("- data:", data)
	}
}

// Подсистема "Мессенджер"
type Messager struct{}

func (messager *Messager) SendMessage(message string) {
	fmt.Println("Message sent:", message)
}

// Фасад - смартфон
type SmartphoneFacade struct {
	camera   *Camera
	phone    *Phone
	messager *Messager
}

func (smartphone *SmartphoneFacade) VideoCall() {
	fmt.Println("Video call...")
	video := smartphone.camera.FilmVideo()
	smartphone.phone.Call(video)
}

func (smartphone *SmartphoneFacade) SendPhoto() {
	fmt.Println("Send photo...")
	photo := smartphone.camera.TakePhoto()
	smartphone.messager.SendMessage(photo)
}

func NewSmartphone() *SmartphoneFacade {
	return &SmartphoneFacade{
		camera:   &Camera{},
		phone:    &Phone{},
		messager: &Messager{},
	}
}

// func main() {
// 	smartphone := NewSmartphone()
// 	smartphone.SendPhoto()
// 	smartphone.VideoCall()
// }
