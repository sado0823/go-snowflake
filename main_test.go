package main

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	// new default
	sf := New()
	syncMap := new(sync.Map)
	wg := new(sync.WaitGroup)
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := sf.NextVal()
			if err == nil {
				_, ok := syncMap.Load(val)
				if ok {
					t.Errorf("gen the same id")
				}
				syncMap.Store(val, val)
			} else {
				t.Errorf("fail to gen id, e:%s", err.Error())
			}
		}()
	}
	wg.Wait()
}

func TestWithEpoch(t *testing.T) {
	// new default
	sf := New(WithEpoch(1577808888888))
	syncMap := new(sync.Map)
	wg := new(sync.WaitGroup)
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := sf.NextVal()
			if err == nil {
				_, ok := syncMap.Load(val)
				if ok {
					t.Errorf("gen the same id")
				}
				syncMap.Store(val, val)
			} else {
				t.Errorf("fail to gen id, e:%s", err.Error())
			}
		}()
	}
	wg.Wait()
}

func BenchmarkNewParallel(b *testing.B) {

	// b.SetParallelism(1) // 设置使用的CPU数
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sf := New()
			syncMap := new(sync.Map)
			wg := new(sync.WaitGroup)
			for i := 0; i < 1000000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					val, err := sf.NextVal()
					if err == nil {
						_, ok := syncMap.Load(val)
						if ok {
							b.Errorf("gen the same id")
						}
						syncMap.Store(val, val)
					} else {
						b.Errorf("fail to gen id, e:%s", err.Error())
					}
				}()
			}
			wg.Wait()
		}
	})
}