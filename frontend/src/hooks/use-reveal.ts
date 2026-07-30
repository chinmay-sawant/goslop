import { useEffect, useRef, useState } from 'react'

/** Adds a class when the element enters the viewport (once). */
export function useReveal<T extends HTMLElement = HTMLDivElement>(
  rootMargin = '0px 0px -8% 0px',
) {
  const ref = useRef<T | null>(null)
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    const el = ref.current
    if (!el || visible) return

    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      setVisible(true)
      return
    }

    const obs = new IntersectionObserver(
      ([entry]) => {
        if (entry?.isIntersecting) {
          setVisible(true)
          obs.disconnect()
        }
      },
      { root: null, rootMargin, threshold: 0.12 },
    )
    obs.observe(el)
    return () => obs.disconnect()
  }, [rootMargin, visible])

  return { ref, visible }
}
