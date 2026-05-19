//go:build android

package skirk

// Memory limits for Android environments.
// Android apps typically receive 200-500MB heap budget. These conservative
// values prevent OOM kills on low-to-mid range devices (2-4GB RAM).
// Worst-case with 4 lanes: ~4×8 + 48 + 48 = ~128MB for buffers.
const (
	muxNormalLaneQueueBytes     = 8 * 1024 * 1024  // 8MB per lane (vs 64MB on desktop)
	muxNormalStreamQueueBytes   = 2 * 1024 * 1024  // 2MB per stream
	muxNormalReceiveQueueBytes  = 8 * 1024 * 1024  // 8MB receive queue
	muxNormalReceiveGlobalBytes = 48 * 1024 * 1024 // 48MB global receive cap
	muxPendingStreamBytes       = 8 * 1024 * 1024  // 8MB per pending stream
	muxPendingGlobalBytes       = 48 * 1024 * 1024 // 48MB global pending cap
	muxStreamPendingBytes       = 8 * 1024 * 1024  // 8MB per stream pending
	// muxStreamPauseBytes must be strictly less than muxNormalLaneQueueBytes (8MB)
	// otherwise streams pause immediately on Android and never send data.
	muxStreamPauseBytes = 1 * 1024 * 1024 // 1MB pause threshold
)
