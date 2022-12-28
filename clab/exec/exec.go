package exec

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/google/shlex"
)

type ExecFormat string

const (
	ExecFormatJSON  ExecFormat = "json"
	ExecFormatPlain ExecFormat = "plain"
)

var ErrRunExecNotSupported = errors.New("exec not supported for this kind")

type ExecResultHolderCreateFn func(*ExecCmd) ExecResultHolderSetter

type ExecResultHolderSetter interface {
	GetExecResultHolder() ExecResultHolder
	SetReturnCode(int)
	SetStdOut([]byte)
	SetStdErr([]byte)
}

// ParseExecOutputFormat parses the exec output format user input.
func ParseExecOutputFormat(s string) (ExecFormat, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case string(ExecFormatJSON):
		return ExecFormatJSON, nil
	case string(ExecFormatPlain), "table":
		return ExecFormatPlain, nil
	}
	return "", fmt.Errorf("cannot parse %q as execution output format, supported output formats %q",
		s, []string{string(ExecFormatJSON), string(ExecFormatPlain)})
}

func GetResultHolderCreateFnFor(eof ExecFormat) ExecResultHolderCreateFn {
	switch eof {
	case ExecFormatJSON:
		return NewExecResultJson
	}
	// default is the plain result
	return NewExecResult
}

// ExecCmd represents an exec command.
type ExecCmd struct {
	Cmd []string `json:"cmd"` // Cmd is a slice-based representation of a string command.
}

// NewExecCmdFromString creates ExecCmd for a string-based command.
func NewExecCmdFromString(cmd string) (*ExecCmd, error) {
	result := &ExecCmd{}
	if err := result.SetCmd(cmd); err != nil {
		return nil, err
	}
	return result, nil
}

// NewExecCmdFromSlice creates ExecCmd for a command represented as a slice of strings.
func NewExecCmdFromSlice(cmd []string) *ExecCmd {
	return &ExecCmd{
		Cmd: cmd,
	}
}

type ExecResultHolder interface {
	GetStdOutString() string
	GetStdErrString() string
	GetStdOutByteSlice() []byte
	GetStdErrByteSlice() []byte
	GetReturnCode() int
	GetCmdString() string
	Dump(format ExecFormat) (string, error)
	String() string
}

// ExecResult represents a result of a command execution.
type ExecResult struct {
	Cmd        []string `json:"cmd"`
	ReturnCode int      `json:"return-code"`
	Stdout     string   `json:"stdout"`
	Stderr     string   `json:"stderr"`
}

func NewExecResult(op *ExecCmd) ExecResultHolderSetter {
	er := &ExecResult{Cmd: op.GetCmd()}
	return er
}

type ExecResultJson struct {
	ExecResult
	Stdout json.RawMessage `json:"stdout"`
}

func (erj *ExecResultJson) SetStdOut(b []byte) {
	erj.Stdout = b
}

func (e *ExecResultJson) GetExecResultHolder() ExecResultHolder {
	return e
}

func NewExecResultJson(op *ExecCmd) ExecResultHolderSetter {
	er := &ExecResultJson{
		ExecResult: ExecResult{Cmd: op.GetCmd()},
	}
	return er
}

// SetCmd sets the command that is to be executed.
func (e *ExecCmd) SetCmd(cmd string) error {
	c, err := shlex.Split(cmd)
	if err != nil {
		return err
	}
	e.Cmd = c
	return nil
}

// GetCmd sets the command that is to be executed.
func (e *ExecCmd) GetCmd() []string {
	return e.Cmd
}

// GetCmdString sets the command that is to be executed.
func (e *ExecCmd) GetCmdString() string {
	return strings.Join(e.Cmd, " ")
}

func (e *ExecResult) String() string {
	return fmt.Sprintf("Cmd: %s\nReturnCode: %d\nStdOut:\n%s\nStdErr:\n%s\n", e.GetCmdString(), e.ReturnCode, e.Stdout, e.Stderr)
}

func (e *ExecResult) GetExecResultHolder() ExecResultHolder {
	return e
}

// Dump dumps execution result as a string in one of the provided formats.
func (e *ExecResult) Dump(format ExecFormat) (string, error) {
	var result string
	switch format {
	case ExecFormatJSON:
		byteData, err := json.MarshalIndent(e, "", "  ")
		if err != nil {
			return "", err
		}
		result = string(byteData)
	case ExecFormatPlain:
		result = e.String()
	}
	return result, nil
}

// GetCmdString returns the initially parsed cmd as a string for e.g. log output purpose.
func (e *ExecResult) GetCmdString() string {
	return strings.Join(e.Cmd, " ")
}

func (e *ExecResult) GetReturnCode() int {
	return e.ReturnCode
}

func (e *ExecResult) SetReturnCode(rc int) {
	e.ReturnCode = rc
}

func (e *ExecResult) GetStdOutString() string {
	return string(e.Stdout)
}

func (e *ExecResult) GetStdErrString() string {
	return string(e.Stderr)
}

func (e *ExecResult) GetStdOutByteSlice() []byte {
	return []byte(e.Stdout)
}

func (e *ExecResult) GetStdErrByteSlice() []byte {
	return []byte(e.Stderr)
}

func (e *ExecResult) GetCmd() []string {
	return e.Cmd
}

func (e *ExecResult) SetStdOut(data []byte) {
	e.Stdout = string(data)
}

func (e *ExecResult) SetStdErr(data []byte) {
	e.Stderr = string(data)
}

// execEntries is a map indexed by container IDs storing lists of ExecResultHolder.
// ExecResultHolder is an interface that is backed by the type storing data for the executed command.
type execEntries map[string][]ExecResultHolder

// ExecCollection represents a datastore for exec commands execution results.
type ExecCollection struct {
	execEntries
}

// NewExecCollection initializes the collection of exec command results.
func NewExecCollection() *ExecCollection {
	return &ExecCollection{
		execEntries{},
	}
}

func (ec *ExecCollection) Add(cId string, e ExecResultHolder) {
	ec.execEntries[cId] = append(ec.execEntries[cId], e)
}

func (ec *ExecCollection) AddAll(cId string, e []ExecResultHolder) {
	ec.execEntries[cId] = append(ec.execEntries[cId], e...)
}

// Dump dumps the contents of ExecCollection as a string in one of the provided formats.
func (ec *ExecCollection) Dump(format ExecFormat) (string, error) {
	result := strings.Builder{}
	switch format {
	case ExecFormatJSON:
		byteData, err := json.MarshalIndent(ec.execEntries, "", "  ")
		if err != nil {
			return "", err
		}
		result.Write(byteData)
	case ExecFormatPlain:
		printSep := false
		for k, execResults := range ec.execEntries {
			if len(execResults) == 0 {
				// skip if there is no result
				continue
			}
			// write seperator
			if printSep {
				result.WriteString("\n+++++++++++++++++++++++++++++\n\n")
			}
			// write header for entry
			result.WriteString("Node: ")
			result.WriteString(k)
			result.WriteString("\n")
			for _, er := range execResults {
				// write entry
				result.WriteString(er.String())
			}
			// starting second run, print sep
			printSep = true
		}
	}
	return result.String(), nil
}

// Log writes to the log execution results stored in ExecCollection.
// If execution result contains error, the error log facility is used,
// otherwise it is logged as INFO.
func (ec *ExecCollection) Log() {
	for k, execResults := range ec.execEntries {
		for _, er := range execResults {
			switch {
			case er.GetReturnCode() != 0 || er.GetStdErrString() != "":
				log.Errorf("Failed to execute command '%s' on node %s. rc=%d,\nstdout:\n%s\nstderr:\n%s",
					er.GetCmdString(), k, er.GetReturnCode(), er.GetStdOutString(), er.GetStdErrString())
			default:
				log.Infof("Executed command '%s' on node %s. stdout:\n%s",
					er.GetCmdString(), k, er.GetStdOutString())
			}
		}
	}
}
