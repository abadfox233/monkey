package object

type Environment struct {
	store  map[string]Object
	outer  *Environment
	inLoop bool
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Upsert(name string, val Object) Object {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return val
	}
	if e.outer != nil {
		return e.outer.Upsert(name, val)
	}
	e.store[name] = val
	return val
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) SetInLoop() {
	e.inLoop = true
}

func (e *Environment) SetOutLoop() {
	e.inLoop = false
}

func (e *Environment) IsInLoop() bool {
	if e.inLoop {
		return true
	}
	if e.outer != nil {
		return e.outer.IsInLoop()
	}
	return false
}
