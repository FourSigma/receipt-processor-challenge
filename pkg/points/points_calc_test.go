package points

import (
	"testing"
	"time"

	"github.com/FourSigma/receipt-processor-challenge/pkg/models"
)

func TestRules(t *testing.T) {
	tests := []struct {
		handler RuleHandlerFn
		name    string
		input   models.Receipt
		want    int64
	}{
		{
			handler: RuleAlphanumeric,
			name:    "RuleAlphanumeric: should return 0 when retailer name is empty",
			input:   models.Receipt{Retailer: ""},
			want:    0,
		},
		{
			handler: RuleAlphanumeric,
			name:    "RuleAlphanumeric: should ignore leading and trailing spaces",
			input:   models.Receipt{Retailer: " Target1 "},
			want:    7,
		},
		{
			handler: RuleAlphanumeric,
			name:    "RuleAlphanumeric: should ignore non-alphanumeric characters",
			input:   models.Receipt{Retailer: "Target&$^#/[] "},
			want:    6,
		},
		{
			handler: RuleRoundDollar,
			name:    "RuleRoundDollar: should return 0",
			input:   models.Receipt{Total: 34.45},
			want:    0,
		},
		{
			handler: RuleRoundDollar,
			name:    "RuleRoundDollar: should return 0",
			input:   models.Receipt{Total: 34.00},
			want:    50,
		},
		{
			handler: RuleMultipleOfQuarter,
			name:    "RuleMultipleOfQuarter: should return 0",
			input:   models.Receipt{Total: 34.22},
			want:    0,
		},
		{
			handler: RuleMultipleOfQuarter,
			name:    "RuleMultipleOfQuarter: should return 25",
			input:   models.Receipt{Total: 34.50},
			want:    25,
		},
		{
			handler: RuleMultipleOfQuarter,
			name:    "RuleMultipleOfQuarter: should return 25",
			input:   models.Receipt{Total: 34.75},
			want:    25,
		},
		{
			handler: RuleMultipleOfQuarter,
			name:    "RuleMultipleOfQuarter: should return 25",
			input:   models.Receipt{Total: 34.00},
			want:    25,
		},
		{
			handler: RuleItemPair,
			name:    "RuleItemPair: should return 0",
			input: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "test", Price: 10},
				},
			},
			want: 0,
		},
		{
			handler: RuleItemPair,
			name:    "RuleItemPair: should return 5 because items contain 1 pair",
			input: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "test", Price: 10},
					{ShortDescription: "test", Price: 10},
				},
			},
			want: 5,
		},
		{
			handler: RuleItemPair,
			name:    "RuleItemPair: should return 5 because items contain 1 pair",
			input: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "test", Price: 10},
					{ShortDescription: "test", Price: 10},
					{ShortDescription: "test", Price: 10},
				},
			},
			want: 5,
		},
		{
			handler: RuleItemDescription,
			name:    "RuleItemDescription: should return 0 because len(test) is not a multiple of 3",
			input: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "test", Price: 10},
					{ShortDescription: "test", Price: 10},
					{ShortDescription: "test", Price: 10},
				},
			},
			want: 0,
		},
		{
			handler: RuleItemDescription,
			name:    "RuleItemDescription: should return 5 because len(tes) is a multiple of 3",
			input: models.Receipt{
				Items: []models.Item{
					{ShortDescription: "tes", Price: 10},
					{ShortDescription: "tes", Price: 10},
					{ShortDescription: "test", Price: 10}, // Not a multiple of 3
				},
			},
			want: 4,
		},
		{
			handler: RuleOddDay,
			name:    "RuleOddDay: should return 0 for an even day",
			input: models.Receipt{
				PurchasedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			want: 0,
		},
		{
			handler: RuleOddDay,
			name:    "RuleOddDay: should return 6 for an odd day",
			input: models.Receipt{
				PurchasedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			want: 6,
		},
		{
			handler: RuleTimeOfPurchase,
			name:    "RuleTimeOfPurchase: should return 0 if time is before 2PM",
			input: models.Receipt{
				PurchasedAt: MustParseTime(t, "1:01PM"),
			},
			want: 0,
		},
		{
			handler: RuleTimeOfPurchase,
			name:    "RuleTimeOfPurchase: should return 10 if time is after 2PM",
			input: models.Receipt{
				PurchasedAt: MustParseTime(t, "2:01PM"),
			},
			want: 10,
		},
		{
			handler: RuleTimeOfPurchase,
			name:    "RuleTimeOfPurchase: should return 10 if time is before 4PM",
			input: models.Receipt{
				PurchasedAt: MustParseTime(t, "3:59PM"),
			},
			want: 10,
		},
		{
			handler: RuleTimeOfPurchase,
			name:    "RuleTimeOfPurchase: should return 0 if time is after 4PM",
			input: models.Receipt{
				PurchasedAt: MustParseTime(t, "4:01PM"),
			},
			want: 0,
		},
		{
			handler: RuleTimeOfPurchase,
			name:    "RuleTimeOfPurchase: should return 0 if time is exactly 2PM",
			input: models.Receipt{
				PurchasedAt: MustParseTime(t, "2:00PM"),
			},
			want: 0,
		},

		{
			handler: RuleTimeOfPurchase,
			name:    "RuleTimeOfPurchase: should return 0 if time is exactly 4PM",
			input: models.Receipt{
				PurchasedAt: MustParseTime(t, "4:00PM"),
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.handler(tt.input)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func MustParseTime(t *testing.T, s string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.Kitchen, s)
	if err != nil {
		t.Fatalf("could not parse time: %v", err)
	}

	return parsed
}
