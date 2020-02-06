package rules

import (
	"reflect"
	"testing"

	"k8s.io/gengo/types"
)

func TestListTypeMissing(t *testing.T) {
	tcs := []struct {
		// name of test case
		name string
		t    *types.Type

		// expected list of violation fields
		expected []string

		expectedError error
	}{
		{
			name:     "none",
			t:        &types.Type{},
			expected: []string{},
		},
		{
			name: "simple missing",
			t: &types.Type{
				Kind: types.Struct,
				Members: []types.Member{
					types.Member{
						Name: "Containers",
						Type: &types.Type{
							Kind: types.Slice,
						},
					},
				},
			},
			expected: []string{"Containers"},
		},
		{
			name: "simple passing",
			t: &types.Type{
				Kind: types.Struct,
				Members: []types.Member{
					types.Member{
						Name: "Containers",
						Type: &types.Type{
							Kind: types.Slice,
						},
						CommentLines: []string{"+listType=map"},
					},
				},
			},
			expected: []string{},
		},

		{
			name: "list Items field should not be annotated",
			t: &types.Type{
				Kind: types.Struct,
				Members: []types.Member{
					types.Member{
						Name: "Items",
						Type: &types.Type{
							Kind: types.Slice,
						},
						CommentLines: []string{"+listType=map"},
					},
					types.Member{
						Name:     "ListMeta",
						Embedded: true,
						Type: &types.Type{
							Kind: types.Struct,
						},
					},
				},
			},
			expected: []string{"Items"},
		},

		{
			name: "list Items field without annotation should pass validation",
			t: &types.Type{
				Kind: types.Struct,
				Members: []types.Member{
					types.Member{
						Name: "Items",
						Type: &types.Type{
							Kind: types.Slice,
						},
					},
					types.Member{
						Name:     "ListMeta",
						Embedded: true,
						Type: &types.Type{
							Kind: types.Struct,
						},
					},
				},
			},
			expected: []string{},
		},

		{
			name: "a list that happens to be called Items (i.e. nested, not top-level list) needs annotations",
			t: &types.Type{
				Kind: types.Struct,
				Members: []types.Member{
					types.Member{
						Name: "Items",
						Type: &types.Type{
							Kind: types.Slice,
						},
					},
				},
			},
			expected: []string{"Items"},
		},

		{
			name: "lists with no annotations but members that look like a map should return a useful hint",
			t: &types.Type{
				Name: types.Name{Name: "PodSpec"},
				Kind: types.Struct,
				Members: []types.Member{
					types.Member{
						Name: "InitContainers",
						Type: &types.Type{
							Kind: types.Slice,
							Elem: &types.Type{
								Name: types.Name{Name: "Container"},
								Kind: types.Struct,
								Members: []types.Member{
									types.Member{
										Name: "Name",
										Type: &types.Type{
											Kind: types.Builtin,
										},
									},
									types.Member{
										Name: "Image",
										Type: &types.Type{
											Kind: types.Builtin,
										},
									},
								},
							},
						},
					},
				},
			},
			expected: []string{"InitContainers; should be taged as +listType=map and +listKey=name"},
		},
	}

	rule := &ListTypeMissing{}
	for _, tc := range tcs {
		violations, err := rule.Validate(tc.t)
		if !reflect.DeepEqual(violations, tc.expected) {
			t.Errorf("unexpected validation result: test name %v, want: %v, got: %v",
				tc.name, tc.expected, violations)
		}
		if tc.expectedError != nil {
			if err == nil {
				t.Errorf("Expected to get error %q, got: nil", tc.expectedError)
			} else if tc.expectedError.Error() != err.Error() {
				t.Errorf("Expected to get error %q, got: %q", tc.expectedError, err)
			}
		}
	}
}
