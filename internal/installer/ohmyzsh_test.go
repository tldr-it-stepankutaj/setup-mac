package installer

import (
	"reflect"
	"testing"
)

func TestParsePluginsLine(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
		ok   bool
	}{
		{"basic", "plugins=(git docker)", []string{"git", "docker"}, true},
		{"empty list", "plugins=()", []string{}, true},
		{"extra whitespace inside parens", "plugins=(  git   docker  fzf )", []string{"git", "docker", "fzf"}, true},
		{"leading whitespace", "  plugins=(git)", []string{"git"}, true},
		{"single plugin", "plugins=(git)", []string{"git"}, true},

		// Multi-line forms aren't supported — must return ok=false so the caller
		// can warn instead of corrupting the file.
		{"multi-line opening", "plugins=(", nil, false},
		{"unrelated line", "ZSH_THEME=\"powerlevel10k/powerlevel10k\"", nil, false},
		{"plugin elsewhere in line", "# my plugins=(git)", nil, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := parsePluginsLine(tc.in)
			if ok != tc.ok {
				t.Fatalf("parsePluginsLine(%q) ok=%v, want %v", tc.in, ok, tc.ok)
			}
			if ok && !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parsePluginsLine(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestMergePluginLists(t *testing.T) {
	cases := []struct {
		name             string
		existing, extras []string
		want             []string
	}{
		{
			"adds missing, keeps existing order",
			[]string{"git", "docker", "kubectl"},
			[]string{"git", "fzf", "ssh-agent"},
			[]string{"git", "docker", "kubectl", "fzf", "ssh-agent"},
		},
		{
			"no extras",
			[]string{"git", "docker"},
			nil,
			[]string{"git", "docker"},
		},
		{
			"no existing",
			nil,
			[]string{"git", "docker"},
			[]string{"git", "docker"},
		},
		{
			"identical",
			[]string{"git", "docker"},
			[]string{"git", "docker"},
			[]string{"git", "docker"},
		},
		{
			"existing has duplicates — dedupe",
			[]string{"git", "git", "docker"},
			[]string{"fzf"},
			[]string{"git", "docker", "fzf"},
		},
		{
			// The motivating real-world case: user has 12 plugins, config has 6 (all subset).
			// Merge should leave the existing line completely unchanged.
			"config is subset of existing — no changes",
			[]string{"git", "docker", "kubectl", "fzf", "zsh-autosuggestions", "zsh-syntax-highlighting", "colored-man-pages", "colorize", "pip", "python", "brew", "ssh-agent"},
			[]string{"git", "docker", "kubectl", "fzf", "zsh-autosuggestions", "zsh-syntax-highlighting"},
			[]string{"git", "docker", "kubectl", "fzf", "zsh-autosuggestions", "zsh-syntax-highlighting", "colored-man-pages", "colorize", "pip", "python", "brew", "ssh-agent"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mergePluginLists(tc.existing, tc.extras)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("mergePluginLists(%v, %v) = %v, want %v", tc.existing, tc.extras, got, tc.want)
			}
		})
	}
}

func TestDiffPluginLists(t *testing.T) {
	got := diffPluginLists(
		[]string{"git", "fzf", "ssh-agent", "docker"},
		[]string{"git", "docker"},
	)
	want := []string{"fzf", "ssh-agent"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("diffPluginLists() = %v, want %v", got, want)
	}
}
