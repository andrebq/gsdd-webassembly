package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/andrebq/gsdd-webassembly/hchan"

	"github.com/perlin-network/life/exec"
	"github.com/sirupsen/logrus"
)

// Resolver defines imports for WebAssembly modules ran in Life.
type Resolver struct {
	tempRet0 int64
}

type pingPong struct {
	ID string
}

func timedFn(fnName string, actual exec.FunctionImport) exec.FunctionImport {
	entry := logrus.WithField("fn", fnName)
	return func(vm *exec.VirtualMachine) int64 {
		start := time.Now()
		ret := actual(vm)
		entry.WithField("duration", time.Now().Sub(start)).Info()
		return ret
	}
}

// ResolveFunc defines a set of import functions that may be called within a WebAssembly module.
func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {
	fmt.Printf("Resolve func: %s %s\n", module, field)
	switch module {
	case "env":
		switch field {
		case "__life_ping":
			return timedFn("__life_ping", func(vm *exec.VirtualMachine) int64 {
				return vm.GetCurrentFrame().Locals[0] + 1
			})
		case "__life_log":
			return timedFn("__life_log", func(vm *exec.VirtualMachine) int64 {
				msg := readMemBytes(vm, 0)
				logrus.WithField("app", string(msg)).Info()
				return 0
			})
		case "__life_write_hchan":
			return timedFn("__life_write_hchan", func(vm *exec.VirtualMachine) int64 {
				msg := readMemBytes(vm, 0)
				err := hchan.Write("http://localhost:8082/write/ping", pingPong{ID: string(msg)})
				if err != nil {
					logrus.WithError(err).Error("error sending data to channel")
				}
				return 0
			})

		case "__life_read_hchan":
			return timedFn("__life_read_hchan", func(vm *exec.VirtualMachine) int64 {
				var out pingPong
				err := hchan.Read(&out, "http://localhost:8082/read/ping")
				if err != nil {
					logrus.WithError(err).Error("error sending data to channel")
				}
				buf, _ := json.Marshal(out) // marshall the output to send it back to rust
				if err := writeMemBytes(vm, 0, buf); err != nil {
					logrus.WithError(err).Error("unable to write all the bytes needed")
					return 0
				}
				return int64(len(buf))
			})
		default:
			panic(fmt.Errorf("unknown field: %s", field))
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module))
	}
}

func readMemBytes(vm *exec.VirtualMachine, localIdx int) []byte {
	ptr := int(uint32(vm.GetCurrentFrame().Locals[0+localIdx]))
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1+localIdx]))
	msg := vm.Memory[ptr : ptr+msgLen]
	return msg
}

func writeMemBytes(vm *exec.VirtualMachine, localIdx int, buf []byte) error {
	ptr := int(uint32(vm.GetCurrentFrame().Locals[0+localIdx]))
	msgLen := int(uint32(vm.GetCurrentFrame().Locals[1+localIdx]))

	if len(buf) > msgLen {
		// will not write since there is no space, and I am too lazy
		// to implement proper buffered
		// TODO: improve on this
		return errors.New("buffer to small")
	}

	copy(vm.Memory[ptr:ptr+msgLen], buf)
	return nil
}

// ResolveGlobal defines a set of global variables for use within a WebAssembly module.
func (r *Resolver) ResolveGlobal(module, field string) int64 {
	fmt.Printf("Resolve global: %s %s\n", module, field)
	switch module {
	case "env":
		switch field {
		case "__life_magic":
			return 424
		default:
			panic(fmt.Errorf("unknown field: %s", field))
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module))
	}
}

func main() {
	entryFunctionFlag := flag.String("entry", "app_main", "entry function id")
	jitFlag := flag.Bool("jit", false, "enable jit")
	flag.Parse()

	// Read WebAssembly *.wasm file.
	input, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	// Instantiate a new WebAssembly VM with a few resolved imports.
	vm, err := exec.NewVirtualMachine(input, exec.VMConfig{
		EnableJIT:          *jitFlag,
		DefaultMemoryPages: 128,
		DefaultTableSize:   65536,
	}, new(Resolver), nil)

	if err != nil {
		panic(err)
	}

	// Get the function ID of the entry function to be executed.
	entryID, ok := vm.GetFunctionExport(*entryFunctionFlag)
	if !ok {
		fmt.Printf("Entry function %s not found; starting from 0.\n", *entryFunctionFlag)
		entryID = 0
	}

	start := time.Now()

	// If any function prior to the entry function was declared to be
	// called by the module, run it first.
	if vm.Module.Base.Start != nil {
		startID := int(vm.Module.Base.Start.Index)
		_, err := vm.Run(startID)
		if err != nil {
			vm.PrintStackTrace()
			panic(err)
		}
	}
	var args []int64
	for _, arg := range flag.Args()[1:] {
		fmt.Println(arg)
		if ia, err := strconv.Atoi(arg); err != nil {
			panic(err)
		} else {
			args = append(args, int64(ia))
		}
	}

	// Run the WebAssembly module's entry function.
	ret, err := vm.Run(entryID, args...)
	if err != nil {
		vm.PrintStackTrace()
		panic(err)
	}
	end := time.Now()

	fmt.Printf("return value = %d, duration = %v\n", ret, end.Sub(start))
}
