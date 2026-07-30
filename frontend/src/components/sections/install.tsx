import { useState } from 'react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { CodeBlock } from '@/components/code-block'
import { Reveal } from '@/components/reveal'
import { cn } from '@/lib/utils'

/** Module path from go.mod: github.com/chinmay-sawant/goslop */
const GO_INSTALL =
  'go install github.com/chinmay-sawant/goslop/cmd/goslop@latest'

const PATHS = [
  {
    id: 'go-install',
    label: 'go install',
    code: `# requires Go 1.26.4+ (see go.mod)
# module: github.com/chinmay-sawant/goslop
# binary lands on your GOPATH/bin or GOBIN
CGO_ENABLED=0 ${GO_INSTALL}

# then (ensure $(go env GOPATH)/bin is on PATH)
goslop --profile recommended .
goslop init   # optional starter goslop.toml`,
  },
  {
    id: 'make',
    label: 'make (from clone)',
    code: `# requires Go 1.26.4+ (see go.mod); pure Go, no C toolchain
git clone https://github.com/chinmay-sawant/goslop.git
cd goslop
make build
./bin/goslop --profile recommended .
./bin/goslop init   # optional starter goslop.toml`,
  },
  {
    id: 'go',
    label: 'go build',
    code: `git clone https://github.com/chinmay-sawant/goslop.git
cd goslop
CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop
./bin/goslop --profile recommended .`,
  },
  {
    id: 'windows',
    label: 'Windows',
    code: `# go install (PowerShell / cmd)
set CGO_ENABLED=0
go install github.com/chinmay-sawant/goslop/cmd/goslop@latest
goslop --profile recommended .

# or from a clone
git clone https://github.com/chinmay-sawant/goslop.git
cd goslop
set CGO_ENABLED=0
go build -o bin\\goslop.exe .\\cmd\\goslop
.\\bin\\goslop.exe --profile recommended .`,
  },
] as const

const PRIMARY_CMD = GO_INSTALL

export function InstallSection() {
  const [id, setId] = useState<(typeof PATHS)[number]['id']>('go-install')
  const active = PATHS.find((p) => p.id === id) ?? PATHS[0]

  return (
    <section id="install" className="border-b border-border bg-card py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <Reveal>
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
              Install in one line
            </h2>
            <p className="mt-4 text-muted-foreground text-balance">
              Module path from{' '}
              <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
                go.mod
              </code>
              :{' '}
              <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
                github.com/chinmay-sawant/goslop
              </code>
              . Prefer{' '}
              <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
                CGO_ENABLED=0
              </code>
              . Primary install:
            </p>
            <p className="mt-4 inline-flex max-w-full items-center gap-2 rounded-lg border border-border bg-background px-3 py-2 font-mono text-[11px] sm:text-[13px]">
              <span className="truncate">{PRIMARY_CMD}</span>
            </p>
          </div>
        </Reveal>

        <Reveal className="mx-auto mt-10 max-w-2xl">
          <div
            className="mb-3 flex flex-wrap gap-2"
            role="tablist"
            aria-label="Install method"
          >
            {PATHS.map((p) => (
              <button
                key={p.id}
                type="button"
                role="tab"
                aria-selected={id === p.id}
                onClick={() => setId(p.id)}
                className={cn(
                  'rounded-md border px-3 py-1.5 text-xs font-medium transition-colors focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50',
                  id === p.id
                    ? 'border-foreground/20 bg-background text-foreground shadow-xs'
                    : 'border-border text-muted-foreground hover:text-foreground',
                )}
              >
                {p.label}
              </button>
            ))}
          </div>
          <CodeBlock code={active.code} filename={`install · ${active.label}`} />
        </Reveal>

        <Reveal>
          <div className="mx-auto mt-8 flex max-w-2xl flex-wrap items-center justify-center gap-2">
            <Badge variant="outline">Go 1.26.4+</Badge>
            <Badge variant="outline">CGO_ENABLED=0</Badge>
            <Badge variant="outline">SARIF 2.1.0</Badge>
            <Badge variant="outline">Linux · macOS · Windows</Badge>
            <Badge variant="outline">pure Go</Badge>
          </div>
        </Reveal>

        <Reveal>
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
              <a href="#demo">See demo</a>
            </Button>
            <Button variant="outline" asChild>
              <a href="#docs">Docs journeys</a>
            </Button>
          </div>
        </Reveal>
      </div>
    </section>
  )
}
