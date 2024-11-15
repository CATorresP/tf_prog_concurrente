package safecounts

import (
	"sync"
)

type SafeCounts struct {
	Counts   []int
	CountsMu sync.RWMutex
	Status   []bool
	StatusMu sync.RWMutex
}

func (sd *SafeCounts) ReadCounts() []int {
	sd.CountsMu.RLock()
	defer sd.CountsMu.RUnlock()
	copiedCounts := make([]int, len(sd.Counts))
	copy(copiedCounts, sd.Counts)
	return copiedCounts
}

func (sd *SafeCounts) CompareCounts(counts []int) bool {
	sd.CountsMu.RLock()
	defer sd.CountsMu.RUnlock()
	for i, count := range counts {
		if count != sd.Counts[i] {
			return false
		}
	}
	return true
}

func (sd *SafeCounts) ReadCountByIndex(index int) int {
	sd.CountsMu.RLock()
	defer sd.CountsMu.RUnlock()
	return sd.Counts[index]
}
func (sd *SafeCounts) WriteCountByIndex(value int, index int) {
	sd.CountsMu.Lock()
	defer sd.CountsMu.Unlock()
	sd.Counts[index] = value
}
func (sd *SafeCounts) ReadStatus() []bool {
	sd.StatusMu.RLock()
	defer sd.StatusMu.RUnlock()
	copiedStatus := make([]bool, len(sd.Status))
	copy(copiedStatus, sd.Status)
	return copiedStatus
}

func (sd *SafeCounts) ReadStatustByIndex(index int) bool {
	sd.StatusMu.RLock()
	defer sd.StatusMu.RUnlock()
	return sd.Status[index]
}
func (sd *SafeCounts) WriteStatusByIndex(value bool, index int) {
	sd.StatusMu.Lock()
	defer sd.StatusMu.Unlock()
	sd.Status[index] = value
}
func (sd *SafeCounts) GetMinCountIdByStatus(status bool) int {
	sd.CountsMu.RLock()
	defer sd.CountsMu.RUnlock()
	min := sd.Counts[0]
	minIndex := -1
	for i, count := range sd.Counts {
		if (minIndex == -1 || count < min) && sd.Status[i] == status {
			min = count
			minIndex = i
		}
	}
	return minIndex
}

func (sd *SafeCounts) GetActiveCountNum() int {
	sd.StatusMu.RLock()
	defer sd.StatusMu.RUnlock()
	total := 0
	for _, status := range sd.Status {
		if status {
			total++
		}
	}
	return total
}

func (sd *SafeCounts) GetActiveCountNumByStatus(status bool) int {
	sd.StatusMu.RLock()
	defer sd.StatusMu.RUnlock()
	total := 0
	for _, iStatus := range sd.Status {
		if iStatus == status {
			total++
		}
	}
	return total
}

func (sd *SafeCounts) GetActiveIdsByStatus(status bool) []int {
	sd.StatusMu.RLock()
	defer sd.StatusMu.RUnlock()
	ids := make([]int, 0)
	for i, iStatus := range sd.Status {
		if iStatus == status {
			ids = append(ids, i)
		}
	}
	return ids
}
