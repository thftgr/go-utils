package bean

import (
	"testing"
)

type IA interface {
	Print()
}

type A struct {
	t   *testing.T
	val string
}

func (r *A) Print() {
	r.t.Log(r.val)
}

type B struct {
	t   *testing.T
	val string
}

func (r *B) Print() {
	r.t.Log(r.val)
}

func TestRegisterAndGetBean(t *testing.T) {
	Register[IA](&A{t: t, val: "A"})
	Get[IA]().Print()

	Register[IA](&B{t: t, val: "B"})
	Get[IA]().Print()

	Register[*A](&A{t: t, val: "*A"})
	Get[*A]().Print()
	//Get[*B]().Print()

}
