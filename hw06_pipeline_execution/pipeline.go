package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done Bi, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		intermediateCh := make(Bi)
		go func(outCh Out, intermedCh Bi) {
			defer func() {
				close(intermedCh)
				for range outCh {
				}
			}()
			for {
				select {
				case <-done:
					return
				case v, ok := <-outCh:
					if !ok {
						return
					}
					intermedCh <- v
				}
			}
		}(out, intermediateCh)
		out = stage(intermediateCh)
	}

	return out
}
