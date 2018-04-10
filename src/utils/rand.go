package utils

import (
	"math/rand"
	"strings"

	"github.com/satori/go.uuid"
)

func GenerateUuidV4() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

func GenerateUuidV5(s string) string {
	return strings.Replace(uuid.NewV5(uuid.NamespaceURL, s).String(), "-", "", -1)
}

func UniqRands(l int, n int) []int {
	set := make(map[int]struct{})
	nums := make([]int, 0, l)
	for {
		num := rand.Intn(n)
		if _, ok := set[num]; !ok {
			set[num] = struct{}{}
			nums = append(nums, num)
		}
		if len(nums) == l {
			goto exit
		}
	}
exit:
	return nums
}
