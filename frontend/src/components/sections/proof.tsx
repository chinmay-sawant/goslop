import { Badge } from '@/components/ui/badge'
import { Reveal } from '@/components/reveal'
import { CodeBlock } from '@/components/code-block'

const BEFORE = `// line hit only (typical SAT noise)
runner.go:110:10  PERF-42  info
fmt.Errorf with a static string...`

const AFTER = `// goslop export · whole function context
Finding 24  PERF-42  runner.go:110:10
...
109: if len(durations) == 0 {
>110:   return fmt.Errorf("no successful runs")
111: }
...
// agent sees surrounding control flow`

const CI_YAML = `name: goslop
on: [push, pull_request]
jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26.x'
      - name: Build goslop
        run: |
          git clone --depth 1 https://github.com/chinmay-sawant/goslop.git /tmp/goslop
          make -C /tmp/goslop build
      - name: Scan
        run: /tmp/goslop/bin/goslop --profile recommended --format sarif . > goslop.sarif
      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: goslop.sarif`

const COMPARE = [
  {
    dim: 'Focus',
    linter: 'Style / simple correctness',
    sast: 'Security patterns (often multi-lang)',
    goslop: 'Perf + CWE + hygiene for Go',
  },
  {
    dim: 'Context',
    linter: 'Line or short span',
    sast: 'Varies; often path snippets',
    goslop: 'Whole enclosing function default',
  },
  {
    dim: 'Agent handoff',
    linter: 'Rare',
    sast: 'Usually dashboards',
    goslop: 'Functions + chunks export',
  },
  {
    dim: 'CI packaging',
    linter: 'Common',
    sast: 'Common (heavier)',
    goslop: 'Profiles + SARIF 2.1.0 + fail policy',
  },
  {
    dim: 'Runtime',
    linter: 'Fast local',
    sast: 'Often heavy / hosted',
    goslop: 'Pure Go binary, no CGO',
  },
]

export function ProofSection() {
  return (
    <section id="proof" className="border-b border-border bg-card py-20 md:py-28">
      <div className="mx-auto max-w-6xl px-6">
        <Reveal>
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="font-heading text-3xl tracking-tight sm:text-4xl md:text-5xl">
              Proof, not adjectives
            </h2>
            <p className="mt-4 text-muted-foreground text-balance">
              Line hits are not enough for agents. Export gives control flow.
              Profiles + SARIF make CI real.
            </p>
          </div>
        </Reveal>

        <Reveal>
          <div className="mt-12 grid gap-4 lg:grid-cols-2">
            <div className="rounded-xl border border-border bg-background p-5">
              <div className="mb-3 flex items-center gap-2">
                <Badge variant="muted">before</Badge>
                <span className="text-sm font-medium">Line hit only</span>
              </div>
              <pre className="overflow-x-auto font-mono text-[12px] leading-relaxed text-muted-foreground">
                <code>{BEFORE}</code>
              </pre>
            </div>
            <div className="rounded-xl border border-foreground/15 bg-background p-5 shadow-xs">
              <div className="mb-3 flex items-center gap-2">
                <Badge>after</Badge>
                <span className="text-sm font-medium">Exported function context</span>
              </div>
              <pre className="overflow-x-auto font-mono text-[12px] leading-relaxed text-foreground/90">
                <code>{AFTER}</code>
              </pre>
            </div>
          </div>
        </Reveal>

        <Reveal>
          <div className="mt-12 overflow-x-auto rounded-xl border border-border">
            <table className="w-full min-w-[40rem] border-collapse text-left text-sm">
              <thead>
                <tr className="border-b border-border bg-secondary/40">
                  <th className="px-4 py-3 font-medium">Dimension</th>
                  <th className="px-4 py-3 font-medium text-muted-foreground">
                    Style linter
                  </th>
                  <th className="px-4 py-3 font-medium text-muted-foreground">
                    Generic SAST
                  </th>
                  <th className="px-4 py-3 font-medium">goslop</th>
                </tr>
              </thead>
              <tbody>
                {COMPARE.map((row) => (
                  <tr key={row.dim} className="border-b border-border last:border-0">
                    <td className="px-4 py-3 font-medium">{row.dim}</td>
                    <td className="px-4 py-3 text-muted-foreground">{row.linter}</td>
                    <td className="px-4 py-3 text-muted-foreground">{row.sast}</td>
                    <td className="px-4 py-3 text-foreground">{row.goslop}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Reveal>

        <Reveal className="mt-12">
          <div className="grid gap-6 lg:grid-cols-[1fr_1.1fr] lg:items-start">
            <div>
              <h3 className="text-lg font-semibold">CI in one paste</h3>
              <p className="mt-2 text-sm text-muted-foreground leading-relaxed">
                Build the pure-Go binary, run{' '}
                <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[11px] text-foreground">
                  --profile recommended
                </code>
                , emit SARIF, upload to Code Scanning. Adjust clone source once
                you publish release binaries.
              </p>
              <div className="mt-4 flex flex-wrap gap-2">
                <Badge variant="outline">SARIF 2.1.0</Badge>
                <Badge variant="outline">fail on high</Badge>
                <Badge variant="outline">Code Scanning</Badge>
              </div>
            </div>
            <CodeBlock code={CI_YAML} filename=".github/workflows/goslop.yml" />
          </div>
        </Reveal>
      </div>
    </section>
  )
}
