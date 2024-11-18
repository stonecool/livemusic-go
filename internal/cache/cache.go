package cache

// Refer to the example on P278 of book <The Go Programming Language>
type Func func(id int) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{}
}

type requestType int

const (
	getRequest requestType = iota
	setRequest
)

type request struct {
	id       int
	op       requestType
	value    interface{}
	response chan<- result
}

type Memo struct {
	requests chan request
}

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(id int) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{id, getRequest, nil, response}
	res := <-response

	return res.value, res.err
}

func (memo *Memo) Set(id int, value interface{}) error {
	response := make(chan result)
	memo.requests <- request{
		id:       id,
		op:       setRequest,
		value:    value,
		response: response,
	}
	res := <-response
	return res.err
}

func (memo *Memo) close() {
	close(memo.requests)
}

func (memo *Memo) server(f Func) {
	cache := make(map[int]*entry)
	for req := range memo.requests {
		switch req.op {
		case getRequest:
			e := cache[req.id]
			if e == nil {
				e = &entry{ready: make(chan struct{})}
				cache[req.id] = e
				go e.call(f, req.id)
			}
			go e.deliver(req.response)
		case setRequest:
			e := &entry{ready: make(chan struct{})}
			e.res.value = req.value
			close(e.ready)
			cache[req.id] = e
			req.response <- e.res
		}
	}
}

func (e *entry) call(f Func, id int) {
	e.res.value, e.res.err = f(id)
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready
	response <- e.res
}
