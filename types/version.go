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

// Version is the current version of the software.
var Version = version{1, 0, 4}

type version struct {
	Major  uint64
	Minor  uint64
	Bugfix uint64
}

// CompareVersion performs a comparison to see which version is higher.
// It returns 1 if the object is higher, -1 if the parameter is higher, and 0 if they are the same.
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

// ToString returns a nicely formatted string of the version.
func (v *version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Bugfix)
}

// CheckForUpdate downloads the latest version of this file and checks to see which version number is higher.
// If the downloaded version is higher, execution ends.
func CheckForUpdate(client *http.Client) {
	url := "https://raw.githubusercontent.com/cass-dlcm/peanutbutteredsalmon/main/types/version.go"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Panicln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	byteBody, _ := ioutil.ReadAll(resp.Body)
	search, err := regexp.Compile("(\\d+, \\d+, \\d+)")
	if err != nil {
		log.Panicln(err)
	}
	versionStr := search.FindString(string(byteBody))
	versionSubstrs := strings.Split(versionStr, ", ")
	major, err := strconv.ParseUint(versionSubstrs[0], 10, 32)
	if err != nil {
		log.Println(err)
	}
	minor, err := strconv.ParseUint(versionSubstrs[1], 10, 32)
	if err != nil {
		log.Println(err)
	}
	bugfix, err := strconv.ParseUint(versionSubstrs[2], 10, 32)
	if err != nil {
		log.Println(err)
	}
	testVers := version{major, minor, bugfix}
	if Version.CompareVersion(&testVers) == 1 {
		log.Panicf("A new version is available. Please update to the new version.\nCurrent Version: %s\nNew Version: %s\nExiting.", Version.ToString(), testVers.ToString())
	}
}
