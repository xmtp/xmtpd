package cmd

import "fmt"

type Version struct {
	GitCommit string
}

func (c *Version) Execute(args []string) error {
	fmt.Println(c.GitCommit)
	return nil
}
