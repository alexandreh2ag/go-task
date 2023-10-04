package types

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func Test_WorkerTask_SuccessValidate(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id:      "test",
		Command: "fake",
	}
	err := validate.Struct(worker)

	assert.NoError(t, err)
}

func Test_WorkerTask_SuccessValidateWithOptionalData(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id:        "test",
		Command:   "fake",
		User:      "test",
		Directory: "/tmp/test/",
	}
	err := validate.Struct(worker)

	assert.NoError(t, err)
}

func Test_WorkerTask_SuccessValidateIDWithUnderscore(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id:        "test_foo",
		Command:   "fake",
		User:      "test",
		Directory: "/tmp/test/",
	}
	err := validate.Struct(worker)

	assert.NoError(t, err)
}

func Test_WorkerTask_ErrorValidate(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id: "test",
	}
	err := validate.Struct(worker)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'Command' failed on the 'required' tag")
}

func Test_WorkerTask_ErrorValidateComplex(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id:        "test",
		Command:   "fake",
		User:      "user/test",
		Directory: "wrong",
	}
	err := validate.Struct(worker)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'User' failed on the 'alphanum' tag")
	assert.Contains(t, err.Error(), "Field validation for 'Directory' failed on the 'dirpath' tag")
}

func Test_WorkerTask_ErrorValidateID(t *testing.T) {
	validate := validator.New()
	worker := WorkerTask{
		Id:        "test foo",
		Command:   "fake",
		User:      "user",
		Directory: "/tmp",
	}
	err := validate.Struct(worker)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field validation for 'Id' failed on the 'excludesall' tag")
}

func TestPrepareWorkerTasks(t *testing.T) {
	type args struct {
		tasks      WorkerTasks
		user       string
		workingDir string
		groupName  string
	}
	tests := []struct {
		name string
		args args
		want WorkerTasks
	}{
		{
			name: "SuccessEmptyTasks",
			args: args{
				tasks:      WorkerTasks{},
				user:       "foo",
				groupName:  "bar",
				workingDir: "/app/foo/",
			},
			want: WorkerTasks{},
		},
		{
			name: "SuccessMultipleTasks",
			args: args{
				tasks: WorkerTasks{
					&WorkerTask{Id: "test", Command: "cmd", GroupName: "bar"},
					&WorkerTask{Id: "test2", Command: "cmd", GroupName: "bar", User: "bar", Directory: "/app/bar/"},
				},
				user:       "foo",
				groupName:  "bar",
				workingDir: "/app/foo/",
			},
			want: WorkerTasks{
				&WorkerTask{Id: "test", Command: "cmd", GroupName: "bar", User: "foo", Directory: "/app/foo/"},
				&WorkerTask{Id: "test2", Command: "cmd", GroupName: "bar", User: "bar", Directory: "/app/bar/"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrepareWorkerTasks(tt.args.tasks, tt.args.groupName, tt.args.user, tt.args.workingDir)
			assert.Equal(t, tt.want, tt.args.tasks)
		})
	}
}

func TestPrefixedName(t *testing.T) {
	prefix := "group"
	worker := WorkerTask{
		Id:        "test",
		Command:   "fake",
		GroupName: prefix,
	}

	assert.Equal(t, worker.PrefixedId(), prefix+"-"+worker.Id)
}

func TestGetUniqueExtraGroupsSingleResult(t *testing.T) {
	tasks := WorkerTasks{
		&WorkerTask{Id: "test", Command: "cmd", GroupName: "bar", ExtraGroups: []string{"foo"}},
		&WorkerTask{Id: "test2", Command: "cmd", GroupName: "bar", User: "bar", Directory: "/app/bar/", ExtraGroups: []string{"foo"}},
	}
	extraGroups := tasks.GetUniqueExtraGroups()
	assert.Equal(t, slices.Contains(extraGroups, "foo"), true)
	assert.Equal(t, len(extraGroups), 1)
}

func TestGetUniqueExtraGroupsMultipleResult(t *testing.T) {
	tasks := WorkerTasks{
		&WorkerTask{Id: "test", Command: "cmd", GroupName: "bar", ExtraGroups: []string{"foo", "bar"}},
		&WorkerTask{Id: "test2", Command: "cmd", GroupName: "bar", User: "bar", Directory: "/app/bar/", ExtraGroups: []string{"foo"}},
	}
	extraGroups := tasks.GetUniqueExtraGroups()
	assert.Equal(t, slices.Contains(extraGroups, "foo"), true)
	assert.Equal(t, slices.Contains(extraGroups, "bar"), true)
	assert.Equal(t, len(extraGroups), 2)
}
func TestGetTasksInGroupAll(t *testing.T) {
	tasks := WorkerTasks{
		&WorkerTask{Id: "test", Command: "cmd", GroupName: "bar", ExtraGroups: []string{"foo", "bar"}},
		&WorkerTask{Id: "test2", Command: "cmd", GroupName: "bar", User: "bar", Directory: "/app/bar/", ExtraGroups: []string{"foo"}},
	}
	selectedTask := tasks.GetTasksInGroup("")
	assert.Equal(t, len(selectedTask), 2)
}

func TestGetTasksInGroup(t *testing.T) {
	tasks := WorkerTasks{
		&WorkerTask{Id: "test", Command: "cmd", GroupName: "bar", ExtraGroups: []string{"foo", "bar"}},
		&WorkerTask{Id: "test2", Command: "cmd", GroupName: "bar", User: "bar", Directory: "/app/bar/", ExtraGroups: []string{"foo"}},
	}
	selectedTask := tasks.GetTasksInGroup("bar")
	assert.Equal(t, len(selectedTask), 1)
}

func TestGetTasksInGroupNoResult(t *testing.T) {
	tasks := WorkerTasks{
		&WorkerTask{Id: "test", Command: "cmd", GroupName: "bar", ExtraGroups: []string{"foo"}},
		&WorkerTask{Id: "test2", Command: "cmd", GroupName: "bar", User: "bar", Directory: "/app/bar/", ExtraGroups: []string{"foo"}},
	}
	selectedTask := tasks.GetTasksInGroup("bar")
	assert.Equal(t, len(selectedTask), 0)
}

func TestGetAllPrefixedId(t *testing.T) {
	tasks := WorkerTasks{
		&WorkerTask{Id: "test", Command: "cmd", GroupName: "bar", ExtraGroups: []string{"foo"}},
		&WorkerTask{Id: "test2", Command: "cmd", GroupName: "bar", User: "bar", Directory: "/app/bar/", ExtraGroups: []string{"foo"}},
	}
	ids := tasks.GetAllPrefixedId()
	assert.Contains(t, ids, tasks[0].PrefixedId())
	assert.Contains(t, ids, tasks[1].PrefixedId())
	assert.Equal(t, len(ids), 2)
}
