package async

import (
	"log"
	"sync"

	"github.com/rfruffer/go-musthave-shortener/internal/repository"
)

// DeleteTask предоставляет переменные для очереди.
type DeleteTask struct {
	UserID    string
	ShortURLs []string
}

// DeleteQueue очередь на удаление
var DeleteQueue = make(chan DeleteTask, 100)

// FanIn метод удаления задач из очереди
func FanIn(doneCh chan struct{}, resultChs ...chan DeleteTask) chan DeleteTask {
	finalCh := make(chan DeleteTask)

	var wg sync.WaitGroup

	for _, ch := range resultChs {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-doneCh:
					return
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}

// StartDeleteWorker воркер на удаление задач
func StartDeleteWorker(doneCh chan struct{}, repo repository.StoreRepositoryInterface, inputCh chan DeleteTask) {
	go func() {
		for {
			select {
			case <-doneCh:
				log.Println("Delete worker stopped")
				return
			case task := <-inputCh:
				if len(task.ShortURLs) == 0 {
					continue
				}
				err := repo.MarkURLsDeleted(task.UserID, task.ShortURLs)
				if err != nil {
					log.Printf("failed to delete URLs for user %s: %v", task.UserID, err)
				}
			}
		}
	}()
}
