package searchutils

func LinearSearchParallelInt(data []int, target int, thread int) int {
	length := len(data)
	dataXThread := length / thread
	oddment := length % thread
	if oddment != 0 {
		oddment += thread * dataXThread
		found := LinearSearchInt(data[thread*dataXThread:oddment], target)
		if found != -1 {
			return found + thread*dataXThread
		}
	}
	for i := 0; i < thread; i++ {
		index := LinearSearchInt(data[i*dataXThread:(i+1)*dataXThread], target)
		if index != -1 {
			return index + i*dataXThread
		}
	}

	return -1
}

func LinearSearchInt(data []int, target int) int {
	for i := range data {
		if target == data[i] {
			return i
		}
	}
	return -1
}
