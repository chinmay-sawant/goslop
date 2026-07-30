import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion'

const faqs = [
  {
    q: 'What is a SAT?',
    a: 'A static analysis tool inspects source without running it. goslop is a pure-Go SAT for Go codebases: it walks AST units with go/parser and go/ast (no CGO, no tree-sitter), runs detector packs, and reports findings with optional taint and export paths.',
  },
  {
    q: 'How is this different from a linter?',
    a: 'Linters often focus on style and simple correctness. goslop targets performance heuristics (239 PERF rules), CWE-class security including experimental taint (CWE-22/78/79/89), and project hygiene, with profiles, baselines, SARIF, and agent-oriented context export.',
  },
  {
    q: 'Why export whole functions for agents?',
    a: 'Agents need surrounding logic to propose safe fixes. By default each finding Context is the whole enclosing function ([goslop.export] whole_function = true). Batched chunks under scripts/chunks/ let you hand packs of findings to parallel agents without dumping the entire monorepo into one prompt.',
  },
  {
    q: 'Does scanning the goslop repo itself count as signal?',
    a: 'No. Self-scans of a SAT repo fire heavily on detector sources, rule needles, and intentional fixtures. That is the same class of noise you get when tools like Semgrep analyze their own packs. Use a real application tree (or the product reference corpus) to judge product signal. See documents/overview.md.',
  },
  {
    q: 'Can I use it in CI?',
    a: 'Yes. Profiles carry fail policies (recommended and security fail on high by default). Reporters include SARIF 2.1.0 for GitHub Code Scanning. Cache and baseline support keep incremental runs practical. Start with --profile recommended or a dedicated security/perf pack.',
  },
  {
    q: 'How do suppressions work?',
    a: 'Use inline // goslop-ignore directives (by rule, file, or block), path ignores, and .goslop-baseline.json to ship with known debt. Matching prefers fingerprint, then file/line/column. Full details are in documents/suppressions-and-baselines.md.',
  },
]

export function FaqSection() {
  return (
    <section id="faq" className="border-b border-border py-24 md:py-32">
      <div className="mx-auto max-w-3xl px-6">
        <div className="text-center">
          <p className="font-mono text-xs font-medium uppercase tracking-wider text-muted-foreground">
            FAQ
          </p>
          <h2 className="mt-3 font-heading text-4xl tracking-tight md:text-5xl">
            Straight answers
          </h2>
        </div>

        <Accordion type="single" collapsible className="mt-12">
          {faqs.map((item, i) => (
            <AccordionItem key={item.q} value={`item-${i}`}>
              <AccordionTrigger>{item.q}</AccordionTrigger>
              <AccordionContent>{item.a}</AccordionContent>
            </AccordionItem>
          ))}
        </Accordion>
      </div>
    </section>
  )
}
