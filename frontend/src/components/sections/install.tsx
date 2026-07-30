import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

const steps = [
  {
    title: 'Build',
    code: 'make build\n# or: CGO_ENABLED=0 go build -o bin/goslop ./cmd/goslop',
  },
  {
    title: 'Scan',
    code: './bin/goslop .\n./bin/goslop --profile perf .\n./bin/goslop --profile security ./cmd',
  },
  {
    title: 'Export for agents',
    code: './bin/goslop --profile all \\\n  --export-context --export-chunks --no-cache .',
  },
  {
    title: 'Machine output',
    code: './bin/goslop --format json .\n./bin/goslop --format sarif . > goslop.sarif',
  },
]

export function InstallSection() {
  return (
    <section id="install" className="border-b border-border bg-card py-24 md:py-32">
      <div className="mx-auto max-w-6xl px-6">
        <div className="grid gap-12 lg:grid-cols-[1fr_1.1fr] lg:items-start">
          <div>
            <h2 className="font-heading text-4xl tracking-tight md:text-5xl">
              One binary. Zero C toolchain.
            </h2>
            <p className="mt-4 text-muted-foreground leading-relaxed">
              Requires Go 1.26.4 (see{' '}
              <code className="font-mono text-xs bg-secondary px-1.5 py-0.5 rounded">
                go.mod
              </code>
              ). Works on Linux, macOS, and Windows. Default build is{' '}
              <code className="font-mono text-xs bg-secondary px-1.5 py-0.5 rounded">
                CGO_ENABLED=0
              </code>
              . Write a starter config with{' '}
              <code className="font-mono text-xs bg-secondary px-1.5 py-0.5 rounded">
                ./bin/goslop init
              </code>
              .
            </p>

            <div className="mt-8 flex flex-wrap gap-3">
              <Button asChild>
                <a
                  href="https://github.com/chinmay-sawant/goslop"
                  target="_blank"
                  rel="noreferrer"
                >
                  View on GitHub
                </a>
              </Button>
              <Button variant="outline" asChild>
                <a href="#docs">Browse docs</a>
              </Button>
              <Button variant="outline" asChild>
                <a
                  href="https://github.com/chinmay-sawant/goslop/blob/main/documents/cli-reference.md"
                  target="_blank"
                  rel="noreferrer"
                >
                  CLI reference
                </a>
              </Button>
            </div>

            <div className="mt-10 flex flex-wrap gap-2">
              <Badge variant="outline">Go 1.26.4</Badge>
              <Badge variant="outline">SARIF 2.1.0</Badge>
              <Badge variant="outline">JSON · text</Badge>
              <Badge variant="outline">CI-ready</Badge>
            </div>
          </div>

          <div className="space-y-4">
            {steps.map((step, index) => (
              <div
                key={step.title}
                className="overflow-hidden rounded-xl border border-border bg-background"
              >
                <div className="flex items-center gap-3 border-b border-border px-4 py-2.5">
                  <span className="flex size-6 items-center justify-center rounded-md bg-secondary font-mono text-xs font-medium">
                    {index + 1}
                  </span>
                  <span className="text-sm font-medium">{step.title}</span>
                </div>
                <pre className="overflow-x-auto p-4 font-mono text-[12px] leading-relaxed md:text-[13px]">
                  <code>{step.code}</code>
                </pre>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
