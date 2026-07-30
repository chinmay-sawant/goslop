import { SiteHeader } from '@/components/site-header'
import { SiteFooter } from '@/components/site-footer'
import { HeroSection } from '@/components/sections/hero'
import { DemoSection } from '@/components/sections/demo'
import { WhySection } from '@/components/sections/why'
import { FeaturesSection } from '@/components/sections/features'
import { CiSection } from '@/components/sections/ci'
import { InstallSection } from '@/components/sections/install'
import { DocsSection } from '@/components/sections/docs'
import { FaqSection } from '@/components/sections/faq'
import { CtaSection } from '@/components/sections/cta'

export default function App() {
  return (
    <div className="min-h-svh">
      <a
        href="#main"
        className="fixed left-4 top-4 z-[100] -translate-y-20 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground opacity-0 transition focus:translate-y-0 focus:opacity-100 focus:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50"
      >
        Skip to content
      </a>
      <SiteHeader />
      <main id="main">
        <HeroSection />
        <DemoSection />
        <WhySection />
        <FeaturesSection />
        <CiSection />
        <InstallSection />
        <DocsSection />
        <FaqSection />
        <CtaSection />
      </main>
      <SiteFooter />
    </div>
  )
}
