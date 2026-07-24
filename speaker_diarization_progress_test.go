package sherpa_onnx

import "testing"
import "context"

func TestInvokeOfflineSpeakerDiarizationProgressCallback(t *testing.T) {
	var processed, total int
	invokeOfflineSpeakerDiarizationProgressCallback(
		offlineSpeakerDiarizationCallbackState{
			ctx: context.Background(),
			cb: func(gotProcessed, gotTotal int) {
				processed = gotProcessed
				total = gotTotal
			},
		},
		3,
		10,
	)

	if processed != 3 || total != 10 {
		t.Fatalf("progress = %d/%d, want 3/10", processed, total)
	}
}

func TestInvokeOfflineSpeakerDiarizationProgressCallbackRecoversPanic(
	t *testing.T,
) {
	invokeOfflineSpeakerDiarizationProgressCallback(
		offlineSpeakerDiarizationCallbackState{
			ctx: context.Background(),
			cb: func(_, _ int) {
				panic("callback failed")
			},
		},
		1,
		2,
	)
}

func TestInvokeOfflineSpeakerDiarizationProgressCallbackStopsOnCancellation(
	t *testing.T,
) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if got := invokeOfflineSpeakerDiarizationProgressCallback(
		offlineSpeakerDiarizationCallbackState{ctx: ctx},
		1,
		2,
	); got != -1 {
		t.Fatalf("callback result = %d, want -1", got)
	}
}
