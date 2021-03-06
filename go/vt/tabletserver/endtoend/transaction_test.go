// Copyright 2015, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package endtoend

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/youtube/vitess/go/vt/tabletserver/endtoend/framework"
)

func TestCommit(t *testing.T) {
	client := framework.NewClient()
	defer client.Execute("delete from vitess_test where intval=4", nil)

	catcher := framework.NewTxCatcher()
	defer catcher.Close()
	vstart := framework.DebugVars()

	query := "insert into vitess_test (intval, floatval, charval, binval) " +
		"values(4, null, null, null)"
	err := client.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	_, err = client.Execute(query, nil)
	if err != nil {
		t.Error(err)
		return
	}
	err = client.Commit()
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := catcher.Next()
	if err != nil {
		t.Error(err)
		return
	}
	want := []string{"insert into vitess_test(intval, floatval, charval, binval) values (4, null, null, null) /* _stream vitess_test (intval ) (4 ); */"}
	if !reflect.DeepEqual(tx.Queries, want) {
		t.Errorf("queries: %v, want %v", tx.Queries, want)
	}
	if !reflect.DeepEqual(tx.Conclusion, "commit") {
		t.Errorf("conclusion: %s, want commit", tx.Conclusion)
	}

	qr, err := client.Execute("select * from vitess_test", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if qr.RowsAffected != 4 {
		t.Errorf("rows affected: %d, want 4", qr.RowsAffected)
	}

	_, err = client.Execute("delete from vitess_test where intval=4", nil)
	if err != nil {
		t.Error(err)
		return
	}

	qr, err = client.Execute("select * from vitess_test", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if qr.RowsAffected != 3 {
		t.Errorf("rows affected: %d, want 4", qr.RowsAffected)
	}

	expectedDiffs := []struct {
		tag  string
		diff int
	}{{
		tag:  "Transactions/TotalCount",
		diff: 2,
	}, {
		tag:  "Transactions/Histograms/Completed/Count",
		diff: 2,
	}, {
		tag:  "Queries/TotalCount",
		diff: 6,
	}, {
		tag:  "Queries/Histograms/BEGIN/Count",
		diff: 1,
	}, {
		tag:  "Queries/Histograms/COMMIT/Count",
		diff: 1,
	}, {
		tag:  "Queries/Histograms/INSERT_PK/Count",
		diff: 1,
	}, {
		tag:  "Queries/Histograms/DML_PK/Count",
		diff: 1,
	}, {
		tag:  "Queries/Histograms/PASS_SELECT/Count",
		diff: 2,
	}}
	vend := framework.DebugVars()
	for _, expected := range expectedDiffs {
		if err := compareIntDiff(vend, expected.tag, vstart, expected.diff); err != nil {
			t.Error(err)
		}
	}
}

func TestRollback(t *testing.T) {
	client := framework.NewClient()

	catcher := framework.NewTxCatcher()
	defer catcher.Close()
	vstart := framework.DebugVars()

	query := "insert into vitess_test values(4, null, null, null)"
	err := client.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	_, err = client.Execute(query, nil)
	if err != nil {
		t.Error(err)
		return
	}
	err = client.Rollback()
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := catcher.Next()
	if err != nil {
		t.Error(err)
		return
	}
	want := []string{"insert into vitess_test values (4, null, null, null) /* _stream vitess_test (intval ) (4 ); */"}
	if !reflect.DeepEqual(tx.Queries, want) {
		t.Errorf("queries: %v, want %v", tx.Queries, want)
	}
	if !reflect.DeepEqual(tx.Conclusion, "rollback") {
		t.Errorf("conclusion: %s, want rollback", tx.Conclusion)
	}

	qr, err := client.Execute("select * from vitess_test", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if qr.RowsAffected != 3 {
		t.Errorf("rows affected: %d, want 3", qr.RowsAffected)
	}

	expectedDiffs := []struct {
		tag  string
		diff int
	}{{
		tag:  "Transactions/TotalCount",
		diff: 1,
	}, {
		tag:  "Transactions/Histograms/Aborted/Count",
		diff: 1,
	}, {
		tag:  "Queries/Histograms/BEGIN/Count",
		diff: 1,
	}, {
		tag:  "Queries/Histograms/ROLLBACK/Count",
		diff: 1,
	}, {
		tag:  "Queries/Histograms/INSERT_PK/Count",
		diff: 1,
	}}
	vend := framework.DebugVars()
	for _, expected := range expectedDiffs {
		if err := compareIntDiff(vend, expected.tag, vstart, expected.diff); err != nil {
			t.Error(err)
		}
	}
}

func TestAutoCommit(t *testing.T) {
	client := framework.NewClient()
	defer client.Execute("delete from vitess_test where intval=4", nil)

	catcher := framework.NewTxCatcher()
	defer catcher.Close()
	vstart := framework.DebugVars()

	query := "insert into vitess_test (intval, floatval, charval, binval) " +
		"values(4, null, null, null)"
	_, err := client.Execute(query, nil)
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := catcher.Next()
	if err != nil {
		t.Error(err)
		return
	}
	want := []string{"insert into vitess_test(intval, floatval, charval, binval) values (4, null, null, null) /* _stream vitess_test (intval ) (4 ); */"}
	if !reflect.DeepEqual(tx.Queries, want) {
		t.Errorf("queries: %v, want %v", tx.Queries, want)
	}
	if !reflect.DeepEqual(tx.Conclusion, "commit") {
		t.Errorf("conclusion: %s, want commit", tx.Conclusion)
	}

	qr, err := client.Execute("select * from vitess_test", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if qr.RowsAffected != 4 {
		t.Errorf("rows affected: %d, want 4", qr.RowsAffected)
	}

	_, err = client.Execute("delete from vitess_test where intval=4", nil)
	if err != nil {
		t.Error(err)
		return
	}

	qr, err = client.Execute("select * from vitess_test", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if qr.RowsAffected != 3 {
		t.Errorf("rows affected: %d, want 4", qr.RowsAffected)
	}

	expectedDiffs := []struct {
		tag  string
		diff int
	}{{
		tag:  "Transactions/TotalCount",
		diff: 2,
	}, {
		tag:  "Transactions/Histograms/Completed/Count",
		diff: 2,
	}, {
		tag:  "Queries/TotalCount",
		diff: 4,
	}, {
		tag:  "Queries/Histograms/BEGIN/Count",
		diff: 0,
	}, {
		tag:  "Queries/Histograms/COMMIT/Count",
		diff: 0,
	}, {
		tag:  "Queries/Histograms/INSERT_PK/Count",
		diff: 1,
	}, {
		tag:  "Queries/Histograms/DML_PK/Count",
		diff: 1,
	}, {
		tag:  "Queries/Histograms/PASS_SELECT/Count",
		diff: 2,
	}}
	vend := framework.DebugVars()
	for _, expected := range expectedDiffs {
		if err := compareIntDiff(vend, expected.tag, vstart, expected.diff); err != nil {
			t.Error(err)
		}
	}
}

func TestAutoCommitOff(t *testing.T) {
	framework.Server.SetAutoCommit(false)
	defer framework.Server.SetAutoCommit(true)

	_, err := framework.NewClient().Execute("insert into vitess_test values(4, null, null, null)", nil)
	want := "error: unsupported query"
	if err == nil || !strings.HasPrefix(err.Error(), want) {
		t.Errorf("Error: %v, must start with %s", err, want)
	}
}

func TestTxPoolSize(t *testing.T) {
	vstart := framework.DebugVars()

	client1 := framework.NewClient()
	err := client1.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer client1.Rollback()
	if err := verifyIntValue(framework.DebugVars(), "TransactionPoolAvailable", framework.BaseConfig.TransactionCap-1); err != nil {
		t.Error(err)
	}

	defer framework.Server.SetTxPoolSize(framework.Server.TxPoolSize())
	framework.Server.SetTxPoolSize(1)
	defer framework.Server.BeginTimeout.Set(framework.Server.BeginTimeout.Get())
	timeout := 1 * time.Millisecond
	framework.Server.BeginTimeout.Set(timeout)
	vend := framework.DebugVars()
	if err := verifyIntValue(vend, "TransactionPoolAvailable", 0); err != nil {
		t.Error(err)
	}
	if err := verifyIntValue(vend, "TransactionPoolCapacity", 1); err != nil {
		t.Error(err)
	}
	if err := verifyIntValue(vend, "BeginTimeout", int(timeout)); err != nil {
		t.Error(err)
	}

	client2 := framework.NewClient()
	err = client2.Begin()
	want := "tx_pool_full"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("Error: %v, must contain %s", err, want)
	}
	if err := compareIntDiff(framework.DebugVars(), "Errors/TxPoolFull", vstart, 1); err != nil {
		t.Error(err)
	}
}

func TestTxTimeout(t *testing.T) {
	vstart := framework.DebugVars()

	defer framework.Server.SetTxTimeout(framework.Server.TxTimeout())
	framework.Server.SetTxTimeout(1 * time.Millisecond)
	if err := verifyIntValue(framework.DebugVars(), "TransactionPoolTimeout", int(1*time.Millisecond)); err != nil {
		t.Error(err)
	}

	catcher := framework.NewTxCatcher()
	defer catcher.Close()
	client := framework.NewClient()
	err := client.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := catcher.Next()
	if err != nil {
		t.Error(err)
		return
	}
	if tx.Conclusion != "kill" {
		t.Errorf("Conclusion: %s, want kill", tx.Conclusion)
	}
	if err := compareIntDiff(framework.DebugVars(), "Kills/Transactions", vstart, 1); err != nil {
		t.Error(err)
	}

	// Ensure commit fails.
	err = client.Commit()
	want := "not_in_tx: Transaction"
	if err == nil || !strings.HasPrefix(err.Error(), want) {
		t.Errorf("Error: %v, must contain %s", err, want)
	}
}

func TestForUpdate(t *testing.T) {
	for _, mode := range []string{"for update", "lock in share mode"} {
		client := framework.NewClient()
		query := fmt.Sprintf("select * from vitess_test where intval=2 %s", mode)
		_, err := client.Execute(query, nil)
		want := "error: Disallowed"
		if err == nil || !strings.HasPrefix(err.Error(), want) {
			t.Errorf("Error: %v, must have prefix %s", err, want)
		}

		// We should not get errors here
		err = client.Begin()
		if err != nil {
			t.Error(err)
			return
		}
		_, err = client.Execute(query, nil)
		if err != nil {
			t.Error(err)
			return
		}
		err = client.Commit()
		if err != nil {
			t.Error(err)
			return
		}
	}
}
