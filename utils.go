package main

func calculateIdRanges(idRange *IdRange, readersCount int) map[int]*IdRange {
	totalIds := idRange.Max - idRange.Min + 1

	perWorkerIdsCount := totalIds / readersCount

	readersIdRanges := make(map[int]*IdRange, readersCount)
	for i := 1; i <= readersCount; i++ {
		readersIdRanges[i] = &IdRange{
			Max: idRange.Max,
			Min: (idRange.Max - perWorkerIdsCount) + 1,
		}
		idRange.Max -= perWorkerIdsCount
	}

	readersIdRanges[readersCount].Min = idRange.Min

	return readersIdRanges
}
