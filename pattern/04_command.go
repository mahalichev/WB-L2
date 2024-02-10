package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
Команда - поведенческий паттерн проектирования, представляющий запросы в виде объектов.
Объект используется для инкапсуляции информации, необходимой для отложенного выполнения запросов.

Применимость:
Команда применяется в случаях, когда:
- Необходима параметризация объекта выполняемым действием;
- Необходима очередь операций, отложенный запуск операций;
- Необходима возможность отмены операции.

Плюсы:
- Убирание зависимости между объектом, вызывающим операцию, и объектом, выполняющим операцию;
- Возможность реализации отмены и повтора операции;
- Возможность создания очереди операций, отложенных операций;
- Создание сложных операций на основе простых.

Минусы:
- Множество дополнительных классов, усложняющих код.

Пример использования:
Графические пользовательские интерфейсы, в которых множество кнопок, выполняющих различные операции.
*/

// Интерфейс получателя
type AudioDevice interface {
	VolumeUp()
	VolumeDown()
	PrintVolume()
}

// Конкретный получатель - радио
type Radio struct {
	Volume int
}

func (radio *Radio) VolumeUp() {
	radio.Volume++
}

func (radio *Radio) VolumeDown() {
	radio.Volume--
}

func (radio Radio) PrintVolume() {
	fmt.Printf("Radio volume - %d\n", radio.Volume)
}

// Интерфейс команды
type Command interface {
	Execute()
}

// Конкретная команда - увеличение громкости
type VolumeUpCommand struct {
	AudioDevice
}

func (volumeUpCommand VolumeUpCommand) Execute() {
	volumeUpCommand.AudioDevice.VolumeUp()
	volumeUpCommand.AudioDevice.PrintVolume()
}

// Конкретная команда - уменьшение громкости
type VolumeDownCommand struct {
	AudioDevice
}

func (volumeUpCommand VolumeDownCommand) Execute() {
	volumeUpCommand.AudioDevice.VolumeDown()
	volumeUpCommand.AudioDevice.PrintVolume()
}

// Отправитель команды - кнопка
type Button struct {
	Command
}

func (button Button) Press() {
	button.Command.Execute()
}

// func main() {
// 	radio := &Radio{}
// 	volumeUpCommand := &VolumeUpCommand{radio}
// 	volumeDownCommand := &VolumeDownCommand{radio}

// 	volumeUpButton := &Button{volumeUpCommand}
// 	volumeDownButton := &Button{volumeDownCommand}
// 	volumeUpButton.Press()
// 	volumeDownButton.Press()
// }
