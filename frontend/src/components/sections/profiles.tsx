import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'

const profiles = [
  {
    name: 'recommended',
    aliases: 'ci · default',
    summary:
      'S-tier PERF + taint-core CWE IDs. Fail on high. Default for everyday use and CI.',
    taint: false,
    bp: false,
    fail: 'high',
  },
  {
    name: 'perf',
    aliases: 'performance',
    summary: 'S + A tier performance catalogue. Fail on high severity findings.',
    taint: false,
    bp: false,
    fail: 'high',
  },
  {
    name: 'security',
    aliases: 'sec',
    summary:
      'Security CWE pack with taint on (depth 4) for injection-class issues.',
    taint: true,
    bp: false,
    fail: 'high',
  },
  {
    name: 'style',
    aliases: 'bp · bad-practices',
    summary:
      'BP-* hygiene rules (skips BP-21/28/30 by default). Soft gate: fail policy none.',
    taint: false,
    bp: true,
    fail: 'none',
  },
  {
    name: 'all',
    aliases: 'full',
    summary: 'Full catalogue for deep audits. BP on, medium fail threshold.',
    taint: false,
    bp: true,
    fail: 'medium',
  },
]

export function ProfilesSection() {
  return (
    <section id="profiles" className="border-b border-border py-24 md:py-32">
      <div className="mx-auto max-w-6xl px-6">
        <div className="mx-auto max-w-2xl text-center">
          <p className="font-mono text-xs font-medium uppercase tracking-wider text-muted-foreground">
            Profiles
          </p>
          <h2 className="mt-3 font-heading text-4xl tracking-tight md:text-5xl">
            Pick a pack, ship a gate
          </h2>
          <p className="mt-4 text-muted-foreground text-balance">
            Curated rule surfaces and fail policies so CI and local runs stay
            intentional, not a firehose. Default is{' '}
            <code className="font-mono text-xs bg-secondary px-1.5 py-0.5 rounded">
              --profile recommended
            </code>
            .
          </p>
        </div>

        <div className="mt-14 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {profiles.map((profile, i) => (
            <Card
              key={profile.name}
              className={i === 0 ? 'border-foreground/15' : ''}
            >
              <CardHeader>
                <div className="flex flex-wrap items-center gap-2">
                  <CardTitle className="font-mono text-base font-semibold">
                    {profile.name}
                  </CardTitle>
                  {i === 0 && <Badge>default</Badge>}
                </div>
                <p className="font-mono text-[11px] text-muted-foreground">
                  {profile.aliases}
                </p>
                <CardDescription className="mt-2">{profile.summary}</CardDescription>
                <div className="mt-4 flex flex-wrap gap-2">
                  <Badge variant={profile.taint ? 'info' : 'muted'}>
                    taint {profile.taint ? 'on' : 'off'}
                  </Badge>
                  <Badge variant={profile.bp ? 'warning' : 'muted'}>
                    BP {profile.bp ? 'on' : 'off'}
                  </Badge>
                  <Badge variant="outline" className="font-mono">
                    fail · {profile.fail}
                  </Badge>
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>

        <div className="mt-8 flex justify-center">
          <Button variant="outline" asChild>
            <a
              href="https://github.com/chinmay-sawant/goslop/blob/main/documents/go-recommended-pack.md"
              target="_blank"
              rel="noreferrer"
            >
              Recommended pack details
            </a>
          </Button>
        </div>
      </div>
    </section>
  )
}
