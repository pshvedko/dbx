package collection

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRow(t *testing.T) {
	const N = 5
	var tx [N]Tx
	var row [N]Row[int]
	rowN := make(map[*Row[int]]int, N)
	txN := make(map[*Tx]int, N)
	for i := 0; i < N; i++ {
		rowN[&row[i]], txN[&tx[i]] = i, i
	}

	txRow(t, &row[1], &tx[1], rowN, txN, true)
	txRow(t, &row[2], &tx[2], rowN, txN, true)
	txRow(t, &row[3], &tx[3], rowN, txN, true)
	rowTx(t, &row[1], rowN, txN)
	rowTx(t, &row[2], rowN, txN)
	rowTx(t, &row[3], rowN, txN)
	require.Same(t, row[1].tx, &tx[1])
	require.Same(t, row[2].tx, &tx[2])
	require.Same(t, row[3].tx, &tx[3])
	require.Zero(t, row[1].tx.tx)
	require.Zero(t, row[2].tx.tx)
	require.Zero(t, row[3].tx.tx)

	txRow(t, &row[2], &tx[1], rowN, txN, true)
	rowTx(t, &row[1], rowN, txN)
	rowTx(t, &row[2], rowN, txN)
	rowTx(t, &row[3], rowN, txN)
	require.Same(t, row[1].tx, &tx[1])
	require.Same(t, row[2].tx, &tx[1])
	require.Same(t, row[3].tx, &tx[3])
	require.Same(t, row[1].tx.tx, &tx[2])
	require.Same(t, row[2].tx.tx, &tx[2])
	require.Zero(t, row[1].tx.tx.tx)
	require.Zero(t, row[2].tx.tx.tx)
	require.Zero(t, row[3].tx.tx)

	txRow(t, &row[3], &tx[2], rowN, txN, true)
	rowTx(t, &row[1], rowN, txN)
	rowTx(t, &row[2], rowN, txN)
	rowTx(t, &row[3], rowN, txN)
	require.Same(t, row[1].tx, &tx[1])
	require.Same(t, row[2].tx, &tx[1])
	require.Same(t, row[3].tx, &tx[2])
	require.Same(t, row[1].tx.tx, &tx[2])
	require.Same(t, row[2].tx.tx, &tx[2])
	require.Same(t, row[3].tx.tx, &tx[3])
	require.Same(t, row[1].tx.tx.tx, &tx[3])
	require.Same(t, row[2].tx.tx.tx, &tx[3])
	require.Zero(t, row[1].tx.tx.tx.tx)
	require.Zero(t, row[2].tx.tx.tx.tx)
	require.Zero(t, row[3].tx.tx.tx)

	txRow(t, &row[3], &tx[2], rowN, txN, false)
	rowTx(t, &row[1], rowN, txN)
	rowTx(t, &row[2], rowN, txN)
	rowTx(t, &row[3], rowN, txN)
	require.Same(t, row[1].tx, &tx[1])
	require.Same(t, row[2].tx, &tx[1])
	require.Same(t, row[3].tx, &tx[3])
	require.Same(t, row[1].tx.tx, &tx[2])
	require.Same(t, row[2].tx.tx, &tx[2])
	require.Zero(t, row[1].tx.tx.tx)
	require.Zero(t, row[2].tx.tx.tx)
	require.Zero(t, row[3].tx.tx)

	txRow(t, &row[2], &tx[1], rowN, txN, false)
	rowTx(t, &row[1], rowN, txN)
	rowTx(t, &row[2], rowN, txN)
	rowTx(t, &row[3], rowN, txN)
	require.Same(t, row[1].tx, &tx[1])
	require.Same(t, row[2].tx, &tx[2])
	require.Same(t, row[3].tx, &tx[3])
	require.Zero(t, row[1].tx.tx)
	require.Zero(t, row[2].tx.tx)
	require.Zero(t, row[3].tx.tx)

	txRow(t, &row[3], &tx[3], rowN, txN, false)
	txRow(t, &row[2], &tx[2], rowN, txN, false)
	txRow(t, &row[1], &tx[1], rowN, txN, false)
	rowTx(t, &row[1], rowN, txN)
	rowTx(t, &row[2], rowN, txN)
	rowTx(t, &row[3], rowN, txN)
	require.Zero(t, row[1].tx)
	require.Zero(t, row[2].tx)
	require.Zero(t, row[3].tx)
}

func txRow[T any](t *testing.T, row *Row[T], tx *Tx, rows map[*Row[T]]int, txs map[*Tx]int, ok bool) {
	t.Helper()
	if ok {
		t.Logf("row%d + tx%d = %t", rows[row], txs[tx], row.acquire(tx))
	} else {
		row.release(tx)
		t.Logf("row%d - tx%d", rows[row], txs[tx])
	}
}

func rowTx[T any](t *testing.T, row *Row[T], rows map[*Row[T]]int, txs map[*Tx]int) {
	t.Helper()
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "row%d :", rows[row])
	for tx := row.tx; tx != nil; tx = tx.tx {
		_, _ = fmt.Fprintf(&b, " tx%d", txs[tx])
	}
	t.Log(b.String())
}
