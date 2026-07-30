import { useEffect, useState } from 'react'

/**
 * Tracks which section id is most visible. Uses IntersectionObserver
 * (no scroll listeners) so it stays cheap on long pages.
 */
export function useActiveSection(ids: string[], rootMargin = '-30% 0px -55% 0px') {
  const [active, setActive] = useState<string>(ids[0] ?? '')

  useEffect(() => {
    if (ids.length === 0) return

    const ratios = new Map<string, number>()
    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          const id = (entry.target as HTMLElement).id
          if (!id) continue
          ratios.set(id, entry.isIntersecting ? entry.intersectionRatio : 0)
        }
        let bestId = ids[0]
        let best = -1
        for (const id of ids) {
          const r = ratios.get(id) ?? 0
          if (r > best) {
            best = r
            bestId = id
          }
        }
        if (best > 0) setActive(bestId)
      },
      { root: null, rootMargin, threshold: [0, 0.15, 0.35, 0.55, 0.75, 1] },
    )

    for (const id of ids) {
      const el = document.getElementById(id)
      if (el) observer.observe(el)
    }

    return () => observer.disconnect()
  }, [ids, rootMargin])

  return active
}
