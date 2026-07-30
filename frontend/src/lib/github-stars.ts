/** GitHub stars with a hard 2-minute client cache to avoid API abuse. */

export const GITHUB_OWNER = 'chinmay-sawant'
export const GITHUB_REPO = 'goslop'
export const GITHUB_URL = `https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}`
export const GITHUB_API_URL = `https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}`

const CACHE_KEY = 'goslop:github-stars:v1'
/** Minimum time between network fetches (client-side rate limit). */
export const STARS_CACHE_TTL_MS = 2 * 60 * 1000

export type StarsCache = {
  stars: number
  fetchedAt: number
}

export type StarsResult = {
  stars: number | null
  fromCache: boolean
  error?: string
}

/** In-flight promise shared across callers in the same tab. */
let inFlight: Promise<StarsResult> | null = null

function readCache(): StarsCache | null {
  try {
    const raw = localStorage.getItem(CACHE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as StarsCache
    if (
      typeof parsed.stars !== 'number' ||
      typeof parsed.fetchedAt !== 'number' ||
      !Number.isFinite(parsed.stars) ||
      !Number.isFinite(parsed.fetchedAt)
    ) {
      return null
    }
    return parsed
  } catch {
    return null
  }
}

function writeCache(entry: StarsCache) {
  try {
    localStorage.setItem(CACHE_KEY, JSON.stringify(entry))
  } catch {
    // private mode / quota
  }
}

function isFresh(entry: StarsCache, now = Date.now()): boolean {
  return now - entry.fetchedAt < STARS_CACHE_TTL_MS
}

/**
 * Returns star count. Uses localStorage cache for 2 minutes and
 * dedupes concurrent requests so remounts cannot hammer the API.
 */
export async function getGitHubStars(): Promise<StarsResult> {
  const cached = readCache()
  if (cached && isFresh(cached)) {
    return { stars: cached.stars, fromCache: true }
  }

  if (inFlight) return inFlight

  inFlight = (async (): Promise<StarsResult> => {
    // Re-check cache in case another tab wrote while we waited.
    const again = readCache()
    if (again && isFresh(again)) {
      return { stars: again.stars, fromCache: true }
    }

    try {
      const res = await fetch(GITHUB_API_URL, {
        headers: {
          Accept: 'application/vnd.github+json',
        },
        // Prefer cache when the browser still has a valid HTTP cache entry.
        cache: 'default',
      })

      if (res.status === 403 || res.status === 429) {
        // Rate limited: serve stale cache if any, otherwise soft-fail.
        if (again) {
          return {
            stars: again.stars,
            fromCache: true,
            error: 'rate_limited',
          }
        }
        return { stars: null, fromCache: false, error: 'rate_limited' }
      }

      if (!res.ok) {
        if (again) {
          return { stars: again.stars, fromCache: true, error: `http_${res.status}` }
        }
        return { stars: null, fromCache: false, error: `http_${res.status}` }
      }

      const data = (await res.json()) as { stargazers_count?: number }
      const stars = data.stargazers_count
      if (typeof stars !== 'number') {
        if (again) {
          return { stars: again.stars, fromCache: true, error: 'bad_payload' }
        }
        return { stars: null, fromCache: false, error: 'bad_payload' }
      }

      writeCache({ stars, fetchedAt: Date.now() })
      return { stars, fromCache: false }
    } catch {
      if (again) {
        return { stars: again.stars, fromCache: true, error: 'network' }
      }
      return { stars: null, fromCache: false, error: 'network' }
    } finally {
      inFlight = null
    }
  })()

  return inFlight
}

export function formatStarCount(n: number): string {
  if (n >= 1000) {
    const k = n / 1000
    return k >= 10 ? `${Math.round(k)}k` : `${k.toFixed(1).replace(/\.0$/, '')}k`
  }
  return String(n)
}
