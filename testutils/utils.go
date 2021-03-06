package testutils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	"github.com/akutz/gotil"
)

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func rootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get current working dir:", err)
		os.Exit(1)
	}
	return path.Dir(cwd)
}

func dumpPath(dumpName string) string {
	cwd := rootDir()
	return cwd + "/resources/dumps/" + dumpName
}

func StartGremlinServerWithDump(confFile string, dumpName string) *exec.Cmd {
	CopyFile(dumpPath(dumpName), "/tmp/dump.json")
	return StartGremlinServer(confFile)
}

// StartGremlinServer starts the gremlin-server
func StartGremlinServer(confFile string) *exec.Cmd {
	gremlinServerPath := os.Getenv("GREMLIN_HOME")
	if gremlinServerPath == "" {
		fmt.Fprintln(os.Stderr, "GREMLIN_HOME env variable not set")
		os.Exit(1)
	}
	cwd := rootDir()
	for _, file := range []string{
		"conf/gremlin-contrail.properties",
		"conf/gremlin-contrail.yml",
		"conf/gremlin-neutron.properties",
		"conf/gremlin-neutron.yml",
		"scripts/gremlin-contrail.groovy",
		"bin/foreground.sh",
	} {
		src := fmt.Sprintf("%s/resources/%s", cwd, file)
		dst := fmt.Sprintf("%s/%s", gremlinServerPath, file)
		err := CopyFile(src, dst)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to copy conf", src, "to", dst)
			os.Exit(1)
		}
	}
	cmd := exec.Command("/bin/bash", "bin/foreground.sh", fmt.Sprintf("conf/%s", confFile))
	cmd.Dir = gremlinServerPath
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start process:", err)
		os.Exit(1)
	}
	for gotil.IsTCPPortAvailable(8182) {
		time.Sleep(1 * time.Second)
	}
	time.Sleep(3 * time.Second)
	return cmd
}

func StopGremlinServer(cmd *exec.Cmd) error {
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to kill process:", err)
		return err
	}
	return cmd.Wait()
}
