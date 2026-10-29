// Package yqbuilder builds yq expressions without executing them.
package yqbuilder

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Expression struct {
	text string
}

// Raw creates an expression from trusted yq syntax.
func Raw(expression string) Expression {
	return Expression{text: expression}
}

// Quote creates a safely escaped yq string literal.
func Quote(value string) Expression {
	encoded, _ := json.Marshal(value)

	return Raw(string(encoded))
}

func Literal(value any) (Expression, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return Expression{}, fmt.Errorf("encode yq literal for %T: %w", value, err)
	}

	return Raw(string(encoded)), nil
}

func (e Expression) String() string {
	return e.text
}

// MapKey returns a path that safely addresses a string key in a map.
func MapKey(path, key string) string {
	return path + "[" + Quote(key).String() + "]"
}

type Builder struct {
	steps []Expression
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Build() string {
	steps := make([]string, len(b.steps))
	for i, step := range b.steps {
		steps[i] = step.String()
	}

	return strings.Join(steps, " | ")
}

// Steps returns a copy of the pipeline steps.
func (b *Builder) Steps() []Expression {
	return append([]Expression{}, b.steps...)
}

func (b *Builder) Reset() *Builder {
	b.steps = b.steps[:0]

	return b
}

func (b *Builder) Pipe(expression Expression) *Builder {
	b.steps = append(b.steps, expression)

	return b
}

func (b *Builder) Set(path string, value Expression) *Builder {
	return b.Pipe(Raw(path + " = " + value.String()))
}

func (b *Builder) Append(path string, value Expression) *Builder {
	return b.Pipe(Raw(path + " += [" + value.String() + "]"))
}

func (b *Builder) Extend(path string, values Expression) *Builder {
	return b.Pipe(Raw(path + " += " + values.String()))
}

func (b *Builder) Update(path string, expression Expression) *Builder {
	return b.Pipe(Raw(path + " |= (" + expression.String() + ")"))
}

func (b *Builder) Delete(path string) *Builder {
	return b.Pipe(Raw("del(" + path + ")"))
}
