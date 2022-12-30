package assets

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/go-containerregistry/cmd/crane/cmd"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/sirupsen/logrus"
)

var AssetDir = "assets"

const (
	K0sBinary = "k0s"
)

func Download() error {
	err := mkAssetDir()
	if err != nil {
		return err
	}

	err = downloadK0sBinary()
	if err != nil {
		return err
	}

	err = downloadAirgapImages()
	if err != nil {
		return err
	}

	ok := OK()
	if !ok {
		return fmt.Errorf("assets not downloaded")
	}

	return nil
}

func downloadK0sBinary() error {
	// skip when file already exists
	if _, err := os.Stat(K0sBinaryPath()); err == nil {
		logrus.Warnf("k0s binary already exists at %s", K0sBinaryPath())
		return nil
	}

	// latest url
	latestURL := "https://docs.k0sproject.io/stable.txt"
	resp, err := http.Get(latestURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// download url
	K0S_VERSION := strings.ReplaceAll(string(body), "\n", "")
	k0sArch := "amd64"
	downloadURL := fmt.Sprintf("https://github.com/k0sproject/k0s/releases/download/%s/%s-%s-%s", K0S_VERSION,
		K0sBinary,
		K0S_VERSION,
		k0sArch)

	logrus.Infof("Downloading k0s binary from %s", downloadURL)

	// download file
	resp, err = http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// save it to a file
	err = os.WriteFile(K0sBinaryPath(), body, 0755)
	if err != nil {
		return err
	}

	return nil
}

func downloadAirgapImages() error {
	images, err := getAirgapImages()
	if err != nil {
		return err
	}

	for _, image := range images {
		err = pullAndExportImage(image)
		if err != nil {
			return err
		}
	}

	return nil
}

func pullAndExportImage(image string) error {
	dest := imageToFilePath(image)

	logrus.Infof("pulling and exporting image %s to %s", image, dest)

	// skip when file already exists and is valid
	if imageTarOK(dest) {
		logrus.Warnf("skipping image %s, already exists", image)
		return nil
	}

	pull := cmd.NewCmdPull(&[]crane.Option{})
	pull.SetArgs([]string{image, dest})
	pull.SetOut(os.Stdout)
	pull.SetErr(os.Stderr)
	err := pull.Execute()
	if err != nil {
		logrus.Errorf("failed to pull image %s: %v", image, err)
		return err
	}
	// out, _, err := k0s("ctr", "image", "pull", image)
	// if err != nil {
	// 	logrus.Errorf("failed to pull image %s: %v", image, err)
	// 	return err
	// }
	// logrus.Infof("pulled image %s: %s", image, string(out))
	return nil
}

func imageTarOK(file string) bool {
	// check if file exists
	if _, err := os.Stat(file); err != nil {
		return false
	}

	// check if file is not empty
	fileInfo, err := os.Stat(file)
	if err != nil {
		return false
	}
	if fileInfo.Size() == 0 {
		return false
	}

	// check if file is valid
	validate := cmd.NewCmdValidate(&[]crane.Option{})
	validate.SetArgs([]string{"--tarball", file})
	validate.SetOut(os.Stdout)
	validate.SetErr(os.Stderr)
	err = validate.Execute()
	if err != nil {
		logrus.Errorf("failed to validate image %s: %v", file, err)
		return false
	}

	return true
}

func getAirgapImages() ([]string, error) {
	out, _, err := k0s("airgap", "list-images")
	if err != nil {
		logrus.Errorf("failed to get airgap images: %v", err)
		return nil, err
	}
	images := []string{}
	for _, image := range strings.Split(string(out), "\n") {
		if image != "" {
			images = append(images, image)
		}
	}
	return images, nil
}

func k0s(args ...string) (string, string, error) {
	stdoutWriter := &bytes.Buffer{}
	stderrWriter := &bytes.Buffer{}

	cmd := exec.Command(K0sBinaryPath(), args...)
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("failed to run k0s: %v", err)
		return "", stderrWriter.String(), err
	}

	return stdoutWriter.String(), stderrWriter.String(), nil
}

func mkAssetDir() error {
	path := path.Join(AssetDir, "images")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func OK() bool {
	_, err := os.Stat(AssetDir)
	if err != nil {
		return false
	}

	paths, err := Paths()
	if err != nil {
		logrus.Errorf("failed to get asset paths: %v", err)
		return false
	}

	for _, path := range paths {
		_, err = os.Stat(path)
		if err != nil {
			logrus.Errorf("failed to stat asset path %s: %v", path, err)
			return false
		}
	}

	return true
}

func Paths() ([]string, error) {
	images, err := getAirgapImages()
	if err != nil {
		return nil, err
	}
	paths := []string{
		K0sBinaryPath(),
	}
	for _, image := range images {
		paths = append(paths, imageToFilePath(image))
	}
	return paths, nil
}

func K0sBinaryPath() string {
	return fmt.Sprintf("%s/%s", AssetDir, K0sBinary)
}

func imageToFilePath(image string) string {
	return fmt.Sprintf("%s/images/%s.tar", AssetDir, strings.ReplaceAll(image, "/", "_"))
}

func ImagesPath() string {
	return fmt.Sprintf("%s/images/", AssetDir)
}
