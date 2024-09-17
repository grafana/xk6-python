package py

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func isRunCommand(args []string) (bool, int) {
	argn := len(args)

	scriptIndex := argn - 1
	if scriptIndex < 0 {
		return false, scriptIndex
	}

	var runIndex int

	for idx := 0; idx < argn; idx++ {
		arg := args[idx]
		if arg == "run" && runIndex == 0 {
			runIndex = idx

			break
		}
	}

	if runIndex == 0 {
		return false, -1
	}

	return true, scriptIndex
}

//nolint:forbidigo
func redirectStdin() {
	isRun, scriptIndex := isRunCommand(os.Args)
	if !isRun {
		return
	}

	script := os.Args[scriptIndex]
	ext := filepath.Ext(script)
	if script == "-" || (ext != ".star" && ext != ".py") {
		return
	}

	if err := os.Setenv(envScript, script); err != nil {
		logrus.WithError(err).Fatal()
	}

	os.Args[scriptIndex] = "-"

	reader, writer, err := os.Pipe()
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	origStdin := os.Stdin

	os.Stdin = reader

	_, err = writer.Write([]byte(jsScript))
	if err != nil {
		os.Stdin = origStdin

		logrus.WithError(err).Fatal()
	}

	if err := writer.Close(); err != nil {
		os.Stdin = origStdin

		logrus.WithError(err).Fatal()
	}
}

const jsScript = `
export * from 'k6/x/python'
export { default } from 'k6/x/python'
`
