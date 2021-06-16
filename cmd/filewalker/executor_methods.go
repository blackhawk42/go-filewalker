package main

import (
	"sort"

	"github.com/blackhak42/go-filewalker/executor"
)

type ExecutorMethods map[string]executor.Executor

func (fm ExecutorMethods) GetMethods() []string {
	methods := make([]string, 0, len(fm))
	for m := range fm {
		methods = append(methods, m)
	}

	sort.Strings(methods)

	return methods
}

var AvaiableExecutorMethods = ExecutorMethods(map[string]executor.Executor{
	"report": executor.ReportExecutor,
	"move":   executor.MoveExecutor,
	"copy":   executor.CopyExecutor,
})
