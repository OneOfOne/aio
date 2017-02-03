package aio

import (
	"os"
	"runtime"
)

func newThread(rq chan interface{}) *thread {
	return &thread{rq}
}

type thread struct {
	rq chan interface{}
}

func (t *thread) listen() {
	runtime.LockOSThread()

	for req := range t.rq {
		switch r := req.(type) {
		case *openRequest:
			t.open(r)
		case *readRequest:
			t.read(r)
		case *writeRequest:
			t.write(r)
		case *deleteRequest:
			t.delete(r)
		case *closeRequest:
			t.close(r)

		default:
			panic("invalid type")
		}
	}

	runtime.UnlockOSThread()
}

func (t *thread) open(r *openRequest) {
	var resp OpenResp
	resp.f, resp.err = newFile(r, t.rq)
	r.resp <- &resp
}

func (t *thread) read(r *readRequest) {
	var resp RWResp
	resp.n, resp.err = r.f.Read(r.b)
	r.resp <- &resp
}

func (t *thread) write(r *writeRequest) {
	var resp RWResp
	resp.n, resp.err = r.f.Write(r.b)
	r.resp <- &resp
}

func (t *thread) close(r *closeRequest) {
	r.errCh <- r.f.Close()
}

func (t *thread) delete(r *deleteRequest) {
	r.errCh <- os.Remove(r.key)
}