package main

import "testing"

func TestSplitTextToLines(t *testing.T) {
	if len(SplitTextToLines("")) != 1 {
		t.Error("Failed to split an empty string")
	}
	if len(SplitTextToLines("This is one line.")) != 1 {
		t.Error("Failed to split one line")
	}
	if len(SplitTextToLines("This is one line with a newline at the end\n")) != 2 {
		t.Error("Failed to split one line with a newline at the end")
	}
	if len(SplitTextToLines("This is the first of two lines\nThis is the second of the two lines")) != 2 {
		t.Error("Failed to split two lines")
	}
	if len(SplitTextToLines("This is the first of three lines\n\nThis is the second of the two lines")) != 3 {
		t.Error("Failed to split three lines (second one blank)")
	}
	if len(SplitTextToLines("This is the first of three lines\n\nThis is the second of the two lines")[1]) != 0 {
		t.Error("Middle line should have been empty")
	}
	lines := SplitTextToLines("\n\n1\n2\n345")
	if len(lines) != 5 {
		t.Error("Should have had 5 lines")
	}
	if lines[2] != "1" {
		t.Error("Invalid result on line 3")
	}
	if lines[3] != "2" {
		t.Error("Invalid result on line 4")
	}
	if lines[4] != "345" {
		t.Error("Invalid result on line 5")
	}
}
