import { useState } from 'react'
import { Bot, FileCode2, Layers } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

/**
 * Real export from Chunk_1_25 / functions/24.txt (gopdfsuit reference scan).
 *
 * Finding 24 · PERF-42 · true positive
 * Evidence: runner.go returns fmt.Errorf("no successful runs") with no
 * format verbs; PERF-42 correctly flags Sprintf allocation via fmt.Errorf.
 * Fix: errors.New("no successful runs").
 *
 * Context is trimmed for the UI but lines, message, and hit marker match the
 * on-disk export (whole-function context in the full files).
 */
const FINDING = {
  index: 24,
  total: 915,
  source: 'gopdfsuit/internal/benchmarktemplates/runner.go:110:10',
  rule: 'PERF-42',
  title: 'fmt.Errorf Without Format Verbs',
  severity: 'info',
  message:
    'fmt.Errorf with a static string allocates a Sprintf; use errors.New instead',
  fix: 'Use errors.New or a sentinel error when the message has no format verbs.',
  functionFile: 'scripts/findings/functions/24.txt',
  chunkFile: 'scripts/chunks/Chunk_1_25.txt',
} as const

type Line = {
  n: number
  code: string
  hit?: boolean
  ellipsis?: boolean
}

/** Authentic context lines around the PERF-42 hit (export line numbers). */
const CONTEXT_LINES: Line[] = [
  { n: 56, code: 'func RunSingleDocumentBenchmark(name string) error {' },
  { n: 57, code: '	template, err := BuildZerodhaRetailTemplate()' },
  { n: 58, code: '	if err != nil {' },
  { n: 59, code: '		return err' },
  { n: 60, code: '	}' },
  { n: 0, code: '// ... workers, goroutines, PDF runs ...', ellipsis: true },
  { n: 106, code: '	memDone <- true' },
  { n: 107, code: '	memWg.Wait()' },
  { n: 108, code: '' },
  { n: 109, code: '	if len(durations) == 0 {' },
  {
    n: 110,
    code: '		return fmt.Errorf("no successful runs")',
    hit: true,
  },
  { n: 111, code: '	}' },
  { n: 112, code: '' },
  { n: 113, code: '	sort.Float64s(durations)' },
  { n: 0, code: '// ... print stats, return nil ...', ellipsis: true },
  { n: 131, code: '}' },
]

type View = 'chunk' | 'function'

function ExportMeta() {
  return (
    <div className="space-y-1 font-mono text-[12px] leading-relaxed md:text-[13px]">
      <p>
        <span className="text-muted-foreground">Finding </span>
        <span className="text-foreground">
          {FINDING.index}/{FINDING.total}
        </span>
      </p>
      <p>
        <span className="text-muted-foreground">Source: </span>
        <span className="text-foreground break-all">{FINDING.source}</span>
      </p>
      <p>
        <span className="text-muted-foreground">Rule: </span>
        <span className="text-warning font-medium">{FINDING.rule}</span>
        <span className="text-muted-foreground"> · </span>
        <span className="text-foreground">{FINDING.title}</span>
      </p>
      <p>
        <span className="text-muted-foreground">Severity: </span>
        <span className="text-info">{FINDING.severity}</span>
      </p>
      <p>
        <span className="text-muted-foreground">Message: </span>
        <span className="text-foreground">{FINDING.message}</span>
      </p>
      <p>
        <span className="text-muted-foreground">Fix: </span>
        <span className="text-foreground">{FINDING.fix}</span>
      </p>
    </div>
  )
}

function ContextBlock() {
  return (
    <div className="mt-4 space-y-1">
      <p className="font-mono text-[12px] text-muted-foreground md:text-[13px]">
        Context:
      </p>
      <div className="overflow-x-auto rounded-lg border border-border bg-card/40">
        <table className="w-full border-collapse font-mono text-[12px] leading-relaxed md:text-[13px]">
          <tbody>
            {CONTEXT_LINES.map((line, i) => (
              <tr
                key={`${line.n}-${i}`}
                className={cn(
                  line.hit && 'bg-danger/10',
                  line.ellipsis && 'text-muted-foreground/80',
                )}
              >
                <td
                  className={cn(
                    'w-10 select-none py-0.5 pr-2 pl-3 text-right align-top tabular-nums text-muted-foreground',
                    line.hit && 'text-danger font-medium',
                  )}
                >
                  {line.hit ? '>' : line.ellipsis ? '' : line.n}
                </td>
                <td
                  className={cn(
                    'py-0.5 pr-3 whitespace-pre',
                    line.hit && 'text-danger',
                    !line.hit && !line.ellipsis && 'text-foreground/90',
                  )}
                >
                  {line.code}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

function ExportPanel() {
  const [view, setView] = useState<View>('chunk')
  const fileLabel =
    view === 'chunk' ? FINDING.chunkFile : FINDING.functionFile

  return (
    <div className="overflow-hidden rounded-xl border border-border bg-background shadow-xs">
      <div className="flex flex-wrap items-center justify-between gap-2 border-b border-border px-3 py-2.5 sm:px-4">
        <div
          className="inline-flex rounded-md border border-border bg-secondary/40 p-0.5"
          role="tablist"
          aria-label="Export view"
        >
          <button
            type="button"
            role="tab"
            aria-selected={view === 'chunk'}
            onClick={() => setView('chunk')}
            className={cn(
              'rounded-[5px] px-3 py-1.5 font-mono text-[11px] transition-colors',
              view === 'chunk'
                ? 'bg-card text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            )}
          >
            chunks
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={view === 'function'}
            onClick={() => setView('function')}
            className={cn(
              'rounded-[5px] px-3 py-1.5 font-mono text-[11px] transition-colors',
              view === 'function'
                ? 'bg-card text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground',
            )}
          >
            functions
          </button>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="warning" className="font-mono text-[10px]">
            {FINDING.rule}
          </Badge>
          <Badge variant="muted" className="font-mono text-[10px]">
            export
          </Badge>
        </div>
      </div>

      <div className="border-b border-border px-4 py-2">
        <p className="truncate font-mono text-[11px] text-muted-foreground">
          {fileLabel}
        </p>
      </div>

      <div className="max-h-[28rem] overflow-auto p-4 sm:p-5">
        {view === 'chunk' && (
          <div className="mb-4 space-y-2 border-b border-border pb-4">
            <p className="font-mono text-[12px] text-muted-foreground md:text-[13px]">
              Findings 1-25 of {FINDING.total}
            </p>
            <p className="font-mono text-[11px] text-muted-foreground">
              Batch for agent delegation · finding {FINDING.index} of 25 in this
              chunk (true positive)
            </p>
          </div>
        )}

        {view === 'function' && (
          <div className="mb-4 space-y-2 border-b border-border pb-4">
            <p className="font-mono text-[12px] text-muted-foreground md:text-[13px]">
              One finding · one file · whole enclosing function
            </p>
            <p className="font-mono text-[11px] text-muted-foreground">
              Stable ref for tickets and single-issue agent prompts
            </p>
          </div>
        )}

        <ExportMeta />
        <ContextBlock />

        {view === 'chunk' && (
          <p className="mt-4 font-mono text-[11px] text-muted-foreground">
            {'='.repeat(48)}
            <br />
            <span className="opacity-70">
              … remaining findings 25 of 25 continue in the same chunk file
            </span>
          </p>
        )}
      </div>
    </div>
  )
}

export function AgentsSection() {
  return (
    <section id="agents" className="border-b border-border bg-card py-24 md:py-32">
      <div className="mx-auto max-w-6xl px-6">
        <div className="grid items-center gap-14 lg:grid-cols-2">
          <div>
            <h2 className="font-heading text-4xl tracking-tight md:text-5xl">
              Findings that agents can act on
            </h2>
            <p className="mt-4 text-muted-foreground leading-relaxed">
              Most SATs dump line hits and leave the model to hunt context.
              goslop exports the whole enclosing function by default, plus
              batched chunks so you can hand work to coding agents without
              stuffing the entire monorepo into a prompt. See{' '}
              <span className="font-mono text-xs">export-context-and-chunks.md</span>.
            </p>

            <ul className="mt-8 space-y-5">
              <li className="flex gap-4">
                <span className="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-md border border-border bg-background">
                  <FileCode2 className="size-4" strokeWidth={1.75} />
                </span>
                <div>
                  <p className="font-medium">Function-level context</p>
                  <p className="text-sm text-muted-foreground">
                    One finding, one file under{' '}
                    <code className="font-mono text-xs bg-secondary px-1.5 py-0.5 rounded">
                      scripts/findings/functions/
                    </code>{' '}
                    for single-issue deep dives and ticket links.
                  </p>
                </div>
              </li>
              <li className="flex gap-4">
                <span className="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-md border border-border bg-background">
                  <Layers className="size-4" strokeWidth={1.75} />
                </span>
                <div>
                  <p className="font-medium">Batched chunks</p>
                  <p className="text-sm text-muted-foreground">
                    Default 25 findings per{' '}
                    <code className="font-mono text-xs bg-secondary px-1.5 py-0.5 rounded">
                      scripts/chunks/Chunk_*.txt
                    </code>{' '}
                    for parallel agent delegation.
                  </p>
                </div>
              </li>
              <li className="flex gap-4">
                <span className="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-md border border-border bg-background">
                  <Bot className="size-4" strokeWidth={1.75} />
                </span>
                <div>
                  <p className="font-medium">Human or agent triage</p>
                  <p className="text-sm text-muted-foreground">
                    Same export path works for manual review queues and
                    multi-agent fix loops in CI or local sessions.
                  </p>
                </div>
              </li>
            </ul>

            <div className="mt-8">
              <Button variant="outline" asChild>
                <a
                  href="https://github.com/chinmay-sawant/goslop/blob/main/documents/export-context-and-chunks.md"
                  target="_blank"
                  rel="noreferrer"
                >
                  Read export docs
                </a>
              </Button>
            </div>
          </div>

          <ExportPanel />
        </div>
      </div>
    </section>
  )
}
