package sorting

import (
	"reflect"
	"sort"
	"testing"

	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
)

func TestDependencySortLenCharacterizesExactSliceLength(t *testing.T) {
	tests := []struct {
		name string
		deps []pom.Dependency
		want int
	}{
		{name: "nil slice", deps: nil, want: 0},
		{name: "empty slice", deps: []pom.Dependency{}, want: 0},
		{name: "populated slice", deps: []pom.Dependency{{}, {}, {}}, want: 3},
	}
	if len(tests) == 0 {
		t.Fatal("dependency length characterization population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sorter := DependencySort{Deps: test.deps, SortKey: "unchanged length sort key"}

			if got := sorter.Len(); got != test.want {
				t.Fatalf("DependencySort.Len() = %d, want exact slice length %d", got, test.want)
			}
			if sorter.SortKey != "unchanged length sort key" {
				t.Fatalf("DependencySort.Len changed SortKey to %q", sorter.SortKey)
			}
		})
	}
}

func TestDependencySortSwapMutatesSharedSliceThroughValueReceiver(t *testing.T) {
	dependencies := []pom.Dependency{
		{GroupId: "group.first", ArtifactId: "first"},
		{GroupId: "group.middle", ArtifactId: "middle"},
		{GroupId: "group.last", ArtifactId: "last"},
	}
	sorter := DependencySort{Deps: dependencies, SortKey: "preserved swap sort key"}
	valueCopy := sorter

	valueCopy.Swap(0, 2)

	want := []pom.Dependency{
		{GroupId: "group.last", ArtifactId: "last"},
		{GroupId: "group.middle", ArtifactId: "middle"},
		{GroupId: "group.first", ArtifactId: "first"},
	}
	if !reflect.DeepEqual(dependencies, want) {
		t.Fatalf("DependencySort.Swap underlying slice = %#v, want %#v", dependencies, want)
	}
	if !reflect.DeepEqual(sorter.Deps, want) || !reflect.DeepEqual(valueCopy.Deps, want) {
		t.Fatalf("value-receiver views do not share the exact swapped slice: original=%#v copy=%#v",
			sorter.Deps, valueCopy.Deps)
	}
	if sorter.Len() != len(want) || valueCopy.Len() != len(want) {
		t.Fatalf("DependencySort.Swap changed length: original=%d copy=%d want=%d",
			sorter.Len(), valueCopy.Len(), len(want))
	}
	if sorter.SortKey != "preserved swap sort key" || valueCopy.SortKey != sorter.SortKey {
		t.Fatalf("DependencySort.Swap changed SortKey: original=%q copy=%q", sorter.SortKey, valueCopy.SortKey)
	}
}

func TestScopeWeightCharacterizesEveryExistingScopeClass(t *testing.T) {
	tests := []struct {
		scope string
		want  int
	}{
		{scope: "", want: 10},
		{scope: "compile", want: 20},
		{scope: "provided", want: 30},
		{scope: "runtime", want: 40},
		{scope: "system", want: 50},
		{scope: "import", want: 60},
		{scope: "test", want: 70},
		{scope: "unknown", want: 100},
		{scope: "Compile", want: 100},
	}
	if len(tests) == 0 {
		t.Fatal("scope-weight characterization population is empty")
	}

	for _, test := range tests {
		t.Run(test.scope, func(t *testing.T) {
			if got := scopeWeight(test.scope); got != test.want {
				t.Fatalf("scopeWeight(%q) = %d, want %d", test.scope, got, test.want)
			}
		})
	}
}

func TestGroupIDWeightCharacterizesContainsCaseAndEmptySortKey(t *testing.T) {
	tests := []struct {
		name    string
		groupID string
		sortKey string
		want    string
	}{
		{
			name:    "matching substring",
			groupID: "org.example.acme.service",
			sortKey: "acme",
			want:    "1-org.example.acme.service",
		},
		{
			name:    "matching substring need not be a segment",
			groupID: "org.concatenate.service",
			sortKey: "cat",
			want:    "1-org.concatenate.service",
		},
		{
			name:    "nonmatching substring",
			groupID: "org.example.external",
			sortKey: "acme",
			want:    "100-org.example.external",
		},
		{
			name:    "contains is case sensitive",
			groupID: "org.example.Acme.service",
			sortKey: "acme",
			want:    "100-org.example.Acme.service",
		},
		{
			name:    "empty sort key matches nonempty group",
			groupID: "org.example.external",
			sortKey: "",
			want:    "1-org.example.external",
		},
		{
			name:    "empty sort key matches empty group",
			groupID: "",
			sortKey: "",
			want:    "1-",
		},
	}
	if len(tests) == 0 {
		t.Fatal("group-weight characterization population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := groupIdWeight(test.groupID, test.sortKey); got != test.want {
				t.Fatalf("groupIdWeight(%q, %q) = %q, want %q",
					test.groupID, test.sortKey, got, test.want)
			}
		})
	}
}

func TestConcatCharacterizesExactScopeGroupAndArtifactKey(t *testing.T) {
	tests := []struct {
		name    string
		dep     pom.Dependency
		sortKey string
		want    string
	}{
		{
			name: "known scope and matching group",
			dep: pom.Dependency{
				Scope: "provided", GroupId: "com.acme.library", ArtifactId: "core",
			},
			sortKey: "acme",
			want:    "30:1-com.acme.library:core",
		},
		{
			name: "unknown scope and nonmatching group",
			dep: pom.Dependency{
				Scope: "custom", GroupId: "org.external", ArtifactId: "tool",
			},
			sortKey: "acme",
			want:    "100:100-org.external:tool",
		},
		{
			name: "empty sort key still uses matching group prefix",
			dep: pom.Dependency{
				GroupId: "org.external", ArtifactId: "default-scope",
			},
			sortKey: "",
			want:    "10:1-org.external:default-scope",
		},
		{
			name:    "all compared fields empty",
			dep:     pom.Dependency{},
			sortKey: "",
			want:    "10:1-:",
		},
	}
	if len(tests) == 0 {
		t.Fatal("concatenated-key characterization population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := concat(test.dep, test.sortKey); got != test.want {
				t.Fatalf("concat(%#v, %q) = %q, want %q", test.dep, test.sortKey, got, test.want)
			}
		})
	}
}

func TestDependencySortLessCharacterizesScopeGroupAndArtifactTieBreaks(t *testing.T) {
	tests := []struct {
		name    string
		left    pom.Dependency
		right   pom.Dependency
		sortKey string
		want    bool
	}{
		{
			name:    "known scope weight",
			left:    pom.Dependency{Scope: "compile", GroupId: "org.external", ArtifactId: "same"},
			right:   pom.Dependency{Scope: "runtime", GroupId: "org.external", ArtifactId: "same"},
			sortKey: "acme",
			want:    true,
		},
		{
			name:    "matching group prefix before nonmatching prefix",
			left:    pom.Dependency{Scope: "compile", GroupId: "z.acme", ArtifactId: "same"},
			right:   pom.Dependency{Scope: "compile", GroupId: "a.external", ArtifactId: "same"},
			sortKey: "acme",
			want:    true,
		},
		{
			name:    "group ID breaks matching-group tie",
			left:    pom.Dependency{Scope: "compile", GroupId: "com.acme.alpha", ArtifactId: "same"},
			right:   pom.Dependency{Scope: "compile", GroupId: "com.acme.beta", ArtifactId: "same"},
			sortKey: "acme",
			want:    true,
		},
		{
			name:    "artifact ID breaks group tie",
			left:    pom.Dependency{Scope: "compile", GroupId: "com.acme.same", ArtifactId: "alpha"},
			right:   pom.Dependency{Scope: "compile", GroupId: "com.acme.same", ArtifactId: "beta"},
			sortKey: "acme",
			want:    true,
		},
		{
			name:    "empty sort key gives both groups matching prefixes",
			left:    pom.Dependency{Scope: "compile", GroupId: "a.external", ArtifactId: "same"},
			right:   pom.Dependency{Scope: "compile", GroupId: "z.external", ArtifactId: "same"},
			sortKey: "",
			want:    true,
		},
		{
			name:    "reverse artifact tie break is not less",
			left:    pom.Dependency{Scope: "compile", GroupId: "com.acme.same", ArtifactId: "zeta"},
			right:   pom.Dependency{Scope: "compile", GroupId: "com.acme.same", ArtifactId: "alpha"},
			sortKey: "acme",
			want:    false,
		},
	}
	if len(tests) == 0 {
		t.Fatal("dependency comparison characterization population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := dependencySort(test.left, test.right, test.sortKey); got != test.want {
				t.Fatalf("dependencySort left=%q right=%q = %t, want %t",
					concat(test.left, test.sortKey), concat(test.right, test.sortKey), got, test.want)
			}
			sorter := DependencySort{
				Deps:    []pom.Dependency{test.left, test.right},
				SortKey: test.sortKey,
			}
			if got := sorter.Less(0, 1); got != test.want {
				t.Fatalf("DependencySort.Less(0, 1) = %t, want %t", got, test.want)
			}
		})
	}
}

func TestDependencySortUsesStrictLessThanForEqualKeys(t *testing.T) {
	left := pom.Dependency{
		GroupId: "com.acme.same", ArtifactId: "artifact", Scope: "compile", Version: "1.0.0",
	}
	right := pom.Dependency{
		Comment: "ignored comment", GroupId: "com.acme.same", ArtifactId: "artifact",
		Version: "9.9.9", Type_: "test-jar", Classifier: "tests", Scope: "compile",
		SystemPath: "/ignored/system/path", Optional: "true",
		Exclusions: &pom.Exclusions{Exclusion: []pom.Exclusion{{
			GroupId: "ignored.exclusion", ArtifactId: "ignored-artifact",
		}}},
	}
	sorter := DependencySort{Deps: []pom.Dependency{left, right}, SortKey: "acme"}

	if dependencySort(left, right, sorter.SortKey) || dependencySort(right, left, sorter.SortKey) {
		t.Fatalf("equal dependency keys compared less in one direction: left=%q right=%q",
			concat(left, sorter.SortKey), concat(right, sorter.SortKey))
	}
	if sorter.Less(0, 1) || sorter.Less(1, 0) {
		t.Fatal("DependencySort.Less is not strict for dependencies with equal comparison keys")
	}
}

func TestDependencySortPreservesLegacyLexicalWeight100Order(t *testing.T) {
	unknown := pom.Dependency{Scope: "custom", GroupId: "org.external", ArtifactId: "unknown"}
	empty := pom.Dependency{Scope: "", GroupId: "org.external", ArtifactId: "empty"}
	compile := pom.Dependency{Scope: "compile", GroupId: "org.external", ArtifactId: "compile"}
	const sortKey = "acme"

	if got := concat(unknown, sortKey); got != "100:100-org.external:unknown" {
		t.Fatalf("unknown-scope comparison key = %q", got)
	}
	if got := concat(empty, sortKey); got != "10:100-org.external:empty" {
		t.Fatalf("empty-scope comparison key = %q", got)
	}
	if !dependencySort(unknown, empty, sortKey) {
		t.Fatalf("legacy lexical comparison no longer orders %q before %q",
			concat(unknown, sortKey), concat(empty, sortKey))
	}

	dependencies := []pom.Dependency{compile, empty, unknown}
	want := []pom.Dependency{unknown, empty, compile}
	if len(dependencies) == 0 || len(want) == 0 {
		t.Fatal("legacy lexical full-sort population is empty")
	}
	sort.Sort(DependencySort{Deps: dependencies, SortKey: sortKey})
	if !reflect.DeepEqual(dependencies, want) {
		t.Fatalf("legacy lexical scope order = %#v, want %#v", dependencies, want)
	}
}

func TestDependencySortProducesExactInPlaceFullOrder(t *testing.T) {
	unknown := pom.Dependency{Scope: "custom", GroupId: "org.external", ArtifactId: "unknown"}
	empty := pom.Dependency{Scope: "", GroupId: "org.external", ArtifactId: "empty"}
	compileMatchAlphaA := pom.Dependency{Scope: "compile", GroupId: "com.acme.alpha", ArtifactId: "alpha"}
	compileMatchAlphaB := pom.Dependency{Scope: "compile", GroupId: "com.acme.alpha", ArtifactId: "beta"}
	compileMatchZeta := pom.Dependency{Scope: "compile", GroupId: "com.acme.zeta", ArtifactId: "same"}
	compileExternal := pom.Dependency{Scope: "compile", GroupId: "a.external", ArtifactId: "external"}
	provided := pom.Dependency{Scope: "provided", GroupId: "com.acme.scope", ArtifactId: "provided"}
	runtime := pom.Dependency{Scope: "runtime", GroupId: "com.acme.scope", ArtifactId: "runtime"}
	system := pom.Dependency{Scope: "system", GroupId: "com.acme.scope", ArtifactId: "system"}
	imported := pom.Dependency{Scope: "import", GroupId: "com.acme.scope", ArtifactId: "import"}
	testScope := pom.Dependency{Scope: "test", GroupId: "com.acme.scope", ArtifactId: "test"}

	dependencies := []pom.Dependency{
		testScope,
		compileExternal,
		runtime,
		compileMatchAlphaB,
		empty,
		system,
		compileMatchZeta,
		unknown,
		provided,
		compileMatchAlphaA,
		imported,
	}
	want := []pom.Dependency{
		unknown,
		empty,
		compileMatchAlphaA,
		compileMatchAlphaB,
		compileMatchZeta,
		compileExternal,
		provided,
		runtime,
		system,
		imported,
		testScope,
	}
	if len(dependencies) == 0 || len(want) == 0 {
		t.Fatal("complete dependency sort population is empty")
	}
	sorter := DependencySort{Deps: dependencies, SortKey: "acme"}

	sort.Sort(sorter)

	if !reflect.DeepEqual(dependencies, want) {
		t.Fatalf("DependencySort in-place result differs:\n got: %#v\nwant: %#v", dependencies, want)
	}
	if !reflect.DeepEqual(sorter.Deps, want) {
		t.Fatalf("DependencySort value-receiver slice view = %#v, want %#v", sorter.Deps, want)
	}
	if sorter.Len() != len(want) || sorter.SortKey != "acme" {
		t.Fatalf("DependencySort full sort changed API state: len=%d SortKey=%q", sorter.Len(), sorter.SortKey)
	}
}
