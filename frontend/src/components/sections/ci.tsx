import { useMemo, useState } from 'react'
import { Badge } from '@/components/ui/badge'
import { CodeBlock } from '@/components/code-block'
import { cn } from '@/lib/utils'

const personas = [
  {
    id: 'ci',
    label: 'Everyday CI',
    profile: 'recommended',
    fail: 'high',
    taint: false,
    blurb: 'S-tier PERF + core security IDs. Default fail on high.',
    cmd: './bin/goslop --profile recommended .',
  },
  {
    id: 'sec',
    label: 'Security review',
    profile: 'security',
    fail: 'high',
    taint: true,
    blurb: 'CWE pack with taint on (depth 4) for injection-class issues.',
    cmd: './bin/goslop --profile security ./cmd\n./bin/goslop --format sarif . > goslop.sarif',
  },
  {
    id: 'perf',
    label: 'Perf pass',
    profile: 'perf',
    fail: 'high',
    taint: false,
    blurb: 'S + A tier performance catalogue for hot-path audits.',
    cmd: './bin/goslop --profile perf .',
  },
  {
    id: 'style',
    label: 'Hygiene only',
    profile: 'style',
    fail: 'none',
    taint: false,
    blurb: 'BP-* soft gate. Good for gradual style cleanup without failing CI.',
    cmd: './bin/goslop --profile style .',
  },
] as const

export function CiSection() {
  const [id, setId] = useState<(typeof personas)[number]['id']>('ci')
  const active = useMemo(
    () => personas.find((p) => p.id === id) ?? personas[0],
    [id],
  )

  return (
    <section id="ci" className="border-b border-border py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <div className="grid gap-12 lg:grid-cols-[1fr_1.05fr] lg:items-start">
          <div>
            <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
              Pick a job, get a command
            </h2>
            <p className="mt-4 text-muted-foreground leading-relaxed">
              Profiles are curated packs with fail policies. Choose what you are
              doing; copy the command. No need to memorize allow-lists on day one.
            </p>

            <div
              className="mt-8 flex flex-col gap-2"
              role="tablist"
              aria-label="CI persona"
            >
              {personas.map((p) => (
                <button
                  key={p.id}
                  type="button"
                  role="tab"
                  aria-selected={id === p.id}
                  onClick={() => setId(p.id)}
                  className={cn(
                    'rounded-lg border px-4 py-3 text-left transition-colors',
                    id === p.id
                      ? 'border-foreground/20 bg-card shadow-xs'
                      : 'border-border bg-background hover:bg-secondary/40',
                  )}
                >
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="text-sm font-medium">{p.label}</span>
                    <Badge variant="outline" className="font-mono text-[10px]">
                      {p.profile}
                    </Badge>
                  </div>
                  <p className="mt-1 text-sm text-muted-foreground">{p.blurb}</p>
                </button>
              ))}
            </div>
          </div>

          <div className="space-y-4">
            <div className="flex flex-wrap gap-2">
              <Badge variant="outline" className="font-mono">
                --profile {active.profile}
              </Badge>
              <Badge variant={active.taint ? 'info' : 'muted'}>
                taint {active.taint ? 'on' : 'off'}
              </Badge>
              <Badge variant="outline" className="font-mono">
                fail · {active.fail}
              </Badge>
            </div>

            <CodeBlock code={active.cmd} filename={`profile: ${active.profile}`} />

            <div className="rounded-xl border border-border bg-secondary/30 px-5 py-4 text-sm text-muted-foreground">
              <p className="font-medium text-foreground">Debt on day one</p>
              <p className="mt-1.5 leading-relaxed">
                Ship with{' '}
                <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[11px] text-foreground">
                  .goslop-baseline.json
                </code>{' '}
                so known findings do not fail the first rollout. Narrow
                exceptions:{' '}
                <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[11px] text-foreground">
                  // goslop-ignore
                </code>
                . Incremental re-scans use{' '}
                <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[11px] text-foreground">
                  .goslop-cache/
                </code>
                .
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
