//go:build !android

package skirk

// Memory limits for desktop/server environments.
// These are generous values suitable for machines with 4GB+ RAM.
// For Android, see mux_limits_android.go.
const (
	muxNormalLaneQueueBytes     = 64 * 1024 * 1024  // 64MB per lane
	muxNormalStreamQueueBytes   = 16 * 1024 * 1024  // 16MB per stream
	muxNormalReceiveQueueBytes  = 64 * 1024 * 1024  // 64MB receive queue
	muxNormalReceiveGlobalBytes = 256 * 1024 * 1024 // 256MB global receive cap
	muxPendingStreamBytes       = 64 * 1024 * 1024  // 64MB per pending stream
	muxPendingGlobalBytes       = 256 * 1024 * 1024 // 256MB global pending cap
	muxStreamPendingBytes       = 64 * 1024 * 1024  // 64MB per stream pending
	// muxStreamPauseBytes must be < muxNormalLaneQueueBytes so streams don't
	// get paused before they even start sending.
	muxStreamPauseBytes = 8 * 1024 * 1024 // 8MB pause threshold
)
