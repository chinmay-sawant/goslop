import { useEffect, useState } from 'react'
import { Star } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  formatStarCount,
  getGitHubStars,
  GITHUB_URL,
} from '@/lib/github-stars'
import { cn } from '@/lib/utils'

type Props = {
  className?: string
  size?: 'sm' | 'default'
  showLabel?: boolean
}

export function GitHubStarsButton({
  className,
  size = 'sm',
  showLabel = true,
}: Props) {
  const [stars, setStars] = useState<number | null>(null)
  const [loaded, setLoaded] = useState(false)

  useEffect(() => {
    let cancelled = false

    getGitHubStars().then((result) => {
      if (cancelled) return
      setStars(result.stars)
      setLoaded(true)
    })

    return () => {
      cancelled = true
    }
  }, [])

  return (
    <Button
      variant="outline"
      size={size}
      asChild
      className={cn('gap-1.5 font-normal', className)}
    >
      <a href={GITHUB_URL} target="_blank" rel="noreferrer">
        <svg
          viewBox="0 0 16 16"
          className="size-3.5 shrink-0 fill-current"
          aria-hidden
        >
          <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" />
        </svg>
        {showLabel && <span>Star</span>}
        <span
          className={cn(
            'inline-flex items-center gap-1 border-l border-border pl-1.5 font-mono text-xs tabular-nums',
            !loaded && 'text-muted-foreground',
          )}
          aria-live="polite"
        >
          <Star className="size-3 fill-current opacity-80" aria-hidden />
          {loaded && stars !== null ? formatStarCount(stars) : '...'}
        </span>
      </a>
    </Button>
  )
}
