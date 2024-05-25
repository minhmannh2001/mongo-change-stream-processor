package main

import (
	"fmt"
	"log"
	"time"

	"github.com/minhmannh2001/mongo-change-stream-processor/domain"
	"github.com/minhmannh2001/mongo-change-stream-processor/internal/pipeline"
)

type MultiplyTenSlow struct{}

func (m MultiplyTenSlow) Process(result domain.Message) ([]domain.Message, error) {
	time.Sleep(1 * time.Second)
	number := result.(int)
	return []domain.Message{number * 10, number * 10}, nil
}

type MultiplyHundredSlow struct{}

func (m MultiplyHundredSlow) Process(result domain.Message) ([]domain.Message, error) {
	time.Sleep(time.Duration(1 * time.Second))
	number := result.(int)
	return []domain.Message{number * 100, number * 100}, nil
}

type DivideThreeSlow struct{}

func (m DivideThreeSlow) Process(result domain.Message) ([]domain.Message, error) {
	time.Sleep(time.Duration(1 * time.Second))
	number := result.(int)
	return []domain.Message{number / 3}, nil
}

func pipeline_example() {

	p := pipeline.NewConcurrentPipeline()

	p.AddPipe(MultiplyHundredSlow{}, &pipeline.PipelineOpts{
		Concurrency: 5,
	})
	p.AddPipe(MultiplyTenSlow{}, &pipeline.PipelineOpts{
		Concurrency: 5,
	})
	p.AddPipe(DivideThreeSlow{}, &pipeline.PipelineOpts{
		Concurrency: 5,
	})

	if err := p.Start(); err != nil {
		log.Println(err)
	}

	for i := 1; i <= 3; i++ {
		p.Input() <- i
	}

	go func() {
		count := 0
		for number := range p.Output() {
			fmt.Println(number)
			count++
		}
	}()

	p.Stop()

}

// https://ketansingh.me/posts/pipeline-pattern-in-go-part-1/
