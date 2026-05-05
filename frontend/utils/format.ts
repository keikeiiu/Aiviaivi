export function formatCount(n: number): string {
  if (n >= 1_0000_0000) return (n / 1_0000_0000).toFixed(1) + "B";
  if (n >= 1_0000) return (n / 1_0000).toFixed(1) + "W";
  if (n >= 1000) return (n / 1000).toFixed(1) + "K";
  return String(n);
}

export function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);
  if (h > 0) return `${h}:${pad(m)}:${pad(s)}`;
  return `${m}:${pad(s)}`;
}

export function formatTimeAgo(dateStr: string): string {
  const now = Date.now();
  const then = new Date(dateStr).getTime();
  const diff = now - then;
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);
  if (minutes < 1) return "just now";
  if (minutes < 60) return `${minutes}m ago`;
  if (hours < 24) return `${hours}h ago`;
  if (days < 30) return `${days}d ago`;
  return new Date(dateStr).toLocaleDateString();
}

function pad(n: number): string {
  return n.toString().padStart(2, "0");
}
