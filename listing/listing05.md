Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```
Функция test() возвращает указатель на customError, который вскоре присваивается err - переменной типа error. Как описывалось в листинге 3, интерфейс хранит информацию о типе данных (в данном случае о *customError), значит переменная err не будет равна nil.

Ответ:
```console
error
```
