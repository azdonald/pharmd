import { useNavigate } from "react-router-dom";

export default function Landing() {
  const navigate = useNavigate();

  return (
    <div className="bg-background text-on-surface overflow-x-hidden">
      {/* TopNavBar */}
      <nav className="fixed top-0 w-full z-50 bg-surface-container-lowest shadow-sm transition-all duration-300">
        <div className="max-w-7xl mx-auto px-container-padding py-4 flex justify-between items-center">
          <div className="font-headline-md text-headline-md font-bold text-primary">PharmaFlow</div>
          <div className="hidden md:flex items-center gap-8">
            <a className="text-primary font-bold border-b-2 border-primary pb-1" href="#features">Features</a>
            <a className="text-on-surface-variant hover:text-primary" href="#how-it-works">How It Works</a>
            <a className="text-on-surface-variant hover:text-primary" href="#">Solutions</a>
            <a className="text-on-surface-variant hover:text-primary" href="#">Pricing</a>
          </div>
          <div className="flex items-center gap-4">
            <button onClick={() => navigate("/login")} className="hidden sm:block text-on-surface-variant hover:text-primary">Login</button>
            <button onClick={() => navigate("/app")} className="bg-primary text-on-primary px-6 py-2 rounded-lg font-semibold hover:bg-primary-container shadow-sm">Get Started</button>
          </div>
        </div>
      </nav>

      {/* Hero */}
      <section className="pt-32 pb-20 px-container-padding max-w-7xl mx-auto">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          <div className="space-y-8">
            <div className="inline-flex items-center gap-2 px-3 py-1 bg-primary-container text-on-primary-container rounded-full font-label-caps text-label-caps">
              <span className="material-symbols-outlined text-[14px]">verified</span> NEW: V2.0 CLINICAL ENGINE RELEASED
            </div>
            <h1 className="font-display-lg text-display-lg text-on-surface leading-tight">Pharmacy Management <br/> <span className="text-primary">Reimagined.</span></h1>
            <p className="text-on-surface-variant font-body-lg text-body-lg max-w-xl leading-relaxed">PharmaFlow delivers clinical precision and operational efficiency to modern pharmacies. Streamline your workflow, manage complex inventory, and focus on what matters most: patient care.</p>
            <div className="flex flex-wrap gap-4">
              <button onClick={() => navigate("/register")} className="bg-primary text-on-primary px-8 py-4 rounded-lg font-headline-sm text-headline-sm hover:bg-primary-container shadow-md flex items-center gap-2 group">
                Get Started Free
                <span className="material-symbols-outlined group-hover:translate-x-1 transition-transform">arrow_forward</span>
              </button>
              <button className="border border-outline-variant text-primary px-8 py-4 rounded-lg font-headline-sm text-headline-sm hover:bg-surface-container flex items-center gap-2">
                <span className="material-symbols-outlined">play_circle</span> Watch Demo
              </button>
            </div>
            <div className="flex items-center gap-6 pt-4">
              <div className="flex -space-x-3">
                <div className="w-10 h-10 rounded-full border-2 border-white bg-surface-container-high flex items-center justify-center overflow-hidden"><div className="w-full h-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary">JD</div></div>
                <div className="w-10 h-10 rounded-full border-2 border-white bg-surface-container-high flex items-center justify-center overflow-hidden"><div className="w-full h-full bg-secondary/10 flex items-center justify-center text-xs font-bold text-secondary">SM</div></div>
                <div className="w-10 h-10 rounded-full border-2 border-white bg-surface-container-high flex items-center justify-center overflow-hidden"><div className="w-full h-full bg-tertiary/10 flex items-center justify-center text-xs font-bold text-tertiary">RK</div></div>
              </div>
              <p className="text-on-surface-variant">Trusted by <span className="font-bold text-on-surface">2,400+</span> healthcare professionals</p>
            </div>
          </div>
          <div className="relative">
            <div className="absolute -inset-4 bg-primary/10 blur-3xl rounded-full opacity-30"></div>
            <div className="relative bg-white p-2 rounded-xl shadow-2xl border border-outline-variant overflow-hidden">
              <div className="w-full h-auto rounded-lg aspect-video bg-gradient-to-br from-primary/5 to-surface-container flex items-center justify-center">
                <div className="text-center p-8">
                  <span className="material-symbols-outlined text-6xl text-primary/30">dashboard</span>
                  <p className="text-on-surface-variant text-sm mt-4">Dashboard Preview</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Stats */}
      <div className="bg-surface-container py-12 px-container-padding">
        <div className="max-w-7xl mx-auto flex flex-wrap justify-around gap-8 text-center">
          <div><h3 className="font-headline-md text-headline-md text-primary">500+</h3><p className="font-label-caps text-label-caps text-on-surface-variant">Active Pharmacies</p></div>
          <div><h3 className="font-headline-md text-headline-md text-primary">99.9%</h3><p className="font-label-caps text-label-caps text-on-surface-variant">Inventory Accuracy</p></div>
          <div><h3 className="font-headline-md text-headline-md text-primary">$2B+</h3><p className="font-label-caps text-label-caps text-on-surface-variant">Annual Revenue Tracked</p></div>
          <div><h3 className="font-headline-md text-headline-md text-primary">24/7</h3><p className="font-label-caps text-label-caps text-on-surface-variant">Dedicated Support</p></div>
        </div>
      </div>

      {/* Features */}
      <section id="features" className="py-24 px-container-padding max-w-7xl mx-auto">
        <div className="text-center mb-16 space-y-4">
          <h2 className="font-headline-md text-headline-md text-on-surface">Designed for Clinical Precision</h2>
          <p className="text-on-surface-variant font-body-lg text-body-lg max-w-2xl mx-auto">Our platform provides a comprehensive suite of tools built to handle the rigorous demands of pharmaceutical administration.</p>
        </div>
        <div className="grid md:grid-cols-3 gap-gutter">
          <div className="bg-white p-8 rounded-xl shadow-sm border border-outline-variant transition-all hover:shadow-md flex flex-col h-full group">
            <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-6 text-primary group-hover:bg-primary group-hover:text-white transition-colors"><span className="material-symbols-outlined">insights</span></div>
            <h3 className="font-headline-sm text-headline-sm mb-4">Real-time Sales Insights</h3>
            <p className="text-on-surface-variant font-body-md text-body-md mb-8 flex-grow">Track revenue, orders, and customer growth instantly. Our high-density data visualizations provide clear action points for your pharmacy management.</p>
            <div className="pt-4 border-t border-surface-container flex justify-between items-center">
              <span className="text-primary font-bold text-label-caps">ANALYTICS ENGINE</span>
              <span className="px-2 py-1 bg-green-50 text-green-700 text-[10px] rounded font-bold">+12.4%</span>
            </div>
          </div>
          <div className="bg-white p-8 rounded-xl shadow-sm border border-outline-variant transition-all hover:shadow-md flex flex-col h-full group">
            <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-6 text-primary group-hover:bg-primary group-hover:text-white transition-colors"><span className="material-symbols-outlined">inventory</span></div>
            <h3 className="font-headline-sm text-headline-sm mb-4">Smart Inventory Control</h3>
            <p className="text-on-surface-variant font-body-md text-body-md mb-8 flex-grow">Automated stock alerts and batch tracking. Prevent stockouts and minimize wastage with our predictive procurement algorithms.</p>
            <div className="pt-4 border-t border-surface-container flex justify-between items-center">
              <span className="text-primary font-bold text-label-caps">SUPPLY CHAIN</span>
              <span className="px-2 py-1 bg-blue-50 text-blue-700 text-[10px] rounded font-bold">AUTOMATED</span>
            </div>
          </div>
          <div className="bg-white p-8 rounded-xl shadow-sm border border-outline-variant transition-all hover:shadow-md flex flex-col h-full group">
            <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center mb-6 text-primary group-hover:bg-primary group-hover:text-white transition-colors"><span className="material-symbols-outlined">health_and_safety</span></div>
            <h3 className="font-headline-sm text-headline-sm mb-4">Clinical Administration</h3>
            <p className="text-on-surface-variant font-body-md text-body-md mb-8 flex-grow">Secure, precise management for modern pharmacy operations. Compliant documentation and patient record tracking in a unified interface.</p>
            <div className="pt-4 border-t border-surface-container flex justify-between items-center">
              <span className="text-primary font-bold text-label-caps">HIPAA COMPLIANT</span>
              <span className="px-2 py-1 bg-purple-50 text-purple-700 text-[10px] rounded font-bold">ENCRYPTED</span>
            </div>
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section id="how-it-works" className="py-24 bg-surface-container-lowest">
        <div className="max-w-7xl mx-auto px-container-padding">
          <div className="flex flex-col md:flex-row justify-between items-end mb-16 gap-4">
            <div className="max-w-xl">
              <h2 className="font-headline-md text-headline-md mb-4">Simplified Workflow</h2>
              <p className="text-on-surface-variant font-body-lg text-body-lg">Transitioning to PharmaFlow is seamless. Our 3-step integration process gets your team up and running within days, not weeks.</p>
            </div>
            <a className="text-primary font-bold flex items-center gap-2 hover:underline cursor-pointer">Explore Implementation Guide <span className="material-symbols-outlined">chevron_right</span></a>
          </div>
          <div className="grid md:grid-cols-3 gap-12 relative">
            <div className="relative z-10 text-center">
              <div className="w-20 h-20 bg-white shadow-xl rounded-full flex items-center justify-center mx-auto mb-6 border-4 border-surface-container"><span className="material-symbols-outlined text-primary text-3xl">hub</span></div>
              <h4 className="font-headline-sm text-headline-sm mb-2">1. Connect</h4>
              <p className="text-on-surface-variant font-body-md text-body-md">Integrate your existing hardware and patient databases with our API-first platform.</p>
            </div>
            <div className="relative z-10 text-center">
              <div className="w-20 h-20 bg-white shadow-xl rounded-full flex items-center justify-center mx-auto mb-6 border-4 border-surface-container"><span className="material-symbols-outlined text-primary text-3xl">monitoring</span></div>
              <h4 className="font-headline-sm text-headline-sm mb-2">2. Monitor</h4>
              <p className="text-on-surface-variant font-body-md text-body-md">Visualize real-time throughput and inventory levels on our central dashboard.</p>
            </div>
            <div className="relative z-10 text-center">
              <div className="w-20 h-20 bg-white shadow-xl rounded-full flex items-center justify-center mx-auto mb-6 border-4 border-surface-container"><span className="material-symbols-outlined text-primary text-3xl">precision_manufacturing</span></div>
              <h4 className="font-headline-sm text-headline-sm mb-2">3. Optimize</h4>
              <p className="text-on-surface-variant font-body-md text-body-md">Use AI-driven insights to refine stock orders and staffing schedules for maximum ROI.</p>
            </div>
            <div className="hidden md:block absolute top-10 left-[15%] right-[15%] h-px bg-surface-container z-0"></div>
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-24 px-container-padding">
        <div className="max-w-5xl mx-auto bg-primary rounded-2xl p-12 md:p-20 text-center text-on-primary relative overflow-hidden shadow-2xl">
          <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full blur-3xl -mr-32 -mt-32"></div>
          <div className="absolute bottom-0 left-0 w-64 h-64 bg-white/5 rounded-full blur-3xl -ml-32 -mb-32"></div>
          <div className="relative z-10">
            <h2 className="font-display-lg text-display-lg mb-6">Experience PharmaFlow Today</h2>
            <p className="text-on-primary/80 font-body-lg text-body-lg mb-10 max-w-2xl mx-auto">Join over 500 pharmacies already optimizing their clinical operations with our modern management suite. Start your 30-day trial today.</p>
            <div className="flex flex-col sm:flex-row justify-center gap-4">
              <button onClick={() => navigate("/register")} className="bg-white text-primary px-10 py-4 rounded-lg font-headline-sm text-headline-sm font-bold hover:bg-surface-bright shadow-lg">Sign Up Free</button>
              <button className="border border-white/30 text-white px-10 py-4 rounded-lg font-headline-sm text-headline-sm hover:bg-white/10">Schedule a Tour</button>
            </div>
            <p className="mt-8 text-on-primary/60 font-label-caps text-label-caps">NO CREDIT CARD REQUIRED • CANCEL ANYTIME</p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="w-full py-12 px-container-padding flex flex-col md:flex-row justify-between items-center gap-base bg-surface-container border-t border-outline-variant">
        <div className="flex flex-col gap-2 mb-8 md:mb-0">
          <div className="font-headline-sm text-headline-sm font-bold text-on-surface">PharmaFlow</div>
          <p className="font-body-md text-body-md text-on-surface-variant max-w-xs">© 2024 PharmaFlow. Clinical Precision for Modern Pharmacy.</p>
        </div>
        <div className="flex flex-wrap justify-center gap-x-8 gap-y-4">
          <a className="font-body-md text-body-md text-on-surface-variant hover:underline hover:text-primary" href="#">Privacy Policy</a>
          <a className="font-body-md text-body-md text-on-surface-variant hover:underline hover:text-primary" href="#">Terms of Service</a>
          <a className="font-body-md text-body-md text-on-surface-variant hover:underline hover:text-primary" href="#">Security</a>
          <a className="font-body-md text-body-md text-on-surface-variant hover:underline hover:text-primary" href="#">API Documentation</a>
          <a className="font-body-md text-body-md text-on-surface-variant hover:underline hover:text-primary" href="#">Contact Support</a>
        </div>
        <div className="flex gap-4 mt-8 md:mt-0">
          <div className="w-10 h-10 rounded-full bg-surface-container-highest flex items-center justify-center cursor-pointer hover:bg-primary hover:text-white transition-colors"><span className="material-symbols-outlined text-sm">public</span></div>
          <div className="w-10 h-10 rounded-full bg-surface-container-highest flex items-center justify-center cursor-pointer hover:bg-primary hover:text-white transition-colors"><span className="material-symbols-outlined text-sm">mail</span></div>
        </div>
      </footer>
    </div>
  );
}
