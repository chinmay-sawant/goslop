import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { GitHubStarsButton } from '@/components/github-stars'

export function CtaSection() {
  return (
    <section className="relative overflow-hidden py-24 md:py-28">
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 -z-10"
        style={{
          background:
            'radial-gradient(ellipse 70% 80% at 50% 100%, color-mix(in srgb, var(--info) 10%, transparent), transparent)',
        }}
      />
      <div className="mx-auto max-w-3xl px-6 text-center">
        <h2 className="font-heading text-4xl tracking-tight md:text-5xl text-balance">
          Put a SAT in your agent loop
        </h2>
        <p className="mx-auto mt-4 max-w-xl text-muted-foreground text-balance">
          Scan for perf, security, and hygiene, then hand structured context to
          the agents that fix the code.
        </p>
        <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
          <Button size="lg" asChild>
            <a
              href="https://github.com/chinmay-sawant/goslop"
              target="_blank"
              rel="noreferrer"
            >
              Open on GitHub
              <ArrowRight className="size-4" />
            </a>
          </Button>
          <GitHubStarsButton size="default" className="h-12 px-5" />
          <Button size="lg" variant="outline" asChild>
            <a href="#install">Install locally</a>
          </Button>
        </div>
      </div>
    </section>
  )
}
