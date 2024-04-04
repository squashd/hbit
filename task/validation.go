package task

import "github.com/SQUASHD/hbit"

func validateRepeatSchedule(taskSchedule map[Day]bool) error {
	defaultTaskMap := getDefaultSchedule()
	for day := range defaultTaskMap {
		if _, ok := taskSchedule[day]; !ok {
			return &hbit.Error{Code: hbit.EINVALID, Message: "Invalid schedule"}
		}
	}
	return nil
}

func getDefaultSchedule() map[Day]bool {
	return map[Day]bool{
		MONDAY:    false,
		TUESDAY:   false,
		WEDNESDAY: false,
		THURSDAY:  false,
		FRIDAY:    false,
		SATURDAY:  false,
		SUNDAY:    false,
	}
}
