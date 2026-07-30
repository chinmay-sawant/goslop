const stats = [
  { value: '239', label: 'PERF rules', hint: 'Hot-path & runtime' },
  { value: '175', label: 'CWE rules', hint: 'Structural security' },
  { value: '3', label: 'Reporters', hint: 'Text · JSON · SARIF' },
  { value: '5', label: 'Profiles', hint: 'Curated packs' },
]

export function StatsSection() {
  return (
    <section className="border-b border-border bg-card">
      <div className="mx-auto grid max-w-6xl grid-cols-2 divide-x divide-y divide-border md:grid-cols-4 md:divide-y-0">
        {stats.map((stat) => (
          <div key={stat.label} className="px-6 py-10 text-center md:py-12">
            <p className="font-heading text-4xl tracking-tight md:text-5xl">
              {stat.value}
            </p>
            <p className="mt-2 text-sm font-medium">{stat.label}</p>
            <p className="mt-1 text-xs text-muted-foreground">{stat.hint}</p>
          </div>
        ))}
      </div>
    </section>
  )
}
