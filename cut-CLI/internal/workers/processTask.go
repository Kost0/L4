package workers

import (
	"context"
	"log"

	"github.com/Kost0/L4/internal/cut"
	tasks2 "github.com/Kost0/L4/internal/gen/tasks"
)

type WorkerServer struct {
	tasks2.UnimplementedDistributedCutServer
}

func (w *WorkerServer) ProcessTask(ctx context.Context, req *tasks2.ProcessTaskRequest) (*tasks2.ProcessTaskResponse, error) {
	fields := make([]int, 0, len(req.Options.Fields))
	for _, field := range req.Options.Fields {
		fields = append(fields, int(field))
	}

	opts := cut.CutOptions{
		Fields:    fields,
		Delimiter: req.Options.Delimiter,
		Separated: req.Options.Separated,
	}

	lines := req.Lines

	log.Printf("Получена задача: %d строк, опции: %+v", len(lines), opts)

	res := cut.CutLines(lines, &opts)

	grpcResult := make([]*tasks2.StringArray, 0, len(res))
	for _, row := range res {
		grpcResult = append(grpcResult, &tasks2.StringArray{Items: row})
	}

	resp := &tasks2.ProcessTaskResponse{
		Result: grpcResult,
	}

	return resp, nil
}
