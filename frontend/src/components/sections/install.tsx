import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { CodeBlock } from '@/components/code-block'

const INSTALL = `# requires Go 1.26.4 (see go.mod); pure Go, no C toolchain
git clone https://github.com/chinmay-sawant/goslop.git
cd goslop
make build
./bin/goslop --profile recommended .
./bin/goslop init   # optional starter goslop.toml`

export function InstallSection() {
  return (
    <section id="install" className="border-b border-border bg-card py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
            One path to a binary
          </h2>
          <p className="mt-4 text-muted-foreground text-balance">
            No package manager install yet. Clone,{' '}
            <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
              make build
            </code>{' '}
            (<code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
              CGO_ENABLED=0
            </code>
            ), scan. That is the product install.
          </p>
        </div>

        <div className="mx-auto mt-10 max-w-2xl">
          <CodeBlock code={INSTALL} filename="install · linux / macOS / windows" />
        </div>

        <div className="mx-auto mt-8 flex max-w-2xl flex-wrap items-center justify-center gap-2">
          <Badge variant="outline">Go 1.26.4</Badge>
          <Badge variant="outline">CGO_ENABLED=0</Badge>
          <Badge variant="outline">SARIF 2.1.0</Badge>
          <Badge variant="outline">Linux · macOS · Windows</Badge>
        </div>

        <div className="mx-auto mt-8 flex max-w-2xl flex-wrap justify-center gap-3">
          <Button asChild>
            <a
              href="https://github.com/chinmay-sawant/goslop"
              target="_blank"
              rel="noreferrer"
            >
              Open repository
            </a>
          </Button>
          <Button variant="outline" asChild>
            <a href="#docs">Docs journeys</a>
          </Button>
        </div>
      </div>
    </section>
  )
}
