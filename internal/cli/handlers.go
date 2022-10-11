package cli

import (
	jr "api-tester/internal/jsonreader"
	"api-tester/internal/net"
	"fmt"
	"strconv"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

func handleTest(cctx *cli.Context) error {
	filename := "./test.json"
	if cctx.String("file") != "" {
		filename = cctx.String("file")
	}
	
	task, err := jr.ReadFile(filename)
	if err != nil {
		return err
	}
	results, err := net.Exec(task)
	if err != nil {
		return err
	}
	for i, r := range results {
		if r.Passed {
			fmt.Printf("%d) %s: %s %s\n", i+1,
				color.Green.Render("PASSED"), color.Blue.Render(strconv.FormatInt(r.Time, 10)), color.Blue.Render("ms"))
		} else {
			fmt.Printf("%d) %s: %s \n", i+1, color.Red.Render("FAILED"), color.Yellow.Render(r.Reason))
		}
	}
	return nil
}
