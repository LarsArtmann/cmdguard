package v2

import (
	"context"
	"slices"
	"testing"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestTagsSeq_MatchesTags(t *testing.T) {
	t.Parallel()

	type TestConfig struct {
		Name    string `flag:"name"    help:"Name"`
		Verbose bool   `flag:"verbose" help:"Verbose"`
		Count   int    `flag:"count"   help:"Count"`
	}

	registry, err := NewFlagRegistry(&TestConfig{})
	testutil.AssertNoError(t, err)

	sliceTags := registry.Tags()
	var seqTags []FlagTag

	for tag := range registry.TagsSeq() {
		seqTags = append(seqTags, tag)
	}

	testutil.AssertEqual(t, len(sliceTags), len(seqTags))

	for i := range sliceTags {
		testutil.AssertEqual(t, sliceTags[i].Name, seqTags[i].Name)
		testutil.AssertEqual(t, sliceTags[i].Field, seqTags[i].Field)
	}
}

func TestFlagNamesSeq_MatchesFlagNames(t *testing.T) {
	t.Parallel()

	type TestConfig struct {
		Name    string `flag:"name"    help:"Name"`
		Verbose bool   `flag:"verbose" help:"Verbose"`
	}

	registry, err := NewFlagRegistry(&TestConfig{})
	testutil.AssertNoError(t, err)

	sliceNames := registry.FlagNames()
	var seqNames []string

	for name := range registry.FlagNamesSeq() {
		seqNames = append(seqNames, name)
	}

	testutil.AssertEqual(t, len(sliceNames), len(seqNames))

	for i := range sliceNames {
		testutil.AssertEqual(t, sliceNames[i], seqNames[i])
	}
}

func TestPathSeq_MatchesPath(t *testing.T) {
	t.Parallel()

	ctx := NewBranchingFlowContext(context.Background())
	child, cancel := ctx.Branch("deploy")
	defer cancel()

	grandchild, cancel2 := child.Branch("production")
	defer cancel2()

	slicePath := grandchild.Path()
	var seqPath []string

	for p := range grandchild.PathSeq() {
		seqPath = append(seqPath, p)
	}

	testutil.AssertEqual(t, len(slicePath), len(seqPath))

	for i := range slicePath {
		testutil.AssertEqual(t, slicePath[i], seqPath[i])
	}
}

func TestChildrenSeq_MatchesChildren(t *testing.T) {
	t.Parallel()

	ctx := NewBranchingFlowContext(context.Background())

	_, cancel1 := ctx.Branch("alpha")
	defer cancel1()

	_, cancel2 := ctx.Branch("beta")
	defer cancel2()

	_, cancel3 := ctx.Branch("gamma")
	defer cancel3()

	sliceChildren := ctx.Children()
	seqChildren := slices.Collect(ctx.ChildrenSeq())

	testutil.AssertEqual(t, len(sliceChildren), len(seqChildren))

	for i := range sliceChildren {
		testutil.AssertEqual(t, sliceChildren[i].PathString(), seqChildren[i].PathString())
	}
}

func TestTagsSeq_EmptyRegistry(t *testing.T) {
	t.Parallel()

	type EmptyConfig struct {
		NoFlag string
	}

	registry, err := NewFlagRegistry(&EmptyConfig{})
	testutil.AssertNoError(t, err)

	count := 0

	for range registry.TagsSeq() {
		count++
	}

	testutil.AssertEqual(t, 0, count)
}

func TestFlagNamesSeq_EmptyRegistry(t *testing.T) {
	t.Parallel()

	type EmptyConfig struct {
		NoFlag string
	}

	registry, err := NewFlagRegistry(&EmptyConfig{})
	testutil.AssertNoError(t, err)

	count := 0

	for range registry.FlagNamesSeq() {
		count++
	}

	testutil.AssertEqual(t, 0, count)
}
