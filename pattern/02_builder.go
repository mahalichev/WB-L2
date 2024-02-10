package pattern

/*
	Реализовать паттерн «строитель».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
Строитель - порождающий паттерн проектирования, позволяющий пошагово создавать сложные объекты.
Строитель позволяет отделить построение сложного объекта от его представления.

Применимость:
Строитель применяется в случаях, когда:
- Создание объекта состоит из нескольких шагов и необходимо использовать один процесс строительства несколько раз;
- Код должен создавать разные представления объекта;
- Необходимо инкапсулировать построение сложного объекта.

Плюсы:
- Пошаговое создание объекта;
- Отделение кода сборки от бизнес логики;
- Использование одного и того же кода для создания различных объектов.

Минусы:
- Привязка клиента к конкретным классам-строителям;
- Дополнительные классы, усложняющие код.

Пример использования:
Конфигуратор продукта с клиентскими параметрами, где пользователь может выбрать спецификацию продукта.
*/

// Продукт - компьютер
type Computer struct {
	CPU         string
	GPU         string
	RAM         string
	Motherboard string
	PSU         string
	Drive       string
}

// Интерфейс строителя
type ComputerBuilder interface {
	SetCPU()
	SetGPU()
	SetRAM()
	SetMotherboard()
	SetPSU()
	SetDrive()
	GetComputer() Computer
}

// Конкретный строитель
type CheapComputerBuilder struct {
	Computer
}

func (computer *CheapComputerBuilder) SetCPU() {
	computer.CPU = "AMD Ryzen 5 5600X"
}

func (computer *CheapComputerBuilder) SetGPU() {
	computer.GPU = "Maxsun AMD Radeon RX 550 4GB"
}

func (computer *CheapComputerBuilder) SetRAM() {
	computer.RAM = "Corsair Vengeance LPX 16 GB DDR4-3600 CL18"
}

func (computer *CheapComputerBuilder) SetMotherboard() {
	computer.Motherboard = "Gigabyte B550M K Micro ATX AM4"
}

func (computer *CheapComputerBuilder) SetPSU() {
	computer.PSU = "Thermaltake Smart BM2 550W"
}

func (computer *CheapComputerBuilder) SetDrive() {
	computer.Drive = "Kingston NV2 500G M.2"
}

func (computer *CheapComputerBuilder) GetComputer() Computer {
	return computer.Computer
}

func NewCheapComputerBuilder() *CheapComputerBuilder {
	return &CheapComputerBuilder{}
}

// Конкретный строитель
type ExpensiveComputerBuilder struct {
	Computer
}

func (computer *ExpensiveComputerBuilder) SetCPU() {
	computer.CPU = "Intel Core i9-14900K"
}

func (computer *ExpensiveComputerBuilder) SetGPU() {
	computer.GPU = "MSI GeForce RTX 4090 Suprim Liquid X 24G"
}

func (computer *ExpensiveComputerBuilder) SetRAM() {
	computer.RAM = "G.Skill 2x32GB Trident Z5 RGB DDR5-6800 CL34"
}

func (computer *ExpensiveComputerBuilder) SetMotherboard() {
	computer.Motherboard = "ASUS ROG Strix Z790-E Gaming WiFi II"
}

func (computer *ExpensiveComputerBuilder) SetPSU() {
	computer.PSU = "be quiet! Straight Power 11 Platinum 1500W"
}

func (computer *ExpensiveComputerBuilder) SetDrive() {
	computer.Drive = "Crucial T700 4TB PCIe Gen 5"
}

func (computer *ExpensiveComputerBuilder) GetComputer() Computer {
	return computer.Computer
}

func NewExpensiveComputerBuilder() *ExpensiveComputerBuilder {
	return &ExpensiveComputerBuilder{}
}

// Получение конкретного строителя
func GetComputerBuilder(builderType string) ComputerBuilder {
	switch builderType {
	case "expensive":
		return NewExpensiveComputerBuilder()
	case "cheap":
		return NewCheapComputerBuilder()
	default:
		return nil
	}
}

// Директор, в котором определена последовательность вызова строительных шагов
type Director struct {
	ComputerBuilder
}

func (director *Director) GetComputer() Computer {
	return director.ComputerBuilder.GetComputer()
}

func (director *Director) SetComputerBuilder(builder ComputerBuilder) {
	director.ComputerBuilder = builder
}

func (director *Director) BuildComputer() {
	director.SetCPU()
	director.SetGPU()
	director.SetRAM()
	director.SetMotherboard()
	director.SetPSU()
	director.SetDrive()
}

func NewDirector(builder ComputerBuilder) *Director {
	return &Director{builder}
}

// func main() {
// 	director := NewDirector(GetComputerBuilder("cheap"))
// 	director.BuildComputer()
// 	fmt.Println(director.GetComputer())

// 	director.SetComputerBuilder(GetComputerBuilder("expensive"))
// 	director.BuildComputer()
// 	fmt.Println(director.GetComputer())
// }
