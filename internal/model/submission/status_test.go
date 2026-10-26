package submission_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model/submission"
)

func TestNewStatus(t *testing.T) {
	tests := []struct {
		name      string
		statusStr string
		want      submission.Status
		wantErr   error
	}{
		{
			name:      "success status PENDING",
			statusStr: "PENDING",
			want:      submission.StatusPending,
			wantErr:   nil,
		},
		{
			name:      "success status RUNNING",
			statusStr: "RUNNING",
			want:      submission.StatusRunning,
			wantErr:   nil,
		},
		{
			name:      "success status ACCEPTED",
			statusStr: "ACCEPTED",
			want:      submission.StatusAccepted,
			wantErr:   nil,
		},
		{
			name:      "success status WRONG_ANSWER",
			statusStr: "WRONG_ANSWER",
			want:      submission.StatusWrongAnswer,
			wantErr:   nil,
		},
		{
			name:      "success status TIMEOUT",
			statusStr: "TIMEOUT",
			want:      submission.StatusTimeout,
			wantErr:   nil,
		},
		{
			name:      "success status INTERNAL",
			statusStr: "INTERNAL",
			want:      submission.StatusInternal,
			wantErr:   nil,
		},
		{
			name:      "invalid status",
			statusStr: "XXX",
			want:      "",
			wantErr:   errs.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := submission.NewStatus(tt.statusStr)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Empty(t, status)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, status)
		})
	}
}
