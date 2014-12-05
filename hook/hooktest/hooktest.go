package hooktest

// Find out which commands will be generated by which Context methods.

// Runner is implemention of hook.Runner suitable for use in tests.
// It calls the given RunFunc function whenever the Run method
// and records all the calls in the Record field.
// Any calls to juju-log are logged using Logger, but otherwise ignored.
// Any calls to config-get or unit-get are passed through to RunFunc but
// not recorded.
type Runner struct {
	RunFunc func(string, ...string) ([]byte, error)
	Record  [][]string
	Logger  interface {
		Logf(string, ...interface{})
	}

	// Close records whether the Close method has been called.
	Closed bool
}

// Run implements hook.Runner.Run.
func (r *Runner) Run(cmd string, args ...string) ([]byte, error) {
	if cmd == "juju-log" {
		if len(args) != 1 {
			panic("expected only one argument to juju-log")
		}
		r.Logger.Logf("%s", args[0])
		return nil, nil
	}
	if cmd != "config-get" && cmd != "unit-get" {
		rec := []string{cmd}
		rec = append(rec, args...)
		r.Record = append(r.Record, rec)
	}
	return r.RunFunc(cmd, args...)
}

// Run implements hook.Runner.Close.
// It panics if called more than once.
func (r *Runner) Close() error {
	if r.Closed {
		panic("runner closed twice")
	}
	r.Closed = true
	return nil
}

// MemState implements hook.PersistentState in memory.
// Each element of the map holds the value key stored in the state.
type MemState map[string][]byte

func (s MemState) Save(name string, data []byte) error {
	s[name] = data
	return nil
}

func (s MemState) Load(name string) ([]byte, error) {
	return s[name], nil
}

// UUID holds an arbitrary environment UUID for testing purposes.
const UUID = "373b309b-4a86-4f13-88e2-c213d97075b8"
