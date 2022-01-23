package wpool

import (
	"context"
)

type JobID string
type jobType string

type ExecutionFn func(url string) ([]byte, error)

type JobDescriptor struct {
	ID       JobID
	JType    jobType
	Metadata map[string]string
}

type Result struct {
	Value      []byte
	Err        error
	Descriptor JobDescriptor
}

type Job struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn
	Args       string
}

func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFn(j.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}

	return Result{
		Value:      value,
		Descriptor: j.Descriptor,
	}
}
