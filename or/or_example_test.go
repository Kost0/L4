package or

import (
	"fmt"
	"time"
)

func ExampleOr() {
	// Функция, которая создает канал с отложенным закрытием
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	// Объединяем 5 каналов со временем закрытия 2 часа, 5 минут, 1 секунда, 1 час, 1 минута и ждем закрытия первого из них
	<-Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// объединенный канал закроется, когда закроется канал с наименьшим временем, то есть 1 секунда
	elapsed := time.Since(start)

	fmt.Println("Объединенный канал закрылся")

	if elapsed > 1100*time.Millisecond {
		fmt.Printf("Внимание: прошло слишком много времени: %v", elapsed)
	}

	// Output:
	// Объединенный канал закрылся
}
