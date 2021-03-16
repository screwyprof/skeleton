package renderer

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/screwyprof/golibs/adaptor"
)

type Request interface {
	Bind(*gin.Context) error
}

type Response interface {
	Status() int
}

type GinRenderer struct {
	requests map[reflect.Value]Request
}

func NewGinRenderer() *GinRenderer {
	return &GinRenderer{requests: make(map[reflect.Value]Request)}
}

func (r *GinRenderer) MustAdapt(h adaptor.HandlerFn) gin.HandlerFunc {
	fn, err := r.Adapt(h)
	if err != nil {
		panic(err)
	}
	return fn
}

func (r *GinRenderer) Adapt(h adaptor.HandlerFn) (gin.HandlerFunc, error) {
	rq, err := r.requestForHandler(h)
	if err != nil {
		return nil, err
	}

	fn := func(c *gin.Context) {
		if err := rq.Bind(c); err != nil {
			r.render(nil, err)(c)
			return
		}

		res, err := h(c, rq)
		r.render(res, err)(c)
	}
	return fn, nil
}

func (r *GinRenderer) Register(rq Request, h adaptor.HandlerFn) {
	r.requests[reflect.ValueOf(h)] = rq
}

func (r *GinRenderer) requestForHandler(h adaptor.HandlerFn) (Request, error) {
	v := reflect.ValueOf(h)
	if r, ok := r.requests[v]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("request for %s not found", runtime.FuncForPC(v.Pointer()).Name())
}

func (r *GinRenderer) render(res interface{}, err error) gin.HandlerFunc {
	if err != nil {
		return (&ErrorRenderer{Err: err}).Render
	}
	status := http.StatusOK
	if r, ok := res.(Response); ok {
		status = r.Status()
	}
	return (&ContentRenderer{Content: res, Status: status}).Render
}
