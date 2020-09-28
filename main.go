package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/mattn/go-colorable"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../")
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())

	var count int
	var instanceCmd = &cobra.Command{
		Use:   "new",
		Short: "Create new instances",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			newDir := filepath.Join(Root, "/instance_*")
			files, err := filepath.Glob(newDir)
			if err != nil {
				panic(err)
			}
			for _, f := range files {
				if err := os.RemoveAll(f); err != nil {
					panic(err)
				}
			}
			i := 1
			for i <= count {
				newDir := filepath.Join(Root, "/instance_"+strconv.Itoa(i))
				err := os.Mkdir(newDir, 0755)
				if err != nil {
					panic(err)
				}

				err = copy.Copy(
					filepath.Join(Root, "/the-crypto-project/"),
					newDir,
				)
				if err != nil {
					panic(err)
				}
				logrus.Info(newDir, "\n")
				i++
			}
		},
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete all instance",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			newDir := filepath.Join(Root, "/instance_*")
			files, err := filepath.Glob(newDir)
			if err != nil {
				panic(err)
			}
			for _, f := range files {
				if err := os.RemoveAll(f); err != nil {
					panic(err)
				}
				logrus.Infof("DELETED: %s \n", f)
			}
		},
	}

	var rootCmd = &cobra.Command{
		Use: "instance",
	}

	instanceCmd.Flags().IntVar(&count, "count", 0, "Instance count")
	rootCmd.AddCommand(deleteCmd, instanceCmd)
	rootCmd.Execute()
}
