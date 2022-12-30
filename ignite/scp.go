package ignite

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mheers/k3i/config"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
)

func Copy(vm *ignite.VM, src, dst string) error {
	config := config.GetConfig()

	err := copy(vm, config.SSHPrivateKeyFile, src, dst, constants.SSH_DEFAULT_TIMEOUT_SECONDS)

	if err != nil && err.Error() != "EOF" {
		return err
	}

	return nil

}

func copy(vm *ignite.VM, privateKeyFile, src, dst string, timeout uint32) error {
	logrus.Infof("copying file: %s to %s", src, dst)

	client, err := newSSHClient(vm, privateKeyFile, timeout)
	if err != nil {
		return err
	}

	// Use sftp to copy file from source to destination.
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create new sftp client: %v", err)
	}
	defer sftpClient.Close()

	// Clean the source and destination paths.
	source := filepath.Clean(src)
	dest := filepath.Clean(dst)

	// Determine the copy direction.
	source, dest, copyDirection, err := copyDirection(source, dest)
	if err != nil {
		return err
	}

	// Copy files based on the copy direction.
	switch copyDirection {
	case run.CopyDirectionHostToVM:
		if err := copyToVM(sftpClient, source, dest); err != nil {
			return fmt.Errorf("failed to copy files from host to VM: %v", err)
		}
	case run.CopyDirectionVMToHost:
		if err := copyFromVM(sftpClient, source, dest); err != nil {
			return fmt.Errorf("failed to copy files from VM to host: %v", err)
		}
	}

	return nil
}

func copyDirection(source, dest string) (string, string, run.CopyDirection, error) {
	var vmMatch string

	sourceComponents := strings.Split(source, run.VMFilePathSeparator)
	destComponents := strings.Split(dest, run.VMFilePathSeparator)

	copyDirection := run.CopyDirectionUnknown

	if len(sourceComponents) > 1 {
		copyDirection = run.CopyDirectionVMToHost
		vmMatch = sourceComponents[0]
		source = sourceComponents[1]
	} else {
		source = sourceComponents[0]
	}

	if len(destComponents) > 1 {
		if vmMatch != "" {
			return "", "", copyDirection, fmt.Errorf("only one of source or destination can have VM name/ID")
		}
		copyDirection = run.CopyDirectionHostToVM
		vmMatch = destComponents[0]
		dest = destComponents[1]
	} else {
		dest = destComponents[0]
	}

	// If no copy direction if known due to no VM reference in the source or
	// destination, fail.
	if copyDirection == run.CopyDirectionUnknown {
		return "", "", copyDirection, fmt.Errorf("no VM reference found in source or destination")
	}
	return source, dest, copyDirection, nil
}

// copyToVM copies from host to VM.
func copyToVM(client *sftp.Client, localPath, remotePath string) error {
	// Check if the source exists.
	fi, err := os.Stat(localPath)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return copyFileToVM(client, localPath, remotePath)
	}
	return copyDirToVM(client, localPath, remotePath)
}

// copyFileToVM copies file from host to VM.
func copyFileToVM(client *sftp.Client, localPath, remotePath string) error {
	in, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer in.Close()

	// Check if remote path already exists. If the destination is a directory,
	// the source must be copied into the directory.
	if existsInVM(client, remotePath) {
		// Check if the remote path is a directory.
		rfi, err := client.Stat(remotePath)
		if err != nil {
			return err
		}

		// If the remote destination is a directory, update the remotePath to be
		// moved into the destination directory.
		// For example: if /tmp/foo.txt is copied to /xyz/, the remote should be
		// updated to /xyz/foo.txt, append the filepath base to remote path.
		if rfi.IsDir() {
			remotePath = filepath.Join(remotePath, filepath.Base(localPath))
		}
		// Else, any existing file will be overwritten with the new file.
	}

	out, err := client.Create(remotePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy source to destination.
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Read and apply source file modes and owner info to destination file.
	sfi, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	stat, ok := sfi.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("failed to get raw syscall.Stat_t data for %q", localPath)
	}

	if err := client.Chmod(remotePath, sfi.Mode()); err != nil {
		return err
	}

	if err := client.Chown(remotePath, int(stat.Uid), int(stat.Gid)); err != nil {
		return err
	}

	return nil
}

// copyDirToVM copies directory from host to VM.
func copyDirToVM(client *sftp.Client, localPath, remotePath string) error {
	// Check if remote destination path already exists. If the destination is a
	// directory with a different name, the source must be copied into the
	// directory. If the destination is a file, copying should fail.
	if existsInVM(client, remotePath) && (filepath.Base(localPath) != filepath.Base(remotePath)) {
		isDir, err := isDirInVM(client, remotePath)
		if err != nil {
			return err
		}

		// If the remote destination is a directory, update the remotePath to be
		// moved into the destination directory.
		if isDir {
			remotePath = filepath.Join(remotePath, filepath.Base(localPath))
		} else {
			// Copying directory info file should fail.
			return fmt.Errorf("cannot overwrite non-directory %q(VM) with directory %q(Host)", remotePath, localPath)
		}

		// If the new subdirectory path exists, ensure that it is a directory,
		// and not a file. Copying directory to file should fail.
		if existsInVM(client, remotePath) {
			isDir, err = isDirInVM(client, remotePath)
			if err != nil {
				return err
			}
			if !isDir {
				return fmt.Errorf("cannot overwrite non-directory %q(VM) with directory %q(Host)", remotePath, localPath)
			}
		}
	}

	// Get the source directory fileinfo.
	dInfo, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	// Ensure destination parent dir exists.
	if err := createIfNotExistsInVM(client, remotePath, dInfo); err != nil {
		return err
	}

	// Remove local source directory.
	entries, err := ioutil.ReadDir(localPath)
	if err != nil {
		return err
	}

	// Go through all the items in the directory and copy them to VM.
	for _, entry := range entries {
		lPath := filepath.Join(localPath, entry.Name())
		rPath := filepath.Join(remotePath, entry.Name())

		fileInfo, err := os.Stat(lPath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExistsInVM(client, rPath, fileInfo); err != nil {
				return err
			}
			if err := copyDirToVM(client, lPath, rPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := copySymLinkToVM(client, lPath, rPath); err != nil {
				return err
			}
		default:
			if err := copyFileToVM(client, lPath, rPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// existsInVM checks if a filepath exists in the VM.
func existsInVM(client *sftp.Client, filePath string) bool {
	if _, err := client.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// createIfNotExistsInVM creates dir if it doesn't exists in VM with the given
// source file info.
func createIfNotExistsInVM(client *sftp.Client, dir string, srcInfo os.FileInfo) error {
	if existsInVM(client, dir) {
		return nil
	}

	if err := client.MkdirAll(dir); err != nil {
		return fmt.Errorf("failed to create directory: %q, error: %q", dir, err.Error())
	}

	// Get source owner and permission info and set on destination.
	stat, ok := srcInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("failed to get raw syscall.Stat_t data for %q", srcInfo.Name())
	}

	if err := client.Chmod(dir, srcInfo.Mode()); err != nil {
		return err
	}

	if err := client.Chown(dir, int(stat.Uid), int(stat.Gid)); err != nil {
		return err
	}

	return nil
}

// copySymLinkToVM reads the symlink destination and copies that to VM.
func copySymLinkToVM(client *sftp.Client, localPath, remotePath string) error {
	link, err := os.Readlink(localPath)
	if err != nil {
		return err
	}

	// Check if the file is a directory or a file and copy accordingly.
	fi, err := os.Stat(link)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return copyDirToVM(client, link, remotePath)
	}
	return copyFileToVM(client, link, remotePath)
}

// copyFromVM copies from VM to host.
func copyFromVM(client *sftp.Client, remotePath, localPath string) error {
	// Check if the source exists.
	fi, err := client.Stat(remotePath)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return copyFileFromVM(client, remotePath, localPath)
	}
	return copyDirFromVM(client, remotePath, localPath)
}

// copyFileFromVM copies file from VM to host.
func copyFileFromVM(client *sftp.Client, remotePath, localPath string) error {
	in, err := client.Open(remotePath)
	if err != nil {
		return err
	}
	defer in.Close()

	// Check if local path already exists. If the destination is a directory,
	// the source must be copied into the directory.
	if existsInHost(localPath) {
		// Check if the local path is a directory.
		lfi, err := os.Stat(localPath)
		if err != nil {
			return err
		}

		// If the local destination is a directory, update the localPath to be
		// moved into the destination directory.
		// For example: if /tmp/foo.txt is copied to /xyz/, the local should be
		// updated to /xyz/foo.txt, append the filepath base to local path.
		if lfi.IsDir() {
			localPath = filepath.Join(localPath, filepath.Base(remotePath))
		}
		// Else, any existing file will be overwritten with the new file.
	}

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy source to destination.
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Read and apply source file modes to destination file.
	// Querying Stat_t from remote FileInfo for setting ownership fails.
	sfi, err := client.Stat(remotePath)
	if err != nil {
		return err
	}

	if err := os.Chmod(localPath, sfi.Mode()); err != nil {
		return err
	}

	return nil
}

// copyDirFromVM copies directory from VM to host.
func copyDirFromVM(client *sftp.Client, remotePath, localPath string) error {
	// Check if local destination path already exists. If the destination is a
	// directory with a different name, the source muse be copied into the
	// directory. If the destination is a file, copying should fail.
	if existsInHost(localPath) && (filepath.Base(remotePath) != filepath.Base(localPath)) {
		isDir, err := isDirInHost(localPath)
		if err != nil {
			return err
		}

		// If the local destination is a directory, update the localPath to be
		// moved into the destination directory.
		if isDir {
			localPath = filepath.Join(localPath, filepath.Base(remotePath))
		} else {
			// Copying directory info file should fail.
			return fmt.Errorf("cannot overwrite non-directory %q(Host) with directory %q(VM)", localPath, remotePath)
		}

		// If the new subdirectory path exists, ensure that it is a directory,
		// and not a file. Copying directory to file should fail.
		if existsInHost(localPath) {
			isDir, err = isDirInHost(localPath)
			if err != nil {
				return err
			}
			if !isDir {
				return fmt.Errorf("cannot overwrite non-directory %q(Host) with directory %q(VM)", localPath, remotePath)
			}
		}
	}

	// Get the source directory fileinfo.
	dInfo, err := client.Stat(remotePath)
	if err != nil {
		return err
	}
	// Ensure destination parent dir exists.
	if err := createIfNotExistsInHost(localPath, dInfo); err != nil {
		return err
	}

	// Read remote source directory.
	entries, err := client.ReadDir(remotePath)
	if err != nil {
		return err
	}

	// Go through all the items in the directory and copy them to host.
	for _, entry := range entries {
		rPath := filepath.Join(remotePath, entry.Name())
		lPath := filepath.Join(localPath, entry.Name())

		fileInfo, err := client.Stat(rPath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExistsInHost(lPath, fileInfo); err != nil {
				return err
			}
			if err := copyDirFromVM(client, rPath, lPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := copySymLinkFromVM(client, rPath, lPath); err != nil {
				return err
			}
		default:
			if err := copyFileFromVM(client, rPath, lPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// existsInHost checks if a filepath exists in the host.
func existsInHost(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// createIfNotExistsInHost creates dir if it doesn't exists in host with the
// given source file info.
func createIfNotExistsInHost(dir string, srcInfo os.FileInfo) error {
	if existsInHost(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create dir: %q, error: %q", dir, err.Error())
	}
	return nil
}

// copySymLinkFromVM reads the symlink destination and copies that to host.
func copySymLinkFromVM(client *sftp.Client, remotePath, localPath string) error {
	link, err := client.ReadLink(remotePath)
	if err != nil {
		return err
	}

	// Check if the file is a directory or a file and copy accordingly.
	fi, err := client.Stat(link)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return copyDirFromVM(client, link, localPath)
	}
	return copyFileFromVM(client, link, localPath)
}

// isDirInVM checks if a given path in VM is a directory.
func isDirInVM(client *sftp.Client, path string) (bool, error) {
	fi, err := client.Stat(path)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		return true, nil
	}
	return false, nil
}

// isDirInHost checks if a given path in host is a directory.
func isDirInHost(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		return true, nil
	}
	return false, nil
}
