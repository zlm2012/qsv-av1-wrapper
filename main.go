package main

import (
	"github.com/jessevdk/go-flags"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var opts struct {
	Input   string `short:"i" long:"input" description:"input file" required:"true"`
	Output  string `short:"o" long:"output" description:"output file" required:"true"`
	OrigExe string `long:"orig-exe" description:"original executable path" required:"true"`
	Format  string `long:"format" description:"for replacement, ignored"`
}

func main() {
	parser := flags.NewParser(&opts, flags.None)
	remainedArgs := make([]string, 0)
	parser.UnknownOptionHandler = func(option string, arg flags.SplitArgument, args []string) ([]string, error) {
		log.Println(option)
		if len(option) == 1 {
			remainedArgs = append(remainedArgs, "-"+option)
		} else {
			remainedArgs = append(remainedArgs, "--"+option)
		}
		if value, exists := arg.Value(); exists {
			remainedArgs = append(remainedArgs, value)
		} else if !strings.HasPrefix(args[0], "-") {
			remainedArgs = append(remainedArgs, args[0])
			args = args[1:]
		}
		return args, nil
	}
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		log.Fatalln("error on parse args; ", err)
	}
	var inputReader io.Reader = os.Stdin
	var outputWriter io.Writer = os.Stdout
	if opts.Input != "-" {
		inputReader, err = os.OpenFile(opts.Input, os.O_RDONLY, 0644)
		if err != nil {
			log.Fatalln("failed on read input file; ", err)
		}
	}
	if opts.Output != "-" {
		outputWriter, err = os.OpenFile(opts.Output, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalln("failed on write output file; ", err)
		}
	}

	var cmdArgs = append([]string{"-i", "-", "-o", "-", "--format", "matroska"}, remainedArgs...)
	av1Cmd := exec.Command(opts.OrigExe, cmdArgs...)
	av1Cmd.Stdout = outputWriter
	av1Cmd.Stdin = inputReader
	av1Cmd.Stderr = os.Stderr
	err = av1Cmd.Run()
	if err != nil {
		log.Println("failed on run command: ", err)
	}
	os.Exit(av1Cmd.ProcessState.ExitCode())
}
