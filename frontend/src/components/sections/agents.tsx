import { Bot, FileCode2, Layers } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

export function AgentsSection() {
  return (
    <section id="agents" className="border-b border-border bg-card py-24 md:py-32">
      <div className="mx-auto max-w-6xl px-6">
        <div className="grid items-center gap-14 lg:grid-cols-2">
          <div>
            <Badge variant="info" className="mb-4 uppercase tracking-wider">
              Built for agents
            </Badge>
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

          <div className="overflow-hidden rounded-xl border border-border bg-background shadow-xs">
            <div className="flex items-center justify-between border-b border-border px-4 py-3">
              <span className="font-mono text-xs text-muted-foreground">
                scripts/chunks/Chunk_1_25.txt
              </span>
              <Badge variant="muted" className="font-mono text-[10px]">
                export
              </Badge>
            </div>
            <pre className="overflow-x-auto p-5 font-mono text-[12px] leading-relaxed text-foreground/85 md:text-[13px]">
              <code>
                <span className="text-muted-foreground"># finding 3 · PERF-042 · high</span>
                {'\n'}
                path: internal/api/handler.go:88{'\n'}
                msg: HTTP client missing Timeout{'\n'}
                {'\n'}
                <span className="text-info">func</span> fetchUpstream(ctx context.Context, url <span className="text-info">string</span>) (*http.Response, error) {'{'}
                {'\n'}
                {'  '}client := &amp;http.Client{'{'}{'}'}  <span className="text-danger">// no Timeout</span>
                {'\n'}
                {'  '}req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
                {'\n'}
                {'  '}...
                {'\n'}
                {'}'}
                {'\n\n'}
                <span className="text-muted-foreground"># finding 4 · CWE-89 · high</span>
                {'\n'}
                path: internal/store/query.go:41{'\n'}
                <span className="text-success"># full enclosing function body attached</span>
              </code>
            </pre>
          </div>
        </div>
      </div>
    </section>
  )
}
