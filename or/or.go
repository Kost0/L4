package or

import (
	"sync"
)

func Or(channels ...<-chan interface{}) <-chan interface{} {
	// Обработка крайних случаев
	if len(channels) == 0 {
		c := make(chan interface{})
		close(c)
		return c
	} else if len(channels) == 1 {
		return channels[0]
	}

	// Создаем канал для объединения
	orDone := make(chan interface{})

	var once sync.Once

	// Для каждого канала запускаем горутину
	for _, channel := range channels {
		go func(c <-chan interface{}) {
			select {
			// Если закрывается один из каналов
			case <-c:
				// Используем once для безопасного закрытия
				once.Do(func() {
					close(orDone)
				})
			// После закрытия orDone, завершаем горутину
			case <-orDone:
			}
		}(channel)
	}

	return orDone
}
