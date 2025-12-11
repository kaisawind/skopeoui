package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func skopeoTask(ctx context.Context, source, destination string, cb func(txt string)) (err error) {
	// skopeo copy --multi-arch=all --dest-tls-verify=false --retry-delay=3s docker://m.daocloud.io/docker.io/nginx:alpine docker://192.168.1.118:5000/nginx:alpine
	cmd := exec.CommandContext(ctx, "skopeo", "copy", "--multi-arch=all", "--dest-tls-verify=false", "--retry-delay=3s", fmt.Sprintf("docker://%s", source), fmt.Sprintf("docker://%s", destination))
	logrus.Infof("skopeo copy command: %s", cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.WithError(err).Error("skopeo copy stdout pipe failed")
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.WithError(err).Error("skopeo copy stderr pipe failed")
		return
	}
	err = cmd.Start()
	// check err
	if err != nil {
		logrus.WithError(err).Error("skopeo copy start failed")
		return
	}
	rd := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		txt := scanner.Text()
		cb(txt)
	}
	if err := scanner.Err(); err != nil {
		logrus.WithError(err).Warn("scanner error")
		return err
	}
	// wait for command to finish
	cmd.Wait()
	return
}
