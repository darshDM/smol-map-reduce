package mapreduce

import (
	"sync"
)

type MapReduce[Input, Intermediate, Output any] struct {
	Pool            chan Input
	OutputCollector chan Intermediate
	Mapper          MapperFunc[Input, Intermediate]
	Reducer         ReduceFunc[Intermediate, Output]
	wg              sync.WaitGroup
	reduceWg        sync.WaitGroup
	WorkerCount     int
	Result          Output
}

type MapperFunc[Input, Intermediate any] func(data Input) *Intermediate

type ReduceFunc[Intermediate, Output any] func(intermediates []*Intermediate) Output

func NewMapReduce[Input, Intermediate, Output any](
	worker int,
	mapper MapperFunc[Input, Intermediate],
	reducer ReduceFunc[Intermediate, Output],
) *MapReduce[Input, Intermediate, Output] {
	return &MapReduce[Input, Intermediate, Output]{
		Pool:            make(chan Input, 100),
		OutputCollector: make(chan Intermediate, 100),
		Mapper:          mapper,
		Reducer:         reducer,
		WorkerCount:     worker,
	}
}

func (mr *MapReduce[Input, Intermediate, Output]) Start() {
	for i := 0; i < mr.WorkerCount; i++ {
		mr.wg.Add(1)
		go func() {
			defer mr.wg.Done()
			for data := range mr.Pool {
				result := mr.Mapper(data)
				if result != nil {
					mr.OutputCollector <- *result
				}
			}
		}()
	}

	mr.reduceWg.Add(1)
	go func() {
		defer mr.reduceWg.Done()
		intermediates := make([]*Intermediate, 0)
		for record := range mr.OutputCollector {
			recordCopy := record // Create a copy to avoid pointer issues
			intermediates = append(intermediates, &recordCopy)
		}
		mr.Result = mr.Reducer(intermediates)
	}()
}

func (mr *MapReduce[Input, Intermediate, Output]) Stop() {
	close(mr.Pool)
	mr.wg.Wait()
	close(mr.OutputCollector)
	mr.reduceWg.Wait()
}

func (mr *MapReduce[Input, Intermediate, Output]) GetResult() Output {
	return mr.Result
}
