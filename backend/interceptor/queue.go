package interceptor

import (
	"net/http"
	"sync"
)

type queue struct {
	sync.Mutex // protects the queue items
	items      []queuedRequest
}

func (q *queue) Add(req *http.Request, id int64) chan *InterceptionResponse {
	q.Lock()
	defer q.Unlock()
	resultChan := make(chan *InterceptionResponse, 1)
	q.items = append(q.items, queuedRequest{
		id:  id,
		req: req,
		c:   resultChan,
	})
	return resultChan
}

func (q *queue) Consume(req *http.Request, id int64, resp *http.Response) {
	q.Lock()
	defer q.Unlock()
	for index := range q.items {
		if q.items[index].id == id {
			req.URL.Host = q.items[index].req.URL.Host
			req.URL.Scheme = q.items[index].req.URL.Scheme
			resultChan := q.items[index].c
			q.items = append(q.items[:index], q.items[index+1:]...)
			resultChan <- &InterceptionResponse{
				Req:  req,
				Resp: resp,
			}
			return
		}
	}
}

func (q *queue) Flush() {
	q.Lock()
	defer q.Unlock()
	for index := range q.items {
		resultChan := q.items[index].c
		resultChan <- &InterceptionResponse{
			Req:  q.items[index].req,
			Resp: nil,
		}
	}
	q.items = []queuedRequest{}
}

func (q *queue) Len() int {
	q.Lock()
	defer q.Unlock()
	return len(q.items)
}
