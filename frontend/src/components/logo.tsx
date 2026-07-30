import { cn } from '@/lib/utils'

type Props = {
  className?: string
  markClassName?: string
  showWordmark?: boolean
}

/**
 * goslop mark: stylized "g" with a scan notch — distinct from the Go gopher.
 * Works at favicon size.
 */
export function Logo({ className, markClassName, showWordmark = true }: Props) {
  return (
    <span className={cn('inline-flex items-center gap-2.5', className)}>
      <span
        className={cn(
          'relative flex size-8 shrink-0 items-center justify-center overflow-hidden rounded-md bg-primary text-primary-foreground',
          markClassName,
        )}
        aria-hidden
      >
        <svg viewBox="0 0 32 32" className="size-[22px]" fill="none">
          {/* outer scan frame */}
          <path
            d="M8 10.5h4.5M8 10.5v4.5M24 10.5h-4.5M24 10.5v4.5M8 21.5h4.5M8 21.5v-4.5M24 21.5h-4.5M24 21.5v-4.5"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            opacity="0.45"
          />
          {/* monogram g */}
          <path
            d="M18.2 12.2c-.7-.6-1.6-.9-2.7-.9-2.5 0-4.3 1.8-4.3 4.7s1.8 4.7 4.3 4.7c1.2 0 2.2-.4 2.9-1.1v.4h1.7v-5.9h-4.4v1.5h2.5v1.1c-.4.5-1.1.8-1.9.8-1.4 0-2.4-1-2.4-2.5s1-2.5 2.4-2.5c.8 0 1.5.3 1.9.8l1-.8z"
            fill="currentColor"
          />
          {/* agent tick */}
          <path
            d="M20.5 19.2l1.2 1.2 2.3-2.5"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
            opacity="0.9"
          />
        </svg>
      </span>
      {showWordmark && (
        <span className="text-sm font-semibold tracking-tight">goslop</span>
      )}
    </span>
  )
}
