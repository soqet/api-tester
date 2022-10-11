package cli

import (
	jr "api-tester/internal/jsonreader"
	"api-tester/internal/net"
	"fmt"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"strconv"
)

func handleTest(cctx *cli.Context) error {
	filename := cctx.String("file")
	task, err := jr.ReadFile(filename)
	if err != nil {
		return err
	}
	results := make(chan net.Info, 100)
	sem := make(chan int)
	go func() {
		i := 0
		for r := range results {
			if r.Passed {
				fmt.Printf("%d) %s: %s %s\n", i+1,
					color.Green.Render("PASSED"), color.Blue.Render(strconv.FormatInt(r.Time, 10)), color.Blue.Render("ms"))
			} else {
				fmt.Printf("%d) %s: %s \n", i+1, color.Red.Render("FAILED"), color.Yellow.Render(r.Reason))
			}
			i++
		}
		sem <- 1
	}()
	err = net.Exec(task, results)
	if err != nil {
		return err
	}
	<-sem
	return nil
}
