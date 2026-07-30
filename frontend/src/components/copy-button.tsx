import { useState } from 'react'
import { Check, Copy } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { copyText } from '@/lib/copy'
import { cn } from '@/lib/utils'

type Props = {
  text: string
  className?: string
  label?: string
}

export function CopyButton({ text, className, label = 'Copy' }: Props) {
  const [copied, setCopied] = useState(false)

  async function onCopy() {
    const ok = await copyText(text)
    if (!ok) return
    setCopied(true)
    window.setTimeout(() => setCopied(false), 1600)
  }

  return (
    <Button
      type="button"
      variant="outline"
      size="sm"
      onClick={onCopy}
      className={cn('h-8 gap-1.5 px-2.5 font-mono text-[11px]', className)}
      aria-label={copied ? 'Copied' : label}
    >
      {copied ? <Check className="size-3.5" /> : <Copy className="size-3.5" />}
      {copied ? 'Copied' : label}
    </Button>
  )
}
