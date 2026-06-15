import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { LogOut } from "lucide-react";
const navItems: { to: string; label: string; icon: string }[] = [
  { to: "/app", label: "Dashboard", icon: "dashboard" },
  { to: "/app/users", label: "Users", icon: "group" },
  { to: "/app/roles", label: "Roles", icon: "admin_panel_settings" },
  { to: "/app/permissions", label: "Permissions", icon: "key" },
  { to: "/app/locations", label: "Locations", icon: "location_on" },
  { to: "/app/patients", label: "Patients", icon: "patient_list" },
  { to: "/app/products", label: "Products", icon: "medication" },
  { to: "/app/categories", label: "Categories", icon: "category" },
  { to: "/app/inventory", label: "Inventory", icon: "inventory_2" },
  { to: "/app/suppliers", label: "Suppliers", icon: "local_shipping" },
  { to: "/app/purchases", label: "Purchases", icon: "receipt_long" },
  { to: "/app/pricing", label: "Pricing", icon: "attach_money" },
  { to: "/app/discounts", label: "Discounts", icon: "percent" },
  { to: "/app/prescribers", label: "Prescribers", icon: "stethoscope" },
  { to: "/app/prescriptions", label: "Prescriptions", icon: "description" },
  { to: "/app/dispensing", label: "Dispensing", icon: "pill" },
  { to: "/app/pos", label: "POS", icon: "point_of_sale" },
  { to: "/app/sales", label: "Sales", icon: "bar_chart" },
];

function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

function SidebarIndicator() {
  return (
    <div className="sidebar-indicator bg-primary absolute left-0 top-0 h-full w-1" />
  );
}

export function Layout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  const initials = user
    ? `${user.first_name?.[0] ?? ""}${user.last_name?.[0] ?? ""}`
    : "??";

  const orgInitials = user?.organisation_name
    ? user.organisation_name
        .split(" ")
        .map((w: string) => w[0])
        .join("")
        .slice(0, 2)
        .toUpperCase()
    : "PC";

  return (
    <div className="flex min-h-screen bg-surface font-body-md text-on-surface">
      {/* Sidebar */}
      <aside className="fixed left-0 top-0 z-50 flex h-screen w-[260px] flex-col border-r border-outline-variant bg-surface py-gutter">
        {/* Brand */}
        <div className="mb-10 px-6">
          <h1 className="font-headline-sm text-headline-sm text-primary">
            PharmD
          </h1>
          <p className="font-body-md text-on-surface-variant opacity-70">
            Clinical Admin
          </p>
        </div>

        {/* Navigation */}
        <nav className="flex-1 overflow-y-auto">
          {navItems.map(({ to, label, icon }) => (
            <NavLink
              key={to}
              to={to}
              end={to === "/app"}
              className={({ isActive }) =>
                `relative flex items-center px-6 py-3 transition-colors duration-200 ${
                  isActive
                    ? "bg-surface-container-low font-bold text-primary"
                    : "group text-on-surface-variant hover:bg-surface-container-high hover:text-primary"
                }`
              }
            >
              {({ isActive }) => (
                <>
                  {isActive && <SidebarIndicator />}
                  <Icon
                    name={icon}
                    className={`mr-3 ${isActive ? "" : "group-hover:text-primary"}`}
                  />
                  <span>{label}</span>
                </>
              )}
            </NavLink>
          ))}
        </nav>

        {/* User Profile */}
        <div className="mt-auto px-6">
          <div className="flex items-center rounded-lg bg-surface-container-low p-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full border border-outline-variant bg-surface-container-lowest text-xs font-bold text-primary">
              {initials || orgInitials}
            </div>
            <div className="ml-3 overflow-hidden">
              <p className="truncate text-sm font-bold">
                {user?.first_name} {user?.last_name}
              </p>
              <p className="truncate text-xs text-on-surface-variant">
                {user?.organisation_name ?? "Pharmacy"}
              </p>
            </div>
          </div>
          <button
            onClick={handleLogout}
            className="mt-3 flex w-full items-center justify-center gap-2 rounded-lg border border-outline-variant px-3 py-2 text-sm text-on-surface-variant transition-colors hover:bg-surface-container-high hover:text-primary"
          >
            <LogOut size={16} />
            Logout
          </button>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="ml-[260px] min-h-screen w-full">
        {/* TopNavBar */}
        <header className="fixed right-0 top-0 z-40 flex h-16 w-[calc(100%-260px)] items-center justify-between border-b border-outline-variant bg-surface px-container-padding transition-all duration-200">
          <div className="flex flex-1 items-center">
            {/* <div className="relative w-full max-w-xl">
              <Icon
                name="search"
                className="absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant"
              />
              <input
                className="w-full rounded-lg border border-outline-variant bg-surface-container-low py-2 pl-10 pr-4 text-body-md outline-none transition-all focus:border-primary focus:ring-2 focus:ring-primary"
                placeholder="Search orders, medicines, or patients..."
                type="text"
              />
            </div>*/}
          </div> 
          <div className="flex items-center space-x-6">
            <button className="relative text-on-surface-variant transition-colors hover:text-primary">
              <Icon name="notifications" />
              <span className="absolute right-0 top-0 h-2 w-2 rounded-full border-2 border-surface bg-error" />
            </button>
            <button className="text-on-surface-variant transition-colors hover:text-primary">
              <Icon name="help_outline" />
            </button>
            <div className="h-8 w-px bg-outline-variant" />
            <div className="flex items-center">
              <span className="mr-3 text-sm font-semibold text-on-surface">
                {user?.organisation_name ?? "Pharmacy Central"}
              </span>
              <div className="flex h-8 w-8 items-center justify-center rounded-full border border-outline-variant bg-surface-container text-xs font-bold text-primary">
                {orgInitials}
              </div>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="px-container-padding pb-12 pt-24">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
