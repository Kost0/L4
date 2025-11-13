package client

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Kost0/L4/internal/config"
	"github.com/Kost0/L4/internal/cut"
	tasks2 "github.com/Kost0/L4/internal/gen/tasks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskResult struct {
	Result [][]string
	Err    error
}

func StartClient(inputFiles []string, opts *cut.CutOptions) [][]string {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	quorum := len(cfg.Workers)/2 + 1

	var wg sync.WaitGroup

	resultChan := make(chan TaskResult, len(inputFiles))

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Println("Starting client")

	for i, file := range inputFiles {
		lines, err := cut.ReadLines(file)
		if err != nil {
			log.Println(err)
			continue
		}

		wg.Add(1)

		workerPort := cfg.Workers[i%len(cfg.Workers)]
		go func() {
			defer wg.Done()

			res, err := callWorker(ctx, workerPort, lines, opts)

			resultChan <- TaskResult{res, err}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	log.Println("All workers finished")

	finalResult := make([][]string, 0)
	successCount := 0
	errorCount := 0

ResultLoop:
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				break ResultLoop
			}

			if result.Err != nil {
				log.Println(result.Err)
				errorCount++
			} else {
				successCount++
				finalResult = append(finalResult, result.Result...)
			}

			if successCount >= quorum {
				cancel()
				break ResultLoop
			}
		}
	}

	return finalResult
}

func callWorker(ctx context.Context, addr string, lines []string, opts *cut.CutOptions) ([][]string, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := tasks2.NewDistributedCutClient(conn)

	fields := make([]int32, 0, len(opts.Fields))
	for _, f := range opts.Fields {
		fields = append(fields, int32(f))
	}

	grpcOptions := &tasks2.CutOptions{
		Fields:    fields,
		Delimiter: opts.Delimiter,
		Separated: opts.Separated,
	}

	req := &tasks2.ProcessTaskRequest{
		Lines:   lines,
		Options: grpcOptions,
	}

	log.Println("Starting processing task")

	resp, err := client.ProcessTask(ctx, req)
	if err != nil {
		return nil, err
	}

	result := make([][]string, 0, len(resp.Result))
	for _, r := range resp.Result {
		result = append(result, r.Items)
	}

	return result, nil
}
