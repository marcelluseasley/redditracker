package client

import (
	"testing"
)

func TestPostInUserMap(t *testing.T) {
	posts := []string{"post1", "post2", "post3"}

	testCases := []struct {
		name     string
		postID   string
		expected bool
	}{
		{"postID exists in posts", "post2", true},
		{"postID does not exist in posts", "post4", false},
	}

	for _, tc := range testCases {
		if result := postInUserMap(tc.postID, posts); result != tc.expected {
			t.Errorf("postInUserMap(%s, %v) = %t, expected %t", tc.postID, posts, result, tc.expected)
		}
	}
}

func TestSortUserByNumPosts(t *testing.T) {
	userMap := map[string][]string{
		"user1": {"post1", "post2"},
		"user2": {"post3"},
		"user3": {"post4", "post5", "post6"},
	}

	expected := []UserPostCount{
		{Username: "user3", PostCount: 3},
		{Username: "user1", PostCount: 2},
		{Username: "user2", PostCount: 1},
	}

	result := sortUserByNumPosts(userMap)

	if len(result) != len(expected) {
		t.Errorf("sortUserByNumPosts() returned %d elements, expected %d", len(result), len(expected))
	}

	for i := 0; i < len(result); i++ {
		if result[i].Username != expected[i].Username || result[i].PostCount != expected[i].PostCount {
			t.Errorf("sortUserByNumPosts() element at index %d does not match the expected result", i)
		}
	}
}
