package mute

import "testing"

import "fmt"

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

func TestEscalationDoubles(t *testing.T) {
	world := New(1000, 100000)
	for index, want := range []int64{1000, 2000, 4000, 8000} {
		if got := world.Mute("p1", 0, fmt.Sprintf("r-%d", index)); got != want {
			t.Fatalf("第 %d 次禁言应该是 %d：%v", index+1, want, got)
		}
	}
}

func TestCapLimitsDuration(t *testing.T) {
	world := New(1000, 3000)
	world.Mute("p1", 0, "r1")
	world.Mute("p1", 0, "r2")
	for index := 3; index <= 6; index++ {
		if got := world.Mute("p1", 0, fmt.Sprintf("r-%d", index)); got != 3000 {
			t.Fatalf("第 %d 次禁言应该封顶在 3000：%v", index, got)
		}
	}
}

func TestExpiryEndsMute(t *testing.T) {
	world := New(1000, 8000)
	world.Mute("p1", 0, "r1")
	if !world.Muted("p1", 999) {
		t.Fatal("禁言期内应该被禁言")
	}
	if world.Muted("p1", 1000) {
		t.Fatal("到期时刻不该再被禁言")
	}
}

func TestLiftClearsMuteKeepsOffenses(t *testing.T) {
	world := New(1000, 8000)
	world.Mute("p1", 0, "r1")
	if !world.Lift("p1", 100, "r2") {
		t.Fatal("禁言期内解除应该返回 true")
	}
	if world.Muted("p1", 200) {
		t.Fatal("解除之后不该再被禁言")
	}
	if world.Offenses("p1") != 1 {
		t.Fatal("解除应该保留违规次数")
	}
}

func TestLiftWithoutMuteReturnsFalse(t *testing.T) {
	world := New(1000, 8000)
	if world.Lift("p1", 0, "r1") {
		t.Fatal("没有禁言时解除应该返回 false")
	}
	world.Mute("p1", 0, "r2")
	if world.Lift("p1", 1000, "r3") {
		t.Fatal("禁言已到期时解除应该返回 false")
	}
}

func TestMuteIdempotentByReqID(t *testing.T) {
	world := New(1000, 100000)
	if got := world.Mute("p1", 0, "same"); got != 1000 {
		t.Fatalf("第一次禁言应该是 1000：%v", got)
	}
	if got := world.Mute("p1", 0, "same"); got != 1000 {
		t.Fatalf("重复请求应该返回原时长 1000：%v", got)
	}
	if world.Offenses("p1") != 1 {
		t.Fatal("重复请求不该增加违规次数")
	}
	if world.Muted("p1", 1000) {
		t.Fatal("重复请求不该延长禁言时长")
	}
}

func TestScannedDoesNotGrowWithPlayers(t *testing.T) {
	world := New(1000, 100000)
	for index := 0; index < 3000; index++ {
		world.Mute(fmt.Sprintf("p-%d", index), 0, fmt.Sprintf("r-%d", index))
	}
	for index := 0; index < 3000; index++ {
		world.Muted("p-2999", 0)
	}
	if got := world.Scanned(); got > 6000 {
		t.Fatalf("扫描量不该随玩家数线性增长：%v", got)
	}
}
