package searchutils

func LinearSearchParallelInt(data []int, target int, thread int) int {
	length := len(data)
	dataXThread := length / thread
	oddment := length % thread

	result := make(chan int)
	if oddment != 0 {
		oddment += thread * dataXThread
		go LinearSearchInt(data[thread*dataXThread:oddment], target, result)
		found := <-result
		if found != -1 {
			return found + thread*dataXThread
		}
	}
	for i := 0; i < thread; i++ {
		go LinearSearchInt(data[i*dataXThread:(i+1)*dataXThread], target, result)
		index := <-result
		if index != -1 {
			return index + i*dataXThread
		}
	}
	return -1
}

func LinearSearchInt(data []int, target int, result chan int) {
	found := false
	for i := range data {
		if target == data[i] {
			result <- i
			found = true
			return
		}
	}
	if !found {
		result <- -1
	}
}
