import {
  Gauge,
  ShieldAlert,
  ListChecks,
  Layers3,
  FileJson2,
  DatabaseZap,
  Ban,
  PackageOpen,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

const features = [
  {
    icon: Gauge,
    title: 'Performance detectors',
    description:
      '239 PERF rules for allocations in loops, missing HTTP timeouts, N+1 patterns, and framework footguns that show up under load.',
    badge: 'PERF-*',
    tone: 'warning' as const,
  },
  {
    icon: ShieldAlert,
    title: 'Security & CWE',
    description:
      '175 structural CWE heuristics plus experimental taint for path traversal, command injection, XSS, and SQL injection (CWE-22/78/79/89).',
    badge: 'CWE-*',
    tone: 'danger' as const,
  },
  {
    icon: ListChecks,
    title: 'Bad practices',
    description:
      'Project-level hygiene: error handling, server shutdown, go.mod discipline, and noisy API patterns across style packs.',
    badge: 'BP-*',
    tone: 'info' as const,
  },
  {
    icon: Layers3,
    title: 'Product profiles',
    description:
      'recommended, perf, security, style, and all: curated packs with fail policies ready for local lint and CI gates.',
    badge: 'Packs',
    tone: 'muted' as const,
  },
  {
    icon: FileJson2,
    title: 'Text · JSON · SARIF',
    description:
      'Human text by default, machine JSON for pipelines, and SARIF 2.1.0 for GitHub Code Scanning. Summary always on stderr.',
    badge: 'Reporters',
    tone: 'muted' as const,
  },
  {
    icon: DatabaseZap,
    title: 'Incremental cache',
    description:
      'Per-file cache under .goslop-cache so re-scans stay fast as the tree changes.',
    badge: 'Cache',
    tone: 'success' as const,
  },
  {
    icon: Ban,
    title: 'Baseline & ignore',
    description:
      'Ship with known debt via .goslop-baseline.json, or silence lines with // goslop-ignore directives and path ignores.',
    badge: 'Debt control',
    tone: 'muted' as const,
  },
  {
    icon: PackageOpen,
    title: 'Pure Go, no CGO',
    description:
      'go/parser + go/ast only (no tree-sitter). CGO_ENABLED=0 by default means easy cross-compile and no C toolchain tax.',
    badge: 'Portable',
    tone: 'success' as const,
  },
]

export function FeaturesSection() {
  return (
    <section id="features" className="border-b border-border py-24 md:py-32">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <p className="font-mono text-xs font-medium uppercase tracking-wider text-muted-foreground">
            Capabilities
          </p>
          <h2 className="mt-3 font-heading text-4xl tracking-tight md:text-5xl">
            A full SAT surface for Go
          </h2>
          <p className="mt-4 text-muted-foreground text-balance">
            One binary for performance, security, and hygiene, with the
            reporting and suppression story you need in real repos. See the
            product overview in documentation for detector counts and status.
          </p>
        </div>

        <div className="mt-14 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {features.map((feature) => (
            <Card
              key={feature.title}
              className="transition-shadow duration-200 hover:shadow-[0_2px_8px_rgba(0,0,0,0.04)] dark:hover:shadow-[0_2px_8px_rgba(0,0,0,0.25)]"
            >
              <CardHeader>
                <div className="mb-2 flex items-center justify-between">
                  <span className="flex size-10 items-center justify-center rounded-lg border border-border bg-secondary/60">
                    <feature.icon className="size-4 text-foreground" strokeWidth={1.75} />
                  </span>
                  <Badge variant={feature.tone}>{feature.badge}</Badge>
                </div>
                <CardTitle className="font-sans text-base font-semibold">
                  {feature.title}
                </CardTitle>
                <CardDescription>{feature.description}</CardDescription>
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>
    </section>
  )
}
