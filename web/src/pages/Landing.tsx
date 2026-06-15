import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function Landing() {
  const navigate = useNavigate();

  useEffect(() => {
    const observerOptions = { threshold: 0.1 };
    const observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add("opacity-100");
          entry.target.classList.remove("opacity-0", "translate-y-8");
        }
      });
    }, observerOptions);
    document.querySelectorAll("section").forEach((el) => {
      el.classList.add("transition-all", "duration-700", "opacity-0", "translate-y-8");
      observer.observe(el);
    });
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    const onScroll = () => {
      const nav = document.querySelector("nav");
      if (!nav) return;
      if (window.scrollY > 20) {
        nav.classList.add("shadow-md");
        nav.classList.remove("shadow-sm");
      } else {
        nav.classList.add("shadow-sm");
        nav.classList.remove("shadow-md");
      }
    };
    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  return (
    <div className="bg-background text-on-surface font-body-md overflow-x-hidden">
      {/* TopNavBar */}
      <nav className="fixed top-0 w-full z-50 bg-surface-container-lowest dark:bg-inverse-surface shadow-sm transition-all duration-300">
        <div className="max-w-7xl mx-auto px-container-padding py-4 flex justify-between items-center">
          <div className="font-headline-md text-headline-md font-bold text-primary dark:text-inverse-primary cursor-pointer active:opacity-80">
            PharmD
          </div>
          <div className="hidden md:flex items-center gap-8">
            <a className="text-primary dark:text-inverse-primary font-bold border-b-2 border-primary pb-1 font-body-md text-body-md cursor-pointer transition-colors" href="#features">Features</a>
            <a className="text-on-surface-variant dark:text-surface-variant hover:text-primary dark:hover:text-inverse-primary font-body-md text-body-md transition-colors cursor-pointer" href="#how-it-works">How It Works</a>
            <a className="text-on-surface-variant dark:text-surface-variant hover:text-primary dark:hover:text-inverse-primary font-body-md text-body-md transition-colors cursor-pointer" href="#">Solutions</a>
            <a className="text-on-surface-variant dark:text-surface-variant hover:text-primary dark:hover:text-inverse-primary font-body-md text-body-md transition-colors cursor-pointer" href="#">Pricing</a>
          </div>
          <div className="flex items-center gap-4">
            <button onClick={() => navigate("/login")} className="hidden sm:block text-on-surface-variant font-body-md text-body-md hover:text-primary transition-colors cursor-pointer">Login</button>
            <button onClick={() => navigate("/register")} className="bg-primary text-on-primary px-6 py-2 rounded-lg font-body-md text-body-md font-semibold hover:bg-primary-container transition-all active:opacity-80 shadow-sm">
              Get Started
            </button>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="pt-32 pb-20 px-container-padding max-w-7xl mx-auto overflow-hidden">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          <div className="space-y-8">
            <div className="inline-flex items-center gap-2 px-3 py-1 bg-primary-container text-on-primary-container rounded-full font-label-caps text-label-caps">
              <span className="material-symbols-outlined text-[14px]">verified</span>
              NEW: V2.0 CLINICAL ENGINE RELEASED
            </div>
            <h1 className="font-display-lg text-display-lg text-on-surface leading-tight">
              Pharmacy Management <br/> <span className="text-primary">Reimagined.</span>
            </h1>
            <p className="text-on-surface-variant font-body-lg text-body-lg max-w-xl leading-relaxed">
              PharmD delivers clinical precision and operational efficiency to modern pharmacies. Streamline your workflow, manage complex inventory, and focus on what matters most: patient care.
            </p>
            <div className="flex flex-wrap gap-4">
              <button onClick={() => navigate("/register")} className="bg-primary text-on-primary px-8 py-4 rounded-lg font-headline-sm text-headline-sm hover:bg-primary-container transition-all shadow-md flex items-center gap-2 group">
                Get Started Free
                <span className="material-symbols-outlined group-hover:translate-x-1 transition-transform">arrow_forward</span>
              </button>
              <button className="border border-outline-variant text-primary px-8 py-4 rounded-lg font-headline-sm text-headline-sm hover:bg-surface-container transition-all flex items-center gap-2">
                <span className="material-symbols-outlined fill-icon">play_circle</span>
                Watch Demo
              </button>
            </div>
            <div className="flex items-center gap-6 pt-4">
              <div className="flex -space-x-3">
                <div className="w-10 h-10 rounded-full border-2 border-white bg-surface-container-high flex items-center justify-center overflow-hidden">
                  <img className="w-full h-full object-cover" alt="A professional female pharmacist in a clean white medical coat smiling warmly in a modern clinical setting with soft natural lighting and blurred pharmaceutical shelves in the background, conveying expertise and trust." src="https://lh3.googleusercontent.com/aida-public/AB6AXuDKFmXVZvJhXps4N3BddUptNtXq5PzSe411J_8VvsZcYknKI9pZLCDD0NTPSfsbH5vW2cOa9MWiKpEh0jCL6YmGeR2p7kXNpgQoSTbOE_4zwrOKWsUNinHYsWtS9XS4tfFJ3yW6Zw25_zRviruRjTuIMf9Gf3DtkGNOgJGTHljpfn5E7-k3QDyqFqgXptpUxswpjNsRLUrL7bZ4emHZ8DKwZpcItL0SLXjGMANCoKHsEPvPgO7OzeTj9OyeEnI4JehObcFt1CHQI4E"/>
                </div>
                <div className="w-10 h-10 rounded-full border-2 border-white bg-surface-container-high flex items-center justify-center overflow-hidden">
                  <img className="w-full h-full object-cover" alt="A portrait of a male healthcare professional in a hospital corridor, wearing a light blue stethoscope and a focused yet friendly expression, set against a bright, airy medical environment with professional cool-toned lighting." src="https://lh3.googleusercontent.com/aida-public/AB6AXuD3QHwtCH-b7-E9FwjGMNZirUx7qPVUVgKoveIZypJzyArPW546gpzVhOR_WkOfCUVOkrE2pxjZ6o21KouZBIB4P4QKaZ5XjOL6UxO8EQiDvJpFamlRv_1v5ejG7aETjC777G2po46k9GJyU1aKK4u86ug6c6AIx1TbUezq7hISe-TgJ3aO_arwbv45VDmHLFP6wctq8M3UFYIYUbOBNT2SiWrLMlr_RAOiSZ_IJnHHMn0T39p1-OD1M3KbDpAHP1eHQRSnzNoTP7c"/>
                </div>
                <div className="w-10 h-10 rounded-full border-2 border-white bg-surface-container-high flex items-center justify-center overflow-hidden">
                  <img className="w-full h-full object-cover" alt="A close-up photograph of a medical researcher in a high-tech laboratory setting, wearing a white coat and safety glasses, looking into a microscope with a soft clinical blue and white color palette creating an atmosphere of precision and discovery." src="https://lh3.googleusercontent.com/aida-public/AB6AXuAQCvvNfQgCp2FBNo0O3_7FTbcu5WlnlLF3DNNKXPIIV7Ei25GMGiRtxeHaXhznTJnTj8nuiyW_IqUjIxLDqJK2MTvpVNRT25bNQfOeNDB_-kC97rneUPi164Vg9xuIx4wQ63lDrgpdLfMhNLKW9MzTyVqc17KYXGnzxOzBos6pyMQW7goei2y3XKnUywdy0PkpNNQ9ckbS00gF5kDD67FdX3WbS5xDNO94UQhUKZmqvNvecELXUT8inusX6QjVzMyKXeaQBcRCQM4"/>
                </div>
              </div>
              <p className="text-on-surface-variant font-body-md text-body-md">
                Trusted by <span className="font-bold text-on-surface">2,400+</span> healthcare professionals
              </p>
            </div>
          </div>
          <div className="relative">
            <div className="absolute -inset-4 bg-primary/10 blur-3xl rounded-full opacity-30"></div>
            <div className="relative bg-white p-2 rounded-xl shadow-2xl border border-outline-variant overflow-hidden">
              <img className="w-full h-auto rounded-lg" alt="A high-resolution dashboard mockup for a clinical management software showing detailed data visualizations, revenue charts, and inventory stock alerts on a clean white user interface with vibrant clinical blue accents and sophisticated minimalist design elements." src="https://lh3.googleusercontent.com/aida-public/AB6AXuDDGL0TsaoeuB5LfpKSxpeDZ251MFD8pYyunMB_EPXp5wnrXwfi_lBNanlsoUk3ZJ8iYRXTMg3QMvMj5_S7f2QBd3ZYbJdW34M8U2YWDtPs8IFoDOvxLryhPhZxu8tHlczu_ea5ISG7xhevNKfp33HDUtd-byG-3V9PFNfRv9YIkMKVEQUW-ZlS2b-4sXR2t8AVeP1RjIlthyHiN-a_wx1OoOnaSznd3chuc37THf-8ZenvW1x-KrUkHCQEkRkXgKbJR9wVopNWep4"/>
              {/* Floating UI element */}
                <div className="absolute top-10 -right-8 bg-white p-4 rounded-lg shadow-lg border border-outline-variant animate-bounce-slow">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-green-50 flex items-center justify-center">
                    <span className="material-symbols-outlined text-green-600">inventory_2</span>
                  </div>
                  <div>
                    <p className="font-label-caps text-label-caps text-on-surface-variant">STOCK LEVEL</p>
                    <p className="font-bold text-on-surface">Optimal (98%)</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Stats Bar */}
      <div className="bg-surface-container py-12 px-container-padding">
        <div className="max-w-7xl mx-auto flex flex-wrap justify-around gap-8 text-center">
          <div>
            <h3 className="font-headline-md text-headline-md text-primary">500+</h3>
            <p className="font-label-caps text-label-caps text-on-surface-variant">Active Pharmacies</p>
          </div>
          <div>
            <h3 className="font-headline-md text-headline-md text-primary">99.9%</h3>
            <p className="font-label-caps text-label-caps text-on-surface-variant">Inventory Accuracy</p>
          </div>
          <div>
            <h3 className="font-headline-md text-headline-md text-primary">$2B+</h3>
            <p className="font-label-caps text-label-caps text-on-surface-variant">Annual Revenue Tracked</p>
          </div>
          <div>
            <h3 className="font-headline-md text-headline-md text-primary">24/7</h3>
            <p className="font-label-caps text-label-caps text-on-surface-variant">Dedicated Support</p>
          </div>
        </div>
      </div>

      {/* Key Features Grid */}
      <section className="py-24 px-container-padding max-w-7xl mx-auto">
        <div className="text-center mb-16 space-y-4">
          <h2 className="font-headline-md text-headline-md text-on-surface">Designed for Clinical Precision</h2>
          <p className="text-on-surface-variant font-body-lg text-body-lg max-w-2xl mx-auto">
            Our platform provides a comprehensive suite of tools built to handle the rigorous demands of pharmaceutical administration.
          </p>
        </div>
        <div className="grid md:grid-cols-3 gap-gutter">
          {/* Feature Card 1 */}
          <div className="tonal-elevation-1 tonal-elevation-2 p-8 rounded-xl border border-outline-variant transition-all flex flex-col h-full group">
            <div className="w-12 h-12 rounded-lg bg-primary-container/20 flex items-center justify-center mb-6 text-primary group-hover:bg-primary group-hover:text-white transition-colors">
              <span className="material-symbols-outlined fill-icon">insights</span>
            </div>
            <h3 className="font-headline-sm text-headline-sm mb-4">Real-time Sales Insights</h3>
            <p className="text-on-surface-variant font-body-md text-body-md mb-8 flex-grow leading-relaxed">
              Track revenue, orders, and customer growth instantly. Our high-density data visualizations provide clear action points for your pharmacy management.
            </p>
            <div className="pt-4 border-t border-surface-container-high flex justify-between items-center">
              <span className="text-primary font-bold text-label-caps font-label-caps">ANALYTICS ENGINE</span>
              <span className="px-2 py-1 bg-green-50 text-green-700 text-[10px] rounded font-bold">+12.4%</span>
            </div>
          </div>
          {/* Feature Card 2 */}
          <div className="tonal-elevation-1 tonal-elevation-2 p-8 rounded-xl border border-outline-variant transition-all flex flex-col h-full group">
            <div className="w-12 h-12 rounded-lg bg-primary-container/20 flex items-center justify-center mb-6 text-primary group-hover:bg-primary group-hover:text-white transition-colors">
              <span className="material-symbols-outlined fill-icon">inventory</span>
            </div>
            <h3 className="font-headline-sm text-headline-sm mb-4">Smart Inventory Control</h3>
            <p className="text-on-surface-variant font-body-md text-body-md mb-8 flex-grow leading-relaxed">
              Automated stock alerts and batch tracking. Prevent stockouts and minimize wastage with our predictive procurement algorithms.
            </p>
            <div className="pt-4 border-t border-surface-container-high flex justify-between items-center">
              <span className="text-primary font-bold text-label-caps font-label-caps">SUPPLY CHAIN</span>
              <span className="px-2 py-1 bg-blue-50 text-blue-700 text-[10px] rounded font-bold">AUTOMATED</span>
            </div>
          </div>
          {/* Feature Card 3 */}
          <div className="tonal-elevation-1 tonal-elevation-2 p-8 rounded-xl border border-outline-variant transition-all flex flex-col h-full group">
            <div className="w-12 h-12 rounded-lg bg-primary-container/20 flex items-center justify-center mb-6 text-primary group-hover:bg-primary group-hover:text-white transition-colors">
              <span className="material-symbols-outlined fill-icon">health_and_safety</span>
            </div>
            <h3 className="font-headline-sm text-headline-sm mb-4">Clinical Administration</h3>
            <p className="text-on-surface-variant font-body-md text-body-md mb-8 flex-grow leading-relaxed">
              Secure, precise management for modern pharmacy operations. Compliant documentation and patient record tracking in a unified interface.
            </p>
            <div className="pt-4 border-t border-surface-container-high flex justify-between items-center">
              <span className="text-primary font-bold text-label-caps font-label-caps">HIPAA COMPLIANT</span>
              <span className="px-2 py-1 bg-purple-50 text-purple-700 text-[10px] rounded font-bold">ENCRYPTED</span>
            </div>
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section className="py-24 bg-surface-container-lowest">
        <div className="max-w-7xl mx-auto px-container-padding">
          <div className="flex flex-col md:flex-row justify-between items-end mb-16 gap-4">
            <div className="max-w-xl">
              <h2 className="font-headline-md text-headline-md mb-4">Simplified Workflow</h2>
              <p className="text-on-surface-variant font-body-lg text-body-lg">
                Transitioning to PharmD is seamless. Our 3-step integration process gets your team up and running within days, not weeks.
              </p>
            </div>
            <div className="text-primary font-bold flex items-center gap-2 cursor-pointer hover:underline">
              Explore Implementation Guide <span className="material-symbols-outlined">chevron_right</span>
            </div>
          </div>
          <div className="grid md:grid-cols-3 gap-12 relative">
            {/* Connect */}
            <div className="relative z-10 text-center">
              <div className="w-20 h-20 bg-white shadow-xl rounded-full flex items-center justify-center mx-auto mb-6 border-4 border-surface-container">
                <span className="material-symbols-outlined text-primary text-3xl">hub</span>
              </div>
              <h4 className="font-headline-sm text-headline-sm mb-2">1. Connect</h4>
              <p className="text-on-surface-variant font-body-md text-body-md">Integrate your existing hardware and patient databases with our API-first platform.</p>
            </div>
            {/* Monitor */}
            <div className="relative z-10 text-center">
              <div className="w-20 h-20 bg-white shadow-xl rounded-full flex items-center justify-center mx-auto mb-6 border-4 border-surface-container">
                <span className="material-symbols-outlined text-primary text-3xl">monitoring</span>
              </div>
              <h4 className="font-headline-sm text-headline-sm mb-2">2. Monitor</h4>
              <p className="text-on-surface-variant font-body-md text-body-md">Visualize real-time throughput and inventory levels on our central dashboard.</p>
            </div>
            {/* Optimize */}
            <div className="relative z-10 text-center">
              <div className="w-20 h-20 bg-white shadow-xl rounded-full flex items-center justify-center mx-auto mb-6 border-4 border-surface-container">
                <span className="material-symbols-outlined text-primary text-3xl">precision_manufacturing</span>
              </div>
              <h4 className="font-headline-sm text-headline-sm mb-2">3. Optimize</h4>
              <p className="text-on-surface-variant font-body-md text-body-md">Use AI-driven insights to refine stock orders and staffing schedules for maximum ROI.</p>
            </div>
            {/* Background Connector Line */}
            <div className="hidden md:block absolute top-10 left-[15%] right-[15%] h-1 bg-surface-container z-0"></div>
          </div>
        </div>
      </section>

      {/* Final CTA */}
      <section className="py-24 px-container-padding">
        <div className="max-w-5xl mx-auto bg-primary rounded-2xl p-12 md:p-20 text-center text-on-primary relative overflow-hidden shadow-2xl">
          <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full blur-3xl -mr-32 -mt-32"></div>
          <div className="absolute bottom-0 left-0 w-64 h-64 bg-white/5 rounded-full blur-3xl -ml-32 -mb-32"></div>
          <div className="relative z-10">
            <h2 className="font-display-lg text-display-lg mb-6">Experience PharmD Today</h2>
            <p className="text-on-primary/80 font-body-lg text-body-lg mb-10 max-w-2xl mx-auto leading-relaxed">
              Join over 500 pharmacies already optimizing their clinical operations with our modern management suite. Start your 30-day trial today.
            </p>
            <div className="flex flex-col sm:flex-row justify-center gap-4">
              <button onClick={() => navigate("/register")} className="bg-white text-primary px-10 py-4 rounded-lg font-headline-sm text-headline-sm font-bold hover:bg-surface-bright transition-all shadow-lg">
                Sign Up Free
              </button>
              <button className="border border-white/30 text-white px-10 py-4 rounded-lg font-headline-sm text-headline-sm hover:bg-white/10 transition-all">
                Schedule a Tour
              </button>
            </div>
            <p className="mt-8 text-on-primary/60 font-label-caps text-label-caps">NO CREDIT CARD REQUIRED • CANCEL ANYTIME</p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="w-full py-12 px-container-padding flex flex-col md:flex-row justify-between items-center gap-base bg-surface-container dark:bg-inverse-surface border-t border-outline-variant dark:border-outline">
        <div className="flex flex-col gap-2 mb-8 md:mb-0">
          <div className="font-headline-sm text-headline-sm font-bold text-on-surface dark:text-inverse-on-surface">PharmD</div>
          <p className="font-body-md text-body-md text-on-surface-variant dark:text-surface-variant max-w-xs">© 2024 PharmD. Clinical Precision for Modern Pharmacy.</p>
        </div>
        <div className="flex flex-wrap justify-center gap-x-8 gap-y-4">
          <a className="font-body-md text-body-md text-on-surface-variant dark:text-surface-variant hover:underline hover:text-primary dark:hover:text-inverse-primary transition-all duration-200" href="#">Privacy Policy</a>
          <a className="font-body-md text-body-md text-on-surface-variant dark:text-surface-variant hover:underline hover:text-primary dark:hover:text-inverse-primary transition-all duration-200" href="#">Terms of Service</a>
          <a className="font-body-md text-body-md text-on-surface-variant dark:text-surface-variant hover:underline hover:text-primary dark:hover:text-inverse-primary transition-all duration-200" href="#">Security</a>
          <a className="font-body-md text-body-md text-on-surface-variant dark:text-surface-variant hover:underline hover:text-primary dark:hover:text-inverse-primary transition-all duration-200" href="#">API Documentation</a>
          <a className="font-body-md text-body-md text-on-surface-variant dark:text-surface-variant hover:underline hover:text-primary dark:hover:text-inverse-primary transition-all duration-200" href="#">Contact Support</a>
        </div>
        <div className="flex gap-4 mt-8 md:mt-0">
          <div className="w-10 h-10 rounded-full bg-surface-container-highest flex items-center justify-center cursor-pointer hover:bg-primary hover:text-white transition-colors">
            <span className="material-symbols-outlined text-sm">public</span>
          </div>
          <div className="w-10 h-10 rounded-full bg-surface-container-highest flex items-center justify-center cursor-pointer hover:bg-primary hover:text-white transition-colors">
            <span className="material-symbols-outlined text-sm">mail</span>
          </div>
        </div>
      </footer>
    </div>
  );
}
