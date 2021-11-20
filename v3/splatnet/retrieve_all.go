package splatnet

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/peanutbutteredsalmon/v3/lib"
	"github.com/cass-dlcm/peanutbutteredsalmon/v3/splatnet/iksm"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

func GetAllShifts(client *http.Client) (errs []error) {
	_, timezone := time.Now().Zone()
	timezone = -timezone / 60
	appHead := http.Header{
		"Host":              []string{"app.splatoon2.nintendo.net"},
		"x-unique-id":       []string{"32449507786579989235"},
		"x-requested-with":  []string{"XMLHttpRequest"},
		"x-timezone-offset": []string{fmt.Sprint(timezone)},
		"User-Agent":        []string{"Mozilla/5.0 (Linux; Android 7.1.2; Pixel Build/NJH47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36"},
		"Accept":            []string{"*/*"},
		"Referer":           []string{"https://app.splatoon2.nintendo.net/home"},
		"Accept-Encoding":   []string{"gzip deflate"},
		"Accept-Language":   []string{viper.GetString("user_lang")},
	}

	if _, err := fmt.Println("Pulling Salmon Run data from online..."); err != nil {
		return []error{err}
	}

	url := "https://app.splatoon2.nintendo.net/api/coop_results"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []error{err}
	}

	req.Header = appHead

	if viper.GetString("cookie") == "" {
		iksm.GenNewCookie("blank", "1.6.0", client)
	}

	req.AddCookie(&http.Cookie{Name: "iksm_session", Value: viper.GetString("cookie")})

	resp, err := client.Do(req)
	if err != nil {
		return []error{err}
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			errs = append(errs, err)
		}
	}()

	var data ShiftList

	var jsonLinesWriter *gzip.Writer
	fileIn, err := os.Open("shifts.jl.gz")
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	gzRead, err := gzip.NewReader(fileIn)
	if err != nil {
		errs = append(errs, err)
		if err := fileIn.Close(); err != nil {
			errs = append(errs, err)
		}
		return errs
	}
	bufScan := bufio.NewScanner(gzRead)
	file, err := os.Create("shifts_out.jl.gz")
	if err != nil {
		return []error{err}
	}
	jsonLinesWriter = gzip.NewWriter(file)

	for bufScan.Scan() {
		if _, err := jsonLinesWriter.Write([]byte(bufScan.Text())); err != nil {
			errs = append(errs, err)
			if err := fileIn.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := gzRead.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return errs
		}
		if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
			errs = append(errs, err)
			if err := fileIn.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := gzRead.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return errs
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		iksm.GenNewCookie("auth", "1.6.0", client)
		if err := fileIn.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := gzRead.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			return errs
		}
		if errsRec := GetAllShifts(client); len(errsRec) > 0 {
			errs = append(errs, errsRec...)
			return errs
		}
	}

	if data.Code != nil {
		iksm.GenNewCookie("auth", "1.6.0", client)
		if err := fileIn.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := gzRead.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			return errs
		}
		if errsRec := GetAllShifts(client); len(errsRec) > 0 {
			errs = append(errs, errsRec...)
			return errs
		}
	}

	if err := fileIn.Close(); err != nil {
		errs = append(errs, err)
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if err := gzRead.Close(); err != nil {
		errs = append(errs, err)
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}

	for i := range data.Results {
		fileText, err := json.Marshal(data.Results[i])
		if err != nil {
			errs = append(errs, err)
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return errs
		}
		if _, err := jsonLinesWriter.Write(fileText); err != nil {
			errs = append(errs, err)
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return errs
		}
		if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
			errs = append(errs, err)
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return errs
		}
	}

	if err := jsonLinesWriter.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := file.Close(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errs
	}

	if err := os.Remove("shifts.jl.gz"); err != nil {
		return []error{err}
	}
	if err := os.Rename("shifts_out.jl.gz", "shifts.jl.gz"); err != nil {
		return []error{err}
	}
	return nil
}

func LoadFromFileIterator() (lib.ShiftIterator, error) {
	returnVal := ShiftSplatnetIterator{}
	var err error
	returnVal.f, err = os.Open("shifts.jl.gz")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	returnVal.gzipReader, err = gzip.NewReader(returnVal.f)
	if err != nil {
		return nil, err
	}
	returnVal.buffRead = bufio.NewScanner(returnVal.gzipReader)
	return &returnVal, nil
}
