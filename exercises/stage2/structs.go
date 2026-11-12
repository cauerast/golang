package stage2

// Person is a base type that will be embedded (composition, not inheritance).
type Person struct {
	Name string
}

// Greet returns "Hi, I'm <Name>".
//
// TODO: Implement it.
func (p Person) Greet() string {
	return ""
}

// Employee embeds Person. Person's fields and methods are promoted, so
// e.Name and e.Greet() work without writing a delegating method.
type Employee struct {
	Person
	Role string
}

// Describe returns "<Name> works as <Role>".
//
// TODO: Access the promoted field e.Name directly.
func (e Employee) Describe() string {
	return ""
}

// Manager embeds Employee (two levels deep).
type Manager struct {
	Employee
	Reports int
}

// Greet shadows the promoted Person.Greet.
// It must return "Hi, I'm <Name>, I manage <Reports> people".
//
// TODO: Implement it. Go has no virtual dispatch through embedding:
// Employee.Describe() still calls Person.Greet, never Manager.Greet.
func (m Manager) Greet() string {
	return ""
}

// Stack is a LIFO stack of ints. The zero value is an empty, usable stack.
type Stack struct {
	items []int
}

// Push adds v to the top of the stack.
//
// TODO: Needs a pointer receiver to mutate the struct. Change (s Stack)
// to (s *Stack) and implement it.
func (s Stack) Push(v int) {
}

// Pop removes and returns the top element. If the stack is empty it returns
// (0, false).
//
// TODO: Use a pointer receiver and the comma-ok idiom for the result.
func (s Stack) Pop() (int, bool) {
	return 0, true
}

// Len returns the number of elements.
//
// TODO: Implement it.
func (s Stack) Len() int {
	return -1
}
