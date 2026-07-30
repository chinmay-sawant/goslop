import { Badge } from '@/components/ui/badge'

const rows = [
  {
    label: 'Typical linter',
    items: [
      'Style and simple correctness',
      'Line hits, thin context',
      'Human reads the terminal',
      'Often one fixed rule surface',
    ],
  },
  {
    label: 'goslop (SAT)',
    items: [
      'Perf + CWE + hygiene catalogues',
      'Whole-function context by default',
      'Chunks for agent delegation',
      'Profiles, SARIF, baseline, cache',
    ],
    highlight: true,
  },
]

export function WhySection() {
  return (
    <section id="why" className="border-b border-border py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
            Not another style linter
          </h2>
          <p className="mt-4 text-muted-foreground text-balance">
            Use a linter for format and simple correctness. Use goslop when you
            need hot-path and security signal, CI fail policies, and exports
            shaped for humans and agents.
          </p>
        </div>

        <div className="mx-auto mt-12 grid max-w-4xl gap-4 md:grid-cols-2">
          {rows.map((col) => (
            <div
              key={col.label}
              className={
                col.highlight
                  ? 'rounded-xl border border-foreground/20 bg-card p-6 shadow-xs'
                  : 'rounded-xl border border-border bg-background p-6'
              }
            >
              <div className="mb-4 flex items-center gap-2">
                <h3 className="text-sm font-semibold">{col.label}</h3>
                {col.highlight && (
                  <Badge className="font-mono text-[10px]">this product</Badge>
                )}
              </div>
              <ul className="space-y-3 text-sm text-muted-foreground">
                {col.items.map((item) => (
                  <li key={item} className="flex gap-2">
                    <span className="mt-2 size-1 shrink-0 rounded-full bg-foreground/40" />
                    <span className={col.highlight ? 'text-foreground/90' : ''}>
                      {item}
                    </span>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="mx-auto mt-10 max-w-3xl rounded-xl border border-border bg-secondary/30 px-5 py-4 text-sm text-muted-foreground">
          <p>
            <span className="font-medium text-foreground">Noise honesty: </span>
            scanning the goslop repo itself is a poor FP benchmark (needles,
            detectors, fixtures). Use a real app tree, then{' '}
            <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[11px] text-foreground">
              .goslop-baseline.json
            </code>{' '}
            and{' '}
            <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[11px] text-foreground">
              // goslop-ignore
            </code>{' '}
            for known debt.
          </p>
        </div>
      </div>
    </section>
  )
}
