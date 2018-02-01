package fsm

import "fmt"

// type Event uint64
type State interface{}

type Machine struct {
	State State
	fsm   *FSM
}

func (m *Machine) Goto(s State, args ...interface{}) error {
	fn, ok := m.fsm.GetHandleFunc(m.State, s)
	isSpecial := m.fsm.IsSpecial(s)
	if !ok && !isSpecial { //如果没有，并且不是特殊的函数
		return fmt.Errorf("Transition %v to %v not permitted", m.State, s)
	}
	{
		stateFuncs, ok := m.fsm.GetStateOnFuncs(m.State)
		if ok && stateFuncs.onExit != nil {
			stateFuncs.onExit(args...)
		}
	}
	{
		if fn != nil {
			err := fn(m.State, s, args...)
			if err != nil {
				return err
			}
		}
	}
	{
		stateFuncs, ok := m.fsm.GetStateOnFuncs(s)
		if ok && stateFuncs.onEnter != nil {
			stateFuncs.onEnter(args...)
		}
	}
	m.State = s
	return nil
}

type FSMState struct {
	onEnter func(...interface{})
	onExit  func(...interface{})
}

func (ft *FSMState) SetOnEnter(fn func(...interface{})) *FSMState {
	ft.onEnter = fn
	return ft
}

func (ft *FSMState) SetOnExit(fn func(...interface{})) *FSMState {
	ft.onExit = fn
	return ft
}

type HandleFunc func(State, State, ...interface{}) error

type FSM struct {
	// State State
	rules         map[State]map[State]HandleFunc
	currentState  State
	toState       State
	states        map[State]*FSMState
	specialStates map[State]bool
}

func (fsm *FSM) GetHandleFunc(from State, to State) (HandleFunc, bool) {
	if from == to {
		return nil, true
	}
	maps, ok := fsm.rules[from]
	if !ok {
		return nil, false
	}
	fn, ok := maps[to]
	return fn, ok
	// if !ok {
	// retu
	// }
	// return fn, nil
}

func NewFSM() *FSM {
	f := &FSM{
		rules:         make(map[State]map[State]HandleFunc, 10),
		states:        make(map[State]*FSMState, 10),
		specialStates: make(map[State]bool, 10),
	}
	return f
}

func (fsm *FSM) Machine(s State) *Machine {
	return &Machine{
		State: s,
		fsm:   fsm,
	}
}

func (fsm *FSM) GetStateOnFuncs(s State) (*FSMState, bool) {
	_s, ok := fsm.states[s]
	return _s, ok
	// return nil, nil
}

func (fsm *FSM) SetStateFuncs(s State, onExit func(...interface{}), onEnter func(...interface{})) {
	_s, ok := fsm.states[s]
	if !ok {
		_s = &FSMState{}
		fsm.states[s] = _s
	}
	_s.onEnter = onEnter
	_s.onExit = onExit
}

func (fsm *FSM) From(s State) *FSM {
	_, ok := fsm.rules[s]
	if !ok {
		fsm.rules[s] = make(map[State]HandleFunc, 10)
	}
	fsm.currentState = s
	fsm.toState = s
	return fsm
}

func (fsm *FSM) Special(s State) {
	fsm.specialStates[s] = true
}

func (fsm *FSM) IsSpecial(s State) bool {
	_, ok := fsm.specialStates[s]
	return ok
}

func (fsm *FSM) To(s State) *FSM {
	fsm.toState = s
	fsm.rules[fsm.currentState][s] = nil
	return fsm
	// fsm.rules[fsm.currentState][s]
}

func (fsm *FSM) Then(fn HandleFunc) {
	fsm.rules[fsm.currentState][fsm.toState] = fn
	// return fn(fsm.currentState, fsm.toState)
}

// //-----------------------------------------------------------
// type S int

// type AA struct {
// 	State int
// }

// func (a *AA) Change(from State, to State, args ...interface{}) error {
// 	fmt.Println(from, to, args)
// 	return nil
// }

// func main() {

// 	a := &AA{}

// 	s1 := S(1)
// 	s2 := S(2)
// 	s3 := S(3)
// 	fsm := NewFSM()
// 	fsm.SetStateFuncs(s1, func(args ...interface{}) {
// 		fmt.Println("on exit s1")
// 	}, nil)
// 	fsm.SetStateFuncs(s2, nil, func(args ...interface{}) {
// 		fmt.Println("on enter s2")
// 	})
// 	// fsm.Start(s1)
// 	fsm.From(s1).To(s2).Then(a.Change)

// 	m := fsm.Machine(s1)
// 	m.Goto(s2)
// 	err := m.Goto(s3)
// 	fmt.Println(err)
// }
