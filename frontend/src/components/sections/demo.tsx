import { useEffect, useState } from 'react'
import { Badge } from '@/components/ui/badge'
import { CopyButton } from '@/components/copy-button'
import { Reveal } from '@/components/reveal'
import {
  buildChunkDossierText,
  buildFunctionDossierText,
  buildScanText,
  CONTEXT_LINES,
  FINDING,
  SCAN_COMMAND,
  SCAN_FINDING,
  SCAN_SUMMARY,
} from '@/lib/export-sample'
import { cn } from '@/lib/utils'

type Panel = 'scan' | 'chunk' | 'function'

const TABS: { id: Panel; label: string }[] = [
  { id: 'scan', label: 'scan' },
  { id: 'chunk', label: 'chunks' },
  { id: 'function', label: 'functions' },
]

function parsePanel(raw: string | null): Panel | null {
  if (raw === 'scan' || raw === 'chunk' || raw === 'function') return raw
  if (raw === 'chunks') return 'chunk'
  if (raw === 'functions') return 'function'
  return null
}

function readPanelFromLocation(): Panel {
  if (typeof window === 'undefined') return 'scan'
  const hash = window.location.hash
  // #demo=scan | #demo=chunk | #demo=function
  if (hash.startsWith('#demo=')) {
    return parsePanel(hash.slice('#demo='.length)) ?? 'scan'
  }
  const q = new URLSearchParams(window.location.search).get('demo')
  return parsePanel(q) ?? 'scan'
}

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
                    line.hit && 'font-medium text-danger',
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

function TerminalReplay() {
  return (
    <div className="mb-8 overflow-hidden rounded-xl border border-border bg-background shadow-xs">
      <div className="flex items-center gap-2 border-b border-border bg-secondary/40 px-4 py-2.5">
        <span className="size-2.5 rounded-full bg-muted-foreground/25" />
        <span className="size-2.5 rounded-full bg-muted-foreground/25" />
        <span className="size-2.5 rounded-full bg-muted-foreground/25" />
        <span className="ml-2 font-mono text-[11px] text-muted-foreground">
          replay · gopdfsuit reference scan
        </span>
        <Badge variant="muted" className="ml-auto font-mono text-[10px]">
          live numbers
        </Badge>
      </div>
      <pre className="overflow-x-auto p-4 font-mono text-[12px] leading-relaxed md:text-[13px]">
        <code className="terminal-replay block whitespace-pre-wrap text-foreground/90">
          <span className="text-muted-foreground">$</span> ./bin/goslop --profile all --no-fail \
          {'\n'}
          {'  '}--no-terminal --export-context --export-chunks --no-cache ./gopdfsuit
          {'\n'}
          <span className="line-reveal delay-1 text-muted-foreground">
            scanned 78 files (28042 lines) in 100.4ms
          </span>
          {'\n'}
          <span className="line-reveal delay-2 text-muted-foreground">
            {'  '}cache: 0 hits, 78 misses (full re-analysis)
          </span>
          {'\n'}
          <span className="line-reveal delay-3">
            915 findings · severity: 10 high, 396 medium, 312 low, 197 info
          </span>
          {'\n'}
          <span className="line-reveal delay-4 text-muted-foreground">
            {'  '}top rules: BP-1 ×181, PERF-6 ×94, PERF-32 ×59, BP-5 ×50, PERF-230 ×44
          </span>
          {'\n'}
          <span className="line-reveal delay-5 text-success">
            exported 915 context file(s) to scripts/findings/functions
          </span>
          {'\n'}
          <span className="line-reveal delay-6 text-success">
            exported 37 chunk file(s) to scripts/chunks
          </span>
        </code>
      </pre>
    </div>
  )
}

export function DemoSection() {
  const [panel, setPanel] = useState<Panel>('scan')

  useEffect(() => {
    setPanel(readPanelFromLocation())
    const onHash = () => setPanel(readPanelFromLocation())
    window.addEventListener('hashchange', onHash)
    return () => window.removeEventListener('hashchange', onHash)
  }, [])

  function selectPanel(next: Panel) {
    setPanel(next)
    const url = new URL(window.location.href)
    url.hash = `demo=${next}`
    window.history.replaceState(null, '', `${url.pathname}${url.search}#demo=${next}`)
  }

  const fileLabel =
    panel === 'scan'
      ? 'stderr summary · sample findings (gopdfsuit)'
      : panel === 'chunk'
        ? FINDING.chunkFile
        : FINDING.functionFile

  const copyText =
    panel === 'scan'
      ? buildScanText()
      : panel === 'chunk'
        ? buildChunkDossierText()
        : buildFunctionDossierText()

  return (
    <section id="demo" className="border-b border-border bg-card py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <Reveal>
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
              The product is the output
            </h2>
            <p className="mt-4 text-muted-foreground text-balance">
              Summary on stderr, findings on stdout, then disk export for agents.
              Real PERF-42 true positive from the reference corpus export.
            </p>
          </div>
        </Reveal>

        <Reveal className="mx-auto mt-10 max-w-3xl">
          <TerminalReplay />
        </Reveal>

        <Reveal className="mx-auto max-w-3xl">
          <div className="overflow-hidden rounded-xl border border-border bg-background shadow-xs">
            <div className="flex flex-wrap items-center justify-between gap-2 border-b border-border px-3 py-2.5 sm:px-4">
              <div
                className="inline-flex rounded-md border border-border bg-secondary/40 p-0.5"
                role="tablist"
                aria-label="Demo view"
              >
                {TABS.map((tab) => (
                  <button
                    key={tab.id}
                    type="button"
                    role="tab"
                    id={`demo-tab-${tab.id}`}
                    aria-controls={`demo-panel-${tab.id}`}
                    aria-selected={panel === tab.id}
                    tabIndex={panel === tab.id ? 0 : -1}
                    onClick={() => selectPanel(tab.id)}
                    className={cn(
                      'rounded-[5px] px-3 py-1.5 font-mono text-[11px] transition-colors focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50',
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
                <CopyButton text={copyText} label="Copy dossier" />
              </div>
            </div>

            <div className="flex items-center justify-between gap-2 border-b border-border px-4 py-2">
              <p className="truncate font-mono text-[11px] text-muted-foreground">
                {fileLabel}
              </p>
              <p className="shrink-0 font-mono text-[10px] text-muted-foreground">
                #demo={panel}
              </p>
            </div>

            <div
              id={`demo-panel-${panel}`}
              role="tabpanel"
              aria-labelledby={`demo-tab-${panel}`}
              className="max-h-[28rem] overflow-auto p-4 sm:p-5"
            >
              {panel === 'scan' && (
                <div className="space-y-5 font-mono text-[12px] leading-relaxed md:text-[13px]">
                  <div>
                    <p className="mb-2 text-muted-foreground"># command</p>
                    <pre className="whitespace-pre-wrap break-all text-foreground/90">
                      $ {SCAN_COMMAND}
                    </pre>
                  </div>
                  <div className="border-t border-border pt-4">
                    <p className="mb-2 text-muted-foreground"># stderr · scan summary</p>
                    <pre className="whitespace-pre-wrap text-foreground/90">
                      {SCAN_SUMMARY}
                    </pre>
                  </div>
                  <div className="border-t border-border pt-4">
                    <p className="mb-2 text-muted-foreground">
                      # sample findings (from same run · text shape)
                    </p>
                    <pre className="whitespace-pre-wrap text-foreground/90">
                      {SCAN_FINDING}
                    </pre>
                  </div>
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
                </>
              )}

              {panel === 'function' && (
                <>
                  <div className="mb-4 space-y-2 border-b border-border pb-4">
                    <p className="font-mono text-[12px] text-muted-foreground">
                      One finding · one file · whole enclosing function
                    </p>
                  </div>
                  <ExportMeta />
                  <ContextBlock />
                </>
              )}
            </div>
          </div>
        </Reveal>

        <p className="mx-auto mt-6 max-w-2xl text-center text-xs text-muted-foreground">
          Reference corpus self-scan counts are large on purpose (fixtures +
          needles). Judge signal on real app trees; use baseline for known debt.
        </p>
      </div>
    </section>
  )
}
