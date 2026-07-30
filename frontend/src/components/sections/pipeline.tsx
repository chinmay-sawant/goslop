import { Reveal } from '@/components/reveal'

const STEPS = [
  { id: '01', title: 'Scan', detail: 'Profiles walk Go AST · detectors emit findings' },
  { id: '02', title: 'Filter', detail: 'Ignore directives · baseline · fail policy' },
  { id: '03', title: 'Report', detail: 'stderr summary · stdout text / JSON / SARIF' },
  { id: '04', title: 'Export', detail: 'Functions + chunks · whole-function context' },
  { id: '05', title: 'Agent', detail: 'Humans or coding agents remediate in parallel' },
]

export function PipelineSection() {
  return (
    <section
      id="pipeline"
      className="border-b border-border py-16 md:py-20"
      aria-label="Product pipeline"
    >
      <div className="mx-auto max-w-6xl px-6">
        <Reveal>
          <div className="overflow-hidden rounded-xl border border-border bg-card shadow-xs">
            <div className="border-b border-border px-5 py-3">
              <p className="font-mono text-[11px] text-muted-foreground uppercase tracking-wider">
                Signature flow
              </p>
            </div>
            <ol className="grid gap-0 sm:grid-cols-5">
              {STEPS.map((step, i) => (
                <li
                  key={step.id}
                  className={
                    i < STEPS.length - 1
                      ? 'border-b border-border px-4 py-5 sm:border-b-0 sm:border-r'
                      : 'px-4 py-5'
                  }
                >
                  <p className="font-mono text-[11px] text-muted-foreground">{step.id}</p>
                  <p className="mt-1 text-sm font-semibold">{step.title}</p>
                  <p className="mt-1 text-xs text-muted-foreground leading-relaxed">
                    {step.detail}
                  </p>
                </li>
              ))}
            </ol>
          </div>
        </Reveal>
      </div>
    </section>
  )
}
