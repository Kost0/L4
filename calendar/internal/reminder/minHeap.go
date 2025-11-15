package reminder

import (
	"container/heap"

	"github.com/Kost0/L4/internal/models"
)

type PriorityQueue []*models.Event

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Date.Before(pq[j].Date)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*models.Event)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func CreateHeap() *PriorityQueue {
	pq := make(PriorityQueue, 0)

	heap.Init(&pq)

	return &pq
}
