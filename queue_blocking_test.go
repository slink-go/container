package container_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.slink.ws/container"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func prepareQueue() *container.BlockingQueue[int] {

	q := container.CreateBlockingQueue[int]()

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Enqueue(4)
	q.Enqueue(5)
	q.Enqueue(6)
	q.Enqueue(7)
	q.Enqueue(8)
	q.Enqueue(9)
	q.Enqueue(10)

	return q

}

func TestLlQueueA(t *testing.T) {

	q := prepareQueue()
	res := q.Extract(func(i int) bool {
		return i%2 == 0
	})
	fmt.Println(res)

	for !q.IsEmpty() {
		v, _ := q.Dequeue()
		fmt.Println("-->", v)
	}

	assert.True(t, q.IsEmpty())

}
func TestLlQueueB(t *testing.T) {

	q := prepareQueue()
	res := q.Extract(func(i int) bool {
		return i%2 == 1
	})
	fmt.Println(res)

	for !q.IsEmpty() {
		v, _ := q.Dequeue()
		fmt.Println("-->", v)
	}

	assert.True(t, q.IsEmpty())

}
func TestLlQueueC(t *testing.T) {

	q := prepareQueue()
	res := q.Extract(func(i int) bool {
		return i%3 == 0
	})
	fmt.Println(res)

	for !q.IsEmpty() {
		v, _ := q.Dequeue()
		fmt.Println("-->", v)
	}

	assert.True(t, q.IsEmpty())

}

// BenchmarkLockFreeQueue - бенчмарки для lock-free очереди
func BenchmarkLLQueue(b *testing.B) {
	// Создаем очередь
	queue := container.CreateBlockingQueue[int]()
	// Тестирование Enqueue операций
	b.Run("Enqueue", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
		}
	})
	// Тестирование Dequeue операций
	b.Run("Dequeue", func(b *testing.B) {
		// Сначала заполняем очередь
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Dequeue()
		}
	})
	// Тестирование смешанных операций
	b.Run("Mixed", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
			if i%2 == 0 {
				queue.Dequeue()
			}
		}
	})
}

// BenchmarkLockFreeQueueConcurrent - конкурентные бенчмарки
func BenchmarkLLQueueConcurrent(b *testing.B) {
	queue := container.CreateBlockingQueue[int]()
	numGoroutines := runtime.NumCPU()
	numOperations := b.N / numGoroutines

	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				queue.Enqueue(j)
				queue.Dequeue()
			}
		}()
	}

	wg.Wait()
}

// BenchmarkLockFreeQueueHighConcurrency - высокая конкуренция
func BenchmarkLLQueueHighConcurrency(b *testing.B) {
	queue := container.CreateBlockingQueue[int]()
	numGoroutines := runtime.NumCPU() * 4
	numOperations := b.N / numGoroutines

	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				queue.Enqueue(j)
				queue.Dequeue()
			}
		}()
	}

	wg.Wait()
}

// BenchmarkLockFreeQueueThroughput - тестирование пропускной способности
func BenchmarkLLQueueThroughput(b *testing.B) {
	queue := container.CreateBlockingQueue[int]()

	b.ResetTimer()

	// Измеряем пропускную способность
	start := time.Now()
	var count int64

	var wg sync.WaitGroup
	numGoroutines := runtime.NumCPU()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < b.N/numGoroutines; j++ {
				queue.Enqueue(j)
				queue.Dequeue()
				atomic.AddInt64(&count, 1)
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	b.ReportMetric(float64(count)/duration.Seconds(), "ops/sec")
}

// BenchmarkLockFreeQueueMemoryUsage - тестирование использования памяти
func BenchmarkLLQueueMemoryUsage(b *testing.B) {
	queue := container.CreateBlockingQueue[int]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
		if i%1000 == 0 {
			queue.Dequeue()
		}
	}
}

// BenchmarkStandardQueue - сравнение с стандартной sync.Queue (для сравнения)
//func BenchmarkStandardQueue(b *testing.B) {
//    var queue list.List
//    var wg sync.WaitGroup
//
//    numGoroutines := runtime.NumCPU()
//    numOperations := b.N / numGoroutines
//
//    b.ResetTimer()
//
//    for i := 0; i < numGoroutines; i++ {
//        wg.Add(1)
//        go func() {
//            defer wg.Done()
//            for j := 0; j < numOperations; j++ {
//                queue.PushBack(j)
//                queue.Remove().Front() .PopDequeue()
//            }
//        }()
//    }
//    wg.Wait()
//}

// BenchmarkLockFreeQueueStress - стресс-тест
func BenchmarkLLQueueStress(b *testing.B) {
	queue := container.CreateBlockingQueue[int]()
	b.ResetTimer()
	// Создаем множество горутин для интенсивной нагрузки
	var wg sync.WaitGroup
	numGoroutines := runtime.NumCPU() * 2
	numOperations := b.N / numGoroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				queue.Enqueue(j)
				queue.Dequeue()
			}
		}()
	}
	wg.Wait()
}

// BenchmarkLockFreeQueueWithLargeData - тест с большими данными
func BenchmarkLLQueueWithLargeData(b *testing.B) {
	queue := container.CreateBlockingQueue[[]byte]()
	// Создаем большие данные для тестирования
	largeData := make([]byte, 1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.Enqueue(largeData)
		queue.Dequeue()
	}
}

// BenchmarkLockFreeQueueWithStrings - тест с строками
func BenchmarkLLQueueWithStrings(b *testing.B) {
	queue := container.CreateBlockingQueue[string]()
	// Создаем строки для тестирования
	testString := "Hello, World! This is a test string for benchmarking."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.Enqueue(testString)
		queue.Dequeue()
	}
}

// BenchmarkLockFreeQueueEmptyQueue - тест с пустой очередью
func BenchmarkLLQueueEmptyQueue(b *testing.B) {
	queue := container.CreateBlockingQueue[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.IsEmpty()
	}
}

// BenchmarkLockFreeQueueMultipleOperations - тест с множественными операциями
func BenchmarkLLQueueMultipleOperations(b *testing.B) {
	queue := container.CreateBlockingQueue[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
		queue.Enqueue(i + 1)
		queue.Enqueue(i + 2)
		queue.Dequeue()
		queue.Dequeue()
		queue.Dequeue()
	}
}
