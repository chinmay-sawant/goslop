import type { ReactNode } from 'react'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion'

const DOCS =
  'https://github.com/chinmay-sawant/goslop/blob/main/documents'

function FaqAnswer({ children }: { children: ReactNode }) {
  return <div className="space-y-5">{children}</div>
}

function FaqLead({ children }: { children: ReactNode }) {
  return <p className="leading-relaxed">{children}</p>
}

function FaqBlock({
  title,
  children,
}: {
  title: string
  children: ReactNode
}) {
  return (
    <div className="space-y-2.5">
      <p className="text-sm font-medium text-foreground">{title}</p>
      {children}
    </div>
  )
}

function CodeInline({ children }: { children: ReactNode }) {
  return (
    <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
      {children}
    </code>
  )
}

function CodeRow({
  code,
  description,
}: {
  code: string
  description: string
}) {
  return (
    <li className="flex flex-col gap-1 sm:flex-row sm:items-baseline sm:gap-3">
      <code className="shrink-0 rounded bg-secondary px-1.5 py-0.5 font-mono text-[12px] text-foreground">
        {code}
      </code>
      <span className="leading-relaxed">{description}</span>
    </li>
  )
}

function CodeList({ children }: { children: ReactNode }) {
  return <ul className="list-none space-y-2.5 pl-0">{children}</ul>
}

function CodeBlock({ children }: { children: string }) {
  return (
    <pre className="overflow-x-auto rounded-lg border border-border bg-background px-3.5 py-3 font-mono text-[12px] leading-relaxed text-foreground">
      <code>{children}</code>
    </pre>
  )
}

function DocLink({ href, children }: { href: string; children: ReactNode }) {
  return (
    <p className="leading-relaxed">
      Full details:{' '}
      <a
        href={href}
        target="_blank"
        rel="noreferrer"
        className="font-medium text-foreground underline underline-offset-4 hover:text-info"
      >
        {children}
      </a>
    </p>
  )
}

type FaqItem = {
  q: string
  a: ReactNode
}

const faqs: FaqItem[] = [
  {
    q: 'What is a SAT?',
    a: (
      <FaqAnswer>
        <FaqLead>
          A static analysis tool (SAT) inspects source without running the
          program. goslop is a pure-Go SAT for Go codebases: it walks AST units
          and reports issues across performance, security, and hygiene.
        </FaqLead>

        <FaqBlock title="How goslop analyzes">
          <CodeList>
            <CodeRow
              code="go/parser + go/ast"
              description="pure Go parsing, no CGO and no tree-sitter"
            />
            <CodeRow
              code="walk → parse → detectors"
              description="files are scanned, then rule packs emit findings"
            />
            <CodeRow
              code="ignore → baseline → report"
              description="filters run before text, JSON, or SARIF output"
            />
          </CodeList>
        </FaqBlock>

        <FaqBlock title="What it finds">
          <CodeList>
            <CodeRow
              code="PERF-*"
              description="239 hot-path rules (loops, allocations, HTTP timeouts, frameworks)"
            />
            <CodeRow
              code="CWE-*"
              description="175 structural security heuristics"
            />
            <CodeRow
              code="Taint (experimental)"
              description="CWE-22 / 78 / 79 / 89 data-flow for injection classes"
            />
            <CodeRow
              code="BP-*"
              description="bad practices and project-level hygiene"
            />
          </CodeList>
        </FaqBlock>

        <FaqBlock title="Outputs">
          <p className="leading-relaxed">
            Findings on <CodeInline>stdout</CodeInline> (text, JSON, or SARIF
            2.1.0). Scan summary always on <CodeInline>stderr</CodeInline>.
            Optional disk export for agent handoff.
          </p>
        </FaqBlock>

        <DocLink href={`${DOCS}/overview.md`}>documents/overview.md</DocLink>
      </FaqAnswer>
    ),
  },
  {
    q: 'How is this different from a linter?',
    a: (
      <FaqAnswer>
        <FaqLead>
          Linters usually target style and simple correctness. goslop is a
          product SAT: curated detector catalogues, CI fail policies, machine
          reporters, debt controls, and exports shaped for humans and agents.
        </FaqLead>

        <FaqBlock title="Broader than style">
          <CodeList>
            <CodeRow
              code="Performance"
              description="hot-path and runtime footguns under load"
            />
            <CodeRow
              code="Security"
              description="CWE heuristics plus optional inter-procedural taint"
            />
            <CodeRow
              code="Hygiene"
              description="error handling, servers, go.mod, and project patterns"
            />
          </CodeList>
        </FaqBlock>

        <FaqBlock title="Product surfaces a linter often lacks">
          <CodeList>
            <CodeRow
              code="Profiles"
              description="recommended · perf · security · style · all"
            />
            <CodeRow
              code="Reporters"
              description="text, JSON, SARIF 2.1.0 for Code Scanning"
            />
            <CodeRow
              code="Cache + baseline"
              description="fast re-scans and controlled debt rollout"
            />
            <CodeRow
              code="Agent export"
              description="whole-function context and batched chunks"
            />
          </CodeList>
        </FaqBlock>

        <DocLink href={`${DOCS}/overview.md`}>documents/overview.md</DocLink>
      </FaqAnswer>
    ),
  },
  {
    q: 'Why export whole functions for agents?',
    a: (
      <FaqAnswer>
        <FaqLead>
          Agents need surrounding logic to propose safe fixes. Line-only hits
          force the model to hunt the tree. goslop attaches the enclosing
          function by default and can batch findings for parallel delegation.
        </FaqLead>

        <FaqBlock title="Two export surfaces">
          <CodeList>
            <CodeRow
              code="--export-context"
              description="one finding → scripts/findings/functions/N.txt"
            />
            <CodeRow
              code="--export-chunks"
              description="N findings (default 25) → scripts/chunks/Chunk_*.txt"
            />
          </CodeList>
        </FaqBlock>

        <FaqBlock title="When to use which">
          <CodeList>
            <CodeRow
              code="Single issue / ticket"
              description="open a numbered function ref"
            />
            <CodeRow
              code="Parallel agent work"
              description="hand a chunk batch to a subagent"
            />
          </CodeList>
        </FaqBlock>

        <FaqBlock title="Default context mode">
          <p className="leading-relaxed">
            Context is the <strong className="font-medium text-foreground">whole enclosing function</strong>{' '}
            (<CodeInline>[goslop.export] whole_function = true</CodeInline>).
            Set it to <CodeInline>false</CodeInline> only if you want the older
            nearby ~4-line window.
          </p>
        </FaqBlock>

        <FaqBlock title="Example">
          <CodeBlock>{`./bin/goslop --profile all \\
  --export-context --export-chunks --no-cache .`}</CodeBlock>
        </FaqBlock>

        <DocLink href={`${DOCS}/export-context-and-chunks.md`}>
          documents/export-context-and-chunks.md
        </DocLink>
      </FaqAnswer>
    ),
  },
  {
    q: 'Can I use it in CI?',
    a: (
      <FaqAnswer>
        <FaqLead>
          Yes. Profiles carry fail policies, reporters include SARIF for GitHub
          Code Scanning, and cache plus baseline keep incremental runs practical.
        </FaqLead>

        <FaqBlock title="Suggested profiles">
          <CodeList>
            <CodeRow
              code="recommended (default)"
              description="everyday CI gate · fail on high"
            />
            <CodeRow
              code="security"
              description="CWE pack with taint on · fail on high"
            />
            <CodeRow
              code="perf"
              description="S + A performance catalogue · fail on high"
            />
            <CodeRow
              code="style"
              description="BP hygiene · soft gate (fail none by default)"
            />
          </CodeList>
        </FaqBlock>

        <FaqBlock title="Machine output">
          <CodeList>
            <CodeRow
              code="--format sarif"
              description="GitHub Code Scanning / SARIF viewers"
            />
            <CodeRow
              code="--format json"
              description="scripts, dashboards, automation"
            />
            <CodeRow
              code="stderr summary"
              description="always separate from stdout findings (pipe-clean)"
            />
          </CodeList>
        </FaqBlock>

        <FaqBlock title="CI-friendly controls">
          <CodeList>
            <CodeRow
              code=".goslop-cache/"
              description="incremental per-file cache"
            />
            <CodeRow
              code=".goslop-baseline.json"
              description="ship known debt without failing day one"
            />
            <CodeRow
              code="// goslop-ignore"
              description="narrow reviewed exceptions in source"
            />
          </CodeList>
        </FaqBlock>

        <FaqBlock title="Example">
          <CodeBlock>{`./bin/goslop --profile recommended .
./bin/goslop --format sarif . > goslop.sarif`}</CodeBlock>
        </FaqBlock>

        <DocLink href={`${DOCS}/reporting-formats.md`}>
          documents/reporting-formats.md
        </DocLink>
      </FaqAnswer>
    ),
  },
  {
    q: 'How do suppressions work?',
    a: (
      <FaqAnswer>
        <FaqLead>
          Use a narrow inline suppression for a reviewed local exception. Use a
          baseline for existing debt that should stay visible without failing an
          initial rollout. Neither replaces fixing a confirmed issue.
        </FaqLead>

        <FaqBlock title="Inline directives">
          <CodeList>
            <CodeRow
              code="// goslop-ignore: RULE"
              description="next line, or same-line suffix"
            />
            <CodeRow
              code="// goslop-ignore-start / -end"
              description="block of lines (CSV rule list or all)"
            />
            <CodeRow
              code="// goslop-ignore-file: RULES"
              description="file scope (first 20 lines)"
            />
          </CodeList>
        </FaqBlock>

        <FaqBlock title="Baseline">
          <p className="leading-relaxed">
            Load <CodeInline>.goslop-baseline.json</CodeInline> (discovered
            upward from the scan root, or via{' '}
            <CodeInline>--baseline-file</CodeInline>). Matching prefers
            fingerprint, then file / line / column.
          </p>
        </FaqBlock>

        <FaqBlock title="Inspect">
          <CodeBlock>{`./bin/goslop --show-ignored .
./bin/goslop --show-baselined .`}</CodeBlock>
        </FaqBlock>

        <DocLink href={`${DOCS}/suppressions-and-baselines.md`}>
          documents/suppressions-and-baselines.md
        </DocLink>
      </FaqAnswer>
    ),
  },
]

export function FaqSection() {
  return (
    <section id="faq" className="border-b border-border py-24 md:py-32">
      <div className="mx-auto max-w-3xl px-6">
        <div className="text-center">
          <h2 className="font-heading text-4xl tracking-tight md:text-5xl">
            Straight answers
          </h2>
        </div>

        <Accordion type="single" collapsible className="mt-12">
          {faqs.map((item, i) => (
            <AccordionItem key={item.q} value={`item-${i}`}>
              <AccordionTrigger>{item.q}</AccordionTrigger>
              <AccordionContent className="text-left">{item.a}</AccordionContent>
            </AccordionItem>
          ))}
        </Accordion>
      </div>
    </section>
  )
}
