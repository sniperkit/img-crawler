package main

import (
	"img-crawler/src/adaptor"
	"img-crawler/src/utils"
)

func main() {

	var wg utils.WaitGroupWrapper

	// Add task
	task_ent_qq := adaptor.Ent_qq()

	// Do task in goroutine
	wg.Wrap(task_ent_qq.Do)

	// Wait all tasks completed
	wg.Wait()
}
