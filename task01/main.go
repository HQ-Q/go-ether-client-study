package main

import "eth-client-study/task01/app"

func main() {
	task01 := app.Task01{}
	task01.QueryBlockInfo()
	task01.TransferEth()
	task01.DeployCounterContract()
}
