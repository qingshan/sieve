package parse

import (
	"testing"
	"github.com/qingshan/sieve/ast"
)


func TestGenericCommand(t *testing.T) {
	input := `header :contains "Subject" "subject keyword";`
	p := Parse("TestGenericCommand", input)
	file := &ast.File{p.name, p.List}
	commandList := []ast.Command{
		&ast.GenericCommand{"header", []ast.Argument{
			ast.TagArgument(":contains"),
			ast.StringArgument([]string{"Subject"}),
			ast.StringArgument([]string{"subject keyword"}),
		}},
	}
	expected := &ast.File{
		Name: p.name,
		List: commandList,
	}
	if !ast.Equals(file, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), file.String())
	}
}

func TestIfControlCommand(t *testing.T) {
	input := `if header :contains "Subject" ["keyword1", "keyword2"] { discard :under "test"; stop; }`
	p := Parse("TestIfControlCommand", input)
	file := &ast.File{p.name, p.List}
	commandList := []ast.Command{
		&ast.ControlCommand{
			Name: "if",
			Test: &ast.GenericTest{"header",
				[]ast.Argument{
					ast.TagArgument(":contains"),
					ast.StringArgument([]string{"Subject"}),
					ast.StringArgument([]string{"keyword1", "keyword2"}),
				}},
			Block: []ast.Command{
				&ast.GenericCommand{
					"discard",
					[]ast.Argument{
						ast.TagArgument(":under"),
						ast.StringArgument([]string{"test"}),
					}},
				&ast.StopCommand{},
			},
		},
	}
	expected := &ast.File{
		Name: p.name,
		List: commandList,
	}
	if !ast.Equals(file, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), file.String())
	}
}

func TestIfElseControlCommand(t *testing.T) {
	input := `if header :contains "Subject" ["keyword1", "keyword2"] { discard :under "test"; stop; } else { stop; }`
	p := Parse("TestIfElseControlCommand", input)
	file := &ast.File{p.name, p.List}
	commandList := []ast.Command{
		&ast.ControlCommand{
			Name: "if",
			Test: &ast.GenericTest{"header",
				[]ast.Argument{
					ast.TagArgument(":contains"),
					ast.StringArgument([]string{"Subject"}),
					ast.StringArgument([]string{"keyword1", "keyword2"}),
				}},
			Block: []ast.Command{
				&ast.GenericCommand{
					"discard",
					[]ast.Argument{
						ast.TagArgument(":under"),
						ast.StringArgument([]string{"test"}),
					}},
				&ast.StopCommand{},
			},
		},
		&ast.ControlCommand{
			Name: "else",
			Block: []ast.Command{
				&ast.StopCommand{},
			},
		},
	}
	expected := &ast.File{
		Name: p.name,
		List: commandList,
	}
	if !ast.Equals(file, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), file.String())
	}
}

func TestIfElsifElseControlCommand(t *testing.T) {
	input := `if header :contains "Subject" ["keyword1", "keyword2"] { discard :under "test"; stop; } elsif true { stop; } else { stop; }`
	p := Parse("TestIfElsifElseControlCommand", input)
	file := &ast.File{p.name, p.List}
	commandList := []ast.Command{
		&ast.ControlCommand{
			Name: "if",
			Test: &ast.GenericTest{"header",
				[]ast.Argument{
					ast.TagArgument(":contains"),
					ast.StringArgument([]string{"Subject"}),
					ast.StringArgument([]string{"keyword1", "keyword2"}),
				}},
			Block: []ast.Command{
				&ast.GenericCommand{
					"discard",
					[]ast.Argument{
						ast.TagArgument(":under"),
						ast.StringArgument([]string{"test"}),
					}},
				&ast.StopCommand{},
			},
		},
		&ast.ControlCommand{
			Name: "elsif",
			Test: &ast.TrueTest{},
			Block: []ast.Command{
				&ast.StopCommand{},
			},
		},
		&ast.ControlCommand{
			Name: "else",
			Block: []ast.Command{
				&ast.StopCommand{},
			},
		},
	}
	expected := &ast.File{
		List: commandList,
	}
	if !ast.Equals(file, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), file.String())
	}
}