package main

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	var env Environment = Environment{}
	env.values = make(map[string]interface{})
	return &env
}

func (env *Environment) Get(name Token) interface{} {
	if value, exists := env.values[name.lexeme]; exists {
		return value
	} else {
		panic("Error, undefined variable!")
	}
}

func (env *Environment) Set(name Token, value interface{}) interface{} {
	env.values[name.lexeme] = value
	return value
}
