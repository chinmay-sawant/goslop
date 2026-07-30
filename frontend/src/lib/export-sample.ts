/** Shared authentic export sample (Chunk_1_25 finding 24 · PERF-42 TP). */

export const FINDING = {
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

export type ContextLine = {
  n: number
  code: string
  hit?: boolean
  ellipsis?: boolean
}

export const CONTEXT_LINES: ContextLine[] = [
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

/**
 * Real stderr summary from:
 * ./bin/goslop --profile all --no-fail --no-terminal --export-context --export-chunks --no-cache gopdfsuit
 * (feat/product-website, wall ~100.4ms, 915 findings)
 */
export const SCAN_SUMMARY = `scanned 78 files (28042 lines) in 100.4ms
  cache: 0 hits, 78 misses (full re-analysis)
  skipped 383 files
915 findings
  severity: 10 high, 197 info, 312 low, 396 medium
  top rules: BP-1 ×181, PERF-6 ×94, PERF-32 ×59, BP-5 ×50, PERF-230 ×44
  example findings: 63 (of 915 total)
exported 915 context file(s) to scripts/findings/functions; exported 37 chunk file(s) to scripts/chunks`

/** Full text for clipboard (function dossier shape). */
export function buildFunctionDossierText(): string {
  const lines = CONTEXT_LINES.map((l) => {
    if (l.ellipsis) return `         ${l.code}`
    const mark = l.hit ? '>' : ' '
    const num = String(l.n).padStart(6, ' ')
    return `${mark}${num}: ${l.code}`
  }).join('\n')

  return [
    `Finding ${FINDING.index}/${FINDING.total}`,
    `Source: ${FINDING.source}`,
    `Rule: ${FINDING.rule}`,
    `Rule title: ${FINDING.title}`,
    `Severity: ${FINDING.severity}`,
    `Message: ${FINDING.message}`,
    `Fix: ${FINDING.fix}`,
    'Context:',
    lines,
  ].join('\n')
}

export function buildChunkDossierText(): string {
  return [
    `Findings 1-25 of ${FINDING.total}`,
    '',
    buildFunctionDossierText(),
    '',
    '='.repeat(48),
  ].join('\n')
}

/** Command that produced the demo scan numbers (reference corpus). */
export const SCAN_COMMAND =
  './bin/goslop --profile all --no-fail --no-terminal --export-context --export-chunks --no-cache ./gopdfsuit'

export function buildScanText(): string {
  return [
    `$ ${SCAN_COMMAND}`,
    '',
    '# stderr · scan summary',
    SCAN_SUMMARY,
  ].join('\n')
}
