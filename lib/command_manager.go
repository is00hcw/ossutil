package lib

import (
	"fmt"
    "time"
	"reflect"
    "os"
    "runtime"
)

// ParseAndRunCommand parse command line user input, get command and options, then run command 
func ParseAndRunCommand() (error) {
	ts := time.Now().UnixNano()

    clearEnv()

	args, options, err := ParseArgOptions()
	if err != nil {
        return err
	}

	if len(args) == 0 {
		if val, _ := GetBool(OptionVersion, options); val {
            fmt.Printf("ossutil version: %s\n", Version)
            return nil 
		} 
		args = append(args, "help")
	}

	if val, _ := GetBool(OptionMan, options); val {
		args = append([]string{"help"}, args...)
		options[OptionMan] = nil
	}

	command := args[0]
	args = args[1:]

	cm := CommandManager{}
	cm.Init()
	showElapse, err := cm.RunCommand(command, args, options)
	if err != nil {
        return err
	}
	if showElapse {
		te := time.Now().UnixNano()
		fmt.Printf("%.6f(s) elapsed\n", float64(te-ts)/1e9)
        return nil
	}
    return nil
}

func clearEnv() {
    if runtime.GOOS == "windows" {
        _, renameFilePath := getBinaryPath()
        os.Remove(renameFilePath)
    }
}

// CommandManager is used to manager commands, such as build command map and run command
type CommandManager struct {
	commandMap map[string]interface{}
}

// Init build command map
func (cm *CommandManager) Init() {
	commandList := GetAllCommands()
	cm.commandMap = make(map[string]interface{}, len(commandList))

	for _, cmd := range commandList {
		name := reflect.ValueOf(cmd).Elem().FieldByName("command").FieldByName("name").String()
		cm.commandMap[name] = cmd
	}
}

// RunCommand select command from command map, initialize command and run command
func (cm *CommandManager) RunCommand(commandName string, args []string, options OptionMapType) (bool, error) {
	if cmd, ok := cm.commandMap[commandName]; ok {
		if err := cmd.(Commander).Init(args, options); err != nil {
			return false, err
		}
		if err := cmd.(Commander).RunCommand(); err != nil {
			return false, err
		}
        group := reflect.ValueOf(cmd).Elem().FieldByName("command").FieldByName("group").String()
        return group == GroupTypeNormalCommand, nil
	}
	return false, fmt.Errorf("no such command: \"%s\", please try \"help\" or \"--man\" for more information", commandName)
}
