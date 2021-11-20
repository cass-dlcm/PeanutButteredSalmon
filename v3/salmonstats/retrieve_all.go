package salmonstats

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/peanutbutteredsalmon/v3/lib"
	"github.com/cass-dlcm/peanutbutteredsalmon/v3/types"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"time"
)

func GetAllShifts(server types.Server, client *http.Client) (errs []error) {
	if _, err := fmt.Println("Pulling Salmon Run data from online..."); err != nil {
		panic(err)
	}
	var jsonLinesWriter *gzip.Writer
	file, err := os.Create(fmt.Sprintf("salmonstats_shifts/%s_out.jl.gz", server.ShortName))
	if err != nil {
		return []error{err}
	}
	jsonLinesWriter = gzip.NewWriter(file)
	getShifts := func(page int) (found bool, errs []error) {
		url := fmt.Sprintf("%splayers/%s/results", server.Address, viper.GetString("user_id"))
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return  false, []error{err}
		}
		query := req.URL.Query()
		query.Set("raw", "1")
		query.Set("count", "200")
		query.Set("page", fmt.Sprint(page))
		req.URL.RawQuery = query.Encode()

		log.Println(req.URL)

		resp, err := client.Do(req)
		if err != nil {
			return false, []error{err}
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				errs = append(errs, err)
			}
		}()
		var data ShiftPage
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			errs = append(errs, err)
			return false, errs
		}

		for i := range data.Results {
			if _, err := os.Stat("salmonstats_shifts"); errors.Is(err, os.ErrNotExist) {
				err := os.Mkdir("salmonstats_shifts", os.ModePerm)
				if err != nil {
					errs = append(errs, err)
					return false, errs
				}
			}
			fileText, err := json.Marshal(data.Results[i])
			if err != nil {
				errs = append(errs, err)
				return false, errs
			}
			if _, err := jsonLinesWriter.Write(fileText); err != nil {
				errs = append(errs, err)
				return false, errs
			}
			if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
				errs = append(errs, err)
				return false, errs
			}
		}
		return len(data.Results) > 0, nil
	}
	f, err := os.Open(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName))
	if err != nil {
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		errs = append(errs, err)
		return errs
	}
	count := 0
	gzRead, err := gzip.NewReader(f)
	if err != nil {
		errs = append(errs, err)
		if err := f.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		return errs
	}
	bufScan := bufio.NewScanner(gzRead)
	for bufScan.Scan() {
		count++
		if _, err := jsonLinesWriter.Write([]byte(bufScan.Text())); err != nil {
			errs = append(errs, err)
			if err := f.Close(); err != nil {
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
			if err := f.Close(); err != nil {
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
	if err := f.Close(); err != nil {
		errs = append(errs, err)
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
	if err := gzRead.Close(); err != nil {
		errs = append(errs, err)
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		return errs
	}
	page := count / 200 + 1
	hasPages := true
	for hasPages {
		hasPages, errs = getShifts(page)
		if len(errs) > 0 {
			return errs
		}
		page++
	}
	if err := jsonLinesWriter.Close(); err != nil {
		errs = append(errs, err)
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		return errs
	}
	if err := file.Close(); err != nil {
		return []error{err}
	}
	if err := os.Remove(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName)); err != nil {
		return []error{err}
	}
	if err := os.Rename(fmt.Sprintf("salmonstats_shifts/%s_out.jl.gz", server.ShortName), fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName)); err != nil {
		return []error{err}
	}
	return nil
}

type ShiftSalmonStatsIterator struct {
	serverName string
	f 		   *os.File
	buffRead   *bufio.Scanner
	gzipReader *gzip.Reader
}

func (s *ShiftSalmonStatsIterator) Next() (shift lib.Shift, errs []error) {
	data := ShiftSalmonStats{}
	if s.buffRead.Scan() {
		if err := json.Unmarshal([]byte(s.buffRead.Text()), &data); err != nil {
			errs = append(errs, err)
			if err := s.f.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := s.gzipReader.Close(); err != nil {
				errs = append(errs, err)
			}
			return nil, errs
		}
		return data, nil
	}
	if err := s.f.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.gzipReader.Close(); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, errors.New("no more shifts"))
	return nil, errs
}

func LoadFromFileIterator(server types.Server) (iter lib.ShiftIterator, err error) {
	returnVal := ShiftSalmonStatsIterator{serverName: server.ShortName}
	returnVal.f, err = os.Open(fmt.Sprintf("salmonstats_shifts/%s.jl.gz", server.ShortName))
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
