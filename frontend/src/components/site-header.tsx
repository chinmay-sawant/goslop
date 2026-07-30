import { Button } from '@/components/ui/button'
import { ThemeToggle } from '@/components/theme-toggle'
import { GitHubStarsButton } from '@/components/github-stars'

const links = [
  { href: '#features', label: 'Features' },
  { href: '#agents', label: 'Agents' },
  { href: '#profiles', label: 'Profiles' },
  { href: '#docs', label: 'Docs' },
  { href: '#install', label: 'Install' },
  { href: '#faq', label: 'FAQ' },
]

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-50 border-b border-border/80 bg-background/80 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between gap-4 px-6">
        <a href="#top" className="flex shrink-0 items-center gap-2.5">
          <span className="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground font-mono text-xs font-semibold tracking-tight">
            go
          </span>
          <span className="text-sm font-semibold tracking-tight">goslop</span>
        </a>

        <nav className="hidden items-center gap-6 lg:flex">
          {links.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className="text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
              {link.label}
            </a>
          ))}
        </nav>

        <div className="flex items-center gap-1.5 sm:gap-2">
          <ThemeToggle />
          <GitHubStarsButton className="hidden sm:inline-flex" />
          <Button size="sm" asChild>
            <a href="#install">Get started</a>
          </Button>
        </div>
      </div>
    </header>
  )
}
