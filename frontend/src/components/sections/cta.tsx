import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { GitHubStarsButton } from '@/components/github-stars'

export function CtaSection() {
  return (
    <section className="relative overflow-hidden py-20 md:py-24">
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 -z-10"
        style={{
          background:
            'radial-gradient(ellipse 70% 80% at 50% 100%, color-mix(in srgb, var(--info) 10%, transparent), transparent)',
        }}
      />
      <div className="mx-auto max-w-3xl px-6 text-center">
        <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl text-balance">
          Scan tonight. Export to agents tomorrow.
        </h2>
        <p className="mx-auto mt-4 max-w-xl text-muted-foreground text-balance">
          One binary, pure Go, profiles for CI. The differentiator is context
          your agents can actually use.
        </p>
        <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
          <Button size="lg" asChild>
            <a href="#install">
              Install goslop
              <ArrowRight className="size-4" />
            </a>
          </Button>
          <GitHubStarsButton size="default" className="h-12 px-5" />
        </div>
      </div>
    </section>
  )
}
