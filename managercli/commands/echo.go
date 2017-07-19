package commands

import (
	"github.com/spf13/cobra"
	"fmt"
	"strings"
)

var cmdEcho = &cobra.Command{
	Use: "echo [string to echo]",
	Short: "Echo anything to the screen",
	Long: `bla bla 
	bldlsdfll.
	`,
	Run: echoRun,
}

func echoRun(cmd* cobra.Command. args []string){
	fmt.Println(strings.Join(args, " "))
}

func init() {
	RootCmd.AddCommand(echoCmd)
}
