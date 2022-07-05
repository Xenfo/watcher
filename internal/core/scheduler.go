package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/Masterminds/semver"
	"github.com/Xenfo/watcher/internal/config"
	"github.com/fatih/color"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	Config *config.Config
	Logger *zap.Logger
	Cron   *gocron.Scheduler
}

type npmResponse struct {
	Error string            `json:"error"`
	Time  map[string]string `json:"time"`
}

func (s *Scheduler) run() {
	var wg sync.WaitGroup

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	for npmPackageName := range s.Config.Packages {
		wg.Add(1)

		go func(packageName string, packageMeta config.Package) {
			defer wg.Done()

			logger := s.Logger.Named(packageName)

			logger.Info("fetching package")
			resp, err := client.Get(fmt.Sprintf("https://registry.npmjs.org/%s", packageName))
			if err != nil {
				logger.Error("could not get npm package", zap.Error(err))
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error("couldn't read response", zap.Error(err))
				return
			}

			var res *npmResponse
			err = json.Unmarshal(body, &res)
			if err != nil {
				logger.Error("couldn't unmarshal response", zap.Error(err))
				return
			}

			if res.Error != "" {
				logger.Error("npm registry returned error", zap.Error(fmt.Errorf(res.Error)))
				return
			}

			versions := []string{}
			for version := range res.Time {
				if version == "created" || version == "modified" {
					continue
				}

				versions = append(versions, version)
			}

			vs := make([]*semver.Version, len(versions))
			for i, r := range versions {
				v, err := semver.NewVersion(r)
				if err != nil {
					logger.Error("couldn't parse semver", zap.Error(err))
					continue
				}

				vs[i] = v
			}

			sort.Sort(semver.Collection(vs))

			vs = lo.Filter(vs, func(v *semver.Version, _ int) bool {
				if v.Prerelease() != "" && !packageMeta.Betas {
					return false
				}

				return true
			})

			latestVersion := vs[len(vs)-1]
			currentVersion, err := semver.NewVersion(packageMeta.CurrentVersion)
			if err != nil {
				logger.Error("couldn't parse semver", zap.Error(err))
				return
			}

			var targetVersion *semver.Version
			if packageMeta.TargetVersion != "" {
				tVersion, err := semver.NewVersion(packageMeta.TargetVersion)
				if err != nil {
					logger.Error("couldn't parse semver", zap.Error(err))
					return
				}

				targetVersion = tVersion
			}

			if (targetVersion != nil && lo.Contains(versions, targetVersion.String()) && (latestVersion.Equal(targetVersion) || latestVersion.GreaterThan(targetVersion))) || latestVersion.GreaterThan(currentVersion) {
				cyanBold := color.New(color.FgCyan, color.Bold)

				if targetVersion != nil {
					fmt.Printf("\n%s %s@%s -> %s@%s\n", cyanBold.Sprintf("Target version detected!"), packageName, packageMeta.CurrentVersion, packageName, targetVersion.String())
				} else {
					fmt.Printf("\n%s %s@%s -> %s@%s\n", cyanBold.Sprintf("New version detected!"), packageName, packageMeta.CurrentVersion, packageName, latestVersion.String())
				}

				if packageMeta.Notes != "" {
					fmt.Printf("%s: %s", cyanBold.Sprintf("User notes"), packageMeta.Notes)
				}

				fmt.Printf("\n")
			}
		}(npmPackageName, s.Config.Packages[npmPackageName])
	}

	wg.Wait()

	fmt.Printf("\n")
}

func (s *Scheduler) Start() {
	s.Cron.Every(30).Seconds().Do(s.run)
	s.Cron.StartBlocking()
}

func (s *Scheduler) Stop() {
	s.Cron.Stop()
}

func CreateScheduler(config *config.Config, logger *zap.Logger, cron *gocron.Scheduler) *Scheduler {
	return &Scheduler{
		config,
		logger,
		cron,
	}
}
