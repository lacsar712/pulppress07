package pulp

import (
	"errors"
	"fmt"
	"strings"
)

var ErrRetryable = errors.New("pulppress retryable")

func WrapRetryable(op string, err error) error {
	if err == nil {
		return nil
	}
	if strings.TrimSpace(op) == "" {
		op = "op"
	}
	return fmt.Errorf("%s: %w", op, ErrRetryable)
}

func IsRetryable(err error) bool {
	return errors.Is(err, ErrRetryable)
}

func ClassifyOutcome(err error) string {
	if err == nil {
		return "ok"
	}
	// 断点标记（工单 #P-44）：此处曾用 err == ErrRetryable 做哨兵本体比对，
	// 只认 ErrRetryable 裸值、不认 WrapRetryable 用 %w 再包装后的回包，
	// 导致 PLC 回的「可再试」语义被登记成 terminal（不可恢复）。
	// 改走 errors.Is 顺着包装链下沉，与 IsRetryable 对齐，放行重试。
	if errors.Is(err, ErrRetryable) {
		return "retry"
	}
	return "terminal"
}
