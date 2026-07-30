import { Badge } from '@/components/ui/badge'
import { Reveal } from '@/components/reveal'

/** High-signal teaser rules (not full catalogue). */
const RULES = [
  {
    id: 'PERF-42',
    title: 'fmt.Errorf without verbs',
    why: 'Static fmt.Errorf allocates like Sprintf; use errors.New.',
    family: 'perf' as const,
  },
  {
    id: 'PERF-68',
    title: 'gin.Logger on hot path',
    why: 'Synchronous request logging serializes I/O under load.',
    family: 'perf' as const,
  },
  {
    id: 'PERF-6',
    title: 'HTTP client without Timeout',
    why: 'Missing Timeout can hang goroutines on slow upstreams.',
    family: 'perf' as const,
  },
  {
    id: 'PERF-41',
    title: 'stdlib log on request path',
    why: 'Unstructured log on hot paths; prefer leveled slog/zap.',
    family: 'perf' as const,
  },
  {
    id: 'CWE-89',
    title: 'SQL injection shape',
    why: 'String-built queries into DB sinks (taint optional).',
    family: 'sec' as const,
  },
  {
    id: 'CWE-78',
    title: 'Command injection shape',
    why: 'Untrusted data into exec/command sinks.',
    family: 'sec' as const,
  },
  {
    id: 'CWE-22',
    title: 'Path traversal shape',
    why: 'User input concatenated into filesystem paths.',
    family: 'sec' as const,
  },
  {
    id: 'CWE-79',
    title: 'XSS shape',
    why: 'Unescaped data into HTML/template sinks.',
    family: 'sec' as const,
  },
  {
    id: 'CWE-497',
    title: 'Environment disclosure',
    why: 'Diagnostics endpoints leaking host details.',
    family: 'sec' as const,
  },
  {
    id: 'BP-70',
    title: 'Log and continue',
    why: 'Errors logged without return/handle path.',
    family: 'bp' as const,
  },
  {
    id: 'BP-47',
    title: 'Missing graceful shutdown',
    why: 'Servers that never drain in-flight work.',
    family: 'bp' as const,
  },
  {
    id: 'BP-48',
    title: 'os.Exit in library code',
    why: 'Libraries should return errors, not kill the process.',
    family: 'bp' as const,
  },
]

const tone = {
  perf: 'warning' as const,
  sec: 'danger' as const,
  bp: 'info' as const,
}

export function RulesSection() {
  return (
    <section id="rules" className="border-b border-border py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <Reveal>
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
              What it catches
            </h2>
            <p className="mt-4 text-muted-foreground text-balance">
              Teaser of high-signal rules. Full catalogue: 239 PERF, 175 CWE,
              plus BP hygiene. List live with{' '}
              <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
                ./bin/goslop --list-rules
              </code>
              .
            </p>
          </div>
        </Reveal>

        <Reveal>
          <div className="mt-12 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {RULES.map((rule) => (
              <div
                key={rule.id}
                className="rounded-xl border border-border bg-card px-4 py-4 shadow-xs"
              >
                <div className="flex items-center justify-between gap-2">
                  <code className="font-mono text-[12px] font-medium text-foreground">
                    {rule.id}
                  </code>
                  <Badge variant={tone[rule.family]} className="font-mono text-[10px]">
                    {rule.family}
                  </Badge>
                </div>
                <p className="mt-2 text-sm font-medium">{rule.title}</p>
                <p className="mt-1 text-sm text-muted-foreground leading-relaxed">
                  {rule.why}
                </p>
              </div>
            ))}
          </div>
        </Reveal>
      </div>
    </section>
  )
}
