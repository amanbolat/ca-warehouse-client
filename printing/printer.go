package printing

import (
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/pkg/errors"
	"os/exec"
	"strconv"
	"strings"
)

const MAX_COPIES = 3
const DEFAULT_MEDIA = "Custom.4x6in"

type Printer struct {
	Name string
}

func (p Printer) PrintFiles(copies int, media string, paths ...string) error {
	if copies > MAX_COPIES {
		return errors.New(fmt.Sprintf("can print no more than %d copies", MAX_COPIES))
	}

	if copies == 0 {
		copies = 1
	}

	if media == "" {
		media = DEFAULT_MEDIA
	}

	if strings.TrimSpace(p.Name) == "" {
		return api.NewError(nil, "打印机名称不能为空", "请联系管理员")
	}
	printCmd := exec.Command("lpr", "-P", p.Name, "-#", strconv.Itoa(copies), "-o", fmt.Sprintf("media=%s", media), strings.Join(paths, " "))
	out, err := printCmd.CombinedOutput()
	if err != nil {
		return api.NewError(errors.WithMessage(err, string(out)), "打印遇到错误", "请确保打印机已开机/已连接电脑/耗材足够")
	}

	return nil
}
