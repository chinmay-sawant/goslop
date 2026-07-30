import type { ReactNode } from 'react'
import { CopyButton } from '@/components/copy-button'
import { cn } from '@/lib/utils'

type Props = {
  code: string
  className?: string
  filename?: string
  /** Optional pre-highlighted nodes; falls back to plain code text. */
  children?: ReactNode
}

export function CodeBlock({ code, className, filename, children }: Props) {
  return (
    <div
      className={cn(
        'overflow-hidden rounded-xl border border-border bg-background shadow-xs',
        className,
      )}
    >
      <div className="flex items-center justify-between gap-3 border-b border-border bg-secondary/40 px-3 py-2">
        <span className="truncate font-mono text-[11px] text-muted-foreground">
          {filename ?? 'terminal'}
        </span>
        <CopyButton text={code} />
      </div>
      <pre className="overflow-x-auto p-4 font-mono text-[12px] leading-relaxed text-foreground/90 md:text-[13px]">
        <code>{children ?? code}</code>
      </pre>
    </div>
  )
}
