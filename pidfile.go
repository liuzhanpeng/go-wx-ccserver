package wxccserver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"golang.org/x/sys/unix"
)

// PIDFile pidfile锁
type PIDFile struct {
	filename string
}

// NewPIDFile 创建pidfile对象
func NewPIDFile(filename string) *PIDFile {
	return &PIDFile{
		filename: filename,
	}
}

// GetPID 获取pidfile文件保存的进行id
func (p *PIDFile) GetPID() (int, error) {
	c, err := ioutil.ReadFile(p.filename)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(string(c))
	if err != nil {
		return 0, err
	}

	if runtime.GOOS == "darwin" {
		err := unix.Kill(pid, 0)
		if err != nil {
			return 0, err
		}
	}

	if _, err := os.Stat(filepath.Join("/proc", strconv.Itoa(pid))); err == nil {
		return 0, err
	}

	return pid, nil
}

// SetPID 设置pidfile文件保存的进程id
func (p *PIDFile) SetPID(pid int) error {
	return ioutil.WriteFile(p.filename, []byte(strconv.Itoa(pid)), 0600)
}

// Remove 移除
func (p *PIDFile) Remove() error {
	return os.Remove(p.filename)
}
