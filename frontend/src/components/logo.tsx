import { cn } from '@/lib/utils'

type Props = {
  className?: string
  markClassName?: string
  showWordmark?: boolean
}

/** goslop mark: monogram "gs" (not the Go language logo). */
export function Logo({ className, markClassName, showWordmark = true }: Props) {
  return (
    <span className={cn('inline-flex items-center gap-2.5', className)}>
      <span
        className={cn(
          'flex size-8 shrink-0 items-center justify-center rounded-md bg-primary font-mono text-[11px] font-semibold tracking-tight text-primary-foreground',
          markClassName,
        )}
        aria-hidden
      >
        gs
      </span>
      {showWordmark && (
        <span className="text-sm font-semibold tracking-tight">goslop</span>
      )}
    </span>
  )
}
