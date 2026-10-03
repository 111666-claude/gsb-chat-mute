// Command chat-mute 跑禁言样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/mute"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("chat-mute", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "escalate", "escalate / cap / lift / idem / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "escalate":
		world := mute.New(1000, 100000)
		world.Mute("p1", 0, "r1")
		world.Mute("p1", 0, "r2")
		fmt.Fprintf(stdout, "third=%d\n", world.Mute("p1", 0, "r3"))
	case "cap":
		world := mute.New(1000, 3000)
		for index := 0; index < 5; index++ {
			world.Mute("p1", 0, fmt.Sprintf("r-%d", index))
		}
		fmt.Fprintf(stdout, "fifth=%d\n", world.Mute("p1", 0, "r-5"))
	case "lift":
		world := mute.New(1000, 8000)
		world.Mute("p1", 0, "r1")
		world.Lift("p1", 100, "r2")
		fmt.Fprintf(stdout, "muted=%v\n", world.Muted("p1", 200))
	case "idem":
		world := mute.New(1000, 100000)
		world.Mute("p1", 0, "same")
		fmt.Fprintf(stdout, "second=%d\n", world.Mute("p1", 0, "same"))
	case "scan":
		world := mute.New(1000, 100000)
		for index := 0; index < 3000; index++ {
			world.Mute(fmt.Sprintf("p-%d", index), 0, fmt.Sprintf("r-%d", index))
		}
		for index := 0; index < 3000; index++ {
			world.Muted("p-2999", 0)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", world.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample escalate|cap|lift|idem|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
