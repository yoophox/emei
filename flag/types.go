package flag

import (
  "flag"
  "time"
)

type FlagSet struct {
  *flag.FlagSet
  args []string
}

func (f *FlagSet) Parse() error {
  return f.FlagSet.Parse(f.args)
}

func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
  f.addArgs(name)
  return f.FlagSet.Bool(name, value, usage)
}

func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
  f.addArgs(name)
  f.FlagSet.BoolVar(p, name, value, usage)
}

func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
  f.addArgs(name)
  return f.FlagSet.Duration(name, value, usage)
}

func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
  f.addArgs(name)
  f.FlagSet.DurationVar(p, name, value, usage)
}

func (f *FlagSet) Float64(name string, value float64, usage string) *float64 {
  f.addArgs(name)
  return f.Float64(name, value, usage)
}

func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
  f.addArgs(name)
  f.FlagSet.Float64Var(p, name, value, usage)
}

func (f *FlagSet) Int(name string, value int, usage string) *int {
  f.addArgs(name)
  return f.FlagSet.Int(name, value, usage)
}

func (f *FlagSet) Int64(name string, value int64, usage string) *int64 {
  f.addArgs(name)
  return f.FlagSet.Int64(name, value, usage)
}

func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
  f.addArgs(name)
  f.Int64Var(p, name, value, usage)
}

func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
  f.addArgs(name)
  f.FlagSet.IntVar(p, name, value, usage)
}

func (f *FlagSet) String(name string, value string, usage string) *string {
  f.addArgs(name)
  return f.FlagSet.String(name, value, usage)
}

func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
  f.addArgs(name)
  f.FlagSet.StringVar(p, name, value, usage)
}

func (f *FlagSet) Uint(name string, value uint, usage string) *uint {
  f.addArgs(name)
  return f.FlagSet.Uint(name, value, usage)
}

func (f *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
  f.addArgs(name)
  return f.FlagSet.Uint64(name, value, usage)
}

func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
  f.addArgs(name)
  f.FlagSet.Uint64Var(p, name, value, usage)
}

func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
  f.addArgs(name)
  f.FlagSet.UintVar(p, name, value, usage)
}

func (f *FlagSet) addArgs(flag string) {
  c, ok := _cmds[flag]
  if ok {
    f.args = append(f.args, "-"+flag)
    if c != "" {
      f.args = append(f.args, c)
    }
  }

  _, ok = _cmds["h"]
  if ok {
    f.args = append(f.args, "--help")
  }
  _, ok = _cmds["help"]
  if ok {
    f.args = append(f.args, "--help")
  }
}
