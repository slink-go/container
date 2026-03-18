package container_test

import (
	"fmt"
	"go.slink.ws/container"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

/*
# Примеры запуска бенчмарков:
# go test -bench=BenchmarkExtractSmallQueue -benchmem
# go test -bench=BenchmarkExtractLargeQueue -benchmem
# go test -bench=BenchmarkExtractWithDifferentFilters -benchmem
# go test -bench=BenchmarkExtractComparison -benchmem -count=5
# go test -bench=. -benchmem
*/

const (
	largeSize  = 1_000_000
	smallSize  = 10
	mediumSize = 10000
)

func benchmark(b *testing.B, name string, fn func(b *testing.B, counter *int, start *time.Time)) {
	b.Run(name, func(b *testing.B) {
		b.ResetTimer()
		start := time.Now()
		count := 0
		fn(b, &count, &start)
		duration := time.Since(start)
		b.ReportMetric(float64(count)/duration.Seconds(), "ops/sec")
		b.ReportMetric(float64(duration.Nanoseconds())/float64(count), "ns/op")
	})
}

// region - base benchmarks declaration

func benchmarkBaseEnqueue(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	// Тестирование Enqueue операций
	benchmark(b, fmt.Sprintf("%s:Enqueue", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
			*counter++
		}
	})
	//b.Run(fmt.Sprintf("%s:Enqueue", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		queue.Enqueue(i)
	//	}
	//})
}
func benchmarkBaseDequeue(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	// Тестирование Dequeue операций
	benchmark(b, fmt.Sprintf("%s:Dequeue", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
		}
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Dequeue()
			*counter++
		}
	})
	//b.Run(fmt.Sprintf("%s:Dequeue", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	// Сначала заполняем очередь
	//	for i := 0; i < b.N; i++ {
	//		queue.Enqueue(i)
	//	}
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		queue.Dequeue()
	//	}
	//})
}
func benchmarkBaseMixed(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	// Тестирование смешанных операций
	benchmark(b, fmt.Sprintf("%s:Mixed", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
			if i%2 == 0 {
				queue.Dequeue()
			}
			*counter++
		}
	})
	//b.Run(fmt.Sprintf("%s:Mixed", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		queue.Enqueue(i)
	//		if i%2 == 0 {
	//			queue.Dequeue()
	//		}
	//	}
	//})
}
func benchmarkQueueConcurrent(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:Concurrent", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		numGoroutines := runtime.NumCPU()
		numOperations := b.N / numGoroutines
		var cnt atomic.Int64
		var wg sync.WaitGroup
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					queue.Enqueue(j)
					queue.Dequeue()
					cnt.Add(1)
				}
			}()
		}
		wg.Wait()
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:Concurrent", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	numGoroutines := runtime.NumCPU()
	//	numOperations := b.N / numGoroutines
	//	b.ResetTimer()
	//	var wg sync.WaitGroup
	//	for i := 0; i < numGoroutines; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			for j := 0; j < numOperations; j++ {
	//				queue.Enqueue(j)
	//				queue.Dequeue()
	//			}
	//		}()
	//	}
	//	wg.Wait()
	//})
}
func benchmarkQueueHighConcurrency(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:HighConcurrency", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		numGoroutines := runtime.NumCPU()
		numOperations := b.N / numGoroutines
		var cnt atomic.Int64
		var wg sync.WaitGroup
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					queue.Enqueue(j)
					queue.Dequeue()
					cnt.Add(1)
				}
			}()
		}
		wg.Wait()
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	numGoroutines := runtime.NumCPU() * 4
	//	numOperations := b.N / numGoroutines
	//	b.ResetTimer()
	//	var wg sync.WaitGroup
	//	for i := 0; i < numGoroutines; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			for j := 0; j < numOperations; j++ {
	//				queue.Enqueue(j)
	//				queue.Dequeue()
	//			}
	//		}()
	//	}
	//	wg.Wait()
	//})
}
func benchmarkQueueThroughput(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:Throughput", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		numGoroutines := runtime.NumCPU()
		var cnt atomic.Int64
		var wg sync.WaitGroup
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < b.N/numGoroutines; j++ {
					queue.Enqueue(j)
					queue.Dequeue()
					cnt.Add(1)
				}
			}()
		}
		wg.Wait()
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:Throughput", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	b.ResetTimer()
	//	// Измеряем пропускную способность
	//	start := time.Now()
	//	var count int64
	//	var wg sync.WaitGroup
	//	numGoroutines := runtime.NumCPU()
	//	for i := 0; i < numGoroutines; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			for j := 0; j < b.N/numGoroutines; j++ {
	//				queue.Enqueue(j)
	//				queue.Dequeue()
	//				atomic.AddInt64(&count, 1)
	//			}
	//		}()
	//	}
	//	wg.Wait()
	//	duration := time.Since(start)
	//	b.ReportMetric(float64(count)/duration.Seconds(), "ops/sec")
	//})
}
func benchmarkQueueMemoryUsage(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:MemoryUsage", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		var cnt atomic.Int64
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
			cnt.Add(1)
			if i%1000 == 0 {
				queue.Dequeue()
				cnt.Add(1)
			}
		}
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:MemoryUsage", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		queue.Enqueue(i)
	//		if i%1000 == 0 {
	//			queue.Dequeue()
	//		}
	//	}
	//})
}
func benchmarkQueueStress(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:Stress", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		numGoroutines := runtime.NumCPU() * 5
		numOperations := b.N / numGoroutines
		var cnt atomic.Int64
		var wg sync.WaitGroup
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					queue.Enqueue(j)
					queue.Dequeue()
					cnt.Add(1)
				}
			}()
		}
		wg.Wait()
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:Stress", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	b.ResetTimer()
	//	// Создаем множество горутин для интенсивной нагрузки
	//	var wg sync.WaitGroup
	//	numGoroutines := runtime.NumCPU() * 5
	//	numOperations := b.N / numGoroutines
	//	for i := 0; i < numGoroutines; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			for j := 0; j < numOperations; j++ {
	//				queue.Enqueue(j)
	//				queue.Dequeue()
	//			}
	//		}()
	//	}
	//	wg.Wait()
	//})
}
func benchmarkQueueWithLargeData(b *testing.B, prefix string, queue container.Queue[[]byte]) {
	b.Run(fmt.Sprintf("%s:WithLargeData", prefix), func(b *testing.B) {
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
	})
}
func benchmarkQueueWithStrings(b *testing.B, prefix string, queue container.Queue[string]) {
	b.Run(fmt.Sprintf("%s:WithStrings", prefix), func(b *testing.B) {
		// Создаем строки для тестирования
		testString := "Hello, World! This is a test string for benchmarking."
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(testString)
			queue.Dequeue()
		}
	})
}
func benchmarkQueueEmptyQueue(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:EmptyQueue", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		var cnt atomic.Int64
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.IsEmpty()
			cnt.Add(1)
		}
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:EmptyQueue", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		queue.IsEmpty()
	//	}
	//})
}
func benchmarkQueueMultipleOperations(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:MultipleOperations", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		var cnt atomic.Int64
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
			queue.Enqueue(i + 1)
			queue.Enqueue(i + 2)
			queue.Dequeue()
			queue.Dequeue()
			queue.Dequeue()
			cnt.Add(1)
		}
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:MultipleOperations", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		queue.Enqueue(i)
	//		queue.Enqueue(i + 1)
	//		queue.Enqueue(i + 2)
	//		queue.Dequeue()
	//		queue.Dequeue()
	//		queue.Dequeue()
	//	}
	//})
}
func benchmarkQueueDetailed(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:SingleProducerSingleConsumer", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		var cnt atomic.Int64
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
			queue.Dequeue()
			cnt.Add(1)
		}
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:SingleProducerSingleConsumer", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	b.ResetTimer()
	//	for i := 0; i < b.N; i++ {
	//		queue.Enqueue(i)
	//		queue.Dequeue()
	//	}
	//})
	benchmark(b, fmt.Sprintf("%s:MultipleProducersSingleConsumer", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		numProducers := runtime.NumCPU()
		numOperations := b.N / numProducers
		var cnt atomic.Int64
		var wg sync.WaitGroup
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < numProducers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					queue.Enqueue(j)
					cnt.Add(1)
				}
			}()
		}
		wg.Wait()
		for i := 0; i < b.N; i++ {
			queue.Dequeue()
			cnt.Add(1)
		}
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:MultipleProducersSingleConsumer", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	numProducers := runtime.NumCPU()
	//	numOperations := b.N / numProducers
	//	b.ResetTimer()
	//	var wg sync.WaitGroup
	//	for i := 0; i < numProducers; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			for j := 0; j < numOperations; j++ {
	//				queue.Enqueue(j)
	//			}
	//		}()
	//	}
	//	wg.Wait()
	//	for i := 0; i < b.N; i++ {
	//		queue.Dequeue()
	//	}
	//})
	benchmark(b, fmt.Sprintf("%s:SingleProducerMultipleConsumers", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		numConsumers := runtime.NumCPU()
		numOperations := b.N / numConsumers
		var cnt atomic.Int64
		*start = time.Now()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
			cnt.Add(1)
		}
		var wg sync.WaitGroup
		for i := 0; i < numConsumers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					queue.Dequeue()
					cnt.Add(1)
				}
			}()
		}
		wg.Wait()
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:SingleProducerMultipleConsumers", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	numConsumers := runtime.NumCPU()
	//	numOperations := b.N / numConsumers
	//	for i := 0; i < b.N; i++ {
	//		queue.Enqueue(i)
	//	}
	//	b.ResetTimer()
	//	var wg sync.WaitGroup
	//	for i := 0; i < numConsumers; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			for j := 0; j < numOperations; j++ {
	//				queue.Dequeue()
	//			}
	//		}()
	//	}
	//	wg.Wait()
	//})
	benchmark(b, fmt.Sprintf("%s:HighConcurrency", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		numGoroutines := runtime.NumCPU() * 8
		numOperations := b.N / numGoroutines
		var cnt atomic.Int64
		*start = time.Now()
		b.ResetTimer()
		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					queue.Enqueue(j)
					queue.Dequeue()
					cnt.Add(1)
				}
			}()
		}
		wg.Wait()
		*counter = int(cnt.Load())
	})
	//b.Run(fmt.Sprintf("%s:HighConcurrency", prefix), func(b *testing.B) {
	//	queue := queueFactory()
	//	numGoroutines := runtime.NumCPU() * 8
	//	numOperations := b.N / numGoroutines
	//	b.ResetTimer()
	//	var wg sync.WaitGroup
	//	for i := 0; i < numGoroutines; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			for j := 0; j < numOperations; j++ {
	//				queue.Enqueue(j)
	//				queue.Dequeue()
	//			}
	//		}()
	//	}
	//	wg.Wait()
	//})
}
func benchmarkQueueVsChannel(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	// очередь
	benchmark(b, fmt.Sprintf("%s:VsChannel:Queue", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		b.ResetTimer()
		*start = time.Now()
		for i := 0; i < b.N; i++ {
			queue.Enqueue(i)
			queue.Dequeue()
			*counter++
		}
	})
	//b.Run(fmt.Sprintf("%s:VsChannel:Queue", prefix), func(b *testing.B) {
	//	b.ResetTimer()
	//	start := time.Now()
	//	count := 0
	//	for i := 0; i < b.N; i++ {
	//		queue.Enqueue(i)
	//		queue.Dequeue()
	//		count++
	//	}
	//	duration := time.Since(start)
	//	b.ReportMetric(float64(count)/duration.Seconds(), "ops/sec")
	//	b.ReportMetric(float64(duration.Nanoseconds())/float64(count), "ns/op")
	//})

	// канал
	benchmark(b, fmt.Sprintf("%s:VsChannel:Channel", prefix), func(b *testing.B, counter *int, start *time.Time) {
		ch := make(chan int, 1000)
		b.ResetTimer()
		*start = time.Now()
		for i := 0; i < b.N; i++ {
			ch <- i
			<-ch
			*counter++
		}
	})
	//b.Run(fmt.Sprintf("%s:VsChannel:Channel", prefix), func(b *testing.B) {
	//	ch := make(chan int, 1000)
	//	b.ResetTimer()
	//	start := time.Now()
	//	count := 0
	//	for i := 0; i < b.N; i++ {
	//		ch <- i
	//		<-ch
	//		count++
	//	}
	//	duration := time.Since(start)
	//	b.ReportMetric(float64(count)/duration.Seconds(), "ops/sec")
	//	b.ReportMetric(float64(duration.Nanoseconds())/float64(count), "ns/op")
	//})
}
func benchmarkQueuePerformance(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:Performance", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := queueFactory()
		var wg sync.WaitGroup
		numGoroutines := runtime.NumCPU()
		operationsPerGoroutine := b.N / numGoroutines
		b.ResetTimer()
		*start = time.Now()
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					queue.Enqueue(j)
					queue.Dequeue()
					*counter++
				}
			}()
		}
		wg.Wait()
	})
	//b.Run(fmt.Sprintf("%s:Performance", prefix), func(b *testing.B) {
	//
	//	b.ResetTimer()
	//
	//	// Измеряем производительность
	//	start := time.Now()
	//	count := 0
	//
	//	var wg sync.WaitGroup
	//	numGoroutines := runtime.NumCPU()
	//	operationsPerGoroutine := b.N / numGoroutines
	//
	//	for i := 0; i < numGoroutines; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			for j := 0; j < operationsPerGoroutine; j++ {
	//				queue.Enqueue(j)
	//				queue.Dequeue()
	//				count++
	//			}
	//		}()
	//	}
	//
	//	wg.Wait()
	//	duration := time.Since(start)
	//
	//	b.ReportMetric(float64(count)/duration.Seconds(), "ops/sec")
	//	b.ReportMetric(float64(duration.Nanoseconds())/float64(count), "ns/op")
	//})
}

// endregion
// region - base benchmarks invocation

func benchmarkComparedBaseEnqueue(b *testing.B) {
	benchmarkBaseEnqueue(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkBaseEnqueue(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedBaseDequeue(b *testing.B) {
	benchmarkBaseDequeue(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkBaseDequeue(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedBaseMixed(b *testing.B) {
	benchmarkBaseMixed(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkBaseMixed(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedQueueConcurrent(b *testing.B) {
	benchmarkQueueConcurrent(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueueConcurrent(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedQueueHighConcurrency(b *testing.B) {
	benchmarkQueueHighConcurrency(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueueHighConcurrency(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedQueueThroughput(b *testing.B) {
	benchmarkQueueThroughput(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueueThroughput(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedQueueMemoryUsage(b *testing.B) {
	benchmarkQueueMemoryUsage(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueueMemoryUsage(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedQueueStress(b *testing.B) {
	benchmarkQueueStress(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueueStress(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}

// func BenchmarkComparedQueueWithLargeData(b *testing.B) {}
// func BenchmarkComparedQueueWithStrings(b *testing.B) {}

func benchmarkComparedQueueEmptyQueue(b *testing.B) {
	benchmarkQueueEmptyQueue(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueueEmptyQueue(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedQueueMultipleOperations(b *testing.B) {
	benchmarkQueueMultipleOperations(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueueMultipleOperations(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedQueueDetailed(b *testing.B) {
	benchmarkQueueDetailed(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueueDetailed(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedQueueVsChannel(b *testing.B) {
	benchmarkQueueVsChannel(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueueVsChannel(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}
func benchmarkComparedQueuePerformance(b *testing.B) {
	benchmarkQueuePerformance(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkQueuePerformance(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}

func BenchmarkComparedBase(b *testing.B) {
	benchmarkComparedBaseEnqueue(b)
	benchmarkComparedBaseDequeue(b)
	benchmarkComparedBaseMixed(b)
	//benchmarkComparedQueueConcurrent(b)
	//benchmarkComparedQueueHighConcurrency(b)
	//benchmarkComparedQueueThroughput(b)
	//benchmarkComparedQueueMemoryUsage(b)
	//benchmarkComparedQueueStress(b)
	//benchmarkComparedQueueEmptyQueue(b)
	//benchmarkComparedQueueMultipleOperations(b)
	//benchmarkComparedQueueDetailed(b)
	//benchmarkComparedQueueVsChannel(b)
	//benchmarkComparedQueuePerformance(b)
}

// endregion

// region - extract benchmarks declaration

// Бенчмарки для метода Extract
func benchmarkExtractSimple(b *testing.B, prefix string, queueFactory func() container.Queue[int]) {
	benchmark(b, fmt.Sprintf("%s:Extract:Base", prefix), func(b *testing.B, counter *int, start *time.Time) {
		queue := container.NewLockFreeQueue()
		var cnt atomic.Int64
		// Заполняем очередь тестовыми данными
		for i := 0; i < largeSize; i++ {
			queue.Enqueue(i)
		}
		*start = time.Now()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				queue.ExtractSimple(func(item interface{}) bool {
					cnt.Add(1)
					if num, ok := item.(int); ok {
						return num%2 == 0
					}
					return false
				})
			}
		})
		*counter = int(cnt.Load())
	})
	//queue := NewLockFreeQueue()
	//// Заполняем очередь тестовыми данными
	//for i := 0; i < 1000; i++ {
	//	queue.Enqueue(i)
	//}
	//b.ResetTimer()
	//b.RunParallel(func(pb *testing.PB) {
	//	for pb.Next() {
	//		queue.ExtractSimple(func(item interface{}) bool {
	//			if num, ok := item.(int); ok {
	//				return num%2 == 0
	//			}
	//			return false
	//		})
	//	}
	//})
}
func benchmarkExtractEfficient(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь тестовыми данными
	for i := 0; i < largeSize; i++ {
		queue.Enqueue(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.ExtractEfficient(func(item interface{}) bool {
				if num, ok := item.(int); ok {
					return num%2 == 0
				}
				return false
			})
		}
	})
}
func benchmarkExtractWithState(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь тестовыми данными
	for i := 0; i < largeSize; i++ {
		queue.Enqueue(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.ExtractWithState(func(item interface{}) bool {
				if num, ok := item.(int); ok {
					return num%2 == 0
				}
				return false
			})
		}
	})
}
func benchmarkExtractOptimized(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь тестовыми данными
	for i := 0; i < largeSize; i++ {
		queue.Enqueue(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.ExtractOptimized(func(item interface{}) bool {
				if num, ok := item.(int); ok {
					return num%2 == 0
				}
				return false
			})
		}
	})
}
func benchmarkExtractSmallQueue(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь небольшим количеством элементов
	for i := 0; i < smallSize; i++ {
		queue.Enqueue(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.ExtractSimple(func(item interface{}) bool {
				if num, ok := item.(int); ok {
					return num%2 == 0
				}
				return false
			})
		}
	})
}
func benchmarkExtractLargeQueue(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь большим количеством элементов
	for i := 0; i < mediumSize; i++ {
		queue.Enqueue(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.ExtractSimple(func(item interface{}) bool {
				if num, ok := item.(int); ok {
					return num%2 == 0
				}
				return false
			})
		}
	})
}
func benchmarkExtractWithDifferentFilters(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь тестовыми данными
	for i := 0; i < mediumSize; i++ {
		queue.Enqueue(i)
	}

	// Бенчмарк для фильтра "четные числа"
	b.Run("EvenNumbers", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				queue.ExtractSimple(func(item interface{}) bool {
					if num, ok := item.(int); ok {
						return num%2 == 0
					}
					return false
				})
			}
		})
	})

	// Бенчмарк для фильтра "нечетные числа"
	b.Run("OddNumbers", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				queue.ExtractSimple(func(item interface{}) bool {
					if num, ok := item.(int); ok {
						return num%2 == 1
					}
					return false
				})
			}
		})
	})

	// Бенчмарк для фильтра "числа больше 500"
	b.Run("GreaterThan500", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				queue.ExtractSimple(func(item interface{}) bool {
					if num, ok := item.(int); ok {
						return num > 500
					}
					return false
				})
			}
		})
	})

	// Бенчмарк для фильтра "все числа"
	b.Run("AllNumbers", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				queue.ExtractSimple(func(item interface{}) bool {
					return true
				})
			}
		})
	})

	// Бенчмарк для фильтра "ни одного числа"
	b.Run("NoNumbers", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				queue.ExtractSimple(func(item interface{}) bool {
					return false
				})
			}
		})
	})
}
func benchmarkExtractConcurrent(b *testing.B) {
	queue := container.NewLockFreeQueue()
	// Заполняем очередь тестовыми данными
	for i := 0; i < mediumSize; i++ {
		queue.Enqueue(i)
	}
	// Создаем несколько горутин для конкурентного тестирования
	b.N = 4 //numGoroutines := 4
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Каждая горутина выполняет Extract с разными фильтрами
			queue.ExtractSimple(func(item interface{}) bool {
				if num, ok := item.(int); ok {
					return num%3 == 0
				}
				return false
			})
		}
	})
}
func benchmarkExtractWithStringData(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь строковыми данными
	for i := 0; i < mediumSize; i++ {
		queue.Enqueue("item_" + string(rune(i+'0')))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.ExtractSimple(func(item interface{}) bool {
				if str, ok := item.(string); ok {
					return len(str) > 5
				}
				return false
			})
		}
	})
}
func benchmarkExtractWithStructData(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь структурными данными
	type TestStruct struct {
		ID   int
		Name string
		Age  int
	}

	for i := 0; i < mediumSize; i++ {
		queue.Enqueue(TestStruct{
			ID:   i,
			Name: "test_" + string(rune(i+'0')),
			Age:  i % 100,
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.ExtractSimple(func(item interface{}) bool {
				if ts, ok := item.(TestStruct); ok {
					return ts.Age > 50
				}
				return false
			})
		}
	})
}

// Бенчмарк для сравнения производительности разных реализаций
func benchmarkExtractComparison(b *testing.B) {
	// Подготовка данных
	data := make([]int, mediumSize)
	for i := 0; i < mediumSize; i++ {
		data[i] = i
	}

	b.Run("SimpleImplementation", func(b *testing.B) {
		queue := container.NewLockFreeQueue()
		for _, item := range data {
			queue.Enqueue(item)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.ExtractSimple(func(item interface{}) bool {
				if num, ok := item.(int); ok {
					return num%2 == 0
				}
				return false
			})
		}
	})

	b.Run("EfficientImplementation", func(b *testing.B) {
		queue := container.NewLockFreeQueue()
		for _, item := range data {
			queue.Enqueue(item)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.ExtractEfficient(func(item interface{}) bool {
				if num, ok := item.(int); ok {
					return num%2 == 0
				}
				return false
			})
		}
	})

	b.Run("WithStateImplementation", func(b *testing.B) {
		queue := container.NewLockFreeQueue()
		for _, item := range data {
			queue.Enqueue(item)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			queue.ExtractWithState(func(item interface{}) bool {
				if num, ok := item.(int); ok {
					return num%2 == 0
				}
				return false
			})
		}
	})
}

// Бенчмарк для измерения времени выполнения
func benchmarkExtractPerformance(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь
	for i := 0; i < mediumSize; i++ {
		queue.Enqueue(i)
	}

	b.ResetTimer()

	// Измеряем производительность
	for i := 0; i < b.N; i++ {
		start := time.Now()
		result := queue.ExtractSimple(func(item interface{}) bool {
			if num, ok := item.(int); ok {
				return num%2 == 0
			}
			return false
		})
		duration := time.Since(start)

		// Выводим информацию о времени выполнения
		b.Logf("Extract operation took %v, returned %d items", duration, len(result))
	}
}

// Дополнительный бенчмарк для тестирования памяти
func benchmarkExtractMemoryUsage(b *testing.B) {
	queue := container.NewLockFreeQueue()

	// Заполняем очередь большими данными
	for i := 0; i < largeSize; i++ {
		queue.Enqueue(strings.Repeat("test_string_", 10))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.ExtractSimple(func(item interface{}) bool {
				if str, ok := item.(string); ok {
					return len(str) > 50
				}
				return false
			})
		}
	})
}

// endregion
// region - extract benchmarks invocation

func benchmarkComparedExtract(b *testing.B) {
	benchmarkExtractSimple(b, "Lfree", func() container.Queue[int] { return container.CreateLockFreeGenericQueue[int]() })
	benchmarkExtractSimple(b, "Lcked", func() container.Queue[int] { return container.CreateBlockingQueue[int]() })
}

func BenchmarkComparedExtract(b *testing.B) {
	benchmarkComparedExtract(b)
}

// endregion
