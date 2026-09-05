interface PlaybackControlsProps {
  playing: boolean;
  onTogglePlaying: () => void;
  onRestart: () => void;
  time: number;
  duration: number;
}

/** Play/pause + restart for a Signals-mode run's playback (feature 010, FR-006). */
export function PlaybackControls({ playing, onTogglePlaying, onRestart, time, duration }: PlaybackControlsProps) {
  return (
    <div style={{ display: "flex", alignItems: "center", gap: "var(--space-3)" }}>
      <button type="button" onClick={onTogglePlaying}>
        {playing ? "Pause" : "Play"}
      </button>
      <button type="button" onClick={onRestart}>
        Restart
      </button>
      <span className="text-muted">
        t = {time.toFixed(1)}s / {duration.toFixed(0)}s
      </span>
    </div>
  );
}
