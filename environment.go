package main

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func NewEnv() *Environment {
	var env Environment = Environment{}
	env.values = make(map[string]interface{})
	env.enclosing = nil
	return &env
}

func NewEnclosingEnv(enclosing *Environment) *Environment {
	var env Environment = Environment{}
	env.values = make(map[string]interface{})
	env.enclosing = enclosing
	return &env
}

func (env *Environment) Get(name Token) interface{} {
	if value, exists := env.values[name.lexeme]; exists {
		return value
	} else if env.enclosing != nil {
		if enclosedValue, enclosingExists := env.enclosing.values[name.lexeme]; enclosingExists {
			return enclosedValue
		}
	}
	RuntimeError(name.line, name.lexeme, "undefined variable.")
	return 0
}

func (env *Environment) Set(name Token, value interface{}) interface{} {
	if _, exists := env.values[name.lexeme]; !exists {
		if env.enclosing != nil {
			if _, enclosedExists := env.enclosing.values[name.lexeme]; enclosedExists {
				env.enclosing.values[name.lexeme] = value
			} else {
				env.values[name.lexeme] = value
			}
			return value
		}
	}

	env.values[name.lexeme] = value
	return value
}
