package rules

import (
	"testing"
	"time"

	"github.com/ankurO7/compilance-monitor/internal/models"
)

func TestWatchlistRule_FlagsBlockedCounterparty(t *testing.T) {
	rule := NewWatchlistRule([]string{"blocked-entity-1"})

	tx := models.Transaction{
		ID:           "tx1",
		UserID:       "test-user",
		Amount:       500,
		Counterparty: "blocked-entity-1",
		CreatedAt:    time.Now(),
	}

	result := rule.Evaluate(tx, nil)

	if !result.Flagged {
		t.Errorf("expected transaction to blocked-entity-1 to be flagged, but it wasn't")
	}
}

func TestWatchlistRule_DoesNotFlagCleanCounterparty(t *testing.T) {
	rule := NewWatchlistRule([]string{"blocked-entity-1"})

	tx := models.Transaction{
		ID:           "tx2",
		UserID:       "test-user",
		Amount:       500,
		Counterparty: "corner-store",
		CreatedAt:    time.Now(),
	}

	result := rule.Evaluate(tx, nil)

	if result.Flagged {
		t.Errorf("expected transaction to corner-store to NOT be flagged, but it was")
	}
}

func TestStructuringRule_FlagsMultipleNearThresholdTransactions(t *testing.T) {
	rule := NewStructuringRule() // Threshold: 10000, NearBand: 0.9, MinMatches: 3

	now := time.Now()
	history := []models.Transaction{
		{ID: "h1", UserID: "u1", Amount: 9200, CreatedAt: now.Add(-1 * time.Hour)},
		{ID: "h2", UserID: "u1", Amount: 9400, CreatedAt: now.Add(-2 * time.Hour)},
	}
	// current transaction is the 3rd near-threshold transaction -> should trigger MinMatches: 3
	tx := models.Transaction{ID: "tx3", UserID: "u1", Amount: 9600, CreatedAt: now}

	result := rule.Evaluate(tx, history)

	if !result.Flagged {
		t.Errorf("expected structuring to be flagged with 3 near-threshold transactions, but it wasn't")
	}
}

func TestStructuringRule_DoesNotFlagBelowMinMatches(t *testing.T) {
	rule := NewStructuringRule()

	now := time.Now()
	history := []models.Transaction{
		{ID: "h1", UserID: "u1", Amount: 9200, CreatedAt: now.Add(-1 * time.Hour)},
	}
	// only 2 near-threshold transactions total -> below MinMatches: 3
	tx := models.Transaction{ID: "tx2", UserID: "u1", Amount: 9400, CreatedAt: now}

	result := rule.Evaluate(tx, history)

	if result.Flagged {
		t.Errorf("expected structuring NOT to be flagged with only 2 near-threshold transactions, but it was")
	}
}

func TestVelocityRule_FlagsClearOutlier(t *testing.T) {
	rule := NewVelocityRule() // MinHistory: 5, ZThreshold: 3.0

	now := time.Now()
	var history []models.Transaction
	// 8 stable transactions around $100
	for i := 0; i < 8; i++ {
		history = append(history, models.Transaction{
			ID: "h", UserID: "u1", Amount: 100, CreatedAt: now.Add(-time.Duration(i+1) * time.Hour),
		})
	}
	// wildly larger transaction, should be flagged as an outlier
	tx := models.Transaction{ID: "tx-outlier", UserID: "u1", Amount: 15000, CreatedAt: now}

	result := rule.Evaluate(tx, history)

	if !result.Flagged {
		t.Errorf("expected a $15000 transaction against a $100 baseline to be flagged as a velocity anomaly, but it wasn't")
	}
}

func TestVelocityRule_DoesNotFlagWithoutEnoughHistory(t *testing.T) {
	rule := NewVelocityRule() // MinHistory: 5

	now := time.Now()
	history := []models.Transaction{
		{ID: "h1", UserID: "u1", Amount: 100, CreatedAt: now.Add(-1 * time.Hour)},
	}
	// only 1 prior transaction, below MinHistory: 5 -> rule should not fire, even on a big jump
	tx := models.Transaction{ID: "tx-early", UserID: "u1", Amount: 15000, CreatedAt: now}

	result := rule.Evaluate(tx, history)

	if result.Flagged {
		t.Errorf("expected rule NOT to fire with insufficient history, but it did")
	}
}