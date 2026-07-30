import { ArrowUpRight, Rocket, Terminal, FileOutput, GitBranch } from 'lucide-react'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'

const DOCS =
  'https://github.com/chinmay-sawant/goslop/blob/main/documents'

/** Four journeys only (not a 12-card sitemap). */
const journeys = [
  {
    icon: Rocket,
    title: 'Get started',
    file: 'overview.md',
    description:
      'What goslop does, detector families, profiles, cache, and suppressions at a glance.',
  },
  {
    icon: Terminal,
    title: 'CLI',
    file: 'cli-reference.md',
    description:
      'Every flag, profile alias, exit code, and config merge rule for ./bin/goslop.',
  },
  {
    icon: FileOutput,
    title: 'Export for agents',
    file: 'export-context-and-chunks.md',
    description:
      'Function refs vs chunks, whole-function Context default, and delegation workflow.',
  },
  {
    icon: GitBranch,
    title: 'CI & reporting',
    file: 'reporting-formats.md',
    description:
      'Text, JSON, SARIF 2.1.0, dual stdout/stderr streams, and Code Scanning shape.',
  },
]

export function DocsSection() {
  return (
    <section id="docs" className="border-b border-border py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
            Four paths into the docs
          </h2>
          <p className="mt-4 text-muted-foreground text-balance">
            Not a second homepage. Pick a journey; deeper guides live under{' '}
            <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
              documents/
            </code>
            .
          </p>
        </div>

        <div className="mt-12 grid gap-4 sm:grid-cols-2">
          {journeys.map((doc) => (
            <a
              key={doc.file}
              href={`${DOCS}/${doc.file}`}
              target="_blank"
              rel="noreferrer"
              className="group"
            >
              <Card className="h-full transition-shadow duration-200 group-hover:shadow-[0_2px_8px_rgba(0,0,0,0.04)] dark:group-hover:shadow-[0_2px_8px_rgba(0,0,0,0.25)]">
                <CardHeader>
                  <div className="mb-2 flex items-start justify-between gap-2">
                    <span className="flex size-10 items-center justify-center rounded-lg border border-border bg-secondary/60">
                      <doc.icon className="size-4" strokeWidth={1.75} />
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

        <div className="mt-8 flex justify-center">
          <Button variant="outline" asChild>
            <a href={`${DOCS}/README.md`} target="_blank" rel="noreferrer">
              Full documents index
              <ArrowUpRight className="size-4" />
            </a>
          </Button>
        </div>
      </div>
    </section>
  )
}
