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

type request struct {
	id       int
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
	memo.requests <- request{id, response}
	res := <-response

	return res.value, res.err
}

func (memo *Memo) close() {
	close(memo.requests)
}

func (memo *Memo) server(f Func) {
	cache := make(map[int]*entry)
	for req := range memo.requests {
		e := cache[req.id]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.id] = e
			go e.call(f, req.id)
		}
		go e.deliver(req.response)
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
