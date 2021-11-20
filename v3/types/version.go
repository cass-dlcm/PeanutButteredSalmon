package types

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Version = version{3, 0, 0}

type version struct {
	Major  uint64
	Minor  uint64
	Bugfix uint64
}

func (v *version) CompareVersion(v2 *version) int {
	if v.Major > v2.Major {
		return -1
	}
	if v.Major < v2.Major {
		return 1
	}
	if v.Minor > v2.Minor {
		return -1
	}
	if v.Minor < v2.Minor {
		return 1
	}
	if v.Bugfix > v2.Bugfix {
		return -1
	}
	if v.Bugfix < v2.Bugfix {
		return 1
	}
	return 0
}

func (v *version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Bugfix)
}

func CheckForUpdate(client *http.Client) (errs []error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/cass-dlcm/PeanutButteredSalmon/main/v%d/types/version.go", Version.Major + 1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	resp, err := client.Do(req)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			errs = append(errs, err)
		}
	}()
	if resp.StatusCode == 404 {
		url := fmt.Sprintf("https://raw.githubusercontent.com/cass-dlcm/PeanutButteredSalmon/main/v%d/types/version.go", Version.Major)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		resp, err = client.Do(req)
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				errs = append(errs, err)
			}
		}()
		if resp.StatusCode == 404 {
			url := fmt.Sprintf("https://raw.githubusercontent.com/cass-dlcm/PeanutButteredSalmon/main/v%d/types/version.go", Version.Major - 1)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				errs = append(errs, err)
				return errs
			}
			resp, err = client.Do(req)
			if err != nil {
				errs = append(errs, err)
				return errs
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					errs = append(errs, err)
				}
			}()
		}
	}
	byteBody, _ := ioutil.ReadAll(resp.Body)
	search, err := regexp.Compile("(\\d+, \\d+, \\d+)")
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	versionStr := search.FindString(string(byteBody))
	versionSubstrs := strings.Split(versionStr, ", ")
	major, err := strconv.ParseUint(versionSubstrs[0], 10, 32)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	minor, err := strconv.ParseUint(versionSubstrs[1], 10, 32)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	bugfix, err := strconv.ParseUint(versionSubstrs[2], 10, 32)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	testVers := version{major, minor, bugfix}
	versionComparison := Version.CompareVersion(&testVers)
	if versionComparison == 1 {
		log.Panicf("A new version is available. Please update to the new version.\nCurrent Version: %s\nNew Version: %s\nExiting.", Version.ToString(), testVers.ToString())
	} else if versionComparison == -1 {
		log.Printf("You are running a unreleased version.\nLatest released version:%s\nCurrent version:%s\n", testVers.ToString(), Version.ToString())
	}
	return errs
}
