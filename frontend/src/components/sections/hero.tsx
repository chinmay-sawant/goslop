import { ArrowRight } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { CodeBlock } from '@/components/code-block'

const HERO_CMD = `CGO_ENABLED=0 go install github.com/chinmay-sawant/goslop/cmd/goslop@latest
goslop --profile recommended .
goslop --export-context --export-chunks .`

export function HeroSection() {
  return (
    <section id="top" className="relative overflow-hidden border-b border-border">
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 -z-10 motion-safe:opacity-100"
        style={{
          background:
            'radial-gradient(ellipse 80% 50% at 50% -20%, color-mix(in srgb, var(--info) 12%, transparent), transparent), radial-gradient(ellipse 60% 40% at 100% 0%, color-mix(in srgb, var(--warning) 8%, transparent), transparent)',
        }}
      />

      <div className="mx-auto max-w-6xl px-6 pt-16 pb-20 md:pt-24 md:pb-28">
        <div className="grid items-center gap-12 lg:grid-cols-[1.05fr_0.95fr] lg:gap-16">
          <div>
            <Badge
              variant="outline"
              className="mb-5 gap-2 px-3 py-1 font-mono text-[11px] uppercase tracking-wider"
            >
              <span className="size-1.5 rounded-full bg-success" />
              Pure-Go SAT · CGO_ENABLED=0
            </Badge>

            <h1 className="font-heading text-4xl leading-[1.08] tracking-tight text-balance sm:text-5xl lg:text-6xl">
              Find it. Export it. Hand it to an agent.
            </h1>

            <p className="mt-5 max-w-xl text-base text-muted-foreground text-balance sm:text-lg">
              <strong className="font-medium text-foreground">goslop</strong> is
              a static analysis tool for Go: performance, CWE security, and
              hygiene. The difference is the export path: whole-function context
              and batched chunks so agents fix issues without spelunking the
              monorepo.
            </p>

            <ol className="mt-8 max-w-md space-y-3 text-sm">
              {[
                { step: '1', label: 'Scan', detail: 'profiles for CI, perf, security, style' },
                { step: '2', label: 'Export', detail: 'functions + chunks with full context' },
                { step: '3', label: 'Delegate', detail: 'humans or coding agents remediate' },
              ].map((item) => (
                <li key={item.step} className="flex gap-3">
                  <span className="flex size-7 shrink-0 items-center justify-center rounded-md border border-border bg-card font-mono text-xs font-medium">
                    {item.step}
                  </span>
                  <span>
                    <span className="font-medium text-foreground">{item.label}</span>
                    <span className="text-muted-foreground"> · {item.detail}</span>
                  </span>
                </li>
              ))}
            </ol>

            <div className="mt-9 flex flex-col gap-3 sm:flex-row sm:items-center">
              <Button size="lg" asChild>
                <a href="#install">
                  Install goslop
                  <ArrowRight className="size-4" />
                </a>
              </Button>
              <Button size="lg" variant="outline" asChild>
                <a href="#demo">See the export</a>
              </Button>
            </div>

            <p className="mt-5 font-mono text-[11px] text-muted-foreground">
              Default gate: --profile recommended · SARIF 2.1.0 · pure Go
            </p>
          </div>

          <div className="motion-safe:animate-fade-up">
            <CodeBlock code={HERO_CMD} filename="get started">
              <span className="text-muted-foreground"># install (module from go.mod)</span>
              {'\n'}
              CGO_ENABLED=0 go install github.com/chinmay-sawant/goslop/cmd/goslop@latest
              {'\n\n'}
              <span className="text-muted-foreground"># scan</span>
              {'\n'}
              goslop --profile recommended .
              {'\n\n'}
              <span className="text-muted-foreground"># export for agents</span>
              {'\n'}
              goslop --export-context --export-chunks .
            </CodeBlock>
          </div>
        </div>
      </div>
    </section>
  )
}
