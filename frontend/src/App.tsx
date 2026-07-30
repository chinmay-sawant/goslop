import { SiteHeader } from '@/components/site-header'
import { SiteFooter } from '@/components/site-footer'
import { HeroSection } from '@/components/sections/hero'
import { StatsSection } from '@/components/sections/stats'
import { FeaturesSection } from '@/components/sections/features'
import { AgentsSection } from '@/components/sections/agents'
import { ProfilesSection } from '@/components/sections/profiles'
import { DocsSection } from '@/components/sections/docs'
import { InstallSection } from '@/components/sections/install'
import { FaqSection } from '@/components/sections/faq'
import { CtaSection } from '@/components/sections/cta'

export default function App() {
  return (
    <div className="min-h-svh">
      <SiteHeader />
      <main>
        <HeroSection />
        <StatsSection />
        <FeaturesSection />
        <AgentsSection />
        <ProfilesSection />
        <DocsSection />
        <InstallSection />
        <FaqSection />
        <CtaSection />
      </main>
      <SiteFooter />
    </div>
  )
}
