Что выведет программа? Объяснить вывод программы.

```go
package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		println(n)
	}
}
```
Цикл for получает данные из канала до того момента, пока канал не будет закрыт. Так как не происходит вызова close(ch), программа не выйдет из ожидания получения данных из канала, которые не будут в него поступать. Произойдет deadlock.

Ответ:
```
0
1
2
3
4
5
6
7
8
9
fatal error: all goroutines are asleep - deadlock!
```
