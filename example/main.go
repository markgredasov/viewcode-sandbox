package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox"
	sandbox_config "github.com/markgredasov/viewcode-sandbox/internal/integration/sandbox/config"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func main() {
	ctx := context.Background()

	sandboxConfig := sandbox_config.NewConfigMust()
	adapter, err := sandbox.NewAdapter(sandboxConfig)
	if err != nil {
		log.Fatal(err)
	}

	exec, err := model.NewExecution("", "go", `
	package solution

	func twoSum(nums []int, target int) []int {
		m := make(map[int]int)
		for k, v := range nums {
			if idx, ok := m[target-v]; ok {
				return []int{idx, k}
			}
			m[v] = k
		}
		return nil
	}
	`)
	if err != nil {
		log.Fatal(err)
	}

	tests, err := model.NewVerification(`
	package solution

	import "testing"

	func isValidSolution(nums []int, target int, result []int) bool {
		if result == nil {
			return false
		}
		
		if len(result) != 2 {
			return false
		}
		
		i, j := result[0], result[1]
		if i < 0 || i >= len(nums) || j < 0 || j >= len(nums) {
			return false
		}
		
		if i == j {
			return false
		}
		
		return nums[i]+nums[j] == target
	}

	func TestTwoSum(t *testing.T) {
		tests := []struct {
			name     string
			nums     []int
			target   int
			hasSolution bool
		}{
			{
				name:     "Example 1: Basic case",
				nums:     []int{2, 7, 11, 15},
				target:   9,
				hasSolution: true,
			},
			{
				name:     "Example 2: Numbers not sorted",
				nums:     []int{3, 2, 4},
				target:   6,
				hasSolution: true,
			},
			{
				name:     "Example 3: Duplicate numbers",
				nums:     []int{3, 3},
				target:   6,
				hasSolution: true,
			},
			{
				name:     "Negative numbers",
				nums:     []int{-1, -2, -3, -4, -5},
				target:   -8,
				hasSolution: true,
			},
			{
				name:     "Zero target",
				nums:     []int{0, 4, 3, 0},
				target:   0,
				hasSolution: true,
			},
			{
				name:     "Large numbers",
				nums:     []int{1000000, 2000000, 3000000},
				target:   3000000,
				hasSolution: true,
			},
			{
				name:     "First and last elements",
				nums:     []int{1, 2, 3, 4, 5},
				target:   6,
				hasSolution: true,
			},
			{
				name:     "Middle elements",
				nums:     []int{1, 2, 3, 4, 5},
				target:   7,
				hasSolution: true,
			},
			{
				name:     "Negative and positive",
				nums:     []int{-3, 4, 3, 90},
				target:   0,
				hasSolution: true,
			},
			{
				name:     "No solution",
				nums:     []int{1, 2, 3, 4},
				target:   10,
				hasSolution: false,
			},
			{
				name:     "Single element (no solution)",
				nums:     []int{5},
				target:   10,
				hasSolution: false,
			},
			{
				name:     "Empty array (no solution)",
				nums:     []int{},
				target:   10,
				hasSolution: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := twoSum(tt.nums, tt.target)
				
				if tt.hasSolution {
					if !isValidSolution(tt.nums, tt.target, result) {
						t.Errorf("twoSum(%v, %d) = %v is not a valid solution", tt.nums, tt.target, result)
					}
				} else {
					if result != nil {
						t.Errorf("twoSum(%v, %d) = %v, expected nil", tt.nums, tt.target, result)
					}
				}
			})
		}
	}

	func BenchmarkTwoSum(b *testing.B) {
		nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		target := 15
		for i := 0; i < b.N; i++ {
			twoSum(nums, target)
		}
	}
	`)
	if err != nil {
		log.Fatal(err)
	}

	res, err := adapter.RunTests(ctx, exec, tests)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println(res)
}
