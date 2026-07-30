import { useEffect, useId, useRef, useState } from 'react'
import { Menu, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ThemeToggle } from '@/components/theme-toggle'
import { GitHubStarsButton } from '@/components/github-stars'
import { Logo } from '@/components/logo'
import { useActiveSection } from '@/hooks/use-active-section'
import { cn } from '@/lib/utils'

export const NAV_LINKS = [
  { href: '#demo', id: 'demo', label: 'Demo' },
  { href: '#rules', id: 'rules', label: 'Rules' },
  { href: '#why', id: 'why', label: 'Why' },
  { href: '#proof', id: 'proof', label: 'Proof' },
  { href: '#ci', id: 'ci', label: 'CI' },
  { href: '#install', id: 'install', label: 'Install' },
  { href: '#docs', id: 'docs', label: 'Docs' },
] as const

const SECTION_IDS = NAV_LINKS.map((l) => l.id)

export function SiteHeader() {
  const [open, setOpen] = useState(false)
  const menuId = useId()
  const menuButtonRef = useRef<HTMLButtonElement>(null)
  const panelRef = useRef<HTMLDivElement>(null)
  const active = useActiveSection([...SECTION_IDS, 'faq', 'features'])

  useEffect(() => {
    if (!open) return

    const panel = panelRef.current
    const focusable = panel?.querySelectorAll<HTMLElement>(
      'a[href], button:not([disabled])',
    )
    focusable?.[0]?.focus()

    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        e.preventDefault()
        setOpen(false)
        menuButtonRef.current?.focus()
        return
      }
      if (e.key !== 'Tab' || !focusable?.length) return
      const list = Array.from(focusable)
      const first = list[0]
      const last = list[list.length - 1]
      if (e.shiftKey && document.activeElement === first) {
        e.preventDefault()
        last.focus()
      } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault()
        first.focus()
      }
    }

    document.addEventListener('keydown', onKey)
    document.body.style.overflow = 'hidden'
    return () => {
      document.removeEventListener('keydown', onKey)
      document.body.style.overflow = ''
    }
  }, [open])

  function closeMenu(restoreFocus = false) {
    setOpen(false)
    if (restoreFocus) menuButtonRef.current?.focus()
  }

  return (
    <header className="sticky top-0 z-50 border-b border-border/80 bg-background/85 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between gap-3 px-6">
        <a
          href="#top"
          className="shrink-0 rounded-md focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50"
          onClick={() => closeMenu()}
        >
          <Logo />
          <span className="sr-only">goslop home</span>
        </a>

        <nav className="hidden items-center gap-0.5 xl:flex" aria-label="Primary">
          {NAV_LINKS.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className={cn(
                'rounded-md px-2 py-1.5 text-sm transition-colors',
                active === link.id
                  ? 'bg-secondary font-medium text-foreground'
                  : 'text-muted-foreground hover:text-foreground',
              )}
            >
              {link.label}
            </a>
          ))}
        </nav>

        <div className="flex items-center gap-1.5 sm:gap-2">
          <ThemeToggle />
          <GitHubStarsButton className="hidden sm:inline-flex" />
          <Button size="sm" asChild className="hidden sm:inline-flex">
            <a href="#install">Install</a>
          </Button>
          <Button
            ref={menuButtonRef}
            type="button"
            variant="outline"
            size="icon"
            className="size-9 xl:hidden"
            aria-expanded={open}
            aria-controls={menuId}
            aria-label={open ? 'Close menu' : 'Open menu'}
            onClick={() => {
              if (open) closeMenu(true)
              else setOpen(true)
            }}
          >
            {open ? <X className="size-4" /> : <Menu className="size-4" />}
          </Button>
        </div>
      </div>

      {open && (
        <div
          id={menuId}
          ref={panelRef}
          className="border-t border-border bg-background xl:hidden"
          role="dialog"
          aria-modal="true"
          aria-label="Mobile navigation"
        >
          <nav
            className="mx-auto flex max-w-6xl flex-col gap-1 px-4 py-3"
            aria-label="Mobile"
          >
            {NAV_LINKS.map((link) => (
              <a
                key={link.href}
                href={link.href}
                onClick={() => closeMenu()}
                className={cn(
                  'rounded-md px-3 py-3 text-sm',
                  active === link.id
                    ? 'bg-secondary font-medium text-foreground'
                    : 'text-muted-foreground hover:bg-secondary/60 hover:text-foreground',
                )}
              >
                {link.label}
              </a>
            ))}
            <a
              href="#faq"
              onClick={() => closeMenu()}
              className="rounded-md px-3 py-3 text-sm text-muted-foreground hover:bg-secondary/60 hover:text-foreground"
            >
              FAQ
            </a>
            <div className="mt-2 flex flex-col gap-2 border-t border-border pt-3 sm:hidden">
              <GitHubStarsButton className="w-full justify-center" />
              <Button size="sm" asChild className="w-full">
                <a href="#install" onClick={() => closeMenu()}>
                  Install
                </a>
              </Button>
            </div>
          </nav>
        </div>
      )}
    </header>
  )
}
