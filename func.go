package main

type ArgList struct {
	args []*Expression
}

func NewArgList(expr *Expression) *ArgList {
	a := &ArgList{}
	a.args = make([]*Expression, 1)
	a.args[0] = expr

	return a
}

func (args *ArgList) appendArg(expr *Expression) *ArgList {
	args.args = append(args.args, expr)
	return args
}

type Function struct {
	builtin func(args ...Node) Node
}

func (f *Function) call(args *ArgList) Node {
	ev := make([]Node, 0)
	for _, v := range args.args {
		ev = append(ev, v.eval())
	}
	return f.builtin(ev...)
}

func builtinFunctions() map[string]*Function {
	f := map[string]*Function{
		"exit": &Function{
			// exit(int) void
			builtin: func(args ...Node) Node {
				exitCode := args[0]
				execContext.exitCode = int(exitCode.asNumber())
				execContext.doExit = true
				return nil
			},
		},
	}

	return f
}
