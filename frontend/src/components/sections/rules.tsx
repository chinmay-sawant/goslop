import { ArrowUpRight } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Reveal } from '@/components/reveal'

const REPO = 'https://github.com/chinmay-sawant/goslop/blob/main'

type Rule = {
  id: string
  title: string
  why: string
  family: 'perf' | 'sec' | 'bp'
  /** documents/ path when prose docs exist; otherwise omit. */
  doc?: string
  /** Authoritative ruleset JSON path under ruleset/. */
  ruleset: string
}

/**
 * High-signal teaser rules. Docs only when documents/ already covers them;
 * otherwise the card links to ruleset/ metadata (no new docs fabricated).
 */
const RULES: Rule[] = [
  {
    id: 'PERF-42',
    title: 'fmt.Errorf without verbs',
    why: 'Static fmt.Errorf allocates like Sprintf; use errors.New.',
    family: 'perf',
    ruleset: 'ruleset/golang/chunks/perf-001-050.json',
  },
  {
    id: 'PERF-68',
    title: 'gin.Logger on hot path',
    why: 'Synchronous request logging serializes I/O under load.',
    family: 'perf',
    ruleset: 'ruleset/golang/chunks/perf-051-100.json',
  },
  {
    id: 'PERF-6',
    title: 'fmt formatting inside loop',
    why: 'fmt.Sprintf / Fprintf in a loop body allocates on every iteration.',
    family: 'perf',
    ruleset: 'ruleset/golang/chunks/perf-001-050.json',
  },
  {
    id: 'PERF-41',
    title: 'stdlib log on request path',
    why: 'Unstructured log on hot paths; prefer leveled slog/zap.',
    family: 'perf',
    ruleset: 'ruleset/golang/chunks/perf-001-050.json',
  },
  {
    id: 'CWE-89',
    title: 'SQL injection shape',
    why: 'String-built queries into DB sinks (taint optional).',
    family: 'sec',
    doc: 'documents/taint.md',
    ruleset: 'ruleset/golang/chunks/cwe-051-100.json',
  },
  {
    id: 'CWE-78',
    title: 'Command injection shape',
    why: 'Untrusted data into exec/command sinks.',
    family: 'sec',
    doc: 'documents/taint.md',
    ruleset: 'ruleset/golang/chunks/cwe-051-100.json',
  },
  {
    id: 'CWE-22',
    title: 'Path traversal shape',
    why: 'User input concatenated into filesystem paths.',
    family: 'sec',
    doc: 'documents/taint.md',
    ruleset: 'ruleset/golang/chunks/cwe-001-050.json',
  },
  {
    id: 'CWE-79',
    title: 'XSS shape',
    why: 'Unescaped data into HTML/template sinks.',
    family: 'sec',
    doc: 'documents/taint.md',
    ruleset: 'ruleset/golang/chunks/cwe-051-100.json',
  },
  {
    id: 'CWE-497',
    title: 'Environment disclosure',
    why: 'Diagnostics endpoints leaking host details.',
    family: 'sec',
    ruleset: 'ruleset/golang/chunks/cwe-201-9999.json',
  },
  {
    id: 'BP-70',
    title: 'Log and continue',
    why: 'Errors logged without return/handle path.',
    family: 'bp',
    ruleset: 'ruleset/golang/bad-practices.json',
  },
  {
    id: 'BP-47',
    title: 'Missing graceful shutdown',
    why: 'Servers that never drain in-flight work.',
    family: 'bp',
    ruleset: 'ruleset/golang/bad-practices.json',
  },
  {
    id: 'BP-48',
    title: 'os.Exit in library code',
    why: 'Libraries should return errors, not kill the process.',
    family: 'bp',
    ruleset: 'ruleset/golang/bad-practices.json',
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
              High-signal teaser. Full catalogue: 239 PERF, 175 CWE, plus BP.
              Each card links to{' '}
              <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
                documents/
              </code>{' '}
              when prose exists, otherwise to the{' '}
              <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
                ruleset/
              </code>{' '}
              JSON definition. Live:{' '}
              <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
                goslop --list-rules
              </code>
              .
            </p>
          </div>
        </Reveal>

        <Reveal>
          <div className="mt-12 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {RULES.map((rule) => {
              const href = rule.doc
                ? `${REPO}/${rule.doc}`
                : `${REPO}/${rule.ruleset}`
              const pathLabel = rule.doc ?? rule.ruleset
              const linkKind = rule.doc ? 'docs' : 'ruleset'

              return (
                <a
                  key={rule.id}
                  href={href}
                  target="_blank"
                  rel="noreferrer"
                  className="group rounded-xl border border-border bg-card px-4 py-4 shadow-xs transition-colors hover:border-foreground/20"
                >
                  <div className="flex items-center justify-between gap-2">
                    <code className="font-mono text-[12px] font-medium text-foreground">
                      {rule.id}
                    </code>
                    <div className="flex items-center gap-1.5">
                      <Badge
                        variant={tone[rule.family]}
                        className="font-mono text-[10px]"
                      >
                        {rule.family}
                      </Badge>
                      <ArrowUpRight className="size-3.5 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100" />
                    </div>
                  </div>
                  <p className="mt-2 text-sm font-medium">{rule.title}</p>
                  <p className="mt-1 text-sm text-muted-foreground leading-relaxed">
                    {rule.why}
                  </p>
                  <p className="mt-3 font-mono text-[10px] leading-relaxed text-muted-foreground break-all">
                    <span className="text-foreground/70">{linkKind}: </span>
                    {pathLabel}
                  </p>
                </a>
              )
            })}
          </div>
        </Reveal>
      </div>
    </section>
  )
}
