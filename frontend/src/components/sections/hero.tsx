import { ArrowRight, Terminal } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { GitHubStarsButton } from '@/components/github-stars'

export function HeroSection() {
  return (
    <section id="top" className="relative overflow-hidden border-b border-border">
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 -z-10"
        style={{
          background:
            'radial-gradient(ellipse 80% 50% at 50% -20%, color-mix(in srgb, var(--info) 12%, transparent), transparent), radial-gradient(ellipse 60% 40% at 100% 0%, color-mix(in srgb, var(--warning) 8%, transparent), transparent)',
        }}
      />

      <div className="mx-auto max-w-6xl px-6 pt-20 pb-24 md:pt-28 md:pb-32">
        <div className="mx-auto max-w-3xl text-center">
          <div className="animate-fade-up mb-6 flex justify-center">
            <Badge variant="outline" className="gap-2 px-3 py-1 font-mono text-[11px] uppercase tracking-wider">
              <span className="size-1.5 rounded-full bg-success" />
              Pure-Go SAT · No CGO
            </Badge>
          </div>

          <h1 className="animate-fade-up delay-100 font-heading text-5xl leading-[1.08] tracking-tight text-balance md:text-6xl lg:text-7xl">
            Static analysis for the{' '}
            <em className="not-italic text-info">agentic era</em> of software
          </h1>

          <p className="animate-fade-up delay-200 mx-auto mt-6 max-w-2xl text-lg text-muted-foreground text-balance md:text-xl">
            <strong className="font-medium text-foreground">goslop</strong> is a
            pure-Go SAT that finds performance, security, and hygiene issues,
            then exports whole-function context so agents can fix them without
            spelunking your tree.
          </p>

          <div className="animate-fade-up delay-300 mt-10 flex flex-col items-center justify-center gap-3 sm:flex-row">
            <Button size="lg" asChild>
              <a href="#install">
                Install goslop
                <ArrowRight className="size-4" />
              </a>
            </Button>
            <Button size="lg" variant="outline" asChild>
              <a href="#docs">
                <Terminal className="size-4" />
                Documentation
              </a>
            </Button>
            <GitHubStarsButton size="default" className="h-12 px-5" />
          </div>
        </div>

        <div className="animate-fade-up delay-400 mx-auto mt-16 max-w-3xl">
          <div className="overflow-hidden rounded-xl border border-border bg-card shadow-xs">
            <div className="flex items-center gap-2 border-b border-border bg-secondary/50 px-4 py-3">
              <span className="size-2.5 rounded-full bg-muted-foreground/25" />
              <span className="size-2.5 rounded-full bg-muted-foreground/25" />
              <span className="size-2.5 rounded-full bg-muted-foreground/25" />
              <span className="ml-2 font-mono text-xs text-muted-foreground">
                terminal
              </span>
            </div>
            <pre className="overflow-x-auto p-5 font-mono text-[13px] leading-relaxed text-foreground/90 md:text-sm">
              <code>
                <span className="text-muted-foreground">$</span> ./bin/goslop --profile recommended .{'\n'}
                <span className="text-muted-foreground">
                  # text on stdout · summary on stderr · exit by severity
                </span>
                {'\n\n'}
                <span className="text-muted-foreground">$</span> ./bin/goslop --format sarif . {'>'} goslop.sarif{'\n'}
                <span className="text-muted-foreground">$</span> ./bin/goslop --export-context --export-chunks .{'\n'}
                <span className="text-success">
                  # → scripts/findings/functions · scripts/chunks for agents
                </span>
              </code>
            </pre>
          </div>
        </div>
      </div>
    </section>
  )
}
