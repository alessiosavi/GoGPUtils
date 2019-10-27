package searchutils

import (
	"sync"
)

// LinearSearchParallelInt is delegated to parallelize the execution of search method
func LinearSearchParallelInt(data []int, target int, thread int) int {
	length := len(data)
	dataXThread := length / thread
	oddment := length % thread

	if oddment != 0 {
		oddment += thread * dataXThread
		found := LinearSearchInt(data[thread*dataXThread:oddment], target)
		if found != -1 {
			return found + thread*dataXThread
		} // else
		return -1
	}
	wg := sync.WaitGroup{}
	result := make([]int, 0)
	//var result []int
	wg.Add(thread)
	for i := 0; i < thread; i++ {
		go LinearSearchParallelIntHelper(&wg, data[i*dataXThread:(i+1)*dataXThread], target, &result)
	}
	wg.Wait()
	// log.Println(result)
	for i := range result {
		if i != -1 {
			return result[i] + i*dataXThread
		}
	}
	return -1
}

// LinearSearchParallelIntHelper is delegated to search the number and append to the given result array
func LinearSearchParallelIntHelper(wg *sync.WaitGroup, data []int, target int, result *[]int) {
	defer wg.Done()
	*result = append(*result, LinearSearchInt(data, target))
}

// LinearSearchInt is a simple for delegated to find the target value
func LinearSearchInt(data []int, target int) int {
	for i := range data {
		if target == data[i] {
			return i
		}
	}
	return -1
}
