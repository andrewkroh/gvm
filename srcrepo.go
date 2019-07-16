package gvm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type srcCacheInfo struct {
	Updated time.Time
}

func (m *Manager) srcCacheDir() string {
	return filepath.Join(m.cacheDir, "go")
}

func (m *Manager) hasSrcCache() bool {
	localGoSrc := m.srcCacheDir()
	exists, err := existsDir(localGoSrc)
	return err == nil && exists
}

func (m *Manager) ensureSrcCache() error {
	localGoSrc := m.srcCacheDir()
	exists, err := existsDir(localGoSrc)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return m.updateSrcCache()
}

func (m *Manager) updateSrcCache() error {
	localGoSrc := filepath.Join(m.cacheDir, "go")
	exists, err := existsDir(localGoSrc)
	if err != nil {
		return err
	}

	if !exists {
		err = gitClone(m.Logger, localGoSrc, m.GoSourceURL, false)
	} else {
		err = gitPull(m.Logger, localGoSrc)
	}
	if err != nil {
		return err
	}

	return writeJsonFile(filepath.Join(m.cacheDir, "go.meta"), srcCacheInfo{
		Updated: time.Now(),
	})
}

func (m *Manager) installSrc(version *GoVersion) (string, error) {
	log := m.Logger

	if err := m.ensureSrcCache(); err != nil {
		return "", err
	}

	tag := "master"
	to := m.VersionGoROOT(version)

	exists, err := existsDir(to)
	if err != nil {
		return "", err
	}

	if !version.IsTip() {
		tag = fmt.Sprintf("go%v", version)
		if err := m.ensureSrcVersionAvail(version); err != nil {
			return "", err
		}
	}

	godir := m.versionDir(version)

	log.Println("create temp directory")
	tmpRoot, err := ioutil.TempDir("", godir)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpRoot)

	err = buildGo(log, tmpRoot, m.srcCacheDir(), version.String(), tag)
	if err != nil {
		return "", err
	}

	if exists {
		log.Println("remove old installation")
		if err := os.RemoveAll(to); err != nil {
			return "", err
		}
	}

	// move final build result into destination
	log.Println("rename")
	if err = os.Rename(filepath.Join(tmpRoot, "go"), to); err != nil {
		return "", err
	}
	return to, nil
}

func buildGo(log logrus.FieldLogger, buildDir, repo, version, tag string) error {
	log.Println("copy cache")
	tmp := filepath.Join(buildDir, "go")
	if err := gitClone(log, tmp, repo, false); err != nil {
		return err
	}
	log.Println("checkout tag:", tag)
	if err := gitCheckout(log, tmp, tag); err != nil {
		return err
	}

	bootstrap := os.Getenv("GOROOT_BOOTSTRAP")
	if bootstrap == "" {
		bootstrap = os.Getenv("GOROOT")
		if bootstrap == "" {
			return errors.New("GOROOT or GOROOT_BOOTSTRAP must be set")
		}
	}

	if version != "tip" {
		// write VERSION file
		versionFile := filepath.Join(tmp, "VERSION")
		err := ioutil.WriteFile(versionFile, []byte(version), 0644)
		if err != nil {
			return err
		}
	}

	if err := os.RemoveAll(filepath.Join(tmp, "go", ".git")); err != nil {
		return err
	}

	log.Println("build")
	srcDir := filepath.Join(tmp, "src")

	var cmd *command
	if runtime.GOOS == "windows" {
		cmd = makeCommand("cmd", "/C", "make.bat")
	} else {
		cmd = makeCommand("bash", "make.bash")
	}
	cmd.Env = []string{
		"GOROOT_BOOTSTRAP=" + bootstrap,
	}
	return cmd.WithDir(srcDir).WithLogger(log).Exec()
}

func (m *Manager) hasSrcVersion(version *GoVersion) (bool, error) {
	log := m.Logger
	localGoSrc := filepath.Join(m.cacheDir, "go")

	tag := fmt.Sprintf("go%s", version)
	log.Println("check version tag")
	hasTag := false
	err := gitListTags(log, localGoSrc, func(t string) { hasTag = hasTag || t == tag })
	return hasTag, err
}

func (m *Manager) ensureSrcVersionAvail(version *GoVersion) error {
	has, err := m.hasSrcVersion(version)
	if err != nil {
		return err
	}

	if !has {
		if err := m.updateSrcCache(); err != nil {
			return err
		}
		has, err = m.hasSrcVersion(version)
		if err != nil {
			return err
		}
	}
	if !has {
		return fmt.Errorf("unknown version %s", version)
	}
	return nil
}

func (m *Manager) tryRefreshSrcCache() (bool, error) {
	log := m.Logger

	localGoSrc := m.srcCacheDir()
	exists, err := existsDir(localGoSrc)
	if err != nil {
		return false, err
	}
	if !exists {
		log.Println("Go cache not found")
		err := m.updateSrcCache()
		return err == nil, err
	}

	info := srcCacheInfo{}
	if err := readJsonFile(filepath.Join(m.cacheDir, "go.meta"), &info); err != nil {
		return false, err
	}

	updTS := info.Updated
	now := time.Now()
	if now.Before(updTS) {
		return false, nil
	}

	// don't refresh cache if still same day
	if updTS.Day() == now.Day() && updTS.Month() == now.Month() && updTS.Year() == updTS.Year() {
		return false, nil
	}

	log.Println("Fetch updates")
	if err := m.updateSrcCache(); err != nil {
		return false, err // update cache failed
	}

	// check for updates ;)
	cTS, err := gitLastCommitTs(log, m.srcCacheDir())
	if err != nil {
		return false, err
	}

	log.Printf("last update ts=%v, last commit ts=%v\n", updTS, cTS)

	// check new commits have been added since last refresh
	updates := updTS.Before(cTS)
	if updates {
		log.Println("New commits since last build")
	} else {
		log.Println("No New commits since last build")
	}

	return updates, nil
}

func (m *Manager) AvailableSource() ([]*GoVersion, error) {
	localGoSrc := m.srcCacheDir()
	var versions []*GoVersion
	err := gitListTags(m.Logger, localGoSrc, func(tag string) {
		if !strings.HasPrefix(tag, "go") {
			return
		}

		ver, err := ParseVersion(tag[2:])
		if err != nil {
			return
		}

		versions = append(versions, ver)
	})

	if err != nil {
		return nil, err
	}

	tip, _ := ParseVersion("tip")
	versions = append(versions, tip)
	sortVersions(versions)
	return versions, err
}

func gitClone(logger logrus.FieldLogger, to string, url string, bare bool) error {
	tmpDir := to + ".tmp"
	if err := os.Mkdir(tmpDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	args := []string{"clone"}
	if bare {
		args = append(args, "--bare")
	}
	args = append(args, url, tmpDir)

	logger.Println("git clone:")
	cmd := makeCommand("git", args...).WithLogger(logger)
	if err := cmd.Exec(); err != nil {
		return err
	}

	// Move into the final location.
	return os.Rename(tmpDir, to)
}

func gitLastCommitTs(logger logrus.FieldLogger, path string) (time.Time, error) {
	var tsLine string

	logger.Println("git log:")
	cmd := makeCommand("git", "log", "-n", "1", "--pretty=format:%ct")
	cmd.Stdout = func(l string) { tsLine = l }
	err := cmd.WithDir(path).WithLogger(logger).Exec()
	if err != nil {
		return time.Time{}, err
	}

	i, err := strconv.ParseInt(tsLine, 10, 64)
	return time.Unix(i, 0), nil
}

func gitFetch(logger logrus.FieldLogger, path string) error {
	logger.Println("git fetch:")
	return makeCommand("git", "fetch").WithDir(path).WithLogger(logger).Exec()
}

func gitPull(logger logrus.FieldLogger, path string) error {
	logger.Println("git pull:")
	return makeCommand("git", "pull").WithDir(path).WithLogger(logger).Exec()
}

func gitCheckout(logger logrus.FieldLogger, path, tag string) error {
	logger.Println("git checkout:")
	return makeCommand("git", "checkout", tag).WithDir(path).WithLogger(logger).Exec()
}

func gitListTags(logger logrus.FieldLogger, path string, fn func(string)) error {
	logger.Println("git tag:")
	cmd := makeCommand("git", "tag").WithDir(path).WithLogger(logger)
	cmd.Stdout = fn
	return cmd.Exec()
}
