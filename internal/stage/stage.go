package stage

import (
	"log"
	"sync"

	"github.com/minhmannh2001/mongo-change-stream-processor/domain"
)

type Stage interface {
	Process(message domain.Message) ([]domain.Message, error)
}

type StageWorker struct {
	wg          *sync.WaitGroup
	input       chan domain.Message
	output      chan domain.Message
	concurrency int
	pipe        Stage
}

func (w *StageWorker) Start() error {

	for i := 0; i < w.concurrency; i++ {
		w.wg.Add(1)

		go func() {
			defer w.wg.Done()
			for i := range w.Input() {
				result, err := w.pipe.Process(i)
				if err != nil {
					log.Println(err)
					continue
				}
				for _, r := range result {
					w.Output() <- r
				}
			}
		}()
	}

	return nil
}

func (w *StageWorker) WaitStop() error {
	w.wg.Wait()
	return nil
}

func (w *StageWorker) Input() chan domain.Message {
	return w.input
}

func (w *StageWorker) Output() chan domain.Message {
	return w.output
}

func NewWorkerGroup(concurrency int, pipe Stage, input chan domain.Message, output chan domain.Message) StageWorker {

	return StageWorker{
		wg:          &sync.WaitGroup{},
		input:       input,
		output:      output,
		concurrency: concurrency,
		pipe:        pipe,
	}

}
