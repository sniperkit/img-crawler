package main

import (
	"img-crawler/src/adaptor"
	"img-crawler/src/utils"
)

func main() {

	var wg utils.WaitGroupWrapper

	// Add task
	//task_ent_qq := adaptor.Ent_qq()
	//wg.Wrap(task_ent_qq.Do)

	//task_pic_699 := adaptor.Pic_699()
	//wg.Wrap(task_pic_699.Do)

	task_renren := adaptor.RenRen()
	wg.Wrap(task_renren.Do)

	// Wait all tasks completed
	wg.Wait()
}
