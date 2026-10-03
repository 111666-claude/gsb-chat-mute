package mute

import "testing"

func TestNotMutedAtStart(t *testing.T) {
	if New(1000, 8000).Muted("p1", 0) {
		t.Fatal("没有违规时不该被禁言")
	}
}

func TestFirstMuteIsBase(t *testing.T) {
	if got := New(1000, 8000).Mute("p1", 0, "r1"); got != 1000 {
		t.Fatalf("第一次禁言应该是基础时长：%v", got)
	}
}

func TestOffensesStartsAtZero(t *testing.T) {
	if New(1000, 8000).Offenses("p1") != 0 {
		t.Fatal("违规次数初始应该是 0")
	}
}

func TestCapIsStored(t *testing.T) {
	if New(1000, 8000).capMs != 8000 {
		t.Fatal("封顶时长应该能读回来")
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New(1000, 8000).Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}
