package chdist_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/p9595jh/chdist"
)

func TestDistribution(t *testing.T) {
	d := chdist.NewDistributor(make(chan int))

	s1 := d.Subscribe(func(value int) {
		t.Log(1, value)
	})
	d.Subscribe(func(value int) {
		t.Log(2, value)
	})

	_ = s1

	go func() {
		for i := 0; i < 3; i++ {
			d.In() <- i
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 3)
	d.Close()
	time.Sleep(time.Second * 5)
}

func TestAsync(t *testing.T) {
	d := chdist.NewDistributor(make(chan int))

	s1 := d.AsyncSubscribe()
	go func() {
		for a := range s1.Out() {
			t.Log(1, a)
		}
	}()

	s2 := d.AsyncSubscribe()
	go func() {
		for a := range s2.Out() {
			t.Log(2, a)
		}
	}()

	d.Subscribe(func(value int) {
		t.Log(3, value)
	})

	go func() {
		for i := 0; i < 3; i++ {
			d.In() <- i
			time.Sleep(time.Second)
		}
	}()

	// d.Close()
	time.Sleep(time.Second * 5)
}

type Model struct {
	Num int    `json:"num"`
	Str string `json:"str"`
}

func TestJsonString(t *testing.T) {
	jsonify := func(i int, s string) string {
		return fmt.Sprintf(`{"num":%d,"str":"%s"}`, i, s)
	}

	ch := make(chan string)
	d := chdist.NewJsonStringDistributor[Model](chdist.NewDistributor(ch))

	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Second * 2)
			ch <- jsonify(i, "hello")
		}
	}()

	d.Subscribe(func(i chdist.Item[Model]) {
		t.Log(i.Value)
	})
	time.Sleep(time.Second * 7)
}

func TestJsonStringAsync(t *testing.T) {
	jsonify := func(i int, s string) string {
		return fmt.Sprintf(`{"num":%d,"str":"%s"}`, i, s)
	}

	ch := make(chan string)
	d := chdist.NewJsonStringDistributor[Model](chdist.NewDistributor(ch))

	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Second * 2)
			ch <- jsonify(i, "hello")
		}
	}()

	go func() {
		for item := range d.AsyncSubscribe() {
			t.Log(item)
		}
	}()
	time.Sleep(time.Second * 7)
}

func TestJsonBytes(t *testing.T) {
	jsonify := func(i int, s string) []byte {
		return []byte(fmt.Sprintf(`{"num":%d,"str":"%s"}`, i, s))
	}

	ch := make(chan []byte)
	d := chdist.NewJsonBytesDistributor[Model](chdist.NewDistributor(ch))

	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Second * 2)
			ch <- jsonify(i, "hello")
		}
	}()

	d.Subscribe(func(item chdist.Item[Model]) {
		t.Log(item)
	})
	time.Sleep(time.Second * 7)
}

func TestClose(t *testing.T) {
	d := chdist.NewDistributor(make(chan int))
	sub1 := d.Subscribe(func(value int) {
		t.Log(1, value)
	})
	sub2 := d.Subscribe(func(value int) {
		t.Log(2, value)
	})
	d.Subscribe(func(value int) {
		t.Log(3, value)
	})
	// _ = sub2

	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Second * 2)
			d.In() <- i
		}
	}()

	for _, sub := range d.Subscriptions() {
		// t.Log(sub, sub.IsAlive())
		t.Log(sub)
	}
	time.Sleep(time.Second * 6)

	sub1.Close()
	t.Log(d.Subscriptions())

	time.Sleep(time.Second * 6)
	_ = sub2
}
