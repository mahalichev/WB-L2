Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}
```
При получении данных из закрытого канала, переменной присваивается нулевое значение типа. Для проверки того, что данные были получены не из закрытого канала, используется проверка с помощью флага:
```go
value, ok := <-channel
if ok {
	...
}
```
В select функции merge не реализована данная проверка, поэтому после закрытия каналов a и b, в канал c будут передаваться 0.

Ответ:
```console
1
2
3
4
5
6
7
8
0
0
...
```
