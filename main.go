package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strings"
	"time"
)

type Result struct {
	code    int
	elapsed time.Duration
}

type Queue struct {
	queueChan chan struct{}
}

func NewQueue(concurrency int) *Queue {
	return &Queue{
		queueChan: make(chan struct{}, concurrency),
	}
}

func (q *Queue) lock() {
	q.queueChan <- struct{}{}
}

func (q *Queue) unlock() {
	<-q.queueChan
}

func main() {
	var (
		url         string
		requests    int
		concurrency int
	)

	// Definindo as flags da linha de comando
	flag.StringVar(&url, "url", "", "URL do serviço a ser testado.")
	flag.IntVar(&requests, "requests", 10, "Número total de requests.")
	flag.IntVar(&concurrency, "concurrency", 1, "Número de chamadas simultâneas.")

	flag.Parse()

	if url == "" || requests <= 0 || concurrency <= 0 {
		flag.PrintDefaults()
		return
	}

	// Iniciando o teste de carga
	fmt.Printf("Iniciando testes de carga em %s com %d requests e %d chamadas simultâneas...\n", url, requests, concurrency)
	start := time.Now()

	queue := NewQueue(concurrency)
	reqChan := make(chan Result)
	doneChan := make(chan struct{}, requests)

	for i := 0; i < requests; i++ {
		go func() {
			queue.lock()
			defer queue.unlock()
			var req *http.Request
			var err error
			for {
				req, err = http.NewRequest("GET", url, nil)
				if err != nil {
					continue
				}
				break
			}
			for {
				start := time.Now()
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					continue
				}
				reqChan <- Result{
					code:    resp.StatusCode,
					elapsed: time.Since(start),
				}
				doneChan <- struct{}{}
				resp.Body.Close()
				break
			}
		}()
	}
	statusCodeCounts := make(map[int]int)

	counter := 1.0
	timeSlice := make([]float64, 0)

	fmt.Print("\033[s")

LOOP:
	for {
		select {
		case result := <-reqChan:
			percentage := 100 * counter / float64(requests)
			printBar := strings.Repeat("|", int(percentage))
			if percentage < 100 {
				printBar = printBar + strings.Repeat(" ", 100-int(percentage)) + "|"
			}

			fmt.Printf("\033[uProgresso: %d/%d %.2f%% %s\n", int(counter), requests, percentage, printBar)
			// fmt.Println("len(reqDone)", len(reqDone))
			timeSlice = append(timeSlice, float64(result.elapsed.Milliseconds()))
			counter++
			statusCodeCounts[result.code]++
		default:
			if len(doneChan) == requests {
				break LOOP
			}
		}
	}

	// Geração do relatório
	fmt.Println("Relatório de Teste:")
	fmt.Println("Tempo total:", time.Since(start))
	fmt.Printf("Tempo médio: %.2fms\n", Avarage(timeSlice))
	fmt.Printf("Desvio padão: %.2fms\n", StandardDeviation(timeSlice))
	fmt.Printf("Percentil 0.75: %.2fms\n", Percentile(timeSlice, 0.75))
	fmt.Printf("Percentil 0.90: %.2fms\n", Percentile(timeSlice, 0.90))
	fmt.Printf("Percentil 0.95: %.2fms\n", Percentile(timeSlice, 0.95))
	fmt.Printf("Percentil 0.99: %.2fms\n", Percentile(timeSlice, 0.99))
	fmt.Println("Quantidade total de requests:", requests)
	fmt.Println("Quantidade de requests com status HTTP 200:", statusCodeCounts[http.StatusOK])
	for code, count := range statusCodeCounts {
		if code != http.StatusOK {
			fmt.Printf("Quantidade de requests com status HTTP %d: %d\n", code, count)
		}
	}
}

func Avarage(slice []float64) float64 {
	sum := 0.0
	for _, num := range slice {
		sum += num
	}
	return sum / float64(len(slice))
}

func StandardDeviation(slice []float64) float64 {
	avarage := Avarage(slice)
	sd := 0.0
	for _, num := range slice {
		sd += (num - avarage) * (num - avarage)
	}
	return math.Sqrt(sd / float64(len(slice)))
}

func Percentile(slice []float64, percent float64) (percentile float64) {
	sorted := slices.Clone(slice)
	slices.Sort(sorted)
	index := int(math.Ceil(float64(len(sorted)) * percent))
	return sorted[index]
}
