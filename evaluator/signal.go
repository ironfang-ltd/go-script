package evaluator

type BreakSignal struct{}

func (b *BreakSignal) Debug() string  { return "break" }
func (b *BreakSignal) Type() ObjectType { return BreakSignalObject }

type ContinueSignal struct{}

func (c *ContinueSignal) Debug() string  { return "continue" }
func (c *ContinueSignal) Type() ObjectType { return ContinueSignalObject }
