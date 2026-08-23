// Package yqbuilder builds yq expressions without executing them.
package yqbuilder

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Expression is a fragment that can be embedded in a yq expression.
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

// Literal converts a Go value to a yq-compatible JSON literal.
func Literal(value any) (Expression, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return Expression{}, fmt.Errorf("encode yq literal for %T: %w", value, err)
	}
	return Raw(string(encoded)), nil
}

// String returns the expression text.
func (e Expression) String() string {
	return e.text
}

// MapKey returns a path that safely addresses a string key in a map.
func MapKey(path, key string) string {
	return path + "[" + Quote(key).String() + "]"
}

// Builder joins yq operations into a pipeline.
type Builder struct {
	steps []Expression
}

// New creates an empty builder.
func New() *Builder {
	return &Builder{}
}

// Build returns the complete yq pipeline.
func (b *Builder) Build() string {
	steps := make([]string, len(b.steps))
	for i, step := range b.steps {
		steps[i] = step.String()
	}
	return strings.Join(steps, " | ")
}

// Steps returns a copy of the pipeline steps.
func (b *Builder) Steps() []Expression {
	steps := make([]Expression, len(b.steps))
	copy(steps, b.steps)
	return steps
}

// Reset removes all pipeline steps.
func (b *Builder) Reset() *Builder {
	b.steps = b.steps[:0]
	return b
}

// Pipe appends an arbitrary expression to the pipeline.
func (b *Builder) Pipe(expression Expression) *Builder {
	b.steps = append(b.steps, expression)
	return b
}

// Set assigns a value to a path.
func (b *Builder) Set(path string, value Expression) *Builder {
	return b.Pipe(Raw(path + " = " + value.String()))
}

// Append appends one value to an array.
func (b *Builder) Append(path string, value Expression) *Builder {
	return b.Pipe(Raw(path + " += [" + value.String() + "]"))
}

// Extend appends an array expression to an array.
func (b *Builder) Extend(path string, values Expression) *Builder {
	return b.Pipe(Raw(path + " += " + values.String()))
}

// Update applies an expression to the value at a path.
func (b *Builder) Update(path string, expression Expression) *Builder {
	return b.Pipe(Raw(path + " |= (" + expression.String() + ")"))
}

// Delete removes a path.
func (b *Builder) Delete(path string) *Builder {
	return b.Pipe(Raw("del(" + path + ")"))
}
