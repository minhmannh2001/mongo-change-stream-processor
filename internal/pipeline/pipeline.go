package pipeline

import (
	"github.com/minhmannh2001/mongo-change-stream-processor/domain"
	"github.com/minhmannh2001/mongo-change-stream-processor/internal/common/errors"
	"github.com/minhmannh2001/mongo-change-stream-processor/internal/stage"
)

type PipelineOpts struct {
	Concurrency int
}

type Pipeline interface {
	AddPipe(pipe stage.Stage, opt *PipelineOpts)
	Start() error
	Stop() error
	Input() chan<- domain.Message
	Output() <-chan domain.Message
}

type ConcurrentPipeline struct {
	stageWorkers []stage.StageWorker
}

func NewConcurrentPipeline() Pipeline {
	return &ConcurrentPipeline{}
}

func (c *ConcurrentPipeline) AddPipe(_stage stage.Stage, opt *PipelineOpts) {

	if opt == nil {
		opt = &PipelineOpts{Concurrency: 1}
	}

	var input = make(chan domain.Message, 10)
	var output = make(chan domain.Message, 10)

	for _, i := range c.stageWorkers {
		input = i.Output()
	}

	worker := stage.NewWorkerGroup(opt.Concurrency, _stage, input, output)
	c.stageWorkers = append(c.stageWorkers, worker)
}

func (c *ConcurrentPipeline) Input() chan<- domain.Message {
	return c.stageWorkers[0].Input()
}

func (c *ConcurrentPipeline) Output() <-chan domain.Message {
	sz := len(c.stageWorkers)
	return c.stageWorkers[sz-1].Output()
}

func (c *ConcurrentPipeline) Start() error {

	if len(c.stageWorkers) == 0 {
		return errors.ErrConcurrentPipelineEmpty
	}

	for i := 0; i < len(c.stageWorkers); i++ {
		g := c.stageWorkers[i]
		g.Start()
	}

	return nil
}

func (c *ConcurrentPipeline) Stop() error {

	for _, i := range c.stageWorkers {
		close(i.Input())
		i.WaitStop()
	}

	sz := len(c.stageWorkers)
	close(c.stageWorkers[sz-1].Output())
	return nil
}
