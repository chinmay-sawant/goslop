import {
  BookOpen,
  Terminal,
  FileOutput,
  Shield,
  ListTree,
  Waves,
  Gauge,
  Wrench,
  Package,
  LayoutDashboard,
  ArrowUpRight,
} from 'lucide-react'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'

const DOCS_BASE =
  'https://github.com/chinmay-sawant/goslop/blob/main/documents'

/** Mirrors documents/README.md start-here + deeper guides. */
const startHere = [
  {
    icon: LayoutDashboard,
    title: 'Product overview',
    file: 'overview.md',
    description:
      'Detector families, profiles, suppressions, cache, baseline, and why self-scan of this repo is noisy.',
  },
  {
    icon: Terminal,
    title: 'CLI reference',
    file: 'cli-reference.md',
    description:
      'Every flag, subcommand, exit code, profile alias, and config merge rule for ./bin/goslop.',
  },
  {
    icon: Package,
    title: 'make run & metrics',
    file: 'make-run.md',
    description:
      'Product make run, reference-metrics, bench targets, and how reference corpus differs from self-scan.',
  },
  {
    icon: FileOutput,
    title: 'Reporting formats',
    file: 'reporting-formats.md',
    description:
      'Text, JSON, and SARIF 2.1.0 with full examples for humans and CI Code Scanning.',
  },
  {
    icon: BookOpen,
    title: 'Export context & chunks',
    file: 'export-context-and-chunks.md',
    description:
      'Function refs vs batched chunks, whole-function Context (default), and agent delegation workflow.',
  },
  {
    icon: Shield,
    title: 'Suppressions & baselines',
    file: 'suppressions-and-baselines.md',
    description:
      'Inline // goslop-ignore, path ignores, .goslop-baseline.json discovery, and matching rules.',
  },
]

const deeper = [
  {
    icon: ListTree,
    title: 'Rule catalog & maturity',
    file: 'rule-catalog-and-maturity.md',
    description: '--list-rules, --explain, runtime counts, and maturity labels.',
  },
  {
    icon: Gauge,
    title: 'PERF rules notes',
    file: 'perf-rules.md',
    description: 'Human notes for many PERF-* rules (partial catalog).',
  },
  {
    icon: Waves,
    title: 'Taint engine',
    file: 'taint.md',
    description: 'Experimental inter-procedural taint for CWE-22 / 78 / 79 / 89.',
  },
  {
    icon: LayoutDashboard,
    title: 'Recommended pack',
    file: 'go-recommended-pack.md',
    description: 'Default recommended pack allow-lists and CI guidance.',
  },
  {
    icon: Wrench,
    title: 'Development',
    file: 'development.md',
    description: 'Build, validation, CI, product scan, and benchmark commands.',
  },
  {
    icon: Package,
    title: 'Architecture & performance',
    file: 'architecture-performance.md',
    description: 'Engine pipeline, cache layout, and performance design notes.',
  },
]

const quickCommands = [
  { label: 'Build (pure Go)', code: 'make build' },
  { label: 'Everyday scan', code: './bin/goslop .' },
  { label: 'List rules', code: './bin/goslop --list-rules' },
  { label: 'Explain a rule', code: './bin/goslop --explain PERF-6' },
  { label: 'Starter config', code: './bin/goslop init' },
  {
    label: 'Product-style export',
    code: 'make run SCAN_PATH=./your/project',
  },
]

export function DocsSection() {
  return (
    <section id="docs" className="border-b border-border py-24 md:py-32">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-heading text-4xl tracking-tight md:text-5xl">
            From scan to agent handoff
          </h2>
          <p className="mt-4 text-muted-foreground text-balance">
            User-facing guides from the{' '}
            <code className="font-mono text-xs bg-secondary px-1.5 py-0.5 rounded">
              documents/
            </code>{' '}
            tree: product overview, CLI, reporters, export, suppressions, and
            deeper engine notes.
          </p>
        </div>

        <div className="mt-12">
          <h3 className="mb-4 text-sm font-semibold">Start here</h3>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {startHere.map((doc) => (
              <a
                key={doc.file}
                href={`${DOCS_BASE}/${doc.file}`}
                target="_blank"
                rel="noreferrer"
                className="group"
              >
                <Card className="h-full transition-shadow duration-200 group-hover:shadow-[0_2px_8px_rgba(0,0,0,0.04)] dark:group-hover:shadow-[0_2px_8px_rgba(0,0,0,0.25)]">
                  <CardHeader>
                    <div className="mb-2 flex items-start justify-between gap-2">
                      <span className="flex size-10 items-center justify-center rounded-lg border border-border bg-secondary/60">
                        <doc.icon className="size-4 text-foreground" strokeWidth={1.75} />
                      </span>
                      <ArrowUpRight className="size-4 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100" />
                    </div>
                    <CardTitle className="font-sans text-base font-semibold">
                      {doc.title}
                    </CardTitle>
                    <CardDescription>{doc.description}</CardDescription>
                    <p className="mt-2 font-mono text-[11px] text-muted-foreground">
                      documents/{doc.file}
                    </p>
                  </CardHeader>
                </Card>
              </a>
            ))}
          </div>
        </div>

        <div className="mt-14">
          <h3 className="mb-4 text-sm font-semibold">Deeper guides</h3>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {deeper.map((doc) => (
              <a
                key={doc.file}
                href={`${DOCS_BASE}/${doc.file}`}
                target="_blank"
                rel="noreferrer"
                className="group"
              >
                <Card className="h-full transition-shadow duration-200 group-hover:shadow-[0_2px_8px_rgba(0,0,0,0.04)] dark:group-hover:shadow-[0_2px_8px_rgba(0,0,0,0.25)]">
                  <CardHeader>
                    <div className="mb-2 flex items-start justify-between gap-2">
                      <span className="flex size-10 items-center justify-center rounded-lg border border-border bg-secondary/60">
                        <doc.icon className="size-4 text-foreground" strokeWidth={1.75} />
                      </span>
                      <ArrowUpRight className="size-4 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100" />
                    </div>
                    <CardTitle className="font-sans text-base font-semibold">
                      {doc.title}
                    </CardTitle>
                    <CardDescription>{doc.description}</CardDescription>
                    <p className="mt-2 font-mono text-[11px] text-muted-foreground">
                      documents/{doc.file}
                    </p>
                  </CardHeader>
                </Card>
              </a>
            ))}
          </div>
        </div>

        <div className="mt-14 overflow-hidden rounded-xl border border-border bg-card">
          <div className="border-b border-border px-6 py-4">
            <h3 className="text-sm font-semibold">Quick commands</h3>
            <p className="mt-1 text-sm text-muted-foreground">
              From the docs index. Scan summary always prints to stderr; findings
              go to stdout (text, JSON, or SARIF).
            </p>
          </div>
          <div className="grid gap-0 sm:grid-cols-2 lg:grid-cols-3">
            {quickCommands.map((cmd) => (
              <div
                key={cmd.label}
                className="border-b border-border px-6 py-4 last:border-b-0 sm:odd:border-r lg:border-r lg:[&:nth-child(3n)]:border-r-0"
              >
                <p className="text-xs font-medium text-muted-foreground">
                  {cmd.label}
                </p>
                <pre className="mt-2 overflow-x-auto font-mono text-[12px] leading-relaxed">
                  <code>{cmd.code}</code>
                </pre>
              </div>
            ))}
          </div>
        </div>

        <div className="mt-8 flex flex-wrap items-center justify-center gap-3">
          <Button variant="outline" asChild>
            <a
              href={`${DOCS_BASE}/README.md`}
              target="_blank"
              rel="noreferrer"
            >
              Open documents/README.md
              <ArrowUpRight className="size-4" />
            </a>
          </Button>
        </div>
      </div>
    </section>
  )
}
