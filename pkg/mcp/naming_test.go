package mcp

import "testing"

func TestQualifiedToolName_SanitizesAndPrefixes(t *testing.T) {
	got := QualifiedToolName("My Server", "Read-File!")
	want := "mcp_my_server__read_file"
	if got != want {
		t.Fatalf("QualifiedToolName() = %q, want %q", got, want)
	}
}

func TestQualifiedToolName_TrimToMaxLen(t *testing.T) {
	longToolName := "tool_name_with_many_segments_and_extra_text_that_exceeds_the_limit_significantly"
	got := QualifiedToolName("server", longToolName)
	if len(got) > qualifiedNameMaxLen {
		t.Fatalf("qualified name length = %d, want <= %d", len(got), qualifiedNameMaxLen)
	}
}
