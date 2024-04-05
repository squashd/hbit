package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/SQUASHD/hbit/feat"
	"github.com/SQUASHD/hbit/rpg"
	"github.com/SQUASHD/hbit/task"
	"github.com/SQUASHD/hbit/user"
)

func main() {
	var errs []error
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		rpgDownErr := rpg.DatabaseDown()
		if rpgDownErr != nil {
			errs = append(errs, rpgDownErr)
		}
	}()
	go func() {
		defer wg.Done()
		userDownErr := user.DatabaseDown()
		if userDownErr != nil {
			errs = append(errs, userDownErr)
		}
	}()
	go func() {
		defer wg.Done()
		featDownErr := feat.DatabaseDown()
		if featDownErr != nil {
			errs = append(errs, featDownErr)
		}
	}()
	go func() {
		defer wg.Done()
		taskDownErr := task.DatabaseDown()
		if taskDownErr != nil {
			errs = append(errs, taskDownErr)
		}
	}()

	for _, err := range errs {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		fmt.Println(errors.Join(errs...))
	} else {
		fmt.Println("All databases down")
	}

}
