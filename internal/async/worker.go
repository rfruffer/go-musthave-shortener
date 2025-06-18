package async

import (
	"log"
	"sync"

	"github.com/rfruffer/go-musthave-shortener/internal/repository"
)

type DeleteTask struct {
	UserID    string
	ShortURLs []string
}

var DeleteQueue = make(chan DeleteTask, 100)

func fanIn(doneCh chan struct{}, resultChs ...chan DeleteTask) chan DeleteTask {
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
