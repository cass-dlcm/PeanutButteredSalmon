package statink

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/peanutbutteredsalmon/v3/lib"
	"github.com/cass-dlcm/peanutbutteredsalmon/v3/types"
	"log"
	"net/http"
	"os"
	"time"
)

func GetAllShifts(statInkServer types.Server, client *http.Client) (errs []error) {
	var jsonLinesWriter *gzip.Writer
	file, err := os.Create(fmt.Sprintf("statink_shifts/%s_out.jl.gz", statInkServer.ShortName))
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	jsonLinesWriter = gzip.NewWriter(file)
	shift := ShiftStatInk{}
	getShift := func(id int) (data []ShiftStatInk, errs []error) {
		url := fmt.Sprintf("%suser-salmon", statInkServer.Address)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			errs = append(errs, err)
			return nil, errs
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", statInkServer.ApiKey))
		query := req.URL.Query()
		query.Set("newer_than", fmt.Sprint(id))
		query.Set("order", "asc")
		req.URL.RawQuery = query.Encode()
		log.Println(req.URL)
		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, err)
			return nil, errs
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				errs = append(errs, err)
			}
		}()
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			errs = append(errs, err)
			return nil, errs
		}
		for i := range data {
			if _, err := os.Stat("statink_shifts"); errors.Is(err, os.ErrNotExist) {
				if err := os.Mkdir("statink_shifts", os.ModePerm); err != nil {
					errs = append(errs, err)
					return nil, errs
				}
			}
			fileText, err := json.Marshal(data[i])
			if err != nil {
				errs = append(errs, err)
				return nil, errs
			}
			if _, err := jsonLinesWriter.Write(fileText); err != nil {
				errs = append(errs, err)
				return nil, errs
			}
			if _, err := jsonLinesWriter.Write([]byte("\n")); err != nil {
				errs = append(errs, err)
				return nil, errs
			}
		}
		return data, nil
	}
	fileIn, err := os.Open(fmt.Sprintf("statink_shifts/%s.jl.gz", statInkServer.ShortName))
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
	gzipReader, err := gzip.NewReader(fileIn)
	if err != nil {
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := fileIn.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		errs = append(errs, err)
		return errs
	}
	bufioScan := bufio.NewScanner(gzipReader)
	id := 1
	for bufioScan.Scan() {
		if err := json.Unmarshal([]byte(bufioScan.Text()), &shift); err != nil {
			errs = append(errs, err)
			if err := fileIn.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := gzipReader.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return errs
		}
		id = shift.ID
		if _, err := jsonLinesWriter.Write([]byte(bufioScan.Text())); err != nil {
			errs = append(errs, err)
			if err := fileIn.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := gzipReader.Close(); err != nil {
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
			if err := gzipReader.Close(); err != nil {
				errs = append(errs, err)
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
			}
			return errs
		}
	}
	if err := fileIn.Close(); err != nil {
		errs = append(errs, err)
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := gzipReader.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		return errs
	}
	if err := gzipReader.Close(); err != nil {
		errs = append(errs, err)
		if err := jsonLinesWriter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
		return errs
	}
	for true {
		tempData, errs2 := getShift(id)
		if errs2 != nil && len(errs2) > 0 {
			errs = append(errs, errs2...)
			return errs
		}
		if len(tempData) == 0 {
			if err := jsonLinesWriter.Close(); err != nil {
				errs = append(errs, err)
				if err := file.Close(); err != nil {
					errs = append(errs, err)
				}
				return errs
			}
			if err := file.Close(); err != nil {
				errs = append(errs, err)
				return errs
			}
			if err := os.Remove(fmt.Sprintf("statink_shifts/%s.jl.gz", statInkServer.ShortName)); err != nil {
				errs = append(errs, err)
				return errs
			}
			if err := os.Rename(fmt.Sprintf("statink_shifts/%s_out.jl.gz", statInkServer.ShortName), fmt.Sprintf("statink_shifts/%s.jl.gz", statInkServer.ShortName)); err != nil {
				errs = append(errs, err)
				return errs
			}
			return nil
		}
		id = tempData[len(tempData)-1].ID
	}
	return nil
}

type ShiftStatInkIterator struct {
	serverName string
	f 		   *os.File
	buffRead   *bufio.Scanner
	gzipReader *gzip.Reader
}

func (s *ShiftStatInkIterator) Next() (shift lib.Shift, errs []error) {
	data := ShiftStatInk{}
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
		return &data, nil
	}
	if err := s.gzipReader.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.f.Close(); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, errors.New("no more shifts"))
	return nil, errs
}

func LoadFromFileIterator(server types.Server) (lib.ShiftIterator, error) {
	returnVal := ShiftStatInkIterator{serverName: server.ShortName}
	var err error
	returnVal.f, err = os.Open(fmt.Sprintf("statink_shifts/%s.jl.gz", server.ShortName))
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
