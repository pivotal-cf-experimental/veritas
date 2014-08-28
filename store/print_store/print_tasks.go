package print_store

import (
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry-incubator/veritas/say"
	"github.com/cloudfoundry-incubator/veritas/veritas_models"
	"github.com/onsi/gomega/format"
)

func printTasks(verbose bool, tasks veritas_models.VeritasTasks) {
	taskTypes := tasks.OrderedTaskTypes()

	say.PrintBanner(say.Green("Tasks"), "~")

	for _, taskType := range taskTypes {
		say.Println(0, say.Green(taskType))
		for _, task := range tasks[taskType] {
			if verbose {
				printVerboseTask(task)
			} else {
				printTask(task)
			}
		}
	}
}

func printVerboseTask(task models.Task) {
	say.Println(0, format.Object(task, 1))
}

func printTask(task models.Task) {
	say.Println(1,
		"%s [%s on %s@%s(%s)] U:%s C:%s (%d MB, %d MB)",
		taskState(task),
		task.Guid,
		task.ContainerHandle,
		task.ExecutorID,
		task.Stack,
		time.Since(time.Unix(0, task.UpdatedAt)).String(),
		time.Since(time.Unix(0, task.CreatedAt)).String(),
		task.MemoryMB,
		task.DiskMB,
	)
}

func taskState(task models.Task) string {
	switch task.State {
	case models.TaskStatePending:
		return say.LightGray("PENDING  ")
	case models.TaskStateClaimed:
		return say.Yellow("CLAIMED  ")
	case models.TaskStateRunning:
		return say.Cyan("RUNNING  ")
	case models.TaskStateCompleted:
		return colorByTaskSuccess(task, "COMPLETED")
	case models.TaskStateResolving:
		return colorByTaskSuccess(task, "RESOLVING")
	default:
		return say.Red("INVALID")
	}
}

func colorByTaskSuccess(task models.Task, format string, args ...interface{}) string {
	if task.Failed {
		return say.Red(format, args...)
	} else {
		return say.Green(format, args...)
	}
}