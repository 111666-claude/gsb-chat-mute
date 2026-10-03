package mute

import (
	"fmt"
	"testing"
)

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
	want := []int64{1000, 2000, 4000, 8000}
	for index, duration := range want {
		if got := world.Mute("p1", 0, fmt.Sprintf("r%d", index)); got != duration {
			t.Fatalf("第 %d 次违规时长应该是 %d：%v", index+1, duration, got)
		}
	}
	if got := world.Offenses("p1"); got != 4 {
		t.Fatalf("违规次数应该是 4：%v", got)
	}
}

func TestDurationCapped(t *testing.T) {
	world := New(1000, 3000)
	want := []int64{1000, 2000, 3000, 3000, 3000, 3000}
	for index, duration := range want {
		if got := world.Mute("p1", 0, fmt.Sprintf("r%d", index)); got != duration {
			t.Fatalf("第 %d 次违规时长应该是 %d（封顶 3000）：%v", index+1, duration, got)
		}
	}
}

func TestExpiredMuteLifts(t *testing.T) {
	world := New(1000, 8000)
	world.Mute("p1", 0, "r1")
	if !world.Muted("p1", 999) {
		t.Fatal("禁言结束前应该还在禁言")
	}
	if world.Muted("p1", 1000) {
		t.Fatal("当前时刻不小于结束时刻时不该再禁言")
	}
}

func TestLiftClearsMuteKeepsOffenses(t *testing.T) {
	world := New(1000, 8000)
	world.Mute("p1", 0, "r1")
	world.Mute("p1", 0, "r2")
	if !world.Lift("p1", 100, "lift1") {
		t.Fatal("禁言中 Lift 应该返回 true")
	}
	if world.Muted("p1", 200) {
		t.Fatal("Lift 之后不该再被禁言")
	}
	if world.Lift("p1", 300, "lift2") {
		t.Fatal("没有禁言时 Lift 应该返回 false")
	}
	if got := world.Offenses("p1"); got != 2 {
		t.Fatalf("Lift 应保留违规次数 2：%v", got)
	}
	if next := world.Mute("p1", 0, "r3"); next != 4000 {
		t.Fatalf("Lift 后第 3 次违规仍应按次数翻倍到 4000：%v", next)
	}
}

func TestMuteIdempotentByReqID(t *testing.T) {
	world := New(1000, 100000)
	first := world.Mute("p1", 0, "same")
	second := world.Mute("p1", 500, "same")
	if first != 1000 || second != 1000 {
		t.Fatalf("重复 reqID 应返回首次时长 1000：first=%v second=%v", first, second)
	}
	if got := world.Offenses("p1"); got != 1 {
		t.Fatalf("重复 reqID 不该增加违规次数：%v", got)
	}
	if !world.Muted("p1", 999) {
		t.Fatal("重复 reqID 不该延长或重置禁言，结束时刻仍应是 1000")
	}
	if world.Muted("p1", 1000) {
		t.Fatal("禁言应在 1000 到期")
	}
}

func TestScannedIndependentOfPlayerCount(t *testing.T) {
	world := New(1000, 100000)
	for index := 0; index < 1000; index++ {
		world.Mute(fmt.Sprintf("p-%d", index), 0, fmt.Sprintf("r-%d", index))
	}
	for index := 0; index < 1000; index++ {
		world.Muted("p-999", 0)
	}
	if got := world.Scanned(); got > 2000 {
		t.Fatalf("扫描量不应随玩家数线性增长：%v", got)
	}
}
