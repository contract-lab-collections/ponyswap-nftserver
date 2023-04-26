package utils

import (
	"crypto/md5"
	"fmt"
	"math/big"
	"time"
)

func CalPage(count int64, page, size int) (int, int) {
	if size <= 0 {
		size = 10
	}
	if page < 1 {
		page = 1
	}

	total := (int(count) + size - 1) / size
	if page > total {
		page = total
	}

	offset := (page - 1) * size

	return size, offset
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GetNowHours() time.Time {
	now := time.Now()
	timestamp := now.Unix() - int64(now.Second()) - int64(60*now.Minute())
	return time.Unix(timestamp, 0)
}

func GetNowDay() time.Time {
	now := time.Now()
	timestamp := now.Unix() - int64(now.Second()) - int64(60*now.Minute()) - int64(3600*now.Hour())
	return time.Unix(timestamp, 0)
}

func String2Bigint(s string) (*big.Int, error) {
	i := big.NewInt(0)
	i, isOK := i.SetString(s, 10)
	if !isOK {
		return nil, fmt.Errorf("string to Bigint err: %s", s)
	}
	return i, nil
}
