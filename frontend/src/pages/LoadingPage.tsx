import './LoadingPage.css';

/**
 * Shown while App.tsx is checking whether a saved session token is still valid.
 * Reuses the same spinning-dot ring animation as the Login page for visual consistency.
 */

const SPANS = Array.from({ length: 50 }, (_, i) => i + 1);

export default function LoadingPage() {
  return (
    <div className="ldp-body">
      <div className="ldp-container">
        {SPANS.map((i) => (
          <span key={i} style={{ '--i': i } as React.CSSProperties} />
        ))}
        <p className="ldp-text">Checking session…</p>
      </div>
    </div>
  );
}
