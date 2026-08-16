package model

import (
	"fmt"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

type Category string

const (
	CategoryArray            = Category("array")
	CategoryHashMap          = Category("hash map")
	CategoryTwoPointers      = Category("two pointers")
	CategoryStack            = Category("stack")
	CategoryBinarySearch     = Category("binary search")
	CategorySlidingWindow    = Category("sliding window")
	CategoryLinkedList       = Category("linked list")
	CategoryTree             = Category("tree")
	CategoryTrie             = Category("trie")
	CategoryBacktracking     = Category("backtracking")
	CategoryHeap             = Category("heap")
	CategoryOneDimensionalDP = Category("1d dynamic programming")
	CategoryTwoDimensionalDP = Category("2d dynamic programming")
	CategoryBitManipulation  = Category("bit manipulation")
	CategoryMath             = Category("math")
	CategoryGraph            = Category("graph")
	CategoryAdvancedGraph    = Category("advanced graph")
	CategoryGreedy           = Category("greedy")
	CategoryInterval         = Category("interval")
)

func NewCategory(category string) (Category, error) {
	categoryValueObj := Category(category)
	if !categoryValueObj.IsValid() {
		return "", fmt.Errorf("invalid category: %s: %w", category, errs.ErrInvalidArgument)
	}

	return categoryValueObj, nil
}

func (c Category) String() string {
	return string(c)
}

func (c Category) IsValid() bool {
	switch c {
	case CategoryArray, CategoryHashMap, CategoryTwoPointers, CategoryStack, CategoryBinarySearch,
		CategorySlidingWindow, CategoryLinkedList, CategoryTree, CategoryTrie, CategoryBacktracking,
		CategoryHeap, CategoryOneDimensionalDP, CategoryTwoDimensionalDP, CategoryBitManipulation,
		CategoryMath, CategoryGraph, CategoryAdvancedGraph, CategoryGreedy, CategoryInterval:
		return true
	}
	return false
}
