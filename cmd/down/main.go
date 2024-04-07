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
	var rpgDownErr, userDownErr, featDownErr, taskDownErr error
	wg.Add(4)
	go func() {
		defer wg.Done()
		rpgDownErr = rpg.DatabaseDown()
	}()
	go func() {
		defer wg.Done()
		userDownErr = user.DatabaseDown()
	}()
	go func() {
		defer wg.Done()
		featDownErr = feat.DatabaseDown()
	}()
	go func() {
		defer wg.Done()
		taskDownErr = task.DatabaseDown()
	}()
	wg.Wait()

	for _, err := range []error{rpgDownErr, userDownErr, featDownErr, taskDownErr} {
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
