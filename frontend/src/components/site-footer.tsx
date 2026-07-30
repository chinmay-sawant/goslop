import { Logo } from '@/components/logo'
import { GitHubStarsButton } from '@/components/github-stars'

export function SiteFooter() {
  return (
    <footer className="border-t border-border bg-card">
      <div className="mx-auto flex max-w-6xl flex-col gap-8 px-6 py-12 md:flex-row md:items-start md:justify-between">
        <div className="max-w-sm space-y-3">
          <Logo />
          <p className="text-sm text-muted-foreground leading-relaxed">
            Pure-Go static analysis for performance, security, and hygiene.
            Built to feed findings into humans and agents.
          </p>
          <GitHubStarsButton size="sm" />
        </div>

        <div className="grid grid-cols-2 gap-10 sm:grid-cols-3">
          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Product
            </p>
            <ul className="space-y-2 text-sm">
              <li>
                <a href="#demo" className="text-muted-foreground hover:text-foreground">
                  Demo
                </a>
              </li>
              <li>
                <a href="#rules" className="text-muted-foreground hover:text-foreground">
                  Rules
                </a>
              </li>
              <li>
                <a href="#proof" className="text-muted-foreground hover:text-foreground">
                  Proof
                </a>
              </li>
              <li>
                <a href="#install" className="text-muted-foreground hover:text-foreground">
                  Install
                </a>
              </li>
            </ul>
          </div>
          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Docs
            </p>
            <ul className="space-y-2 text-sm">
              <li>
                <a
                  href="https://github.com/chinmay-sawant/goslop/blob/main/documents/overview.md"
                  target="_blank"
                  rel="noreferrer"
                  className="text-muted-foreground hover:text-foreground"
                >
                  Overview
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/chinmay-sawant/goslop/blob/main/documents/cli-reference.md"
                  target="_blank"
                  rel="noreferrer"
                  className="text-muted-foreground hover:text-foreground"
                >
                  CLI
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/chinmay-sawant/goslop/blob/main/documents/export-context-and-chunks.md"
                  target="_blank"
                  rel="noreferrer"
                  className="text-muted-foreground hover:text-foreground"
                >
                  Export
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/chinmay-sawant/goslop/blob/main/documents/reporting-formats.md"
                  target="_blank"
                  rel="noreferrer"
                  className="text-muted-foreground hover:text-foreground"
                >
                  Reporting
                </a>
              </li>
            </ul>
          </div>
          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Source
            </p>
            <ul className="space-y-2 text-sm">
              <li>
                <a
                  href="https://github.com/chinmay-sawant/goslop"
                  target="_blank"
                  rel="noreferrer"
                  className="text-muted-foreground hover:text-foreground"
                >
                  GitHub
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/chinmay-sawant/goslop/blob/main/LICENSE"
                  target="_blank"
                  rel="noreferrer"
                  className="text-muted-foreground hover:text-foreground"
                >
                  License
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <div className="border-t border-border">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-5 text-xs text-muted-foreground">
          <span>Static analysis tool for Go</span>
          <span className="font-mono">goslop</span>
        </div>
      </div>
    </footer>
  )
}
