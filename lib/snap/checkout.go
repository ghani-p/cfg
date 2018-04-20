package snap

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Confbase/cfg/lib/dotcfg"
)

func Checkout(name string) {
	cfg := dotcfg.MustLoadCfg()
	if cfg.NoGit {
		fmt.Fprintf(os.Stderr, "error: checkout is not a valid command in a non-git base")
		os.Exit(1)
	}

	snapExists := false
	snaps := dotcfg.MustLoadSnaps()
	for _, s := range snaps.Snapshots {
		if s.Name == name {
			snapExists = true
			break
		}
	}
	if !snapExists {
		fmt.Fprintf(os.Stderr, "error: there is no snapshot named '%v'\n", name)
		os.Exit(1)
	}

	statusCmd := exec.Command("git", "status", "-s")
	stsBytes, err := statusCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to run 'git status -s'\n")
		os.Exit(1)
	}
	sts := string(stsBytes)
	if sts != "" {
		for _, line := range strings.Split(sts, "\n") {
			if len(line) >= 2 && line[0] != '?' && line[1] != '?' {
				fmt.Fprintf(os.Stderr, "error: there are uncommitted files\n")
				fmt.Fprintf(os.Stderr, "'git status -s' output:\n%v", sts)
				os.Exit(1)
			}
		}
	}

	gitCmd := exec.Command("git", "checkout", name)
	if err := gitCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to checkout snapshot\n%v\n", err)
		os.Exit(1)
	}
}
