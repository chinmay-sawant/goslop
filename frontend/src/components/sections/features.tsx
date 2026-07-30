import {
  Gauge,
  ShieldAlert,
  Bot,
  GitBranch,
  FileJson2,
  Ban,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

/** Outcome-first features (not an 8-card taxonomy wall). */
const primary = [
  {
    icon: Gauge,
    title: 'Catch hot-path cost before prod',
    description:
      '239 PERF rules: allocations in loops, missing HTTP timeouts, framework footguns, and other runtime taxes.',
    badge: '239 PERF',
    tone: 'warning' as const,
  },
  {
    icon: ShieldAlert,
    title: 'Security signal that maps to CWE',
    description:
      '175 structural CWE heuristics plus experimental taint for path traversal, command injection, XSS, SQLi.',
    badge: '175 CWE',
    tone: 'danger' as const,
  },
  {
    icon: Bot,
    title: 'Exports agents can act on',
    description:
      'Whole enclosing function by default. Per-finding refs and 25-wide chunks for parallel agent work.',
    badge: 'Export',
    tone: 'info' as const,
  },
]

const secondary = [
  {
    icon: GitBranch,
    title: 'CI gates with profiles',
    description:
      'recommended, perf, security, style, all — curated packs and fail policies.',
  },
  {
    icon: FileJson2,
    title: 'Text · JSON · SARIF 2.1.0',
    description:
      'Human terminal output, machine JSON, GitHub Code Scanning-ready SARIF. Summary on stderr.',
  },
  {
    icon: Ban,
    title: 'Debt you can control',
    description:
      'Baseline known findings, path ignores, and // goslop-ignore for reviewed exceptions.',
  },
]

export function FeaturesSection() {
  return (
    <section id="features" className="border-b border-border bg-card py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
            Outcomes, not a catalogue dump
          </h2>
          <p className="mt-4 text-muted-foreground text-balance">
            Pure Go parsing (go/parser + go/ast, no CGO). Ship as a local linter,
            a CI gate, or a triage pipeline for agent remediations.
          </p>
        </div>

        <div className="mt-12 grid gap-4 md:grid-cols-3">
          {primary.map((f) => (
            <Card key={f.title} className="h-full">
              <CardHeader>
                <div className="mb-2 flex items-center justify-between">
                  <span className="flex size-10 items-center justify-center rounded-lg border border-border bg-secondary/60">
                    <f.icon className="size-4" strokeWidth={1.75} />
                  </span>
                  <Badge variant={f.tone}>{f.badge}</Badge>
                </div>
                <CardTitle className="font-sans text-base font-semibold">
                  {f.title}
                </CardTitle>
                <CardDescription>{f.description}</CardDescription>
              </CardHeader>
            </Card>
          ))}
        </div>

        <div className="mt-4 grid gap-4 sm:grid-cols-3">
          {secondary.map((f) => (
            <div
              key={f.title}
              className="rounded-xl border border-border bg-background px-5 py-5"
            >
              <div className="mb-3 flex size-9 items-center justify-center rounded-md border border-border bg-secondary/50">
                <f.icon className="size-4" strokeWidth={1.75} />
              </div>
              <h3 className="text-sm font-semibold">{f.title}</h3>
              <p className="mt-1.5 text-sm text-muted-foreground leading-relaxed">
                {f.description}
              </p>
            </div>
          ))}
        </div>

        <div className="mt-10 flex flex-wrap justify-center gap-2">
          {[
            '239 PERF',
            '175 CWE',
            'BP hygiene',
            '5 profiles',
            'SARIF 2.1.0',
            'incremental cache',
          ].map((chip) => (
            <Badge key={chip} variant="outline" className="font-mono text-[11px]">
              {chip}
            </Badge>
          ))}
        </div>
      </div>
    </section>
  )
}
