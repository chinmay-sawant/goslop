import { useState } from 'react'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

/**
 * Real export from Chunk_1_25 / functions/24.txt (gopdfsuit reference scan).
 * Finding 24 · PERF-42 · true positive (fmt.Errorf static string).
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

const CONTEXT_LINES: Line[] = [
  { n: 56, code: 'func RunSingleDocumentBenchmark(name string) error {' },
  { n: 57, code: '\ttemplate, err := BuildZerodhaRetailTemplate()' },
  { n: 58, code: '\tif err != nil {' },
  { n: 59, code: '\t\treturn err' },
  { n: 60, code: '\t}' },
  { n: 0, code: '// ... workers, goroutines, PDF runs ...', ellipsis: true },
  { n: 106, code: '\tmemDone <- true' },
  { n: 107, code: '\tmemWg.Wait()' },
  { n: 108, code: '' },
  { n: 109, code: '\tif len(durations) == 0 {' },
  { n: 110, code: '\t\treturn fmt.Errorf("no successful runs")', hit: true },
  { n: 111, code: '\t}' },
  { n: 112, code: '' },
  { n: 113, code: '\tsort.Float64s(durations)' },
  { n: 0, code: '// ... print stats, return nil ...', ellipsis: true },
  { n: 131, code: '}' },
]

type Panel = 'scan' | 'chunk' | 'function'

const SCAN_SUMMARY = `goslop scan summary
  files: 142   lines: 48.2k   wall: 1.8s
  cache: warm   findings: 915
  high: 41   medium: 312   low: 401   info: 161
  top: PERF-42, PERF-68, CWE-497, BP-70`

const SCAN_FINDING = `PERF-42  info  runner.go:110:10
fmt.Errorf with a static string allocates a Sprintf; use errors.New instead

PERF-68  medium  main.go:66:14
gin.Logger() performs synchronous I/O on the request path

CWE-497  high  main.go:72:19
diagnostics endpoint exposes host environment details`

function ContextBlock() {
  return (
    <div className="mt-4 space-y-1">
      <p className="font-mono text-[12px] text-muted-foreground">Context:</p>
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
        <span className="break-all text-foreground">{FINDING.source}</span>
      </p>
      <p>
        <span className="text-muted-foreground">Rule: </span>
        <span className="font-medium text-warning">{FINDING.rule}</span>
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

export function DemoSection() {
  const [panel, setPanel] = useState<Panel>('scan')

  const fileLabel =
    panel === 'scan'
      ? 'stderr summary · stdout findings'
      : panel === 'chunk'
        ? FINDING.chunkFile
        : FINDING.functionFile

  return (
    <section id="demo" className="border-b border-border bg-card py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
            The product is the output
          </h2>
          <p className="mt-4 text-muted-foreground text-balance">
            Summary on stderr, findings on stdout, then optional disk export for
            agents. Below: a real PERF-42 true positive from the reference
            corpus export (Chunk_1_25 / functions/24.txt).
          </p>
        </div>

        <div className="mx-auto mt-12 max-w-3xl overflow-hidden rounded-xl border border-border bg-background shadow-xs">
          <div className="flex flex-wrap items-center justify-between gap-2 border-b border-border px-3 py-2.5 sm:px-4">
            <div
              className="inline-flex rounded-md border border-border bg-secondary/40 p-0.5"
              role="tablist"
              aria-label="Demo view"
            >
              {(
                [
                  { id: 'scan', label: 'scan' },
                  { id: 'chunk', label: 'chunks' },
                  { id: 'function', label: 'functions' },
                ] as const
              ).map((tab) => (
                <button
                  key={tab.id}
                  type="button"
                  role="tab"
                  aria-selected={panel === tab.id}
                  onClick={() => setPanel(tab.id)}
                  className={cn(
                    'rounded-[5px] px-3 py-1.5 font-mono text-[11px] transition-colors',
                    panel === tab.id
                      ? 'bg-card text-foreground shadow-xs'
                      : 'text-muted-foreground hover:text-foreground',
                  )}
                >
                  {tab.label}
                </button>
              ))}
            </div>
            <div className="flex items-center gap-2">
              {panel !== 'scan' && (
                <Badge variant="warning" className="font-mono text-[10px]">
                  {FINDING.rule}
                </Badge>
              )}
              <Badge variant="muted" className="font-mono text-[10px]">
                {panel === 'scan' ? 'live shape' : 'export'}
              </Badge>
            </div>
          </div>

          <div className="border-b border-border px-4 py-2">
            <p className="truncate font-mono text-[11px] text-muted-foreground">
              {fileLabel}
            </p>
          </div>

          <div className="max-h-[28rem] overflow-auto p-4 sm:p-5">
            {panel === 'scan' && (
              <div className="space-y-5 font-mono text-[12px] leading-relaxed md:text-[13px]">
                <div>
                  <p className="mb-2 text-muted-foreground"># stderr · scan summary</p>
                  <pre className="whitespace-pre-wrap text-foreground/90">
                    {SCAN_SUMMARY}
                  </pre>
                </div>
                <div className="border-t border-border pt-4">
                  <p className="mb-2 text-muted-foreground"># stdout · findings (text)</p>
                  <pre className="whitespace-pre-wrap text-foreground/90">
                    {SCAN_FINDING}
                  </pre>
                </div>
                <p className="text-[11px] text-muted-foreground">
                  Summary stays pipe-clean for JSON/SARIF: findings on stdout only.
                </p>
              </div>
            )}

            {panel === 'chunk' && (
              <>
                <div className="mb-4 space-y-2 border-b border-border pb-4">
                  <p className="font-mono text-[12px] text-muted-foreground">
                    Findings 1-25 of {FINDING.total}
                  </p>
                  <p className="font-mono text-[11px] text-muted-foreground">
                    Batch for agent delegation · finding {FINDING.index} (true
                    positive)
                  </p>
                </div>
                <ExportMeta />
                <ContextBlock />
                <p className="mt-4 font-mono text-[11px] text-muted-foreground">
                  {'='.repeat(40)}
                  <br />
                  <span className="opacity-70">
                    … remaining findings continue in the same chunk file
                  </span>
                </p>
              </>
            )}

            {panel === 'function' && (
              <>
                <div className="mb-4 space-y-2 border-b border-border pb-4">
                  <p className="font-mono text-[12px] text-muted-foreground">
                    One finding · one file · whole enclosing function
                  </p>
                  <p className="font-mono text-[11px] text-muted-foreground">
                    Stable ref for tickets and single-issue agent prompts
                  </p>
                </div>
                <ExportMeta />
                <ContextBlock />
              </>
            )}
          </div>
        </div>

        <p className="mx-auto mt-6 max-w-2xl text-center text-xs text-muted-foreground">
          Reference corpus self-scan counts are large on purpose (fixtures +
          needles). Judge signal on real app trees; use baseline for known debt.
        </p>
      </div>
    </section>
  )
}
