package utils

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestIsTagList(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1,2,3", true},
		{"1,2,3,", true},
		{"1", true},
		{"1-2", false},
		{"1 2 3", false},
	}

	for _, test := range tests {
		result := IsTagList(test.input)
		require.Equal(t, test.expected, result)
	}
}

func TestTagsToIntSlice(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput []int32
		expectedError  error
	}{
		{input: "1,2,3", expectedOutput: []int32{1, 2, 3}, expectedError: nil},
		{input: "10,20,30", expectedOutput: []int32{10, 20, 30}, expectedError: nil},
		{input: "100,200,300", expectedOutput: []int32{100, 200, 300}, expectedError: nil},
		{input: "", expectedOutput: []int32{}, expectedError: nil},              // Empty input should return an empty slice
		{input: "a,b,c", expectedOutput: nil, expectedError: strconv.ErrSyntax}, // Invalid number should return an error
	}

	for _, testCase := range testCases {
		output, err := TagsToIntSlice(testCase.input)
		require.Equal(t, output, testCase.expectedOutput)
		if testCase.expectedError == nil {
			require.NoError(t, err)
		} else {
			require.Error(t, err, testCase.expectedError)
		}

	}
}

func TestCompareTagLists(t *testing.T) {
	// Test case 1: Both slices are empty
	slice1 := []string{}
	slice2 := []string{}
	expectedUniqueSlice1 := []string{}
	expectedUniqueSlice2 := []string{}
	actualUniqueSlice1, actualUniqueSlice2 := CompareTagLists(slice1, slice2)
	require.Equal(t, expectedUniqueSlice1, actualUniqueSlice1, "UniqueSlice1 should be empty")
	require.Equal(t, expectedUniqueSlice2, actualUniqueSlice2, "UniqueSlice2 should be empty")

	// Test case 2: slice1 is empty, slice2 has elements
	slice1 = []string{}
	slice2 = []string{"a", "b", "c"}
	expectedUniqueSlice1 = []string{}
	expectedUniqueSlice2 = []string{"a", "b", "c"}
	actualUniqueSlice1, actualUniqueSlice2 = CompareTagLists(slice1, slice2)
	require.Equal(t, expectedUniqueSlice1, actualUniqueSlice1, "UniqueSlice1 should be empty")
	require.Equal(t, expectedUniqueSlice2, actualUniqueSlice2, "UniqueSlice2 should contain all elements of slice2")

	// Test case 3: slice1 has elements, slice2 is empty
	slice1 = []string{"a", "b", "c"}
	slice2 = []string{}
	expectedUniqueSlice1 = []string{"a", "b", "c"}
	expectedUniqueSlice2 = []string{}
	actualUniqueSlice1, actualUniqueSlice2 = CompareTagLists(slice1, slice2)
	require.Equal(t, expectedUniqueSlice1, actualUniqueSlice1, "UniqueSlice1 should contain all elements of slice1")
	require.Equal(t, expectedUniqueSlice2, actualUniqueSlice2, "UniqueSlice2 should be empty")

	// Test case 4: Both slices have some common elements
	slice1 = []string{"a", "b", "c", "d"}
	slice2 = []string{"c", "d", "e", "f"}
	expectedUniqueSlice1 = []string{"a", "b"}
	expectedUniqueSlice2 = []string{"e", "f"}
	actualUniqueSlice1, actualUniqueSlice2 = CompareTagLists(slice1, slice2)
	require.ElementsMatch(t, expectedUniqueSlice1, actualUniqueSlice1, "UniqueSlice1 should contain elements 'a' and 'b'")
	require.ElementsMatch(t, expectedUniqueSlice2, actualUniqueSlice2, "UniqueSlice2 should contain elements 'e' and 'f'")
}
