package util

import "testing"

func TestIsJosaAttachable(t *testing.T) {
	if !isJosaAttachable('플', '은') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('플', '이') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('플', '을') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('플', '과') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('플', '은') {
		t.Error("Expected result should match.")
	}

	if isJosaAttachable('플', '는') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('플', '가') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('플', '를') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('플', '와') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('플', '야') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('플', '여') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('플', '라') {
		t.Error("Expected result should match.")
	}

	if isJosaAttachable('프', '은') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('프', '이') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('프', '을') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('프', '과') {
		t.Error("Expected result should match.")
	}
	if isJosaAttachable('프', '아') {
		t.Error("Expected result should match.")
	}

	if !isJosaAttachable('프', '는') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('프', '가') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('프', '를') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('프', '와') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('프', '야') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('프', '여') {
		t.Error("Expected result should match.")
	}
	if !isJosaAttachable('프', '라') {
		t.Error("Expected result should match.")
	}
}

func TestIsNameLengthLessThan3(t *testing.T) {
	if IsName("김") {
		t.Error("Expected result should match.")
	}
	if IsName("관진") {
		t.Error("Expected result should match.")
	}

}

func TestIsName3CharPersonName(t *testing.T) {
	if !IsName("유호현") {
		t.Error("Expected result should match.")
	}
	if !IsName("김혜진") {
		t.Error("Expected result should match.")
	}
	if IsName("개루루") {
		t.Error("Expected result should match.")
	}
	if IsName("사다리") {
		t.Error("Expected result should match.")
	}
}

func TestIsName4CharPersonName(t *testing.T) {
	if !IsName("독고영재") {
		t.Error("Expected result should match.")
	}
	if !IsName("제갈경준") {
		t.Error("Expected result should match.")
	}
	if IsName("유호현진") {
		t.Error("Expected result should match.")
	}
}

func TestIsKoreanNumber(t *testing.T) {
	if IsName("영삼") {
		t.Error("Expected result should match.")
	}
	if IsName("이정") {
		t.Error("Expected result should match.")
	}
	if IsName("조삼모사") {
		t.Error("Expected result should match.")
	}
}

func TestIsKoreanNameVariation(t *testing.T) {
	if !IsKoreanNameVariation("호혀니") {
		t.Error("Expected result should match.")
	}
	if !IsKoreanNameVariation("혜지니") {
		t.Error("Expected result should match.")
	}
	if !IsKoreanNameVariation("은벼리") {
		t.Error("Expected result should match.")
	}
	if !IsKoreanNameVariation("이오니") {
		t.Error("Expected result should match.")
	}

	if IsKoreanNameVariation("이") {
		t.Error("Expected result should match.")
	}
	if IsKoreanNameVariation("장미") {
		t.Error("Expected result should match.")
	}

	if IsKoreanNameVariation("가라찌") {
		t.Error("Expected result should match.")
	}
	if IsKoreanNameVariation("유하기") {
		t.Error("Expected result should match.")
	}
}
